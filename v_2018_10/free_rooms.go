package v_2018_10

import (
	"encoding/xml"

	"github.com/HGV/alpinebits/version"
	"github.com/HGV/x/timex"
)

type HotelAvailNotifRQ struct {
	XMLName             xml.Name            `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelAvailNotifRQ"`
	UniqueID            *UniqueID           `xml:"UniqueID,omitempty"`
	AvailStatusMessages AvailStatusMessages `xml:"AvailStatusMessages"`
}

type UniqueIDType int

const (
	// UniqueIDTypeReservation           UniqueIDType = 14
	// UniqueIDTypeCancellation          UniqueIDType = 15
	// UniqueIDTypeReference             UniqueIDType = 16
	UniqueIDTypePurgedMasterReference UniqueIDType = 35
)

type Instance string

const (
	InstanceCompleteSet Instance = "CompleteSet"
)

type UniqueID struct {
	Type     UniqueIDType `xml:"Type,attr"`
	ID       string       `xml:"ID,attr"`
	Instance Instance     `xml:"Instance,attr,omitempty"`
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
	response

	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelAvailNotifRS"`
	Version string   `xml:"Version,attr"`
}
