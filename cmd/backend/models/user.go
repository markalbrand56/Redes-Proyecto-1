package models

import (
	"gosrc.io/xmpp"
	"os"
)

type User struct {
	Client   *xmpp.Client
	UserName string               // Nombre de usuario
	Contacts []string             // Lista de contactos
	Messages map[string][]Message // Mensajes recibidos
}

// NewUser crea un nuevo usuario dado un Cliente XMPP previamente conectado y un nombre de usuario
func NewUser(client *xmpp.Client, username string) *User {
	return &User{
		UserName: username,
		Client:   client,
		Contacts: make([]string, 0),
		Messages: make(map[string][]Message),
	}
}

func (u *User) LoadConfig() {
	// TODO Cargar la configuración
	_, _ = os.UserConfigDir() // La carpeta de configuración del usuario
}

func (u *User) SaveConfig() {
	// TODO Guardar la configuración
	_, _ = os.UserConfigDir() // La carpeta de configuración del usuario
}
