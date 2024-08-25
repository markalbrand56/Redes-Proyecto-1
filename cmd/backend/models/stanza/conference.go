package stanza

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

// Ejemplo de la Stanza a la que se refiere este archivo:
// <message xmlns="jabber:client" to="alb21004@alumchat.lol" id="9946c7cb-b8fe-4214-9314-cc7ac91e1ab9" from="alb21005@alumchat.lol/gajim.0O3D5ZZ0">
//		<x xmlns="jabber:x:conference" jid="ogivox@conference.alumchat.lol"></x>
// </message>

// Conference es una estructura que representa un mensaje de tipo Message con namespace jabber:x:conference para manejar las invitaciones a salas de chat
type Conference struct {
	XMLName xml.Name `xml:"jabber:x:conference x"`
	JID     string   `xml:"jid,attr"`
}

func (c Conference) Name() string {
	return "Conference"
}

func (c Conference) Namespace() string {
	return c.XMLName.Space
}

func (c Conference) GetSet() *stanza.ResultSet {
	return nil
}

// NewConferenceInvitation crea un nuevo mensaje de tipo Message con namespace jabber:x:conference para invitar a un usuario a una sala de chat
func (c Conference) NewConferenceInvitation(jid string) Conference {
	return Conference{
		XMLName: xml.Name{Space: "jabber:x:conference", Local: "x"},
		JID:     jid,
	}
}

func init() {
	stanza.TypeRegistry.MapExtension(stanza.PKTMessage, xml.Name{Space: "jabber:x:conference", Local: "x"}, Conference{})
}
