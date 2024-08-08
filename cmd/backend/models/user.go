package models

import "gosrc.io/xmpp"

type User struct {
	Client   *xmpp.Client
	UserName string
	Contacts []string
}

func NewUser(client *xmpp.Client, username string) *User {
	return &User{
		UserName: username,
		Client:   client,
	}
}
