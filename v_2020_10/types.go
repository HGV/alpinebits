package v_2020_10

import (
	"encoding/xml"

	"github.com/HGV/x/timex"
)

type ErrorWarningType int

const (
	ErrorWarningTypeAdvisory         ErrorWarningType = 11
	ErrorWarningTypeApplicationError ErrorWarningType = 13
)

type Status string

const (
	StatusSendInventory Status = "ALPINEBITS_SEND_INVENTORY"
	StatusSendFreeRooms Status = "ALPINEBITS_SEND_FREEROOMS"
	StatusSendRatePlans Status = "ALPINEBITS_SEND_RATEPLANS"
)

type Success struct{}

type Warning struct {
	Type   ErrorWarningType `xml:"Type,attr"`
	Code   int              `xml:"Code,attr,omitempty"`
	Status Status           `xml:"Status,attr,omitempty"`
	Value  string           `xml:",innerxml"`
}

type Error struct {
	Type   ErrorWarningType `xml:"Type,attr"`
	Code   int              `xml:"Code,attr,omitempty"`
	Status Status           `xml:"Status,attr,omitempty"`
	Value  string           `xml:",innerxml"`
}

type response struct {
	Success  *Success   `xml:"Success"`
	Warnings *[]Warning `xml:"Warnings>Warning"`
	Errors   *[]Error   `xml:"Errors>Error"`
}

type PingRQ struct {
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRQ"`
	Version  string   `xml:"Version,attr"`
	EchoData EchoData `xml:"EchoData"`
}

type PingRS struct {
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRS"`
	Version  string   `xml:"Version,attr"`
	Success  Success  `xml:"Success"`
	Warnings Warning  `xml:"Warnings>Warning"`
	EchoData EchoData `xml:"EchoData"`
}

type EchoData struct {
	Value string `xml:",innerxml"`
}

type UniqueIDType uint8

const (
	// UniqueIDTypeReservation  UniqueIDType = 14
	// UniqueIDTypeCancellation UniqueIDType = 15
	// UniqueIDTypeReference    UniqueIDType = 16
	UniqueIDTypePurgedMasterReference UniqueIDType = 35
)

type UniqueIDInstance string

const (
	UniqueIDInstanceCompleteSet UniqueIDInstance = "CompleteSet"
)

type UniqueID struct {
	Type     UniqueIDType     `xml:"Type,attr"`
	ID       string           `xml:"ID,attr"`
	Instance UniqueIDInstance `xml:"Instance,attr,omitempty"`
}

type HotelInvCountNotifRQ struct {
	XMLName     xml.Name    `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelInvCountNotifRQ"`
	Version     string      `xml:"Version,attr"`
	UniqueID    *UniqueID   `xml:"UniqueID,omitempty"`
	Inventories Inventories `xml:"Inventories"`
}

type Inventories struct {
	HotelCode   string      `xml:"HotelCode,attr"`
	HotelName   string      `xml:"HotelName,attr"`
	Inventories []Inventory `xml:"Inventory"`
}

type Inventory struct {
	StatusApplicationControl *StatusApplicationControl `xml:"StatusApplicationControl,omitempty"`
	InvCounts                *[]InvCount               `xml:"InvCounts>InvCount"`
}

type StatusApplicationControl struct {
	Start       timex.Date `xml:"Start,attr"`
	End         timex.Date `xml:"End,attr"`
	InvTypeCode string     `xml:"InvTypeCode,attr,omitempty"`
	InvCode     string     `xml:"InvCode,attr,omitempty"`
	AllInvCode  bool       `xml:"AllInvCode,attr,omitempty"`
}

type CountType int

const (
	CountTypeBookable   CountType = 2
	CountTypeOutOfOrder CountType = 6
	CountTypeFree       CountType = 9
)

type InvCount struct {
	CountType CountType `xml:"CountType,attr"`
	Count     int       `xml:"Count,attr"`
}

type HotelInvCountNotifRS struct {
	response

	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelInvCountNotifRS"`
	Version string   `xml:"Version,attr"`
}
