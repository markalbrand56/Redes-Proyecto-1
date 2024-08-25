package stanza

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

/*
Example 1. Entity Requests Registration Fields from Host¶
<iq type='get' id='reg1' to='shakespeare.lit'>
  <query xmlns='jabber:iq:register'/>
</iq>

*/

type RegisterQuery struct {
	XMLName xml.Name `xml:"query"`
	XMLNS   string   `xml:"xmlns,attr"`
}

func (r RegisterQuery) Name() string {
	return "query"
}

func (r RegisterQuery) Namespace() string {
	return r.XMLName.Space
}

func (r RegisterQuery) GetSet() *stanza.ResultSet {
	return nil
}

// NewRegisterQuery creates a new IQ message with namespace jabber:iq:register to request registration fields
func NewRegisterQuery() RegisterQuery {
	return RegisterQuery{
		XMLNS: "jabber:iq:register",
	}
}

/*
<iq type='set' id='reg2'>
  <query xmlns='jabber:iq:register'>
    <username>bill</username>
    <password>Calliope</password>
    <email>bard@shakespeare.lit</email>
  </query>
</iq>
*/

type RegisterQueryWithUser struct {
	XMLName  xml.Name `xml:"query"`
	XMLNS    string   `xml:"xmlns,attr"`
	Username string   `xml:"username"`
	Password string   `xml:"password"`
	Email    string   `xml:"email"`
}

func (r RegisterQueryWithUser) Name() string {
	return "query"
}

func (r RegisterQueryWithUser) Namespace() string {
	return r.XMLName.Space
}

func (r RegisterQueryWithUser) GetSet() *stanza.ResultSet {
	return nil
}

// NewRegisterQueryWithUser creates a new IQ message with namespace jabber:iq:register to register a new user
func NewRegisterQueryWithUser(username string, password string, email string) RegisterQueryWithUser {
	return RegisterQueryWithUser{
		XMLNS:    "jabber:iq:register",
		Username: username,
		Password: password,
		Email:    email,
	}
}
