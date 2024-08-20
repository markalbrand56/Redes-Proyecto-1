package models

import (
	"encoding/json"
	"fmt"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"log"
	"os"
	"path/filepath"
)

const (
	StatusOnline       = "online"
	StatusAway         = "away"
	StatusBusy         = "busy"
	StatusNotAvailable = "not-available"
	StatusOffline      = "offline"
)

type User struct {
	Client      *xmpp.Client           `json:"-,omitempty"`        // Client cliente XMPP
	UserName    string                 `json:"username"`           // UserName nombre de usuario
	Contacts    []string               `json:"contacts"`           // Contacts lista de contactos
	Conferences map[string]*Conference `json:"conferences"`        // Conferences lista de salas de chat
	Messages    map[string][]Message   `json:"messages,omitempty"` // Messages mensajes
	Status      string                 `json:"status"`             // Status estado del usuario
	Show        stanza.PresenceShow    `json:"show,omitempty"`
}

// NewUser crea un nuevo usuario dado un Cliente XMPP previamente conectado y un nombre de usuario
func NewUser(client *xmpp.Client, username string) *User {
	return &User{
		Client:      client,
		UserName:    username,
		Contacts:    make([]string, 0),
		Conferences: make(map[string]*Conference, 0),
		Messages:    make(map[string][]Message),
		Status:      StatusOnline,
		Show:        stanza.PresenceShowChat,
	}
}

func (u *User) InsertConference(conference *Conference) {
	u.Conferences[conference.JID] = conference
}

func (u *User) InsertMessage(message Message) {
	u.Messages[message.From] = append(u.Messages[message.From], message)
}

func (u *User) String() string {
	return fmt.Sprintf("{UserName: %s, Contacts: %v, Conferences: %v, Messages: %v}", u.UserName, u.Contacts, u.Conferences, u.Messages)
}

func (u *User) LoadConfig() error {
	// TODO Cargar la configuración
	dir, err := os.UserConfigDir() // La carpeta de configuración del usuario

	if err != nil {
		log.Printf("Error obteniendo el directorio de configuración del usuario: %+v\n", err)
		return err
	}

	dir = filepath.Join(dir, "alumchat") // La carpeta de configuración de la aplicación

	// Crear el directorio si no existe
	err = os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		log.Printf("Error creando el directorio de configuración del usuario: %+v", err)
		return err
	}

	// Definir la ruta completa del archivo de configuración
	configFile := filepath.Join(dir, "xmpp_user_config.json")

	// Leer el archivo de configuración
	configData, err := os.ReadFile(configFile)

	if err != nil {
		log.Printf("Error leyendo el archivo de configuración del usuario: %+v", err)
		return err
	}

	// Convertir el archivo JSON a la estructura User
	userFile := User{}

	err = json.Unmarshal(configData, &userFile)

	if err != nil {
		log.Printf("Error deserializando la configuración del usuario desde JSON: %+v", err)
		return err
	}

	if userFile.UserName != u.UserName {
		log.Printf("El nombre de usuario no coincide con el archivo de configuración")
		return fmt.Errorf("el nombre de usuario no coincide con el archivo de configuración")
	}

	// Actualizar el estado del usuario
	if userFile.Status != "" {
		u.Status = userFile.Status
	}

	if userFile.Show != "" {
		u.Show = userFile.Show
	}

	return nil
}

func (u *User) SaveConfig() error {
	// Obtener la carpeta de configuración del usuario
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Error obteniendo el directorio de configuración del usuario: %+v\n", err)
		return err
	}

	dir = filepath.Join(dir, "alumchat")

	fmt.Println("Saving config to: ", dir)

	// Crear el directorio si no existe
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Printf("Error creando el directorio de configuración del usuario: %+v", err)
		return err
	}

	// Definir la ruta completa del archivo de configuración
	configFile := filepath.Join(dir, "xmpp_user_config.json")

	userSave := User{
		UserName: u.UserName,
		Status:   u.Status,
		Show:     u.Show,
	}

	// Convertir la estructura User a JSON
	configData, err := json.MarshalIndent(userSave, "", "  ")
	if err != nil {
		log.Fatalf("Error serializando la configuración del usuario a JSON: %+v", err)
		return err
	}

	// Guardar el archivo JSON
	err = os.WriteFile(configFile, configData, 0644)
	if err != nil {
		log.Fatalf("Error guardando el archivo de configuración del usuario: %+v", err)
		return err
	}

	log.Printf("Configuración guardada exitosamente en %s", configFile)

	return nil
}

func (u *User) ShowConversations() {
	for contact, messages := range u.Messages {
		fmt.Printf("Messages from %s (%d): %s\n", contact, len(messages), messages[len(messages)-1].Body)
	}
}
