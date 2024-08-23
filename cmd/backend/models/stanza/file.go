package stanza

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

// File representa una estructura para enviar una URL como archivo
type File struct {
	XMLName xml.Name `xml:"x"`
	XMLNS   string   `xml:"xmlns,attr"`
	URL     URL      `xml:"url"`
}

// URL representa la URL del archivo
type URL struct {
	XMLName xml.Name `xml:"url"`
	Value   string   `xml:",chardata"` // El contenido de la URL
}

// Name devuelve el nombre del tipo de archivo
func (f File) Name() string {
	return "File"
}

// Namespace devuelve el espacio de nombres para el archivo
func (f File) Namespace() string {
	return f.XMLName.Space
}

// GetSet no es necesario en este contexto
func (f File) GetSet() *stanza.ResultSet {
	return nil
}

// NewFile crea una nueva instancia de la estructura File con una URL
func NewFile(url string) File {
	return File{
		XMLName: xml.Name{Space: "jabber:x:oob", Local: "x"},
		XMLNS:   "jabber:x:oob",
		URL: URL{
			XMLName: xml.Name{Space: "jabber:x:oob", Local: "url"},
			Value:   url, // Asigna la URL al campo Value
		},
	}
}
