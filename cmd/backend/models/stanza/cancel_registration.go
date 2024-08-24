package stanza

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

/*
<iq type='set' from='bill@shakespeare.lit/globe' id='unreg1'>
  <query xmlns='jabber:iq:register'>
    <remove/>
  </query>
</iq>
*/

type CancelRegistration struct {
	XMLName xml.Name       `xml:"iq"`
	Type    string         `xml:"type,attr"`
	From    string         `xml:"from,attr"`
	Id      string         `xml:"id,attr"`
	Query   CancelingQuery `xml:"query"`
}

type CancelingQuery struct {
	XMLName xml.Name `xml:"query"`
	Xmlns   string   `xml:"xmlns,attr"`
	Remove  string   `xml:"remove"`
}

func (c CancelRegistration) Name() string {
	return "CancelRegistration"
}

func (c CancelRegistration) Namespace() string {
	return c.XMLName.Space
}

func (c CancelRegistration) GetSet() *stanza.ResultSet {
	return nil
}

func NewCancelRegistration(from string) CancelRegistration {
	return CancelRegistration{
		Type: "set",
		From: from,
		Id:   "unreg1",
		Query: CancelingQuery{
			Xmlns:  "jabber:iq:register",
			Remove: "",
		},
	}
}
