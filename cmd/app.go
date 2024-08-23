package main

import (
	"RedesProyecto/backend/chat"
	"RedesProyecto/backend/models"
	"context"
	"fmt"
	"log"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// DEBUG
	//a.Login("alb21004@alumchat.lol", "mark123")
}

func (a *App) OnShutdown(ctx context.Context) {
	chat.Close()
}

func (a *App) Login(username string, password string) {
	chat.Start(a.ctx, username, password)
}

func (a *App) Logout() {
	chat.LogoutChannel <- true
}

func (a *App) SendMessage(body string, to string, from string) {
	message := models.NewMessage(body, to, from)

	fmt.Printf("Sending message: %s\n", message)

	chat.TextChannel <- message
}

func (a *App) SendConferenceMessage(body string, to string, from string) {
	message := models.NewMessage(body, to, from)

	fmt.Printf("Sending conference message: %s\n", message)

	chat.ConferenceTextChannel <- message
}

func (a *App) GetContacts() []string {
	return chat.User.Contacts
}

func (a *App) GetMessages(username string) []models.Message {
	//chat.User.ShowConversations()

	if _, ok := chat.User.Messages[username]; !ok {
		return []models.Message{}
	}
	r := chat.User.Messages[username]

	log.Printf("Messages from '%s' (%d): %v\n", username, len(r), r)
	return r
}

func (a *App) GetMessagesConference(jid string) []models.Message {
	u := chat.User
	if _, ok := u.Conferences[jid]; !ok {
		log.Println("Conference not found")
		return []models.Message{}
	}

	log.Println("Conference found")
	return chat.User.Conferences[jid].Messages
}

func (a *App) UpdateContacts() {
	chat.FetchContactsChannel <- true
}

func (a *App) RequestContact(username string) {
	chat.SubscribeToChannel <- username
}

func (a *App) AcceptSubscription(username string) {
	// Se ha decidido aceptar la solicitud de suscripción
	chat.SubscriptionRequestChannel <- username
}

func (a *App) AcceptConferenceInvitation(jid string) {
	chat.ConferenceInvitationChannel <- jid
}

func (a *App) CancelSubscription(username string) {
	chat.UnsubscribeFromChannel <- username
}

func (a *App) RejectSubscription(username string) {
	chat.RejectionChannel <- username
}

func (a *App) SetStatus(status int) {
	switch status {
	case 0: // Online
		chat.StatusChannel <- models.StatusOnline
	case 1: // Away
		chat.StatusChannel <- models.StatusAway
	case 2: // Busy
		chat.StatusChannel <- models.StatusBusy
	case 3: // Not Available
		chat.StatusChannel <- models.StatusNotAvailable
	case 4: // Offline
		chat.StatusChannel <- models.StatusOffline
	default:
		chat.StatusChannel <- models.StatusOnline
	}
}

func (a *App) GetArchive(username string) {
	chat.FetchArchiveChannel <- username
}

func (a *App) GetCurrentUser() string {
	return chat.User.UserName
}

func (a *App) ProbeContacts() {
	chat.ProbeContactsChannel <- true
}
