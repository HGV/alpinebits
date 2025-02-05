package freerooms

import (
	"encoding/xml"

	"github.com/HGV/alpinebits/v_2020_10/common"
	"github.com/HGV/x/timex"
)

type UniqueIDType int

const UniqueIDTypePurgedMasterReference UniqueIDType = 35

type UniqueIDInstance string

const UniqueIDInstanceCompleteSet UniqueIDInstance = "CompleteSet"

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
	common.Response

	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelInvCountNotifRS"`
	Version string   `xml:"Version,attr"`
}
