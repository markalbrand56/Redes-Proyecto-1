package chat

import (
	"RedesProyecto/backend/chat/events"
	"RedesProyecto/backend/models"
	cstanza "RedesProyecto/backend/models/stanza"
	"fmt"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
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
			ConferenceInvitationChannel <- invite.JID
			return
		} else if mam, ok := ext.(*cstanza.MAM); ok {
			fmt.Println("MAM message: ", *mam)

			ts := mam.Forwarded.Delay.Stamp

			// parsear ts a time.Time
			tsTime, err := time.Parse(time.RFC3339, ts)

			if err != nil {
				fmt.Println(err)
			}

			if mam.Forwarded.Message.To == strings.Split(User.UserName, "/")[0] {
				fromFormatted := strings.Split(mam.Forwarded.Message.From, "/")[0]
				toFormatted := strings.Split(mam.Forwarded.Message.To, "/")[0]

				User.Messages[fromFormatted] = append(User.Messages[fromFormatted], models.Message{
					Body:      mam.Forwarded.Message.Body,
					From:      fromFormatted,
					To:        toFormatted,
					Timestamp: tsTime,
				})
			} else if mam.Forwarded.Message.From == strings.Split(User.UserName, "/")[0] {
				fromFormatted := strings.Split(mam.Forwarded.Message.From, "/")[0]
				toFormatted := strings.Split(mam.Forwarded.Message.To, "/")[0]

				User.Messages[mam.Forwarded.Message.To] = append(User.Messages[toFormatted], models.Message{
					Body:      mam.Forwarded.Message.Body,
					From:      fromFormatted,
					To:        toFormatted,
					Timestamp: tsTime,
				})
			}

			events.EmitMessages(AppContext)
		}
	}

	switch msg.Type {
	case stanza.MessageTypeNormal:
		_, _ = fmt.Fprintf(os.Stdout, "(N) Message from: %s\n", msg.From)

		//message := models.NewMessage(msg.Body)

		if msg.Body != "" {
			//User.Messages[msg.From] = append(User.Messages[msg.From], *message)
			events.EmitMessage(AppContext, msg.From)
			User.SaveConfig()
		}

	case stanza.MessageTypeChat:
		_, _ = fmt.Fprintf(os.Stdout, "(C) Message from: %s\n", msg.From)

		//message := models.NewMessage(msg.Body)

		if msg.Body != "" {
			//User.Messages[msg.From] = append(User.Messages[msg.From], *message)
			events.EmitMessage(AppContext, msg.From)
			User.SaveConfig()
		}

	case stanza.MessageTypeGroupchat:
		_, _ = fmt.Fprintf(os.Stdout, "(G) Message from: %s\n", msg.From)

		//message := models.NewMessage(msg.Body)

		if msg.Body != "" {
			//User.Messages[msg.From] = append(User.Messages[msg.From], *message)
			events.EmitMessage(AppContext, msg.From)
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
		switch presence.Type {
		case stanza.PresenceTypeSubscribe:
			// Un usuario ha solicitado suscribirse a nuestro estado de presencia.
			_, _ = fmt.Fprintf(os.Stdout, "Subscription request from: %s\n", presence.From)

			if presence.From != "" {
				// TODO Verificar si el usuario está en la lista de contactos
				SubscriptionRequestChannel <- presence.From
			}

			events.EmitSubscription(AppContext, presence.From)

		case stanza.PresenceTypeSubscribed:
			// El usuario al que se solicitó la suscripción ha aceptado.
			_, _ = fmt.Fprintf(os.Stdout, "Subscription accepted from: %s\n", presence.From)

		case stanza.PresenceTypeUnsubscribed:
			// El usuario al que se solicitó la suscripción ha rechazado.
			_, _ = fmt.Fprintf(os.Stdout, "Unsubscribed from: %s\n", presence.From)

		case stanza.PresenceTypeUnavailable:
			// El usuario está desconectado.
			_, _ = fmt.Fprintf(os.Stdout, "User %s is offline\n", presence.From)

		default:
			_, _ = fmt.Fprintf(os.Stdout, "(%s) Presence from: %s\n", presence.Type, presence.From)

			/*
				() Presence from: alb21005@alumchat.lol/gajim.0O3D5ZZ0
				RECV:
				<presence xmlns="jabber:client" id="e6edad8c-ec83-4be6-8102-674f780c7b53" from="alb21005@alumchat.lol/gajim.0O3D5ZZ0" to="alb21004@alumchat.lol"><show>dnd</show><status>kkkk</status><c xmlns="http://jabber.org/protocol/caps" hash="sha-1" node="https://gajim.org" ver="XSQ75zZMlPNIYlOYsfXsvB/0F0g="></c></presence>

				Aquí, hay que extraer el estado de presencia y el mensaje de estado.
			*/

			if presence.Show != "" {
				fmt.Println("Show: ", presence.Show)
			}

			if presence.Status != "" {
				fmt.Println("Status: ", presence.Status)
			}

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
