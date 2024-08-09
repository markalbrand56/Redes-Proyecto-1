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
	chat.Start(ctx)
}

func (a *App) SetCorrespondent(correspondent string) {
	chat.SetCorrespondent(correspondent)
}

func (a *App) SendMessage(message string) {
	chat.SendMessage(message)
}

func (a *App) GetContacts() []string {
	return chat.User.Contacts
}

func (a *App) UpdateContacts() {
	chat.FetchContacts()
}

func (a *App) RequestContact(username string) {
	chat.RequestContact(username)
}

func (a *App) AcceptSubscription(username string) {
	chat.AcceptSubscription(username)
}
