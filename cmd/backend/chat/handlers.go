package chat

import (
	"RedesProyecto/backend/chat/events"
	"RedesProyecto/backend/models"
	cstanza "RedesProyecto/backend/models/stanza"
	"fmt"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"log"
	"os"
	"strings"
	"time"
)

func handleMessage(s xmpp.Sender, p stanza.Packet) {
	msg, ok := p.(stanza.Message)
	if !ok {
		_, _ = fmt.Fprintf(os.Stdout, "Ignoring packet: %T\n", p)
		return
	}

	for _, ext := range msg.Extensions {
		if invite, ok := ext.(*cstanza.Conference); ok {
			// <message xmlns="jabber:client" to="alb21004@alumchat.lol" id="9946c7cb-b8fe-4214-9314-cc7ac91e1ab9" from="alb21005@alumchat.lol/gajim.0O3D5ZZ0"><x xmlns="jabber:x:conference" jid="ogivox@conference.alumchat.lol"></x></message>
			fmt.Println("Conference invitation from: ", invite.JID)
			//ConferenceInvitationChannel <- invite.JID

			events.EmitConferenceInvitation(AppContext, invite.JID, strings.Split(msg.From, "/")[0])

			return
		} else if mam, ok := ext.(*cstanza.MAM); ok {
			fmt.Println("MAM message: ", *mam)

			ts := mam.Forwarded.Delay.Stamp

			// parsear ts a time.Time
			tsTime, err := time.Parse(time.RFC3339, ts)

			if err != nil {
				fmt.Println(err)
			}

			toFormatted := strings.Split(mam.Forwarded.Message.To, "/")[0]
			fromFormatted := strings.Split(mam.Forwarded.Message.From, "/")[0]
			userFormatted := strings.Split(User.UserName, "/")[0]

			fmt.Printf("From: %s\n", fromFormatted)
			fmt.Printf("To: %s\n", toFormatted)
			fmt.Printf("User: %s\n", userFormatted)

			if toFormatted == userFormatted {
				// Un mensaje enviado a este usuario

				User.Messages[fromFormatted] = append(User.Messages[fromFormatted], models.Message{
					Body:      mam.Forwarded.Message.Body,
					From:      fromFormatted,
					To:        toFormatted,
					Timestamp: tsTime,
				})
				events.EmitArchive(AppContext)
				//events.EmitMessage(AppContext, fromFormatted)
			} else if fromFormatted == userFormatted {
				// Un mensaje enviado por este usuario
				User.Messages[toFormatted] = append(User.Messages[toFormatted], models.Message{
					Body:      mam.Forwarded.Message.Body,
					From:      fromFormatted,
					To:        toFormatted,
					Timestamp: tsTime,
				})
				events.EmitArchive(AppContext)
				//events.EmitMessage(AppContext, toFormatted)
			}

		}
	}

	switch msg.Type {
	case stanza.MessageTypeNormal:
		_, _ = fmt.Fprintf(os.Stdout, "(N) Message from: %s\n", msg.From)

		if msg.Body != "" {

			fromFormatted := strings.Split(msg.From, "/")[0]
			toFormatted := strings.Split(msg.To, "/")[0]

			message := models.NewMessage(msg.Body, toFormatted, fromFormatted)

			User.InsertMessage(message)

			events.EmitMessage(AppContext, msg.From)
			User.SaveConfig()
		}

	case stanza.MessageTypeChat:
		_, _ = fmt.Fprintf(os.Stdout, "(C) Message from: %s\n", msg.From)

		if msg.Body != "" {
			fromFormatted := strings.Split(msg.From, "/")[0]
			toFormatted := strings.Split(msg.To, "/")[0]

			message := models.NewMessage(msg.Body, toFormatted, fromFormatted)

			User.InsertMessage(message)

			events.EmitMessage(AppContext, msg.From)
			User.SaveConfig()
		}

	case stanza.MessageTypeGroupchat:
		_, _ = fmt.Fprintf(os.Stdout, "(G) Message from: %s\n", msg.From)

		//message := models.NewMessage(msg.Body)

		if msg.Body != "" {
			// Insertar en las conferencias

			conf := strings.Split(msg.From, "/")[0]

			if conference, ok := User.Conferences[conf]; ok {

				if msg.Body != "" {
					// Compara el Body y el From de los mensajes
					if conference.InsertMessage(models.NewMessage(msg.Body, conference.JID, msg.From)) {
						log.Println("Inserting message in conference: ", conference.JID)
						events.EmitMessage(AppContext, msg.From)
					}
				}
			}
		}

	default:
		if msg.Body != "" {
			events.EmitMessage(AppContext, msg.From)
		}
		_, _ = fmt.Fprintf(os.Stdout, "(%s) Message from: %s\n", msg.Type, msg.From)

	}

}

