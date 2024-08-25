package stanza

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

/*

Al solicitar la configuracion de una sala, el servidor respondera con un formulario

En esta aplicación, solo se necesita que el room sea persistente

Example 159. Owner Submits Configuration Form¶
<iq from='crone1@shakespeare.lit/desktop'
    id='create2'
    to='coven@chat.shakespeare.lit'
    type='set'>
  <query xmlns='http://jabber.org/protocol/muc#owner'>
    <x xmlns='jabber:x:data' type='submit'>
      <field var='FORM_TYPE'>
        <value>http://jabber.org/protocol/muc#roomconfig</value>
      </field>
      <field var='muc#roomconfig_roomname'>
        <value>A Dark Cave</value>
      </field>
      <field var='muc#roomconfig_roomdesc'>
        <value>The place for all good witches!</value>
      </field>
      <field var='muc#roomconfig_enablelogging'>
        <value>0</value>
      </field>
      <field var='muc#roomconfig_changesubject'>
        <value>1</value>
      </field>
      <field var='muc#roomconfig_allowinvites'>
        <value>0</value>
      </field>
      <field var='muc#roomconfig_allowpm'>
        <value>anyone</value>
      </field>
      <field var='muc#roomconfig_maxusers'>
        <value>10</value>
      </field>
      <field var='muc#roomconfig_publicroom'>
        <value>0</value>
      </field>
      <field var='muc#roomconfig_persistentroom'>
        <value>0</value>
      </field>
      <field var='muc#roomconfig_moderatedroom'>
        <value>0</value>
      </field>
      <field var='muc#roomconfig_membersonly'>
        <value>0</value>
      </field>
      <field var='muc#roomconfig_passwordprotectedroom'>
        <value>1</value>
      </field>
      <field var='muc#roomconfig_roomsecret'>
        <value>cauldronburn</value>
      </field>
      <field var='muc#roomconfig_whois'>
        <value>moderators</value>
      </field>
      <field var='muc#maxhistoryfetch'>
        <value>50</value>
      </field>
      <field var='muc#roomconfig_roomadmins'>
        <value>wiccarocks@shakespeare.lit</value>
        <value>hecate@shakespeare.lit</value>
      </field>
    </x>
  </query>
</iq>

de lo anterior, solo nos interesa enviar el campo muc#roomconfig_persistentroom

crear un struct para el formulario de configuracion de la sala con el campo muc#roomconfig_persistentroom

*/

type RoomConfigForm struct {
	XMLName xml.Name `xml:"jabber:x:data x"`
	Type    string   `xml:"type,attr"`
	Field   []struct {
		Var   string `xml:"var,attr"`
		Value string `xml:"value"`
	} `xml:"field"`
}

func (RoomConfigForm) Name() string {
	return "RoomConfigForm"
}

func (m RoomConfigForm) Namespace() string {
	return m.XMLName.Space
}

func (m RoomConfigForm) GetSet() *stanza.ResultSet {
	return nil
}

func NewRoomConfigForm() RoomConfigForm {
	return RoomConfigForm{
		XMLName: xml.Name{Space: "jabber:x:data", Local: "x"},
		Type:    "submit",
		Field: []struct {
			Var   string `xml:"var,attr"`
			Value string `xml:"value"`
		}{
			{Var: "FORM_TYPE", Value: "http://jabber.org/protocol/muc#roomconfig"},
			{Var: "muc#roomconfig_persistentroom", Value: "1"},
		},
	}
}

type MUCOwnerWithForm struct {
	XMLName xml.Name       `xml:"http://jabber.org/protocol/muc#owner query"`
	X       RoomConfigForm `xml:"x"`
}

func (MUCOwnerWithForm) Name() string {
	return "MUCOwnerWithForm"
}

func (m MUCOwnerWithForm) Namespace() string {
	return m.XMLName.Space
}

func (MUCOwnerWithForm) GetSet() *stanza.ResultSet {
	return nil
}

func NewMUCOwnerWithForm() MUCOwnerWithForm {
	return MUCOwnerWithForm{
		XMLName: xml.Name{Space: "http://jabber.org/protocol/muc#owner", Local: "query"},
		X:       NewRoomConfigForm(),
	}
}
