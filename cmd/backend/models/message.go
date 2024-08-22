package models

import "time"

type Message struct {
	Body      string    `json:"body"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Timestamp time.Time `json:"timestamp"`
}

func NewMessage(body string, to string, from string) Message {
	return Message{
		Body:      body,
		From:      from,
		To:        to,
		Timestamp: time.Now(),
	}
}

func (m *Message) String() string {
	return m.From + " -> " + m.To + ": " + m.Body
}

// Equals compara dos mensajes, pero no toma en cuenta el timestamp. Está pensado para comparar mensajes que se obtienen del servidor sin timestamp.
func Equals(m1 Message, m2 Message) bool {
	return m1.Body == m2.Body && m1.From == m2.From && m1.To == m2.To
}