func handlePresence(s xmpp.Sender, p stanza.Packet) {
	//_, _ = fmt.Fprintf(os.Stdout, "Presence = %s\n", p)

	if presence, ok := p.(stanza.Presence); ok {

		/*
			<presence
			    from='coven@chat.shakespeare.lit/firstwitch'
			    to='crone1@shakespeare.lit/desktop'>

			<x xmlns='http://jabber.org/protocol/muc#user'>
			    <item affiliation='owner'
			          role='moderator'/>
			    <status code='110'/>
			    <status code='201'/>
			  </x>
			</presence>
		*/

		if presence.Extensions != nil {
			// https://xmpp.org/extensions/xep-0045.html#createroom

			for _, ext := range presence.Extensions {

				if mucUser, ok := ext.(*cstanza.MUCUser); ok && mucUser != nil && mucUser.Item.Role == "moderator" && mucUser.Item.Affiliation == "owner" && mucUser.Item.JID != "" {
					itemJid := strings.Split(mucUser.Item.JID, "/")[0] // jid="alb21004@alumchat.lol/9790d249-1efd-4341-9635-3dae9314b511"

					if itemJid == User.UserName {
						// Verificar si es de una sala existente de la cuál el usuario ya es dueño, o si es de una nueva sala que se acaba de crear

						fromFormatted := strings.Split(presence.From, "/")[0]

						if _, ok := User.Conferences[fromFormatted]; ok {
							events.EmitSuccess(AppContext, fmt.Sprintf("You are owner of conference: %s", fromFormatted))
							continue
						} else {
							// Crear la sala
							alias := strings.Split(fromFormatted, "@")[0]
							User.Conferences[fromFormatted] = models.NewConference(alias, fromFormatted)
						}

						log.Println("---------------------------------> MUC User extension: ", mucUser, fromFormatted)

						iq := &stanza.IQ{
							Attrs: stanza.Attrs{
								Type: stanza.IQTypeSet,
								From: itemJid,                              // JID of the user
								To:   strings.Split(presence.From, "/")[0], // JID of the conference
							},
							Payload: cstanza.NewMUCOwnerGet(),
						}

						c, err := s.SendIQ(AppContext, iq)

						if err != nil {
							log.Println(err)
							events.EmitError(AppContext, err.Error())
							continue
						}

						go func() {
							serverResponse := <-c

							if serverResponse.Type == stanza.IQTypeResult {
								log.Println("Server response: ", serverResponse)

								// Construir el IQ para enviar la configuración modificada
								acceptConfigIQ := &stanza.IQ{
									Attrs: stanza.Attrs{
										Type: stanza.IQTypeSet,
										From: itemJid,                              // JID del usuario
										To:   strings.Split(presence.From, "/")[0], // JID de la sala de conferencia
										Id:   "accept-config",
									},
									Payload: cstanza.NewMUCOwnerWithForm(), // Crea un MUCOwnerWithForm con el campo muc#roomconfig_persistentroom
								}

								// Enviar el IQ con la configuración modificada
								_, err := s.SendIQ(AppContext, acceptConfigIQ)

								if err != nil {
									log.Println(err)
									events.EmitError(AppContext, err.Error())
								} else {
									log.Println("Configuration updated and submitted for conference")
								}

							} else {
								log.Println("Unexpected server response type:", serverResponse.Type)
							}
						}()

						events.EmitSuccess(AppContext, "Conference created")

					}

				}
			}

		}

		switch presence.Type {
		case stanza.PresenceTypeSubscribe:
			// Un usuario ha solicitado suscribirse a nuestro estado de presencia.
			_, _ = fmt.Fprintf(os.Stdout, "Subscription request from: %s\n", presence.From)

			if presence.From != "" {
				// Emitir evento de suscripción, para decidir si aceptar o no.
				events.EmitSubscription(AppContext, presence.From)
			}

		case stanza.PresenceTypeSubscribed:
			// El usuario al que se solicitó la suscripción ha aceptado.
			_, _ = fmt.Fprintf(os.Stdout, "Subscription accepted from: %s\n", presence.From)

			events.EmitSuccess(AppContext, "Subscription accepted from: "+presence.From)

		case stanza.PresenceTypeUnsubscribed:
			// Un usuario ha cancelado la suscripción a nuestro estado de presencia.
			_, _ = fmt.Fprintf(os.Stdout, "Unsubscribed from: %s\n", presence.From)

			// TODO Notificar al usuario que se ha cancelado la suscripción

		case stanza.PresenceTypeUnsubscribe:
			// Un usuario ha solicitado cancelar la suscripción a nuestro estado de presencia.
			_, _ = fmt.Fprintf(os.Stdout, "Unsubscribe request from: %s\n", presence.From)

			// TODO Notificar al usuario que se ha solicitado cancelar la suscripción

		case stanza.PresenceTypeUnavailable:
			// El usuario está desconectado.
			_, _ = fmt.Fprintf(os.Stdout, "User %s is offline\n", presence.From)

			events.EmitPresenceUpdate(AppContext, presence.From, "Disconnected", "")

		default:
			_, _ = fmt.Fprintf(os.Stdout, "(%s) Presence from: %s\n", presence.Type, presence.From)

			/*
				() Presence from: alb21005@alumchat.lol/gajim.0O3D5ZZ0
				RECV:
				<presence xmlns="jabber:client" id="e6edad8c-ec83-4be6-8102-674f780c7b53" from="alb21005@alumchat.lol/gajim.0O3D5ZZ0" to="alb21004@alumchat.lol"><show>dnd</show><status>kkkk</status><c xmlns="http://jabber.org/protocol/caps" hash="sha-1" node="https://gajim.org" ver="XSQ75zZMlPNIYlOYsfXsvB/0F0g="></c></presence>

				Aquí, hay que extraer el estado de presencia y el mensaje de estado.
			*/

			fmt.Println("Show: ", presence.Show)

			var userPresence string

			switch presence.Show {

			case "dnd":
				userPresence = "Do Not Disturb"
			case "xa":
				userPresence = "Extended Away"
			case "away":
				userPresence = "Away"
			case "chat":
				userPresence = "Online"
			case "invisible":
				userPresence = "Disconnected"

			}

			var username string

			if presence.From != "" {
				if strings.Contains(presence.From, "@conference") {
					// group@conference/jid
					username = strings.Split(presence.From, "/")[1] + " at " + strings.Split(presence.From, "/")[0]
				} else if strings.Contains(presence.From, "/") {
					// username@server/resource
					username = strings.Split(presence.From, "/")[0]
				} else {
					username = presence.From
				}

			}
			//events.EmitNotification(AppContext, fmt.Sprintf("User %s: %s", username, userPresence), "info")
			events.EmitPresenceUpdate(AppContext, username, userPresence, presence.Status)

		}

	}
}

