package events

import (
	"RedesProyecto/backend/models"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strings"
)

func EmitContacts(ctx context.Context, contacts []string) {
	runtime.EventsEmit(ctx, "contacts", contacts)
}

func EmitArchive(ctx context.Context) {
	runtime.EventsEmit(ctx, "archive")
}

func EmitMessage(ctx context.Context, from string) {
	fromFormatted := strings.Split(from, "/")[0]
	runtime.EventsEmit(ctx, "message", fromFormatted)
}

func EmitConferences(ctx context.Context, conferences map[string]*models.Conference) {
	mp := make(map[string]string) // map[alias]jid

	for _, conference := range conferences {
		mp[conference.Alias] = conference.JID
	}

	runtime.EventsEmit(ctx, "conferences", mp)
}

func EmitSuccess(ctx context.Context, message string) {
	runtime.EventsEmit(ctx, "success", message)
}

func EmitNotification(ctx context.Context, message string, notificationType string) {
	runtime.EventsEmit(ctx, "notification", message, notificationType)
}

func EmitSubscription(ctx context.Context, username string) {
	runtime.EventsEmit(ctx, "subscription", username)
}

func EmitError(ctx context.Context, message string) {
	runtime.EventsEmit(ctx, "error", message)
}

func EmitPresenceUpdate(ctx context.Context, user string, presence string) {
	runtime.EventsEmit(ctx, "presence", user, presence)
}
