package events

import (
	"RedesProyecto/backend/models"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strings"
)

// EmitSuccess emite un evento de éxito
func EmitSuccess(ctx context.Context, message string) {
	runtime.EventsEmit(ctx, "success", message)
}

// EmitError emite un evento de error
func EmitError(ctx context.Context, message string) {
	runtime.EventsEmit(ctx, "error", message)
}

// EmitLogin emite un evento de login exitoso
func EmitLogin(ctx context.Context, username string) {
	runtime.EventsEmit(ctx, "login", username)
}

// EmitLogout emite un evento de logout exitoso
func EmitLogout(ctx context.Context) {
	runtime.EventsEmit(ctx, "logout")
}

// EmitContacts emite un evento de contactos, con la lista de contactos del usuario actual
func EmitContacts(ctx context.Context, contacts []string) {
	runtime.EventsEmit(ctx, "contacts", contacts)
}

// EmitArchive emite un evento de que se actualizó el archive de mensajes del usuario
func EmitArchive(ctx context.Context) {
	runtime.EventsEmit(ctx, "archive")
}

// EmitMessage emite un evento de mensaje recibido
func EmitMessage(ctx context.Context, from string) {
	fromFormatted := strings.Split(from, "/")[0]
	runtime.EventsEmit(ctx, "message", fromFormatted)
}

// EmitConferences emite un evento de conferencias, con un mapa de alias de conferencias y sus JID
func EmitConferences(ctx context.Context, conferences map[string]*models.Conference) {
	mp := make(map[string]string) // map[alias]jid

	for _, conference := range conferences {
		mp[conference.Alias] = conference.JID
	}

	runtime.EventsEmit(ctx, "conferences", mp)
}

// EmitSubscription emite un evento acerca de que un usuario ha solicitado suscribirse a nuestro estado de presencia
func EmitSubscription(ctx context.Context, username string) {
	runtime.EventsEmit(ctx, "subscription", username)
}

// EmitPresenceUpdate emite un evento acerca de que un usuario ha actualizado su presencia
func EmitPresenceUpdate(ctx context.Context, user string, presence string, status string) {
	runtime.EventsEmit(ctx, "presence", user, presence, status)
}

// EmitConferenceInvitation emite un evento acerca de que se ha recibido una invitación a una conferencia
func EmitConferenceInvitation(ctx context.Context, jid string) {
	runtime.EventsEmit(ctx, "conference-invitation", jid)
}
