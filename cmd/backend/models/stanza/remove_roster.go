package stanza

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

/*
	<iq type="set" id="roster_1" xmlns="jabber:client">
	  <query xmlns="jabber:iq:roster">
	    <item jid="contact@example.com" subscription="remove"/>
	  </query>
	</iq>
*/

// RosterRemove es una estructura que representa un mensaje de tipo IQ con namespace jabber:iq:roster para eliminar un contacto del roster
type RosterRemove struct {
	XMLName xml.Name `xml:"iq"`
	Type    string   `xml:"type,attr"`
	ID      string   `xml:"id,attr"`
	Query   Query2   `xml:"query"`
}

type Query2 struct {
	XMLName xml.Name `xml:"query"`
	XMLNS   string   `xml:"xmlns,attr"`
	Item    Item     `xml:"item"`
}

type Item struct {
	XMLName      xml.Name `xml:"item"`
	JID          string   `xml:"jid,attr"`
	Subscription string   `xml:"subscription,attr"`
}

func (r RosterRemove) Name() string {
	return "RosterRemove"
}

func (r RosterRemove) Namespace() string {
	return r.XMLName.Space
}

func (r RosterRemove) GetSet() *stanza.ResultSet {
	return nil
}

// NewRosterRemove crea un nuevo mensaje de tipo IQ con namespace jabber:iq:roster para eliminar un contacto del roster
func NewRosterRemove(jid string) RosterRemove {
	return RosterRemove{
		XMLName: xml.Name{Space: "jabber:client", Local: "iq"},
		Type:    "set",
		ID:      "roster_1",
		Query: Query2{
			XMLName: xml.Name{Space: "jabber:iq:roster", Local: "query"},
			XMLNS:   "jabber:iq:roster",
			Item: Item{
				XMLName:      xml.Name{Space: "jabber:iq:roster", Local: "item"},
				JID:          jid,
				Subscription: "remove",
			},
		},
	}
}
