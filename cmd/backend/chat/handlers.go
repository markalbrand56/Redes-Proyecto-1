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

	runtime.EventsEmit(AppContext, "message", msg.Body, msg.From)

	_, _ = fmt.Fprintf(os.Stdout, "Body = %s - from = %s\n", msg.Body, msg.From)
}

func handlePresence(s xmpp.Sender, p stanza.Packet) {
	_, _ = fmt.Fprintf(os.Stdout, "Presence = %s\n", p)
}

func handleIQ(s xmpp.Sender, p stanza.Packet) {
	_, _ = fmt.Fprintf(os.Stdout, "IQ = %s\n", p)

	switch iq := p.(type) {
	case *stanza.IQ:
		switch payload := iq.Payload.(type) {
		case *stanza.Roster:
			fmt.Println("ENTERED ROSTER")

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
