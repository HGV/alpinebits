package v_2024_10

import "encoding/xml"

type PingRQ struct {
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRQ"`
	Version  string   `xml:"Version,attr"`
	EchoData EchoData `xml:"EchoData"`
}

type EchoData struct {
	Value string `xml:",innerxml"`
}
