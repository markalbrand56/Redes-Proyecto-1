package chat

import (
	"RedesProyecto/backend/chat/events"
	"RedesProyecto/backend/models"
	cstanza "RedesProyecto/backend/models/stanza"
	"context"
	"encoding/xml"
	"fmt"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"log"
	"os"
	"strings"
)

var (
	LogoutChannel         = make(chan bool)           // Canal para cerrar la sesión
	TextChannel           = make(chan models.Message) // Canal para enviar mensajes
	ConferenceTextChannel = make(chan models.Message) // Canal para enviar mensajes a salas de chat

	FetchContactsChannel = make(chan bool) // Canal para enviar una solicitud de lista de contactos
	ProbeContactsChannel = make(chan bool) // Canal para enviar una solicitud de estado de contactos

	SubscribeToChannel         = make(chan string) // Canal para enviar una solicitud de suscripción
	SubscriptionRequestChannel = make(chan string) // Canal para recibir solicitudes de suscripción
	UnsubscribeFromChannel     = make(chan string) // Canal para enviar una solicitud de cancelación de suscripción

	ConferenceInvitationChannel = make(chan string) // Canal para recibir invitaciones a salas de chat

	StatusChannel = make(chan string) // Canal para cambiar el estado del usuario

	FetchArchiveChannel = make(chan string) // Canal para solicitar mensajes archivados
)

const (
	address = "ws://alumchat.lol:7070/ws"
)

var (
	User       *models.User    // Usuario actual
	AppContext context.Context // Contexto de la aplicación
)

func Start(ctx context.Context, username string, password string) bool {
	AppContext = ctx
	go startClient(username, password)

	return true
}

func Close() {
	close(TextChannel)
	close(FetchContactsChannel)
	close(SubscribeToChannel)
	close(SubscriptionRequestChannel)
	close(UnsubscribeFromChannel)
	close(ConferenceInvitationChannel)

	err := User.SaveConfig()
	if err != nil {
		log.Println("Error saving user configuration: ", err)
	}
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
		// Login failed
		log.Printf("%+v\n", err)

		events.EmitError(AppContext, "Login failed")
		return
	}

	err = newClient.Connect()

	if err != nil {
		log.Printf("%+v\n", err)
		events.EmitLoginError(AppContext, "Connection failed")
		return
	}

	User = models.NewUser(newClient, username)
	fmt.Println(User)
	err = User.LoadConfig()

	if err != nil {
		log.Println(err)
	}
	fmt.Println(User)

	// Se envía un mensaje de presencia para indicar que el usuario está en línea
	presence := stanza.Presence{
		Attrs: stanza.Attrs{
			From: User.UserName,
		},
		Show: User.Show,
	}

	err = User.Client.Send(presence)

	if err != nil {
		log.Println("Error sending presence: ", err)
		events.EmitError(AppContext, "Error sending presence")
	}

	events.EmitLogin(AppContext, User.UserName)

	startMessaging()
}

