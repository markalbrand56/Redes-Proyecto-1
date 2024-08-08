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
	textChannel           = make(chan string)
	correspondentChannel  = make(chan string)
	requestContactChannel = make(chan bool)
)

const (
	username = "alb21004@alumchat.lol"
	password = "mark123"
	address  = "ws://alumchat.lol:7070/ws"
)

var user *models.User
var AppContext context.Context

func Start(ctx context.Context) {
	AppContext = ctx
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

	user = models.NewUser(newClient, username)
	startMessaging()
}

func startMessaging() {
	var text string
	var correspondent string

	for {
		select {
		case text = <-textChannel:
			fmt.Printf("Correspondent: %s Message: %s\n", correspondent, text)
			msg := stanza.Message{Attrs: stanza.Attrs{To: correspondent, Type: stanza.MessageTypeChat}, Body: text}
			err := user.Client.Send(msg)
			if err != nil {
				log.Fatalf("%+v", err)
			}

		case crrsp := <-correspondentChannel:
			fmt.Println("Correspondent: ", crrsp)
			correspondent = crrsp

		case <-requestContactChannel:
			// fecth contacts from xmpp server
			req, err := stanza.NewIQ(
				stanza.Attrs{
					From: user.UserName,
					Type: stanza.IQTypeGet,
					Id:   "roster_1",
				},
			)

			if err != nil {
				log.Fatalf("%+v", err)
			}

			req.RosterItems()

			c, err := user.Client.SendIQ(AppContext, req)

			if err != nil {
				log.Fatalf("%+v", err)
			}

			go func() {
				serverResp := <-c

				if rosterItems, ok := serverResp.Payload.(*stanza.RosterItems); ok {
					contacts := make([]string, 0)
					for _, item := range rosterItems.Items {
						contacts = append(contacts, item.Jid)
					}
					user.Contacts = contacts
					fmt.Println("Contacts: ", contacts)
					runtime.EventsEmit(AppContext, "contacts", contacts)
				}
			}()

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

func FetchContacts() {
	requestContactChannel <- true
}
