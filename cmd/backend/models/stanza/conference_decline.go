package stanza

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

/*
<message xmlns="jabber:client" from="testingv2@conference.alumchat.lol" to="alb210041@alumchat.lol"><x xmlns="http://jabber.org/protocol/muc#user"><decline from="alb21005@alumchat.lol"/></x></message>
*/

// ConferenceDeclineMessage es una estructura que representa un mensaje de tipo Message con namespace http://jabber.org/protocol/muc#user para declinar una invitación a una sala de chat
type ConferenceDeclineMessage struct {
	XMLName           xml.Name          `xml:"message"`
	XMLNS             string            `xml:"xmlns,attr"`
	From              string            `xml:"from,attr"`
	To                string            `xml:"to,attr"`
	ConferenceDecline ConferenceDecline `xml:"x"`
}

func (c ConferenceDeclineMessage) Name() string {
	return "ConferenceDeclineMessage"
}

func (c ConferenceDeclineMessage) Namespace() string {
	return c.XMLNS
}

func (c ConferenceDeclineMessage) GetSet() *stanza.ResultSet {
	return nil
}

// NewConferenceDeclineMessage crea un nuevo mensaje de tipo Message con namespace http://jabber.org/protocol/muc#user para declinar una invitación a una sala de chat
func NewConferenceDeclineMessage(conference, from, to string) ConferenceDeclineMessage {
	return ConferenceDeclineMessage{
		XMLName: xml.Name{Space: "jabber:client", Local: "message"},
		XMLNS:   "jabber:client",
		From:    conference, // JID de la conferencia a la que se declina
		To:      to,         // JID del usuario que invitó
		ConferenceDecline: ConferenceDecline{
			XMLName: xml.Name{Space: "http://jabber.org/protocol/muc#user", Local: "x"},
			XMLNS:   "http://jabber.org/protocol/muc#user",
			Decline: Decline{
				From: from, // JID del usuario que declina
			},
		},
	}
}

// ConferenceDecline es una estructura que representa un mensaje de tipo Message con namespace http://jabber.org/protocol/muc#user para declinar una invitación a una sala de chat
type ConferenceDecline struct {
	XMLName xml.Name `xml:"x"`
	XMLNS   string   `xml:"xmlns,attr"`
	Decline Decline  `xml:"decline"`
}

type Decline struct {
	From string `xml:"from,attr"`
}

// Name devuelve el nombre del tipo de archivo
func (f ConferenceDecline) Name() string {
	return "ConferenceDecline"
}

// Namespace devuelve el espacio de nombres para el archivo
func (f ConferenceDecline) Namespace() string {
	return f.XMLName.Space
}

// GetSet no es necesario en este contexto
func (f ConferenceDecline) GetSet() *stanza.ResultSet {
	return nil
}
