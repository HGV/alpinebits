// Package freerooms holds the AlpineBits 2018-10 FreeRooms messages
// (OTA_HotelAvailNotif, §4.1).
package freerooms

import (
	"encoding/xml"

	"github.com/HGV/alpinebits"
	"github.com/HGV/x/timex"
)

type HotelAvailNotifRQ struct {
	XMLName             xml.Name            `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelAvailNotifRQ"`
	Version             string              `xml:"Version,attr"`
	UniqueID            *UniqueID           `xml:"UniqueID,omitempty"`
	AvailStatusMessages AvailStatusMessages `xml:"AvailStatusMessages"`
}

// HotelCode implements alpinebits.HotelCoded.
func (r HotelAvailNotifRQ) HotelCode() string {
	return r.AvailStatusMessages.HotelCode
}

// IsCompleteSet reports whether the message holds all availability of the
// hotel; the server replaces what it has on record. Otherwise it is a delta
// (AlpineBits 2018-10 §4.1.1).
func (r HotelAvailNotifRQ) IsCompleteSet() bool {
	return r.UniqueID != nil
}

// IsReset reports whether r resets all availability of the hotel: a
// CompleteSet with a single empty AvailStatusMessage.
func (r HotelAvailNotifRQ) IsReset() bool {
	msgs := r.AvailStatusMessages.AvailStatusMessages
	return r.IsCompleteSet() && len(msgs) == 1 && msgs[0] == AvailStatusMessage{}
}

type UniqueIDType string

const (
	UniqueIDTypeReference UniqueIDType = "16"
	// UniqueIDTypePurgedMasterReference marks a CompleteSet sent because the
	// client considers the data stale. The server may keep the data it has on
	// record, marked as purged, or treat it like UniqueIDTypeReference.
	UniqueIDTypePurgedMasterReference UniqueIDType = "35"
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

// AvailStatusMessages identifies the hotel by HotelCode; HotelName is optional.
type AvailStatusMessages struct {
	HotelCode           string               `xml:"HotelCode,attr,omitempty"`
	HotelName           string               `xml:"HotelName,attr,omitempty"`
	AvailStatusMessages []AvailStatusMessage `xml:"AvailStatusMessage"`
}

type BookingLimitMessageType string

const (
	BookingLimitMessageTypeSetLimit BookingLimitMessageType = "SetLimit"
)

// AvailStatusMessage gives the number of free rooms of a room or category for
// a period. It is empty only in a reset.
type AvailStatusMessage struct {
	BookingLimit            *int                    `xml:"BookingLimit,attr,omitempty"`
	BookingLimitMessageType BookingLimitMessageType `xml:"BookingLimitMessageType,attr,omitempty"`
	// BookingThreshold is the number of free rooms that are not bookable. It
	// is sent if and only if CapAcceptBookingThreshold was negotiated.
	BookingThreshold         *int                      `xml:"BookingThreshold,attr,omitempty"`
	StatusApplicationControl *StatusApplicationControl `xml:"StatusApplicationControl,omitempty"`
}

// DateRange returns the period from Start to End.
func (m AvailStatusMessage) DateRange() timex.DateRange {
	if s := m.StatusApplicationControl; s != nil {
		return timex.DateRange{Start: s.Start, End: s.End}
	}
	return timex.DateRange{}
}

// StatusApplicationControl refers to a specific room if InvCode is set and to
// a category otherwise. End is the last night.
type StatusApplicationControl struct {
	Start       timex.Date `xml:"Start,attr"`
	End         timex.Date `xml:"End,attr"`
	InvTypeCode string     `xml:"InvTypeCode,attr,omitempty"`
	InvCode     string     `xml:"InvCode,attr,omitempty"`
}

func (s StatusApplicationControl) isRoom() bool {
	return s.InvCode != ""
}

type HotelAvailNotifRS struct {
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelAvailNotifRS"`
	alpinebits.Response
}
