package models

import "gosrc.io/xmpp"

type User struct {
	Client   *xmpp.Client
	Contacts []string
}

func NewUser(client *xmpp.Client) *User {
	return &User{
		Client: client,
	}
}
