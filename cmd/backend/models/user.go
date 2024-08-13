package models

import (
	"encoding/json"
	"fmt"
	"gosrc.io/xmpp"
	"log"
	"os"
	"path/filepath"
	"slices"
)

type User struct {
	Client      *xmpp.Client         `json:"-,omitempty"` // Client cliente XMPP
	UserName    string               `json:"username"`    // UserName nombre de usuario
	Contacts    []string             `json:"contacts"`    // Contacts lista de contactos
	Conferences map[string]string    `json:"conferences"` // Conferences lista de salas de chat
	Messages    map[string][]Message `json:"messages"`    // Messages mensajes
}

// NewUser crea un nuevo usuario dado un Cliente XMPP previamente conectado y un nombre de usuario
func NewUser(client *xmpp.Client, username string) *User {
	return &User{
		Client:      client,
		UserName:    username,
		Contacts:    make([]string, 0),
		Conferences: make(map[string]string),
		Messages:    make(map[string][]Message),
	}
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

	// Copiar los datos del archivo de configuración al usuario actual

	for _, contact := range userFile.Contacts {
		if slices.Contains(u.Contacts, contact) == false {
			fmt.Println("Adding contact: ", contact)
			u.Contacts = append(u.Contacts, contact)
		}
	}

	for jid, conference := range userFile.Conferences {
		// Si no existe la sala de chat en la lista, se agrega
		if _, ok := u.Conferences[jid]; !ok {
			u.Conferences[jid] = conference
		}
	}

	for key, messages := range userFile.Messages {
		if _, ok := u.Messages[key]; !ok {
			// Si no existe la clave en el mapa de mensajes, se agrega
			u.Messages[key] = messages
		} else {
			for _, message := range messages {
				// Si ya existe el mensaje en la lista, no se agrega
				if slices.Contains(u.Messages[key], message) == false {
					u.Messages[key] = append(u.Messages[key], message)
				}
			}
		}
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
		UserName:    u.UserName,
		Contacts:    u.Contacts,
		Conferences: u.Conferences,
		Messages:    u.Messages,
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