// startMessaging es una goroutine que escucha los canales de mensajes y solicitudes de suscripción
func startMessaging() {
	sendPresence(User.Show)
	getArchivedMessages(User.UserName)

	listening := true
	for listening {
		select {
		case <-LogoutChannel:
			// Cerrar la sesión
			err := User.Client.Disconnect()
			if err != nil {
				return // No se pudo cerrar la sesión
			}
			User = nil

			events.EmitLogout(AppContext)
			listening = false

		case msg := <-TextChannel:
			// Envío de mensaje a un contacto V2
			fmt.Printf("Correspondent: %s Message: %s\n", msg.To, msg.Body)
			message := stanza.Message{
				Attrs: stanza.Attrs{
					To:   msg.To,
					Type: stanza.MessageTypeChat,
					From: User.UserName,
				}, Body: msg.Body,
			}
			err := User.Client.Send(message)

			if err != nil {
				log.Println("Error sending message: ", err)
				events.EmitError(AppContext, "Error sending message")
				continue
			}

			User.InsertMessage(msg)
			events.EmitSuccess(AppContext, "Message sent")

		case msg := <-ConferenceTextChannel:
			// Envío de mensaje a una sala de chat
			fmt.Printf("Conference: %s Message: %s\n", msg.To, msg.Body)
			message := stanza.Message{Attrs: stanza.Attrs{To: msg.To, Type: stanza.MessageTypeGroupchat}, Body: msg.Body}
			err := User.Client.Send(message)

			if err != nil {
				log.Println("Error sending message: ", err)
				events.EmitError(AppContext, "Error sending message")
			}

		// Obtener la lista de contactos
		case <-FetchContactsChannel:
			// Para obtener la lista de contactos, se debe enviar una solicitud IQ de tipo "get"
			contactReq, err := stanza.NewIQ(
				stanza.Attrs{
					From: User.UserName,
					Type: stanza.IQTypeGet,
					Id:   "roster_1",
				},
			)

			if err != nil {
				log.Println(err)
				events.EmitError(AppContext, "Error fetching contacts")
			}

			contactReq.RosterItems() // Obtiene los contactos

			contactsResp, err := User.Client.SendIQ(AppContext, contactReq)

			if err != nil {
				log.Println(err)
				events.EmitError(AppContext, "Error fetching contacts")
			}

			// Para obtener la respuesta del servidor, Client.SendIQ() devuelve un canal de respuesta que se debe escuchar.
			// Se usa una goroutine para no bloquear el flujo principal, y poder esperar la respuesta del servidor con los contactos
			go func() {
				serverResp := <-contactsResp

				if rosterItems, ok := serverResp.Payload.(*stanza.RosterItems); ok {
					contacts := make([]string, 0)

					for _, item := range rosterItems.Items {
						contacts = append(contacts, item.Jid)
					}

					User.Contacts = contacts

					fmt.Println("Contacts: ", contacts)
					events.EmitContacts(AppContext, contacts)
					err = User.SaveConfig()

					if err != nil {
						log.Println(err)
					}
				}
			}()

			// salas de chat
			discoReq, err := stanza.NewIQ(
				stanza.Attrs{
					From: User.UserName,
					Type: stanza.IQTypeGet,
					Id:   "disco_1",
					To:   "conference.alumchat.lol",
				},
			)

			if err != nil {
				log.Println(err)
				events.EmitError(AppContext, "Error fetching conferences")
			}

			discoReq.Payload = &stanza.DiscoItems{
				XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#items", Local: "query"},
			}

			discoResp, err := User.Client.SendIQ(AppContext, discoReq)

			if err != nil {
				log.Println(err)
				events.EmitError(AppContext, "Error fetching conferences")
			}

			go func() {
				serverResp := <-discoResp

				if discoItems, ok := serverResp.Payload.(*stanza.DiscoItems); ok {
					fmt.Println("Found disco items")

					for _, item := range discoItems.Items {
						//User.InsertConference(models.NewConference(item.Name, item.JID))
						// Si la sala de chat no está en la lista de salas de chat del usuario, se agrega
						if _, ok := User.Conferences[item.JID]; !ok {
							User.InsertConference(models.NewConference(item.Name, item.JID))
						}
					}

					sendPresence(User.Show)

					events.EmitConferences(AppContext, User.Conferences)

					err = User.SaveConfig()

					if err != nil {
						log.Println(err)
					}
				}
			}()

		// Solicitar el estado de los contactos
		case <-ProbeContactsChannel:
			// Se envía una solicitud de presencia a los contactos para obtener su estado actual
			log.Println("Probing contacts...")
			for _, contact := range User.Contacts {
				presence := stanza.Presence{
					Attrs: stanza.Attrs{
						To:   contact,
						From: User.UserName,
						Type: stanza.PresenceTypeProbe,
					},
				}

				err := User.Client.Send(presence)

				if err != nil {
					log.Println("Error sending presence probe: ", err)
				}
			}

		// Suscripción a un contacto
		case u := <-SubscribeToChannel:
			// Para enviar una solicitud de suscripción a un contacto, se debe enviar un mensaje de presencia con el tipo "subscribe"
			fmt.Println("Subscribing to: ", u)
			presence := stanza.Presence{Attrs: stanza.Attrs{To: u, Type: stanza.PresenceTypeSubscribe}}

			err := User.Client.Send(presence)

			if err != nil {
				log.Println("Error sending subscription request: ", err)
				events.EmitError(AppContext, "Error sending subscription request")
			}

			userFormatted := strings.Split(u, "/")[0]

			User.Contacts = append(User.Contacts, userFormatted)

			events.EmitSuccess(AppContext, "Subscription request sent")

		// Aceptar solicitud de suscripción
		case u := <-SubscriptionRequestChannel:
			// Se ha decidido aceptar la solicitud de suscripción
			log.Println("Accepting subscription from: ", u)
			presence := stanza.Presence{
				Attrs: stanza.Attrs{
					To: u, Type: stanza.PresenceTypeSubscribed,
				},
			}

			err := User.Client.Send(presence)

			if err != nil {
				log.Println("Error accepting subscription: ", err)
				events.EmitError(AppContext, "Error accepting subscription")
				continue
			}

			userFormatted := strings.Split(u, "/")[0]

			User.Contacts = append(User.Contacts, userFormatted)

			events.EmitSuccess(AppContext, "Subscription accepted")

		// Cancelar suscripción
		case u := <-UnsubscribeFromChannel:
			// Para cancelar la suscripción a un contacto, se debe enviar un mensaje de presencia con el tipo "unsubscribe"
			fmt.Println("Unsubscribing from: ", u)
			presence := stanza.Presence{
				Attrs: stanza.Attrs{
					To:   u,
					Type: stanza.PresenceTypeUnsubscribed,
				},
			}

			err := User.Client.Send(presence)

			if err != nil {
				log.Println("Error sending unsubscription request: ", err)
				events.EmitError(AppContext, "Error sending unsubscription request")
			}

			// Remove from roster
			rosterRemove := cstanza.NewRosterRemove(u)

			err = User.Client.Send(&rosterRemove)

			if err != nil {
				log.Println("Error removing contact from roster: ", err)
				events.EmitError(AppContext, "Error removing contact from roster")
			}

		// Invitación a sala de chat
		case jid := <-ConferenceInvitationChannel:
			// Se ha aceptado una invitación a una sala de chat
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

				User.InsertConference(models.NewConference(alias, jid))
			}

		// Cambiar el estado del usuario
		case status := <-StatusChannel:
			// Para cambiar el estado del usuario, se debe enviar un mensaje de presencia con el estado deseado
			presence := stanza.Presence{
				Attrs: stanza.Attrs{
					From: User.UserName,
				},
			}

			switch status {
			case models.StatusOnline:
				log.Println("Setting status to online")
				presence.Show = stanza.PresenceShowChat // <show>chat</show>

			case models.StatusAway:
				log.Println("Setting status to away")
				presence.Show = stanza.PresenceShowAway // <show>away</show>

			case models.StatusNotAvailable:
				log.Println("Setting status to not available")
				presence.Show = stanza.PresenceShowXA // <show>xa</show>

			case models.StatusOffline:
				log.Println("Setting status to offline")
				presence.Type = stanza.PresenceTypeUnavailable

			case models.StatusBusy:
				log.Println("Setting status to busy")
				presence.Show = stanza.PresenceShowDND // <show>dnd</show>

			default:
				log.Println("Setting status to online")
				presence.Show = stanza.PresenceShowChat
			}

			err := User.Client.Send(presence)
			sendPresence(presence.Show)
			User.Show = presence.Show

			if err != nil {
				log.Println("Error al enviar presencia para cambiar el estado del usuario:", err)
			}

		// Obtener mensajes archivados
		case username := <-FetchArchiveChannel:
			getArchivedMessages(username)
		default:
			continue
		}
	}

	log.Println("Closing XMPP client")
}

