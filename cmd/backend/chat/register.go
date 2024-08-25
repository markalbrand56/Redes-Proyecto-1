package chat

import (
	cstanza "RedesProyecto/backend/models/stanza"
	"context"
	"fmt"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"log"
	"os"
	"strings"
)

// RegisterNewUser registers a new user with the given username and password
func RegisterNewUser(ctx context.Context, username string, password string) bool {
	AppContext = ctx

	successChannel := make(chan bool) // Channel to return the result of the registration

	go register(successChannel, username, password)

	return <-successChannel
}

// Register registers a new user with the given username and password
func register(successChan chan bool, email string, password string) {
	defer close(successChan)
	done := make(chan struct{}) // Channel to wait for the response from the server before closing the client

	// Register the user{
	config := xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: address,
		},
		Jid:          "alb21004@alumchat.lol", // Admin
		Credential:   xmpp.Password("221756"),
		StreamLogger: os.Stdout,
		Insecure:     true,
	}

	registerRouter := xmpp.NewRouter()
	registerRouter.HandleFunc("message", registerHandleMessage)
	registerRouter.HandleFunc("presence", registerHandlePresence)
	registerRouter.HandleFunc("iq", registerHandleIQ)

	registerClient, err := xmpp.NewClient(&config, registerRouter, errorHandler)

	if err != nil {
		log.Println("Error while starting the client for registration", err)
		return
	}

	err = registerClient.Connect()

	if err != nil {
		log.Println("Error while connecting the client for registration", err)
		return
	}

	username := strings.Split(email, "@")[0]

	reg := &stanza.IQ{
		Attrs: stanza.Attrs{
			Type: stanza.IQTypeSet,
			To:   "alumchat.lol",
			Id:   "reg1",
		},
		Payload: cstanza.NewRegisterQueryWithUser(username, password, email),
	}

	c, err := registerClient.SendIQ(AppContext, reg)

	if err != nil {
		log.Println("Error while sending the IQ for registration", err)
		return
	}

	go func() {
		serverResponse := <-c

		if serverResponse.Type == stanza.IQTypeResult {
			fmt.Println("Server response: ", serverResponse)
			successChan <- true

		} else if serverResponse.Type == stanza.IQTypeError {
			fmt.Println("Error while registering the user")
			successChan <- false
		}

		close(done)
	}()

	select { // Wait for the response from the server
	case <-done:
		return // Close the client
	}
}

func registerHandleIQ(s xmpp.Sender, p stanza.Packet) {
	// Handle the IQ stanza
}

func registerHandlePresence(s xmpp.Sender, p stanza.Packet) {
	// Handle the presence stanza
}

func registerHandleMessage(s xmpp.Sender, p stanza.Packet) {
	// Handle the message stanza
}
