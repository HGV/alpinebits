// Package freerooms holds the AlpineBits 2020-10 FreeRooms messages
// (OTA_HotelInvCountNotif, §4.1).
package freerooms

import (
	"encoding/xml"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rules"
	"github.com/HGV/x/timex"
)

type HotelInvCountNotifRQ struct {
	XMLName     xml.Name    `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelInvCountNotifRQ"`
	Version     string      `xml:"Version,attr"`
	UniqueID    *UniqueID   `xml:"UniqueID,omitempty"`
	Inventories Inventories `xml:"Inventories"`
}

// HotelCode implements alpinebits.HotelCoded.
func (r HotelInvCountNotifRQ) HotelCode() string {
	return r.Inventories.HotelCode
}

// IsCompleteSet reports whether r replaces all availability of the hotel;
// otherwise it is a delta (AlpineBits 2020-10 §4.1.1).
func (r HotelInvCountNotifRQ) IsCompleteSet() bool {
	return r.UniqueID != nil
}

// IsReset reports whether r deletes all availability of the hotel: a
// CompleteSet with a single empty Inventory.
func (r HotelInvCountNotifRQ) IsReset() bool {
	inv := r.Inventories.Inventories
	return r.IsCompleteSet() && len(inv) == 1 && inv[0].isEmpty()
}

// UniqueIDType represents the type of a CompleteSet.
type UniqueIDType string

const (
	UniqueIDTypeReference UniqueIDType = "16"
	// UniqueIDTypePurgedMasterReference is a CompleteSet hinting that the
	// replaced data is stale; the server may keep it flagged as purged.
	UniqueIDTypePurgedMasterReference UniqueIDType = "35"
)

// Instance represents the UniqueID instance.
type Instance string

const (
	InstanceCompleteSet Instance = "CompleteSet"
)

type UniqueID struct {
	Type     UniqueIDType `xml:"Type,attr"`
	ID       string       `xml:"ID,attr"`
	Instance Instance     `xml:"Instance,attr"`
}

type Inventories struct {
	HotelCode   string      `xml:"HotelCode,attr,omitempty"`
	HotelName   string      `xml:"HotelName,attr,omitempty"`
	Inventories []Inventory `xml:"Inventory"`
}

// Inventory holds the counts of a room or room category for a period, or a
// closing season of the hotel. Missing InvCounts mean no free rooms.
type Inventory struct {
	StatusApplicationControl *StatusApplicationControl `xml:"StatusApplicationControl,omitempty"`
	InvCounts                *[]InvCount               `xml:"InvCounts>InvCount"`
}

func (i Inventory) isEmpty() bool {
	return i.StatusApplicationControl == nil && i.InvCounts == nil
}

// IsClosingSeason reports whether i marks a period in which the hotel is
// closed.
func (i Inventory) IsClosingSeason() bool {
	return i.StatusApplicationControl != nil && i.StatusApplicationControl.AllInvCode
}

// IsRoom reports whether i refers to a specific room (InvCode) rather than a
// room category.
func (i Inventory) IsRoom() bool {
	return i.StatusApplicationControl != nil && i.StatusApplicationControl.InvCode != ""
}

// Count returns the Count of the InvCount with CountType t, 0 if there is
// none (§4.1.1).
func (i Inventory) Count(t CountType) int {
	for _, c := range rules.Items(i.InvCounts) {
		if c.CountType == t {
			return c.Count
		}
	}
	return 0
}

func (i Inventory) total() int {
	n := 0
	for _, c := range rules.Items(i.InvCounts) {
		n += c.Count
	}
	return n
}

// DateRange returns the period from Start to End.
func (i Inventory) DateRange() timex.DateRange {
	if s := i.StatusApplicationControl; s != nil {
		return timex.DateRange{Start: s.Start, End: s.End}
	}
	return timex.DateRange{}
}

type StatusApplicationControl struct {
	Start       timex.Date `xml:"Start,attr"`
	End         timex.Date `xml:"End,attr"`
	InvTypeCode string     `xml:"InvTypeCode,attr,omitempty"`
	InvCode     string     `xml:"InvCode,attr,omitempty"`
	AllInvCode  bool       `xml:"AllInvCode,attr,omitempty"`
}

// CountType represents the kind of rooms an InvCount counts.
type CountType string

const (
	CountTypeBookable    CountType = "2"
	CountTypeOutOfOrder  CountType = "6"
	CountTypeOutOfMarket CountType = "9" // available but not bookable
)

type InvCount struct {
	CountType CountType `xml:"CountType,attr"`
	Count     int       `xml:"Count,attr"`
}

type HotelInvCountNotifRS struct {
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelInvCountNotifRS"`
	alpinebits.Response
}
