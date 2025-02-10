package handshake

import (
	"encoding/xml"

	"github.com/HGV/alpinebits/v_2020_10/common"
)

type PingRQ struct {
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRQ"`
	Version  string   `xml:"Version,attr"`
	EchoData EchoData `xml:"EchoData"`
}

type EchoData struct {
	Value string `xml:",innerxml"`
}

type PingRS struct {
	XMLName  xml.Name       `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRS"`
	Version  string         `xml:"Version,attr"`
	Success  common.Success `xml:"Success"`
	Warnings common.Warning `xml:"Warnings>Warning"`
	EchoData EchoData       `xml:"EchoData"`
}
