package freerooms

import (
	"encoding/xml"

	"github.com/HGV/alpinebits/v_2018_10/common"
	"github.com/HGV/alpinebits/version"
	"github.com/HGV/x/timex"
)

type HotelAvailNotifRQ struct {
	XMLName             xml.Name            `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelAvailNotifRQ"`
	UniqueID            *UniqueID           `xml:"UniqueID,omitempty"`
	AvailStatusMessages AvailStatusMessages `xml:"AvailStatusMessages"`
}

var _ version.HotelCodeProvider = (*HotelAvailNotifRQ)(nil)

func (h HotelAvailNotifRQ) HotelCode() string {
	return h.AvailStatusMessages.HotelCode
}

type UniqueIDType int

const (
	UniqueIDTypeReference             UniqueIDType = 16
	UniqueIDTypePurgedMasterReference UniqueIDType = 35
)

type Instance string

const (
	InstanceCompleteSet Instance = "CompleteSet"
)

type UniqueID struct {
	Type     UniqueIDType `xml:"Type,attr"`
	ID       string       `xml:"ID,attr"`
	Instance Instance     `xml:"Instance,attr"`
}

type AvailStatusMessages struct {
	HotelCode           string               `xml:"HotelCode,attr"`
	HotelName           string               `xml:"HotelName,attr"`
	AvailStatusMessages []AvailStatusMessage `xml:"AvailStatusMessage"`
}

func (a AvailStatusMessages) IsReset() bool {
	var zero AvailStatusMessage
	return len(a.AvailStatusMessages) == 1 &&
		a.AvailStatusMessages[0] == zero
}

type BookingLimitMessageType string

const (
	BookingLimitMessageTypeSetLimit BookingLimitMessageType = "SetLimit"
)

type AvailStatusMessage struct {
	BookingLimit             int                      `xml:"BookingLimit,attr"`
	BookingLimitMessageType  BookingLimitMessageType  `xml:"BookingLimitMessageType,attr"`
	BookingThreshold         int                      `xml:"BookingThreshold,attr"`
	StatusApplicationControl StatusApplicationControl `xml:"StatusApplicationControl"`
}

var _ version.DateRangeProvider = (*AvailStatusMessage)(nil)

func (s AvailStatusMessage) DateRange() timex.DateRange {
	return timex.DateRange{
		Start: s.StatusApplicationControl.Start,
		End:   s.StatusApplicationControl.End,
	}
}

type StatusApplicationControl struct {
	Start       timex.Date `xml:"Start,attr"`
	End         timex.Date `xml:"End,attr"`
	InvTypeCode string     `xml:"InvTypeCode,attr,omitempty"`
	InvCode     string     `xml:"InvCode,attr,omitempty"`
}

type HotelAvailNotifRS struct {
	common.Response

	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelAvailNotifRS"`
	Version string   `xml:"Version,attr"`
}
