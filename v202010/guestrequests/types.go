// Package guestrequests holds the AlpineBits 2020-10 GuestRequests messages
// (OTA_Read and OTA_NotifReport, §4.2). ResRetrieveRS and the types leading to
// the new SpecialRequest are its own; the others are aliases of 2018-10.
package guestrequests

import (
	"encoding/xml"

	"github.com/HGV/alpinebits"
	guestrequests201810 "github.com/HGV/alpinebits/v201810/guestrequests"
)

type (
	DateTime                = guestrequests201810.DateTime
	ReadRQ                  = guestrequests201810.ReadRQ
	ReadRequest             = guestrequests201810.ReadRequest
	SelectionCriteria       = guestrequests201810.SelectionCriteria
	ResStatus               = guestrequests201810.ResStatus
	RoomStay                = guestrequests201810.RoomStay
	RoomType                = guestrequests201810.RoomType
	RatePlan                = guestrequests201810.RatePlan
	Commission              = guestrequests201810.Commission
	CommissionPayableAmount = guestrequests201810.CommissionPayableAmount
	MealsIncluded           = guestrequests201810.MealsIncluded
	GuestCount              = guestrequests201810.GuestCount
	TimeSpan                = guestrequests201810.TimeSpan
	StartDateWindow         = guestrequests201810.StartDateWindow
	Guarantee               = guestrequests201810.Guarantee
	PaymentCard             = guestrequests201810.PaymentCard
	CardNumber              = guestrequests201810.CardNumber
	Total                   = guestrequests201810.Total
	ResGuest                = guestrequests201810.ResGuest
	Gender                  = guestrequests201810.Gender
	Customer                = guestrequests201810.Customer
	PersonName              = guestrequests201810.PersonName
	PhoneTechType           = guestrequests201810.PhoneTechType
	Telephone               = guestrequests201810.Telephone
	Email                   = guestrequests201810.Email
	Address                 = guestrequests201810.Address
	CountryName             = guestrequests201810.CountryName
	CommentName             = guestrequests201810.CommentName
	Comment                 = guestrequests201810.Comment
	ListItem                = guestrequests201810.ListItem
	CancelPenalty           = guestrequests201810.CancelPenalty
	PenaltyDescription      = guestrequests201810.PenaltyDescription
	HotelReservationID      = guestrequests201810.HotelReservationID
	GlobalProfile           = guestrequests201810.GlobalProfile
	CompanyInfo             = guestrequests201810.CompanyInfo
	CompanyName             = guestrequests201810.CompanyName
	AddressInfo             = guestrequests201810.AddressInfo
	TelephoneInfo           = guestrequests201810.TelephoneInfo
	BasicPropertyInfo       = guestrequests201810.BasicPropertyInfo
	NotifReportRQ           = guestrequests201810.NotifReportRQ
	Warning                 = guestrequests201810.Warning
	UniqueIDType            = guestrequests201810.UniqueIDType
	NotifReservation        = guestrequests201810.NotifReservation
	UniqueID                = guestrequests201810.UniqueID
	NotifReportRS           = guestrequests201810.NotifReportRS
)

const (
	ResStatusRequested = guestrequests201810.ResStatusRequested
	ResStatusReserved  = guestrequests201810.ResStatusReserved
	ResStatusCancelled = guestrequests201810.ResStatusCancelled
	ResStatusModify    = guestrequests201810.ResStatusModify

	GenderUnknown = guestrequests201810.GenderUnknown
	GenderMale    = guestrequests201810.GenderMale
	GenderFemale  = guestrequests201810.GenderFemale

	PhoneTechTypeVoice  = guestrequests201810.PhoneTechTypeVoice
	PhoneTechTypeFax    = guestrequests201810.PhoneTechTypeFax
	PhoneTechTypeMobile = guestrequests201810.PhoneTechTypeMobile

	CommentNameIncludedServices = guestrequests201810.CommentNameIncludedServices
	CommentNameCustomerComment  = guestrequests201810.CommentNameCustomerComment

	UniqueIDTypeReservation  = guestrequests201810.UniqueIDTypeReservation
	UniqueIDTypeCancellation = guestrequests201810.UniqueIDTypeCancellation
)

// ResRetrieveRS answers OTA_Read:GuestRequests. A successful response always
// has a ReservationsList, empty if there are no HotelReservations.
type ResRetrieveRS struct {
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_ResRetrieveRS"`
	alpinebits.Response
	HotelReservations []HotelReservation `xml:"ReservationsList>HotelReservation"`
}

// MarshalXML writes ReservationsList next to Success, where the XSD requires
// it, and never next to Errors, where the XSD forbids it.
func (rs ResRetrieveRS) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	type list struct {
		HotelReservations []HotelReservation `xml:"HotelReservation"`
	}
	out := struct {
		XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_ResRetrieveRS"`
		alpinebits.Response
		ReservationsList *list `xml:"ReservationsList"`
	}{Response: rs.Response}
	if rs.Success != nil {
		out.ReservationsList = &list{rs.HotelReservations}
	}
	return e.Encode(out)
}

// HotelReservation is a quote request, a booking reservation, a modification
// or a cancellation, as ResStatus says.
type HotelReservation struct {
	UniqueID       UniqueID       `xml:"UniqueID"`
	RoomStays      *[]RoomStay    `xml:"RoomStays>RoomStay"`
	ResGuest       *ResGuest      `xml:"ResGuests>ResGuest"`
	ResGlobalInfo  *ResGlobalInfo `xml:"ResGlobalInfo"`
	CreateDateTime DateTime       `xml:"CreateDateTime,attr"`
	ResStatus      ResStatus      `xml:"ResStatus,attr"`
}

type ResGlobalInfo struct {
	Comments            *[]Comment            `xml:"Comments>Comment"`
	SpecialRequest      *SpecialRequest       `xml:"SpecialRequests>SpecialRequest"`
	CancelPenalty       *CancelPenalty        `xml:"CancelPenalties>CancelPenalty"`
	HotelReservationIDs *[]HotelReservationID `xml:"HotelReservationIDs>HotelReservationID"`
	GlobalProfile       *GlobalProfile        `xml:"Profiles>ProfileInfo>Profile"`
	BasicPropertyInfo   BasicPropertyInfo     `xml:"BasicPropertyInfo"`
}

// SpecialRequestName represents the subject of a special request.
type SpecialRequestName string

const (
	SpecialRequestNamePets SpecialRequestName = "PETS"
)

// SpecialRequest tells whether the guests bring pets: a Text, even an empty
// one, means they do and describes them; no Text means they do not. Without
// SpecialRequest it is unknown (AlpineBits 2020-10 §4.2.2).
type SpecialRequest struct {
	Name SpecialRequestName `xml:"Name,attr"`
	// Text has TextFormat PlainText and no Language.
	Text *alpinebits.Text `xml:"Text"`
}