// sendPresence envía una presencia para unirse a las salas de chat a las que pertenece el usuario
func sendPresence(pres stanza.PresenceShow) {
	for jid, name := range User.Conferences {
		log.Println("Joining conference: ", name)
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
			Show: pres,
		}

		// Aquí enviamos el `presence` usando el cliente XMPP para unirse a la sala de chat.
		err := User.Client.Send(presence)
		if err != nil {
			fmt.Println("Error al enviar presencia para unirse a la sala de chat:", err)
		} else {
			fmt.Println("Presencia enviada para unirse a la sala de chat:", name)
		}
	}
}

func getArchivedMessages(jid string) {
	log.Println("Getting archived messages...", jid, User.UserName)
	// Para obtener los mensajes archivados, se debe enviar una solicitud IQ de tipo "get"
	archiveQuery := cstanza.NewArchiveQuery(jid, 500)

	iq := stanza.IQ{
		Attrs: stanza.Attrs{
			Type: stanza.IQTypeSet,
			Id:   "mam_query_1",
			To:   User.UserName,
		},
		Payload: archiveQuery,
	}

	_, err := User.Client.SendIQ(AppContext, &iq)

	if err != nil {
		log.Printf("Error sending IQ: %+v\n", err)
		events.EmitError(AppContext, "Error fetching archived messages")
	}

	User.Messages = make(map[string][]models.Message)

}

/*
TODO Leído de mensajes

<message xmlns="jabber:client" to="alb21004@alumchat.lol" type="chat" id="e8427b6d-cc79-4771-aaeb-a94cc5dcacc3" from="alb21005@alumchat.lol/gajim.0O3D5ZZ0"><origin-id xmlns="urn:xmpp:sid:0" id="e8427b6d-cc79-4771-aaeb-a94cc5dcacc3"></origin-id><displayed xmlns="urn:xmpp:chat-markers:0" id="837ca35e-53ae-4684-bc63-f2068d188745"></displayed></message>

Esta stanza llega cuando un mensaje ha sido leído por el destinatario. Se debe marcar el mensaje como leído en la interfaz de usuario.
*/
