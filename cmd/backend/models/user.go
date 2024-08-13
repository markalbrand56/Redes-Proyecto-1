package models

import (
	"encoding/json"
	"fmt"
	"gosrc.io/xmpp"
	"log"
	"os"
	"path/filepath"
)

type User struct {
	Client      *xmpp.Client         `json:"-,omitempty"` // Client cliente XMPP
	UserName    string               `json:"username"`    // UserName nombre de usuario
	Contacts    []string             `json:"contacts"`    // Contacts lista de contactos
	Conferences []string             `json:"conferences"` // Conferences lista de salas de chat
	Messages    map[string][]Message `json:"messages"`    // Messages mensajes
}

// NewUser crea un nuevo usuario dado un Cliente XMPP previamente conectado y un nombre de usuario
func NewUser(client *xmpp.Client, username string) *User {
	return &User{
		Client:      client,
		UserName:    username,
		Contacts:    make([]string, 0),
		Conferences: make([]string, 0),
		Messages:    make(map[string][]Message),
	}
}

func (u *User) LoadConfig() {
	// TODO Cargar la configuración
	_, _ = os.UserConfigDir() // La carpeta de configuración del usuario
}

func (u *User) SaveConfig() {
	// Obtener la carpeta de configuración del usuario
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error obteniendo el directorio de configuración del usuario: %+v", err)
	}

	dir = filepath.Join(dir, "alumchat")

	fmt.Println("Saving config to: ", dir)

	// Crear el directorio si no existe
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creando el directorio de configuración del usuario: %+v", err)
	}

	// Definir la ruta completa del archivo de configuración
	configFile := filepath.Join(dir, "xmpp_user_config.json")

	userSave := User{
		UserName:    u.UserName,
		Contacts:    u.Contacts,
		Conferences: u.Conferences,
		Messages:    u.Messages,
	}

	// Convertir la estructura User a JSON
	configData, err := json.MarshalIndent(userSave, "", "  ")
	if err != nil {
		log.Fatalf("Error serializando la configuración del usuario a JSON: %+v", err)
	}

	// Guardar el archivo JSON
	err = os.WriteFile(configFile, configData, 0644)
	if err != nil {
		log.Fatalf("Error guardando el archivo de configuración del usuario: %+v", err)
	}

	log.Printf("Configuración guardada exitosamente en %s", configFile)
}
