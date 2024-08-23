package models

type Invitation struct {
	ConferenceJID string // JID de la sala de chat
	To            string // JID del usuario al que se le envía la invitación
}

// NewInvitation crea una nueva invitación con un JID de sala de chat y un JID de usuario
func NewInvitation(conferenceJID string, to string) Invitation {
	return Invitation{
		ConferenceJID: conferenceJID,
		To:            to,
	}
}

func (i Invitation) String() string {
	return i.ConferenceJID + " -> " + i.To
}
