package chat

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"os"
)

func handleMessage(s xmpp.Sender, p stanza.Packet) {
	msg, ok := p.(stanza.Message)
	if !ok {
		_, _ = fmt.Fprintf(os.Stdout, "Ignoring packet: %T\n", p)
		return
	}

	switch msg.Type {
	case stanza.MessageTypeNormal:
		_, _ = fmt.Fprintf(os.Stdout, "(N) Message from: %s\n", msg.From)

		if msg.Body != "" {
			runtime.EventsEmit(AppContext, "message", msg.Body, msg.From)
		}

	case stanza.MessageTypeChat:
		_, _ = fmt.Fprintf(os.Stdout, "(C) Message from: %s\n", msg.From)

		if msg.Body != "" {
			runtime.EventsEmit(AppContext, "message", msg.Body, msg.From)
		}

	default:
		if msg.Body != "" {
			runtime.EventsEmit(AppContext, "message", msg.Body, msg.From)
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
				AcceptSubscription(presence.From)
			}

			runtime.EventsEmit(AppContext, "subscription-request", presence.From)

		case stanza.PresenceTypeSubscribed:
			// El usuario al que se solicitó la suscripción ha aceptado.
			_, _ = fmt.Fprintf(os.Stdout, "Subscription accepted from: %s\n", presence.From)

		case stanza.PresenceTypeUnsubscribed:
			_, _ = fmt.Fprintf(os.Stdout, "Unsubscribed from: %s\n", presence.From)

		case stanza.PresenceTypeUnavailable:
			_, _ = fmt.Fprintf(os.Stdout, "User %s is offline\n", presence.From)

		default:
			_, _ = fmt.Fprintf(os.Stdout, "(%s) Presence from: %s\n", presence.Type, presence.From)

		}

	}
}

func handleIQ(s xmpp.Sender, p stanza.Packet) {
	_, _ = fmt.Fprintf(os.Stdout, "IQ = %s\n", p)

	switch iq := p.(type) {
	case *stanza.IQ:
		switch payload := iq.Payload.(type) {
		case *stanza.Roster:
			fmt.Println("ENTERED ROSTER")

		case *stanza.RosterItems:
			items := payload.Items

			for _, item := range items {
				fmt.Println("Item: ", item.Jid, item.Name, item.Subscription)
			}

		case *stanza.Version:
			// Aquí puedes manejar la versión del servidor u otros IQs de versión.
			fmt.Printf("Received version request from: %s\n", iq.From)

		default:
			fmt.Printf("Unhandled IQ payload type: %T\n", payload)
		}

	default:
		fmt.Printf("Unhandled packet type: %T\n", p)
	}
}

func errorHandler(err error) {
	fmt.Println(err.Error())
}
