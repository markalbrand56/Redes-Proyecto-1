package stanza

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

/*

Example 201. Owner Submits Room Destruction Request¶
<iq from='crone1@shakespeare.lit/desktop'
    id='begone'
    to='heath@chat.shakespeare.lit'
    type='set'>
  <query xmlns='http://jabber.org/protocol/muc#owner'>
    <destroy jid='coven@chat.shakespeare.lit'>
      <reason>Macbeth doth come.</reason>
    </destroy>
  </query>
</iq>

*/

type RoomDestruction struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#owner query"`
	Destroy struct {
		Jid    string `xml:"jid,attr"`
		Reason string `xml:"reason"`
	} `xml:"destroy"`
}

func (RoomDestruction) Name() string {
	return "RoomDestruction"
}

func (m RoomDestruction) Namespace() string {
	return m.XMLName.Space
}

func (RoomDestruction) GetSet() *stanza.ResultSet {
	return nil
}

// NewRoomDestruction creates a new RoomDestruction stanza
func NewRoomDestruction(jid, reason string) RoomDestruction {
	return RoomDestruction{
		XMLName: xml.Name{Space: "http://jabber.org/protocol/muc#owner", Local: "query"},
		Destroy: struct {
			Jid    string `xml:"jid,attr"`
			Reason string `xml:"reason"`
		}{
			Jid:    jid,
			Reason: reason,
		},
	}
}
