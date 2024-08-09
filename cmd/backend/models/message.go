package models

import "time"

type Message struct {
	Body      string    `json:"body"`
	Timestamp time.Time `json:"timestamp"`
}

func NewMessage(body string) *Message {
	return &Message{
		Body:      body,
		Timestamp: time.Now(),
	}
}
