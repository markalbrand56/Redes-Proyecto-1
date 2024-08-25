package stanza

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

/*

<presence
  xmlns="jabber:client"
  from="testingv3@conference.alumchat.lol/alb21004"
  to="alb21004@alumchat.lol/57a95075-0ca4-4356-b1c9-31d607f02c7e"
  >

  <x xmlns="http://jabber.org/protocol/muc#user">
	  <item
		  jid="alb21004@alumchat.lol/57a95075-0ca4-4356-b1c9-31d607f02c7e"
		  affiliation="owner"
		  role="moderator"/>
	  <status code="110"/>
	  <status code="100"/>
	  <status code="170"/>
	  <status code="201"/>
  </x>
</presence>

dentro de un presence vendra un x con el namespace http://jabber.org/protocol/muc#user
Hacer una estructura para manejarlo

*/

// MUCUser es una stanza utilizada para identificar la respuesta del servidor al enviar la solicitud de creación de una sala
type MUCUser struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#user x"`
	Item    struct {
		JID         string `xml:"jid,attr"`
		Affiliation string `xml:"affiliation,attr"`
		Role        string `xml:"role,attr"`
	} `xml:"item"`
	Status []struct {
		Code string `xml:"code,attr"`
	} `xml:"status"`
}

func (MUCUser) Name() string {
	return "MUCUser"
}

func (MUCUser) Namespace() string {
	return "http://jabber.org/protocol/muc#user"
}

func (MUCUser) GetSet() *stanza.ResultSet {
	return nil
}

func init() {
	stanza.TypeRegistry.MapExtension(stanza.PKTPresence, xml.Name{Space: "http://jabber.org/protocol/muc#user", Local: "x"}, MUCUser{})
}

/*
<iq from='crone1@shakespeare.lit/desktop'
    id='create1'
    to='coven@chat.shakespeare.lit'
    type='get'>
  <query xmlns='http://jabber.org/protocol/muc#owner'/>
</iq>
*/

// MUCOwnerGet es una stanza utilizada para obtener la configuración de una sala en creación
type MUCOwnerGet struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#owner query"`
}

func (MUCOwnerGet) Name() string {
	return "MUCOwnerGet"
}

func (m MUCOwnerGet) Namespace() string {
	return m.XMLName.Space
}

func (MUCOwnerGet) GetSet() *stanza.ResultSet {
	return nil
}

// NewMUCOwnerGet crea una stanza para obtener la configuración de una sala
func NewMUCOwnerGet() MUCOwnerGet {
	return MUCOwnerGet{
		XMLName: xml.Name{Space: "http://jabber.org/protocol/muc#owner", Local: "query"},
	}
}
