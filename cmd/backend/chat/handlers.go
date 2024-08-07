package chat

import (
	"fmt"
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

	_, _ = fmt.Fprintf(os.Stdout, "Body = %s - from = %s\n", msg.Body, msg.From)
}

func handlePresence(s xmpp.Sender, p stanza.Packet) {
	_, _ = fmt.Fprintf(os.Stdout, "Presence = %s\n", p)
}

func handleIQ(s xmpp.Sender, p stanza.Packet) {
	_, _ = fmt.Fprintf(os.Stdout, "IQ = %s\n", p)
}

func errorHandler(err error) {
	fmt.Println(err.Error())
}
