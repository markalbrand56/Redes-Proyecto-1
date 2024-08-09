package main

import (
	"RedesProyecto/backend/chat"
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
	chat.Start(ctx, "alb21004@alumchat.lol", "mark123")
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

func (a *App) UpdateContacts() {
	chat.FetchContactsChannel <- true
}

func (a *App) RequestContact(username string) {
	chat.SubscribeToChannel <- username
}

func (a *App) AcceptSubscription(username string) {
	chat.SubscriptionRequestChannel <- username
}
