package freerooms

import (
	"encoding/xml"

	"github.com/HGV/alpinebits/v_2020_10/common"
	"github.com/HGV/alpinebits/version"
	"github.com/HGV/x/timex"
)

type UniqueIDType int

const (
	UniqueIDTypeReference             UniqueIDType = 16
	UniqueIDTypePurgedMasterReference UniqueIDType = 35
)

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

var _ version.HotelCodeProvider = (*HotelInvCountNotifRQ)(nil)

func (h HotelInvCountNotifRQ) HotelCode() string {
	return h.Inventories.HotelCode
}

type Inventories struct {
	HotelCode   string      `xml:"HotelCode,attr"`
	HotelName   string      `xml:"HotelName,attr"`
	Inventories []Inventory `xml:"Inventory"`
}

func (i Inventories) IsReset() bool {
	var zero Inventory
	return len(i.Inventories) == 1 &&
		i.Inventories[0] == zero
}

type Inventory struct {
	StatusApplicationControl *StatusApplicationControl `xml:"StatusApplicationControl,omitempty"`
	InvCounts                *[]InvCount               `xml:"InvCounts>InvCount"`
}

var _ version.DateRangeProvider = (*Inventory)(nil)

func (i Inventory) DateRange() timex.DateRange {
	return timex.DateRange{
		Start: i.StatusApplicationControl.Start,
		End:   i.StatusApplicationControl.End,
	}
}

func (i Inventory) isAvailability() bool {
	return !i.StatusApplicationControl.AllInvCode
}

func (i Inventory) isClosingSeason() bool {
	return i.StatusApplicationControl.AllInvCode
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