func handleIQ(s xmpp.Sender, p stanza.Packet) {
	switch iq := p.(type) {
	case *stanza.IQ:
		switch payload := iq.Payload.(type) {
		case *stanza.Roster:
			fmt.Println("ENTERED ROSTER")

		case *stanza.RosterItems:
			fmt.Println("Updating contacts from IQ handler")
			items := payload.Items

			contacts := make([]string, 0)

			for _, item := range items {
				fmt.Println("Item: ", item.Jid, item.Name, item.Subscription)
				contacts = append(contacts, item.Jid)
			}

			User.Contacts = contacts
			events.EmitContacts(AppContext, contacts)

		case *stanza.Version:
			// Aquí puedes manejar la versión del servidor u otros IQs de versión.
			fmt.Printf("Received version request from: %s\n", iq.From)

			// Responder con la versión del servidor
			resp := stanza.IQ{
				Attrs: stanza.Attrs{
					Type: stanza.IQTypeResult,
					From: iq.To,
					To:   iq.From,
					Id:   iq.Id,
				},
				Payload: &stanza.Version{
					Name:    "Mark Albrand",
					Version: "0.1",
					OS:      "Windows",
				},
			}

			err := s.Send(&resp)

			if err != nil {
				fmt.Println(err)
			}

			_, _ = fmt.Fprintf(os.Stdout, "Responded Version request from: %s\n", iq.From)

		case *cstanza.Ping:
			log.Println("Received ping from: ", iq.From)

			// Responder
			resp := stanza.IQ{
				Attrs: stanza.Attrs{
					Type: stanza.IQTypeResult,
					From: iq.To,
					To:   iq.From,
				},
				Payload: &cstanza.Ping{},
			}

			err := s.Send(&resp)

			if err != nil {
				log.Println(err)
				events.EmitError(AppContext, err.Error())
			}

		default:
			fmt.Printf("Unhandled IQ payload type: %T\n", payload)
		}

	default:
		fmt.Printf("Unhandled packet type: %T\n", p)
		_, _ = fmt.Fprintf(os.Stdout, "IQ = %s\n", p)

	}
}

func errorHandler(err error) {
	fmt.Println(err.Error())
}
