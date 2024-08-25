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
	LogoutChannel        = make(chan bool) // Canal para cerrar la sesión
	DeleteAccountChannel = make(chan bool) // Canal para eliminar la cuenta

	TextChannel = make(chan models.Message) // Canal para enviar mensajes
	FileChannel = make(chan models.Message) // Canal para enviar archivos

	ConferenceTextChannel = make(chan models.Message) // Canal para enviar mensajes a salas de chat
	ConferenceFileChannel = make(chan models.Message) // Canal para enviar archivos a salas de chat

	FetchContactsChannel = make(chan bool) // Canal para enviar una solicitud de lista de contactos
	ProbeContactsChannel = make(chan bool) // Canal para enviar una solicitud de estado de contactos

	SubscribeToChannel         = make(chan string) // Canal para enviar una solicitud de suscripción
	SubscriptionRequestChannel = make(chan string) // Canal para recibir solicitudes de suscripción
	UnsubscribeFromChannel     = make(chan string) // Canal para enviar una solicitud de cancelación de suscripción (eliminar contacto)
	RejectionChannel           = make(chan string) // Canal para rechazar una solicitud de suscripción

	ConferenceInvitationChannel        = make(chan string) // Canal para recibir invitaciones a salas de chat
	InviteToConferenceChannel          = make(chan models.Invitation)
	ConferenceDeclineInvitationChannel = make(chan models.Invitation) // Canal para rechazar una invitación a sala de chat
	NewConferenceChannel               = make(chan models.Conference)
	DeleteConferenceChannel            = make(chan string)

	ShowChannel   = make(chan string) // Canal para cambiar el estado del usuario
	StatusChannel = make(chan string) // Canal para cambiar el mensaje de estado del usuario

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
	close(LogoutChannel)
	close(DeleteAccountChannel)

	close(TextChannel)
	close(FileChannel)

	close(ConferenceTextChannel)
	close(ConferenceFileChannel)

	close(FetchContactsChannel)
	close(ProbeContactsChannel)

	close(SubscribeToChannel)
	close(SubscriptionRequestChannel)
	close(UnsubscribeFromChannel)
	close(RejectionChannel)

	close(ConferenceInvitationChannel)
	close(InviteToConferenceChannel)

	close(ShowChannel)

	close(FetchArchiveChannel)

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
		Show:   User.Show,
		Status: User.Status,
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
	sendPresence(User.Show, User.Status)
	getArchivedMessages(User.UserName)

	listening := true
	for listening {
		select {

		// Cerrar la sesión
		case <-LogoutChannel:
			err := User.Client.Disconnect()
			if err != nil {
				return // No se pudo cerrar la sesión
			}
			User = nil

			events.EmitLogout(AppContext)
			listening = false

		// Eliminar la cuenta
		case <-DeleteAccountChannel:
			log.Println("Deleting account ", User.UserName)
			cr := cstanza.NewCancelRegistration(User.UserName)

			err := User.Client.Send(&cr)

			if err != nil {
				log.Println("Error deleting account: ", err)
				events.EmitError(AppContext, "Error deleting account")
				continue
			}

			err = User.Client.Disconnect()

			if err != nil {
				return // No se pudo cerrar la sesión
			}

			events.EmitLogout(AppContext)
			listening = false

		// Envío de mensaje a un contacto (chat)
		case msg := <-TextChannel:
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

		// Envío de archivo a un contacto
		case msg := <-FileChannel:

			/*
				ejemplo de como se hace en js, replicar acá
				        if (isFile) {
				            const x = xml('x', { xmlns: 'jabber:x:oob' }, xml('url', {}, message));
				            messageXML.append(x);
			*/

			fmt.Printf("Correspondent: %s File: %s\n", msg.To, msg.Body)
			message := stanza.Message{
				Attrs: stanza.Attrs{
					To:   msg.To,
					Type: stanza.MessageTypeChat,
					From: User.UserName,
				},
				Body: msg.Body,
			}

			f := cstanza.NewFile(msg.Body)

			message.Extensions = append(message.Extensions, f)

			err := User.Client.Send(message)

			if err != nil {
				log.Println("Error sending file: ", err)
				events.EmitError(AppContext, "Error sending file")
				continue
			}

			User.InsertMessage(msg)
			events.EmitSuccess(AppContext, "Message sent")

		// Envío de mensaje a una sala de chat
		case msg := <-ConferenceTextChannel:
			fmt.Printf("Conference: %s Message: %s\n", msg.To, msg.Body)
			message := stanza.Message{Attrs: stanza.Attrs{To: msg.To, Type: stanza.MessageTypeGroupchat}, Body: msg.Body}
			err := User.Client.Send(message)

			if err != nil {
				log.Println("Error sending message: ", err)
				events.EmitError(AppContext, "Error sending message")
				continue
			}

		// Envío de archivo a una sala de chat
		case msg := <-ConferenceFileChannel:
			fmt.Printf("Conference: %s File: %s\n", msg.To, msg.Body)
			message := stanza.Message{
				Attrs: stanza.Attrs{
					To:   msg.To,
					Type: stanza.MessageTypeGroupchat,
					From: User.UserName,
				},
				Body: msg.Body,
			}

			f := cstanza.NewFile(msg.Body)

			message.Extensions = append(message.Extensions, f)

			err := User.Client.Send(message)

			if err != nil {
				log.Println("Error sending file: ", err)
				events.EmitError(AppContext, "Error sending file")
				continue
			}

			User.InsertMessage(msg)
			events.EmitSuccess(AppContext, "Message sent")

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
						log.Println("Conference: ", item.Name, item.JID)
						if _, ok := User.Conferences[item.JID]; !ok {
							log.Println("Inserting conference: ", item.Name)
							User.InsertConference(models.NewConference(item.Name, item.JID))
						} else {
							// Si la sala de chat ya está en la lista de salas de chat del usuario, se actualiza el alias
							User.Conferences[item.JID].Alias = item.Name
						}
					}

					sendPresence(User.Show, User.Status)

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
					continue
				}
			}

		// Suscripción a un contacto (solitud de contacto)
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

		// Aceptar solicitud de suscripción (agregar contacto)
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

			events.EmitContacts(AppContext, User.Contacts)

		// Cancelar suscripción a un contacto (eliminar contacto)
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
				continue
			}

			// Remove from roster
			rosterRemove := cstanza.NewRosterRemove(u)

			err = User.Client.Send(&rosterRemove)

			if err != nil {
				log.Println("Error removing contact from roster: ", err)
				events.EmitError(AppContext, "Error removing contact from roster")
				continue
			}

			userFormatted := strings.Split(u, "/")[0]

			events.EmitSuccess(AppContext, "Unsubscribed from: "+userFormatted)

		// Rechazar solicitud de suscripción (rechazar contacto)
		case u := <-RejectionChannel:
			// Rechazar solicitud de suscripción
			log.Println("Rejecting subscription from: ", u)
			presence := stanza.Presence{
				Attrs: stanza.Attrs{
					To:   u,
					Type: stanza.PresenceTypeUnsubscribe,
				},
			}

			err := User.Client.Send(presence)

			if err != nil {
				log.Println("Error rejecting subscription: ", err)
				events.EmitError(AppContext, "Error rejecting subscription")
				continue
			}

			events.EmitSuccess(AppContext, "Subscription rejected")

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
				events.EmitError(AppContext, "Error joining conference")
			} else {
				fmt.Println("Presencia enviada para unirse a la sala de chat:", jid)
			}

		// Invitar a sala de chat (enviar invitación a un contacto)
		case invite := <-InviteToConferenceChannel:
			log.Println("Inviting to conference: ", invite)

			// Primero se debe enviar una solicitud de afiliación a la sala de chat
			// Sin esto, el usuario no podrá unirse a la sala de chat a menos que ya sea pública
			affiliationRequest := cstanza.NewMUCAffiliationRequest(invite.To, "member")

			affiliationRequestMessage := stanza.IQ{
				Attrs: stanza.Attrs{
					To:   invite.ConferenceJID,
					From: User.UserName,
					Type: stanza.IQTypeSet,
					Id:   "affiliation_request_1",
				},
				Payload: affiliationRequest,
			}

			_, err := User.Client.SendIQ(AppContext, &affiliationRequestMessage)

			if err != nil {
				log.Println("Error sending affiliation request: ", err)
				events.EmitError(AppContext, "Error sending affiliation request")
			}

			// Luego se debe enviar una invitación a la sala de chat
			inviteMessage := stanza.Message{
				Attrs: stanza.Attrs{
					To:   invite.To,
					From: User.UserName,
				},
			}

			inviteExtension := cstanza.NewConferenceInvite(invite.ConferenceJID, fmt.Sprintf("%s has invited you to the conference '%s'", User.UserName, invite.ConferenceJID))
			mucInviteExtension := cstanza.NewMucInvite(User.UserName)

			inviteMessage.Extensions = append(inviteMessage.Extensions, inviteExtension)
			inviteMessage.Extensions = append(inviteMessage.Extensions, mucInviteExtension)

			err = User.Client.Send(inviteMessage)

			if err != nil {
				log.Println("Error sending conference invitation: ", err)
				events.EmitError(AppContext, "Error sending conference invitation")
				continue
			}

			events.EmitSuccess(AppContext, "Conference invitation sent")

		// Rechazar invitación a sala de chat
		case inv := <-ConferenceDeclineInvitationChannel:
			log.Println("Declining conference invitation: ", inv)
			// Se ha rechazado una invitación a una sala de chat

			/*
				<message xmlns="jabber:client" from="testingv2@conference.alumchat.lol" to="alb210041@alumchat.lol"><x xmlns="http://jabber.org/protocol/muc#user"><decline from="alb21005@alumchat.lol"/></x></message>
				message.from = inv.ConferenceJID
				message.to = inv.To  // El dueño de la sala de chat / quien envió la invitación
				message.x.decline.from = User.UserName  // El usuario que rechaza la invitación
			*/

			decline := cstanza.NewConferenceDeclineMessage(inv.ConferenceJID, User.UserName, inv.To)

			err := User.Client.Send(decline)

			if err != nil {
				log.Println("Error sending conference decline: ", err)
				events.EmitError(AppContext, "Error sending conference decline")
				continue
			}

			// Eliminar la sala de chat del disco items del servidor
			// Mandar unsubscribe a la sala de chat
			preference := stanza.Presence{
				Attrs: stanza.Attrs{
					To:   inv.ConferenceJID,
					From: User.UserName,
					Type: stanza.PresenceTypeUnsubscribe,
				},
			}

			err = User.Client.Send(preference)

			if err != nil {
				log.Println("Error sending conference unsubscribe: ", err)
				events.EmitError(AppContext, "Error sending conference unsubscribe")
				continue
			}

		// Crear sala de chat
		case conference := <-NewConferenceChannel:
			log.Println("Creating conference: ", conference.Alias)

			// Se envía un mensaje de presencia para crear la sala de chat
			p := stanza.Presence{
				Attrs: stanza.Attrs{
					From: User.UserName,
					To:   fmt.Sprintf("%s/%s", conference.JID, User.UserName[:strings.Index(User.UserName, "@")]), // JID de la sala de chat + alias del usuario
				},
				Extensions: []stanza.PresExtension{
					&stanza.MucPresence{},
				},
			}

			err := User.Client.Send(p)

			if err != nil {
				log.Println("Error creating conference: ", err)
				events.EmitError(AppContext, "Error creating conference")
				continue
			}

		// Eliminar sala de chat
		case jid := <-DeleteConferenceChannel:
			log.Println("Deleting conference: ", jid)

			jid = strings.Split(jid, "/")[0]

			iq := stanza.IQ{
				Attrs: stanza.Attrs{
					From: User.UserName,
					Type: stanza.IQTypeSet,
					Id:   "destroy_1",
					To:   jid,
				},
				Payload: cstanza.NewRoomDestruction(jid, "Room destruction request"),
			}

			res, err := User.Client.SendIQ(AppContext, &iq)

			if err != nil {
				log.Println("Error destroying conference: ", err)
				events.EmitError(AppContext, "Error destroying conference")
				continue
			}

			go func() {
				serverResp := <-res

				if serverResp.Type == stanza.IQTypeResult {
					log.Println("Conference destroyed: ", jid)
					events.EmitSuccess(AppContext, "Conference destroyed")
				} else {
					log.Println("Error destroying conference: ", serverResp.Error)
					events.EmitError(AppContext, "Error destroying conference")
				}
			}()

		// Cambiar el estado del usuario
		case show := <-ShowChannel:
			// Para cambiar el estado del usuario, se debe enviar un mensaje de presencia con el estado deseado
			presence := stanza.Presence{
				Attrs: stanza.Attrs{
					From: User.UserName,
				},
			}

			switch show {
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
			sendPresence(presence.Show, User.Status)
			User.Show = presence.Show

			if err != nil {
				log.Println("Error al enviar presencia para cambiar el estado del usuario:", err)
			}

		// Cambiar el mensaje de estado del usuario
		case status := <-StatusChannel:
			// Para cambiar el mensaje de estado del usuario, se debe enviar un mensaje de presencia con el mensaje deseado
			presence := stanza.Presence{
				Attrs: stanza.Attrs{
					From: User.UserName,
				},
				Status: status,
				Show:   User.Show,
			}

			err := User.Client.Send(presence)

			if err != nil {
				log.Println("Error al enviar presencia para cambiar el mensaje de estado del usuario:", err)
				events.EmitError(AppContext, "Error changing status message")
				continue
			}

			User.Status = status
			sendPresence(User.Show, status)

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
func sendPresence(pres stanza.PresenceShow, status string) {
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
			Show:   pres,   // Estado de presencia
			Status: status, // Mensaje de estado
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
