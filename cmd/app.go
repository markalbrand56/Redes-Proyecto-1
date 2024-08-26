package main

import (
	"RedesProyecto/backend/chat"
	"RedesProyecto/backend/models"
	"context"
	"fmt"
	"gosrc.io/xmpp/stanza"
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

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// OnShutdown is called when the app is shutting down
func (a *App) OnShutdown(ctx context.Context) {
	chat.Close()
}

// Login logs in the user with the given username and password
func (a *App) Login(username string, password string) {
	chat.Start(a.ctx, username, password)
	//chat.RegisterNewUser(a.ctx, username, password)
}

// Register registers a new user with the given username, password and email
func (a *App) Register(email string, password string) bool {
	return chat.RegisterNewUser(a.ctx, email, password)
}

// Logout logs out the user
func (a *App) Logout() {
	chat.LogoutChannel <- true
}

// DeleteAccount deletes the account of the user
func (a *App) DeleteAccount() {
	chat.DeleteAccountChannel <- true
}

// SendMessage sends a message to the given user
func (a *App) SendMessage(body string, to string, from string) {
	message := models.NewMessage(body, to, from)

	fmt.Printf("Sending message: %s\n", message)

	chat.TextChannel <- message
}

// SendFileMessage sends a file message to the given user
func (a *App) SendFileMessage(file string, to string, from string) {
	message := models.NewMessage(file, to, from)

	fmt.Printf("Sending file message: %s\n", message)

	chat.FileChannel <- message
}

// SendConferenceMessage sends a message to the given conference
func (a *App) SendConferenceMessage(body string, to string, from string) {
	message := models.NewMessage(body, to, from)

	fmt.Printf("Sending conference message: %s\n", message)

	chat.ConferenceTextChannel <- message
}

// SendConferenceFileMessage sends a file message to the given conference
func (a *App) SendConferenceFileMessage(file string, to string, from string) {
	message := models.NewMessage(file, to, from)

	fmt.Printf("Sending conference file message: %s\n", message)

	chat.ConferenceFileChannel <- message
}

// GetContacts returns the contacts of the user
func (a *App) GetContacts() []string {
	return chat.User.Contacts
}

func (a *App) GetConferences() []map[string]string {
	var conferences []map[string]string

	for _, conference := range chat.User.Conferences {
		c := make(map[string]string)
		c["alias"] = conference.Alias
		c["jid"] = conference.JID

		conferences = append(conferences, c)
	}

	return conferences
}

// GetMessages returns the messages of the given user
func (a *App) GetMessages(username string) []models.Message {
	if _, ok := chat.User.Messages[username]; !ok {
		return []models.Message{}
	}
	r := chat.User.Messages[username]

	log.Printf("Messages from '%s' (%d): %v\n", username, len(r), r)
	return r
}

// GetMessagesConference returns the messages of the given conference
func (a *App) GetMessagesConference(jid string) []models.Message {
	u := chat.User
	if _, ok := u.Conferences[jid]; !ok {
		log.Println("Conference not found")
		return []models.Message{}
	}

	log.Println("Conference found")
	return chat.User.Conferences[jid].Messages
}

// UpdateContacts updates the contacts of the user from the server
func (a *App) UpdateContacts() {
	log.Println("Updating contacts")
	chat.FetchContactsChannel <- true
}

// RequestContact sends a contact request to the given username
func (a *App) RequestContact(username string) {
	chat.SubscribeToChannel <- username
}

// AcceptSubscription accepts the subscription request of the given username
func (a *App) AcceptSubscription(username string) {
	// Se ha decidido aceptar la solicitud de suscripción
	chat.SubscriptionRequestChannel <- username
}

// AcceptConferenceInvitation accepts the conference invitation of the given JID
func (a *App) AcceptConferenceInvitation(jid string) {
	chat.ConferenceInvitationChannel <- jid
}

// DeclineConference declines the conference invitation of the given JID
func (a *App) DeclineConference(conferenceJID string, sender string) {
	chat.ConferenceDeclineInvitationChannel <- models.NewInvitation(conferenceJID, sender)
}

// SendInvitation sends an invitation to the given username to join the conference
func (a *App) SendInvitation(conferenceJID string, to string) {
	log.Printf("Inviting %s to conference %s\n", to, conferenceJID)
	chat.InviteToConferenceChannel <- models.NewInvitation(conferenceJID, to)
}

// CreateConference creates a new conference with the given alias
func (a *App) CreateConference(alias string) {
	if alias == "" {
		log.Println("Alias is empty")
		return
	}

	conf := models.NewConference(alias, fmt.Sprintf("%s@conference.alumchat.lol", alias))

	log.Printf("Creating conference: %s {%s}\n", conf.Alias, conf.JID)

	chat.NewConferenceChannel <- *conf
}

// DeleteConference deletes the conference with the given JID
func (a *App) DeleteConference(jid string) {
	chat.DeleteConferenceChannel <- jid
}

// ExitConference edits the conference with the given JID
func (a *App) ExitConference(jid string) {
	chat.ExitConferenceChannel <- jid
}

// CancelSubscription cancels the subscription of the given username
func (a *App) CancelSubscription(username string) {
	chat.UnsubscribeFromChannel <- username
}

// RejectSubscription rejects the subscription of the given username
func (a *App) RejectSubscription(username string) {
	chat.RejectionChannel <- username
}

// SetStatus sets the status of the user
func (a *App) SetStatus(status int) {
	switch status {
	case 0: // Online
		chat.ShowChannel <- models.StatusOnline
	case 1: // Away
		chat.ShowChannel <- models.StatusAway
	case 2: // Busy
		chat.ShowChannel <- models.StatusBusy
	case 3: // Not Available
		chat.ShowChannel <- models.StatusNotAvailable
	case 4: // Offline
		chat.ShowChannel <- models.StatusOffline
	default:
		chat.ShowChannel <- models.StatusOnline
	}
}

// SetStatusMessage sets the status message of the user
func (a *App) SetStatusMessage(message string) {
	chat.StatusChannel <- message
}

// GetArchive gets the archive of the given username
func (a *App) GetArchive(username string) {
	chat.FetchArchiveChannel <- username
}

// GetCurrentUser returns the current user
func (a *App) GetCurrentUser() string {
	return chat.User.UserName
}

// ProbeContacts probes the contacts of the user to fetch their status
func (a *App) ProbeContacts() {
	chat.ProbeContactsChannel <- true
}

// GetCurrentStatus returns the current status of the user
func (a *App) GetCurrentStatus() string {
	return chat.User.Status
}

// GetCurrentShow returns the current status message of the user
func (a *App) GetCurrentShow() int {
	switch chat.User.Show {
	case stanza.PresenceShowChat:
		return 0
	case stanza.PresenceShowAway:
		return 1
	case stanza.PresenceShowDND:
		return 2
	case stanza.PresenceShowXA:
		return 3
	default:
		return 4
	}
}
