package stanza

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

/*
<message
			    from='crone1@shakespeare.lit/desktop'
			    to='hecate@shakespeare.lit'>
			  <x xmlns='jabber:x:conference'
			     jid='darkcave@macbeth.shakespeare.lit'
			     reason='Hey Hecate, this is the place for all good witches!'/>
			</message>
*/

type ConferenceInvite struct {
	XMLName xml.Name `xml:"x"`
	XMLNS   string   `xml:"xmlns,attr"`
	JID     string   `xml:"jid,attr"`
	Reason  string   `xml:"reason,attr"`
}

// Name devuelve el nombre del tipo de archivo
func (f ConferenceInvite) Name() string {
	return "ConferenceInvite"
}

// Namespace devuelve el espacio de nombres para el archivo
func (f ConferenceInvite) Namespace() string {
	return f.XMLName.Space
}

// GetSet no es necesario en este contexto
func (f ConferenceInvite) GetSet() *stanza.ResultSet {
	return nil
}

// NewConferenceInvite crea una nueva instancia de la estructura ConferenceInvite con un JID, contraseña y razón
func NewConferenceInvite(jid string, reason string) ConferenceInvite {
	return ConferenceInvite{
		XMLName: xml.Name{Space: "jabber:x:conference", Local: "x"},
		XMLNS:   "jabber:x:conference",
		JID:     jid,    // Asigna el JID al campo JID
		Reason:  reason, // Asigna la razón al campo Reason
	}
}
