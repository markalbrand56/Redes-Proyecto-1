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
	"strings"
)

var (
	TextChannel                 = make(chan string) // Canal para enviar mensajes
	CorrespondentChannel        = make(chan string) // Canal para guardar el receptor del mensaje
	FetchContactsChannel        = make(chan bool)   // Canal para enviar una solicitud de lista de contactos
	SubscribeToChannel          = make(chan string) // Canal para enviar una solicitud de suscripción
	SubscriptionRequestChannel  = make(chan string) // Canal para recibir solicitudes de suscripción
	UnsubscribeFromChannel      = make(chan string) // Canal para enviar una solicitud de cancelación de suscripción
	ConferenceInvitationChannel = make(chan string) // Canal para recibir invitaciones a salas de chat
)

const (
	address = "ws://alumchat.lol:7070/ws"
)

var (
	User       *models.User    // Usuario actual
	AppContext context.Context // Contexto de la aplicación
)

func Start(ctx context.Context, username string, password string) {
	AppContext = ctx
	go startClient(username, password)
}

func Close() {
	close(TextChannel)
	close(CorrespondentChannel)
	close(FetchContactsChannel)
	close(SubscribeToChannel)
	close(SubscriptionRequestChannel)
	close(UnsubscribeFromChannel)
	close(ConferenceInvitationChannel)

	User.SaveConfig()
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
	fmt.Println(User)
	err = User.LoadConfig()
	fmt.Println(User)

	if err != nil {
		log.Println(err)
	}

	startMessaging()
}

// startMessaging es una goroutine que escucha los canales de mensajes y solicitudes de suscripción
func startMessaging() {
	var text string
	var correspondent string

	sendPresence()

	for {
		select {
		// Envío de mensaje a un contacto
		case text = <-TextChannel:
			fmt.Printf("Correspondent: %s Message: %s\n", correspondent, text)
			msg := stanza.Message{Attrs: stanza.Attrs{To: correspondent, Type: stanza.MessageTypeChat}, Body: text}
			err := User.Client.Send(msg)
			if err != nil {
				log.Fatalf("%+v", err)
			}

		// Guardar el receptor del mensaje
		case csrp := <-CorrespondentChannel:
			correspondent = csrp
			fmt.Println("Correspondent: ", correspondent)

		// Obtener la lista de contactos
		case <-FetchContactsChannel:
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
		case u := <-SubscribeToChannel:
			// Para enviar una solicitud de suscripción a un contacto, se debe enviar un mensaje de presencia con el tipo "subscribe"
			fmt.Println("Subscribing to: ", u)
			presence := stanza.Presence{Attrs: stanza.Attrs{To: u, Type: stanza.PresenceTypeSubscribe}}

			err := User.Client.Send(presence)

			if err != nil {
				log.Fatalf("%+v", err)
			}

			runtime.EventsEmit(AppContext, "success", "Subscription request sent")

		// Aceptar solicitud de suscripción
		case u := <-SubscriptionRequestChannel:
			fmt.Println("Subscription (channel) request from: ", u)

			// aceptar la solicitud de suscripción
			presence := stanza.Presence{Attrs: stanza.Attrs{To: u, Type: stanza.PresenceTypeSubscribed}}

			err := User.Client.Send(presence)

			if err != nil {
				log.Fatalf("%+v", err)
			}

			runtime.EventsEmit(AppContext, "success", "Subscription accepted")

		// Cancelar suscripción
		case u := <-UnsubscribeFromChannel:
			// Para cancelar la suscripción a un contacto, se debe enviar un mensaje de presencia con el tipo "unsubscribe"
			fmt.Println("Unsubscribing from: ", u)
			presence := stanza.Presence{Attrs: stanza.Attrs{To: u, Type: stanza.PresenceTypeUnsubscribe}}

			err := User.Client.Send(presence)

			if err != nil {
				log.Fatalf("%+v", err)
			}
		// Invitación a sala de chat
		case jid := <-ConferenceInvitationChannel:
			fmt.Println("Conference invitation from: ", jid)

			alias := User.UserName[:strings.Index(User.UserName, "@")]

			presence := stanza.Presence{
				Attrs: stanza.Attrs{
					From: User.UserName,                    // El JID del usuario actual
					To:   fmt.Sprintf("%s/%s", jid, alias), // El JID del usuario actual con el recurso
					Id:   "join_1",
				},
				Extensions: []stanza.PresExtension{
					&stanza.MucPresence{},
				},
			}

			// Aquí enviamos el `presence` usando el cliente XMPP para unirse a la sala de chat.
			err := User.Client.Send(presence)
			if err != nil {
				fmt.Println("Error al enviar presencia para unirse a la sala de chat:", err)
			} else {
				fmt.Println("Presencia enviada para unirse a la sala de chat:", jid)
			}
		default:
			continue
		}
	}
}

// sendPresence envía una presencia para unirse a las salas de chat a las que pertenece el usuario
func sendPresence() {
	for _, conference := range User.Conferences {
		alias := User.UserName[:strings.Index(User.UserName, "@")]

		presence := stanza.Presence{
			Attrs: stanza.Attrs{
				From: User.UserName,                           // El JID del usuario actual
				To:   fmt.Sprintf("%s/%s", conference, alias), // El JID del usuario actual con el recurso
				Id:   "join_1",
			},
			Extensions: []stanza.PresExtension{
				&stanza.MucPresence{},
			},
		}

		// Aquí enviamos el `presence` usando el cliente XMPP para unirse a la sala de chat.
		err := User.Client.Send(presence)
		if err != nil {
			fmt.Println("Error al enviar presencia para unirse a la sala de chat:", err)
		} else {
			fmt.Println("Presencia enviada para unirse a la sala de chat:", conference)
		}
	}
}
