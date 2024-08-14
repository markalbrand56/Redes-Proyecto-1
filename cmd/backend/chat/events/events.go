package events

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func EmitContacts(ctx context.Context, contacts []string) {
	runtime.EventsEmit(ctx, "contacts", contacts)
}

func EmitMessages(ctx context.Context, messages []string) {
	runtime.EventsEmit(ctx, "messages", messages)
}

func EmitMessage(ctx context.Context, message string, from string) {
	runtime.EventsEmit(ctx, "message", message, from)
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
