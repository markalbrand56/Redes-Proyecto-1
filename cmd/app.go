package main

import (
	"RedesProyecto/backend/chat"
	"RedesProyecto/backend/models"
	"context"
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

func (a *App) SetCorrespondent(correspondent string) {
	chat.CorrespondentChannel <- correspondent
}

func (a *App) SendMessage(message string) {
	chat.TextChannel <- message
}

func (a *App) GetContacts() []string {
	return chat.User.Contacts
}

func (a *App) GetMessages(username string) []models.Message {
	return chat.User.Messages[username]
}

func (a *App) UpdateContacts() {
	chat.FetchContactsChannel <- true
}

func (a *App) RequestContact(username string) {
	chat.SubscribeToChannel <- username
}

func (a *App) AcceptSubscription(username string) {
	chat.SubscriptionRequestChannel <- username
}

func (a *App) CancelSubscription(username string) {
	chat.UnsubscribeFromChannel <- username
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
