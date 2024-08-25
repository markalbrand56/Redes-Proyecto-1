package stanza

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

//<iq xmlns="jabber:client" type="get" id="271-343276" from="alumchat.lol" to="alb21004@alumchat.lol/50eb7f07-6bf0-47e7-91e2-a9f67ad7711d"><ping xmlns="urn:xmpp:ping"/></iq>

// Ping representa una solicitud de ping en XMPP.
type Ping struct {
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}

// Name devuelve el nombre de la extensión de la interfaz de paquetes.
func (Ping) Name() string {
	return "Ping"
}

// Namespace devuelve el espacio de nombres de la extensión de la interfaz de paquetes.
func (Ping) Namespace() string {
	return "urn:xmpp:ping"
}

// GetSet devuelve el conjunto de resultados de la extensión de la interfaz de paquetes.
func (Ping) GetSet() *stanza.ResultSet {
	return nil
}

// NewPing crea una nueva instancia de Ping.
func NewPing() Ping {
	return Ping{
		XMLName: xml.Name{Space: "urn:xmpp:ping", Local: "ping"},
	}
}

func init() {
	stanza.TypeRegistry.MapExtension(stanza.PKTIQ, xml.Name{Space: "urn:xmpp:ping", Local: "ping"}, Ping{})
}
