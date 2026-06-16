package freerooms

import (
	"encoding/xml"

	"github.com/HGV/alpinebits"
)

type HotelInvCountNotifRQ struct {
	XMLName     xml.Name    `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelInvCountNotifRQ"`
	Version     string      `xml:"Version,attr"`
	Inventories Inventories `xml:"Inventories"`
}

// HotelCode implements alpinebits.HotelCoded.
func (r HotelInvCountNotifRQ) HotelCode() string {
	return r.Inventories.HotelCode
}

type Inventories struct {
	HotelCode string `xml:"HotelCode,attr"`
}

type HotelInvCountNotifRS struct {
	XMLName  xml.Name             `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelInvCountNotifRS"`
	Version  string               `xml:"Version,attr"`
	Success  *alpinebits.Success  `xml:"Success"`
	Warnings *alpinebits.Warnings `xml:"Warnings>Warning"`
	Errors   *alpinebits.Errors   `xml:"Errors>Error"`
}
