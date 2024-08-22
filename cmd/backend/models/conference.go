package models

type Conference struct {
	Alias    string    // Alias del usuario en la sala de chat
	JID      string    // JID de la sala de chat
	Messages []Message // Mensajes de la sala de chat
}

// NewConference crea una nueva sala de chat con un alias y un JID
func NewConference(alias string, jid string) *Conference {
	return &Conference{
		Alias:    alias,
		JID:      jid,
		Messages: make([]Message, 0),
	}
}

func (c *Conference) String() string {
	return c.Alias + " (" + c.JID + ")"
}

func (c *Conference) InsertMessage(message Message) bool {
	//c.Messages = append(c.Messages, message)

	// antes de insertar el mensaje, verificar si ya existe
	for _, m := range c.Messages {
		if Equals(m, message) {
			return false
		}
	}

	c.Messages = append(c.Messages, message)
	return true
}
