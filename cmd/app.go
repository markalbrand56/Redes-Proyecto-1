package main

import (
	"RedesProyecto/backend/chat"
	"context"
	"fmt"
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) SetCorrespondent(correspondent string) {
	chat.SetCorrespondent(correspondent)
}

func (a *App) SendMessage(message string) {
	chat.SendMessage(message)
}
