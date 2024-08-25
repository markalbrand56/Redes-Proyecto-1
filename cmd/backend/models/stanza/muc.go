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

/* If the initial room owner wants to accept the default room configuration (i.e., create an "instant room"), the room owner MUST decline an initial configuration form by sending an IQ set to the <room@service> itself containing a <query/> element qualified by the 'http://jabber.org/protocol/muc#owner' namespace, where the only child of the <query/> is an empty <x/> element that is qualified by the 'jabber:x:data' namespace and that possesses a 'type' attribute whose value is "submit":
Example 155. Owner Requests Instant Room¶
<iq from='crone1@shakespeare.lit/desktop'
    id='create1'
    to='coven@chat.shakespeare.lit'
    type='set'>
  <query xmlns='http://jabber.org/protocol/muc#owner'>
    <x xmlns='jabber:x:data' type='submit'/>
  </query>
*/

//type MUCOwner struct {
//	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#owner query"`
//	X       struct {
//		XMLName xml.Name `xml:"jabber:x:data x"`
//		Type    string   `xml:"type,attr"`
//	} `xml:"x"`
//}
//
//func (MUCOwner) Name() string {
//	return "MUCOwner"
//}
//
//func (m MUCOwner) Namespace() string {
//	return m.XMLName.Space
//}
//
//func (MUCOwner) GetSet() *stanza.ResultSet {
//	return nil
//}
//
//func NewMUCOwner() MUCOwner {
//	return MUCOwner{
//		XMLName: xml.Name{Space: "http://jabber.org/protocol/muc#owner", Local: "query"},
//		X: struct {
//			XMLName xml.Name `xml:"jabber:x:data x"`
//			Type    string   `xml:"type,attr"`
//		}{
//			XMLName: xml.Name{Space: "jabber:x:data", Local: "x"},
//			Type:    "submit",
//		},
//	}
//}

/*
<iq from='crone1@shakespeare.lit/desktop'
    id='create1'
    to='coven@chat.shakespeare.lit'
    type='get'>
  <query xmlns='http://jabber.org/protocol/muc#owner'/>
</iq>
*/

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

func NewMUCOwnerGet() MUCOwnerGet {
	return MUCOwnerGet{
		XMLName: xml.Name{Space: "http://jabber.org/protocol/muc#owner", Local: "query"},
	}
}

type MUCOwnerWithRoomConfig struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#owner query"`
	X       struct {
		XMLName xml.Name `xml:"jabber:x:data x"`
		Type    string   `xml:"type,attr"`
	} `xml:"x"`
}
