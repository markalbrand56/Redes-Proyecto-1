package stanza

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

/*
Las definiciones de este modelo se usan para manejar las invitaciones a conferencias y las afiliaciones en salas de conferencia MUC
Las afiliaciones son necesarias de manejar al momento de invitar a un usuario a una sala de conferencia

<message
	 from="alb21004@alumchat.lol"
	 to="alb210041@alumchat.lol"
	 >
	 <x
		 xmlns="jabber:x:conference"
		 jid="ogivox@conference.alumchat.lol"
		 reason="alb21004@alumchat.lol has invited you to the conference &#39;ogivox@conference.alumchat.lol&#39;">
	 </x>

	 <x xmlns="http://jabber.org/protocol/muc#user">
		 <invite from="alb21004@alumchat.lol"></invite>
	 </x>
 </message>
*/

// ConferenceInvite es una estructura que representa un mensaje de tipo Message con namespace jabber:x:conference para manejar las invitaciones a salas de chat
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

// NewConferenceInvite crea una nueva instancia de la estructura ConferenceInvite con un JID y una razón
func NewConferenceInvite(jid string, reason string) ConferenceInvite {
	return ConferenceInvite{
		XMLName: xml.Name{Space: "jabber:x:conference", Local: "x"},
		XMLNS:   "jabber:x:conference",
		JID:     jid,    // Asigna el JID al campo JID
		Reason:  reason, // Asigna la razón al campo Reason
	}
}

type MucInvite struct {
	XMLName xml.Name `xml:"x"`
	XMLNS   string   `xml:"xmlns,attr"`
	Invite  Invite   `xml:"invite"`
}

type Invite struct {
	From string `xml:"from,attr"`
}

// Name devuelve el nombre del tipo de archivo
func (f MucInvite) Name() string {
	return "MucInvite"
}

// Namespace devuelve el espacio de nombres para el archivo
func (f MucInvite) Namespace() string {
	return f.XMLName.Space
}

// GetSet no es necesario en este contexto
func (f MucInvite) GetSet() *stanza.ResultSet {
	return nil
}

// NewMucInvite crea una nueva instancia de la estructura MucInvite con un JID, contraseña y razón
func NewMucInvite(from string) MucInvite {
	return MucInvite{
		XMLName: xml.Name{Space: "http://jabber.org/protocol/muc#user", Local: "x"},
		XMLNS:   "http://jabber.org/protocol/muc#user",
		Invite: Invite{
			From: from,
		},
	}
}

// MUCAdminAffiliation representa una afiliación en una sala de conferencia MUC
type MUCAdminAffiliation struct {
	XMLName     xml.Name `xml:"item"`
	JID         string   `xml:"jid,attr"`
	Affiliation string   `xml:"affiliation,attr"`
}

// MUCAffiliationRequest es la estructura para enviar una solicitud de afiliación
type MUCAffiliationRequest struct {
	XMLName xml.Name              `xml:"query"`
	XMLNS   string                `xml:"xmlns,attr"`
	Items   []MUCAdminAffiliation `xml:"item"`
}

// Name devuelve el nombre del tipo de archivo
func (f MUCAffiliationRequest) Name() string {
	return "MUCAffiliationRequest"
}

// Namespace devuelve el espacio de nombres para el archivo
func (f MUCAffiliationRequest) Namespace() string {
	return f.XMLName.Space
}

// GetSet no es necesario en este contexto
func (f MUCAffiliationRequest) GetSet() *stanza.ResultSet {
	return nil
}

// NewMUCAffiliationRequest crea una nueva solicitud para cambiar la afiliación en una sala MUC
func NewMUCAffiliationRequest(jid, affiliation string) MUCAffiliationRequest {
	return MUCAffiliationRequest{
		XMLName: xml.Name{Local: "query"},
		XMLNS:   "http://jabber.org/protocol/muc#admin",
		Items: []MUCAdminAffiliation{
			{JID: jid, Affiliation: affiliation},
		},
	}
}
