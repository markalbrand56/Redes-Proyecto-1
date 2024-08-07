package chat

import (
	"fmt"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"log"
	"os"
)

var (
	textChannel          = make(chan string)
	correspondentChannel = make(chan string)
)

const (
	username = "alb21004@alumchat.lol"
	password = "mark123"
	address  = "ws://alumchat.lol:7070/ws"
)

var client *xmpp.Client

func init() {
	go startClient()
}

func startClient() {
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

	client = newClient
	startMessaging(client)
}

func startMessaging(client xmpp.Sender) {
	var text string
	var correspondent string

	for {
		select {
		case text = <-textChannel:
			fmt.Printf("Correspondent: %s Message: %s\n", correspondent, text)
			msg := stanza.Message{Attrs: stanza.Attrs{To: correspondent, Type: stanza.MessageTypeChat}, Body: text}
			err := client.Send(msg)
			if err != nil {
				log.Fatalf("%+v", err)
			}

		case crrsp := <-correspondentChannel:
			fmt.Println("Correspondent: ", crrsp)
			correspondent = crrsp

		default:
			continue
		}
	}
}

func SetCorrespondent(correspondent string) {
	correspondentChannel <- correspondent
}

func SendMessage(message string) {
	textChannel <- message
}
