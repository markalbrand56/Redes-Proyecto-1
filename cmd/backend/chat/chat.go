package chat

import (
	"RedesProyecto/backend/models"
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"log"
	"os"
)

var (
	textChannel                = make(chan string) // Canal para enviar mensajes
	correspondentChannel       = make(chan string) // Canal para guardar el receptor del mensaje
	RequestContactChannel      = make(chan bool)   // Canal para enviar una solicitud de lista de contactos
	subscribeToChannel         = make(chan string) // Canal para enviar una solicitud de suscripción
	subscriptionRequestChannel = make(chan string) // Canal para recibir solicitudes de suscripción
)

const (
	address = "ws://alumchat.lol:7070/ws"
)

var User *models.User
var AppContext context.Context

func Start(ctx context.Context, username string, password string) {
	AppContext = ctx
	go startClient(username, password)
}

func startClient(username string, password string) {
	config := xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: address,
		},
		Jid:          username,
		Credential:   xmpp.Password(password),
		StreamLogger: os.Stdout,
		Insecure:     true,
	}

	router := xmpp.NewRouter()
	router.HandleFunc("message", handleMessage)
	router.HandleFunc("presence", handlePresence)
	router.HandleFunc("iq", handleIQ)

	newClient, err := xmpp.NewClient(&config, router, errorHandler)

	if err != nil {
		log.Fatalf("%+v", err)
	}

	defer func(client *xmpp.Client) {
		if client != nil {
			log.Println("closing client...")
			err := client.Disconnect()
			if err != nil {
				log.Fatalf("%+v", err)
			}
		}
	}(newClient)

	err = newClient.Connect()

	if err != nil {
		log.Fatalf("%+v", err)
	}

	User = models.NewUser(newClient, username)
	startMessaging()
}

// startMessaging es una goroutine que escucha los canales de mensajes y solicitudes de suscripción
func startMessaging() {
	var text string
	var correspondent string

	for {
		select {
		// Envío de mensaje a un contacto
		case text = <-textChannel:
			fmt.Printf("Correspondent: %s Message: %s\n", correspondent, text)
			msg := stanza.Message{Attrs: stanza.Attrs{To: correspondent, Type: stanza.MessageTypeChat}, Body: text}
			err := User.Client.Send(msg)
			if err != nil {
				log.Fatalf("%+v", err)
			}

		// Guardar el receptor del mensaje
		case csrp := <-correspondentChannel:
			correspondent = csrp
			fmt.Println("Correspondent: ", correspondent)

		// Obtener la lista de contactos
		case <-RequestContactChannel:
			// Para obtener la lista de contactos, se debe enviar una solicitud IQ de tipo "get"
			req, err := stanza.NewIQ(
				stanza.Attrs{
					From: User.UserName,
					Type: stanza.IQTypeGet,
					Id:   "roster_1",
				},
			)

			if err != nil {
				log.Fatalf("%+v", err)
			}

			req.RosterItems()

			c, err := User.Client.SendIQ(AppContext, req)

			if err != nil {
				log.Fatalf("%+v", err)
			}

			// Para obtener la respuesta del servidor, Client.SendIQ() devuelve un canal de respuesta que se debe escuchar.
			// Se usa una goroutine para no bloquear el flujo principal, y poder esperar la respuesta del servidor con los contactos
			go func() {
				serverResp := <-c

				if rosterItems, ok := serverResp.Payload.(*stanza.RosterItems); ok {
					contacts := make([]string, 0)

					for _, item := range rosterItems.Items {
						contacts = append(contacts, item.Jid)
					}

					User.Contacts = contacts

					fmt.Println("Contacts: ", contacts)
					runtime.EventsEmit(AppContext, "contacts", contacts)
				}
			}()

		// Suscripción a un contacto
		case u := <-subscribeToChannel:
			// Para enviar una solicitud de suscripción a un contacto, se debe enviar un mensaje de presencia con el tipo "subscribe"
			fmt.Println("Subscribing to: ", u)
			presence := stanza.Presence{Attrs: stanza.Attrs{To: u, Type: stanza.PresenceTypeSubscribe}}

			err := User.Client.Send(presence)

			if err != nil {
				log.Fatalf("%+v", err)
			}

			runtime.EventsEmit(AppContext, "success", "Subscription request sent")

		// Aceptar solicitud de suscripción
		case u := <-subscriptionRequestChannel:
			fmt.Println("Subscription (channel) request from: ", u)

			// aceptar la solicitud de suscripción
			presence := stanza.Presence{Attrs: stanza.Attrs{To: u, Type: stanza.PresenceTypeSubscribed}}

			err := User.Client.Send(presence)

			if err != nil {
				log.Fatalf("%+v", err)
			}

			runtime.EventsEmit(AppContext, "success", "Subscription accepted")

		default:
			continue
		}
	}
}

// SetCorrespondent guarda el receptor del mensaje en correspondentChannel
func SetCorrespondent(correspondent string) {
	correspondentChannel <- correspondent
}

// SendMessage envía un mensaje al servidor, usando el canal textChannel y el receptor guardado en correspondentChannel
func SendMessage(message string) {
	textChannel <- message
}

// FetchContacts envía una solicitud IQ al servidor para obtener la lista de contactos y actualizar la lista de contactos del usuario
func FetchContacts() {
	RequestContactChannel <- true
}

// RequestContact envía una solicitud de suscripción a un contacto
func RequestContact(username string) {
	subscribeToChannel <- username
}

func AcceptSubscription(username string) {
	subscriptionRequestChannel <- username
}
