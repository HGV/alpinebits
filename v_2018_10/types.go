package v_2018_10

import (
	"encoding/xml"

	"github.com/HGV/alpinebits/version"
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

func (err Error) Error() string {
	return err.Value
}

type response struct {
	Success  *Success   `xml:"Success"`
	Warnings *[]Warning `xml:"Warnings>Warning"`
	Errors   *[]Error   `xml:"Errors>Error"`
}

type PingRQ struct {
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRQ"`
	Version  string   `xml:"Version,attr"`
	EchoData string   `xml:",innerxml"`
}

type PingRS struct {
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRS"`
	Version  string   `xml:"Version,attr"`
	Success  Success  `xml:"Success"`
	Warnings Warning  `xml:"Warnings>Warning"`
	EchoData string   `xml:",innerxml"`
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

type HotelAvailNotifRQ struct {
	XMLName             xml.Name            `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelAvailNotifRQ"`
	UniqueID            *UniqueID           `xml:"UniqueID,omitempty"`
	AvailStatusMessages AvailStatusMessages `xml:"AvailStatusMessages"`
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
