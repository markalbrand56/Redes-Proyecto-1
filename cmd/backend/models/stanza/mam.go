package stanza

import (
	"encoding/xml"
	"fmt"
	"gosrc.io/xmpp/stanza"
)

// <message xmlns="jabber:client" to="alb21004@alumchat.lol/82dd3a98-eee3-45bb-9d84-c819728d3fc8"><result xmlns="urn:xmpp:mam:2" id="3516"><forwarded xmlns="urn:xmpp:forward:0"><delay xmlns="urn:xmpp:delay" stamp="2024-08-14T05:23:41.310Z"/><message xmlns="jabber:client" to="alb21004@alumchat.lol" type="chat" id="b5e21f08-8b2c-4854-9227-df084324922b" from="alb21005@alumchat.lol/gajim.0O3D5ZZ0"><body>hola</body><origin-id xmlns="urn:xmpp:sid:0" id="b5e21f08-8b2c-4854-9227-df084324922b"></origin-id><request xmlns="urn:xmpp:receipts"></request><markable xmlns="urn:xmpp:chat-markers:0"></markable></message></forwarded></result></message>

// MAM es una estructura que representa un mensaje de tipo IQ con namespace urn:xmpp:mam:2 para solicitar el historial de mensajes
type MAM struct {
	XMLName   xml.Name  `xml:"urn:xmpp:mam:2 result"`
	ID        string    `xml:"id,attr"`
	Forwarded Forwarded `xml:"urn:xmpp:forward:0 forwarded"`
}

type Forwarded struct {
	XMLName xml.Name       `xml:"forwarded"`
	Delay   Delay          `xml:"urn:xmpp:delay delay"`
	Message stanza.Message `xml:"jabber:client message"`
}

type Delay struct {
	XMLName xml.Name `xml:"delay"`
	Stamp   string   `xml:"stamp,attr"`
}

func (m MAM) String() string {
	return fmt.Sprintf("{From: %s, To: %s, Body: %s, Timestamp: %s}", m.Forwarded.Message.From, m.Forwarded.Message.To, m.Forwarded.Message.Body, m.Forwarded.Delay.Stamp)
}

func (m MAM) Name() string {
	return "MAM"
}

func (m MAM) Namespace() string {
	return m.XMLName.Space
}

func (m MAM) GetSet() *stanza.ResultSet {
	return nil
}

func init() {
	stanza.TypeRegistry.MapExtension(stanza.PKTMessage, xml.Name{Space: "urn:xmpp:mam:2", Local: "result"}, MAM{})
}
