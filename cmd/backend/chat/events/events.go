package events

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strings"
)

func EmitContacts(ctx context.Context, contacts []string) {
	runtime.EventsEmit(ctx, "contacts", contacts)
}

func EmitMessages(ctx context.Context) {
	runtime.EventsEmit(ctx, "update-messages")
}

func EmitMessage(ctx context.Context, from string) {
	fromFormatted := strings.Split(from, "/")[0]
	runtime.EventsEmit(ctx, "message", fromFormatted)
}

func EmitConferences(ctx context.Context, conferences map[string]string) {
	runtime.EventsEmit(ctx, "conferences", conferences)
}

func EmitSuccess(ctx context.Context, message string) {
	runtime.EventsEmit(ctx, "success", message)
}

func EmitSubscription(ctx context.Context, username string) {
	runtime.EventsEmit(ctx, "subscription", username)
}
