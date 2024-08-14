package stanza

import (
	"encoding/xml"
	"fmt"
	"gosrc.io/xmpp/stanza"
)

type Archive struct {
	XMLName xml.Name `xml:"urn:xmpp:mam:2 query"`
	Type    string   `xml:"type,attr"`
	ID      string   `xml:"id,attr"`
	Query   Query    `xml:"urn:xmpp:mam:2 query"`
}

type Query struct {
	XMLName xml.Name `xml:"query"`
	X       X        `xml:"jabber:x:data x"`
	Set     Set      `xml:"http://jabber.org/protocol/rsm set"`
}

type X struct {
	XMLName xml.Name `xml:"x"`
	XMLNS   string   `xml:"xmlns,attr"`
	Type    string   `xml:"type,attr"`
	Field   []Field  `xml:"field"`
}

type Field struct {
	XMLName xml.Name `xml:"field"`
	Var     string   `xml:"var,attr"`
	Type    string   `xml:"type,attr"`
	Value   string   `xml:"value"`
}

type Set struct {
	XMLName xml.Name `xml:"set"`
	XMLNS   string   `xml:"xmlns,attr"`
	Max     Max      `xml:"max"`
}

type Max struct {
	XMLName xml.Name `xml:"max"`
	Value   string   `xml:",chardata"`
}

func (a Archive) Name() string {
	return "Archive"
}

func (a Archive) Namespace() string {
	return a.XMLName.Space
}

func (a Archive) GetSet() *stanza.ResultSet {
	return nil
}

func NewArchiveQuery(jid string, max int) Archive {
	return Archive{
		XMLName: xml.Name{Space: "urn:xmpp:mam:2", Local: "query"},
		Type:    "set",
		ID:      "mam_query_1",
		Query: Query{
			XMLName: xml.Name{Local: "query"},
			X: X{
				XMLName: xml.Name{Local: "x"},
				XMLNS:   "jabber:x:data",
				Type:    "submit",
				Field: []Field{
					{
						XMLName: xml.Name{Local: "field"},
						Var:     "FORM_TYPE",
						Type:    "hidden",
						Value:   "urn:xmpp:mam:2",
					},
					{
						XMLName: xml.Name{Local: "field"},
						Var:     "with",
						Value:   jid,
					},
				},
			},
			Set: Set{
				XMLName: xml.Name{Local: "set"},
				XMLNS:   "http://jabber.org/protocol/rsm",
				Max: Max{
					XMLName: xml.Name{Local: "max"},
					Value:   fmt.Sprintf("%d", max),
				},
			},
		},
	}
}

func init() {
	stanza.TypeRegistry.MapExtension(stanza.PKTIQ, xml.Name{Space: "urn:xmpp:mam:2", Local: "query"}, Archive{})
}
