// Package guestrequests holds the AlpineBits 2022-10 GuestRequests messages
// (OTA_Read and OTA_NotifReport, §4.2). ResRetrieveRS and the types leading to
// RoomStay and SpecialRequest are its own; the others are aliases of 2020-10.
package guestrequests

import (
	"encoding/xml"

	"github.com/HGV/alpinebits"
	guestrequests202010 "github.com/HGV/alpinebits/v202010/guestrequests"
)

type (
	DateTime                = guestrequests202010.DateTime
	ReadRQ                  = guestrequests202010.ReadRQ
	ReadRequest             = guestrequests202010.ReadRequest
	SelectionCriteria       = guestrequests202010.SelectionCriteria
	ResStatus               = guestrequests202010.ResStatus
	RoomType                = guestrequests202010.RoomType
	RatePlan                = guestrequests202010.RatePlan
	Commission              = guestrequests202010.Commission
	CommissionPayableAmount = guestrequests202010.CommissionPayableAmount
	MealsIncluded           = guestrequests202010.MealsIncluded
	GuestCount              = guestrequests202010.GuestCount
	TimeSpan                = guestrequests202010.TimeSpan
	StartDateWindow         = guestrequests202010.StartDateWindow
	Guarantee               = guestrequests202010.Guarantee
	PaymentCard             = guestrequests202010.PaymentCard
	CardNumber              = guestrequests202010.CardNumber
	Total                   = guestrequests202010.Total
	ResGuest                = guestrequests202010.ResGuest
	Gender                  = guestrequests202010.Gender
	Customer                = guestrequests202010.Customer
	PersonName              = guestrequests202010.PersonName
	PhoneTechType           = guestrequests202010.PhoneTechType
	Telephone               = guestrequests202010.Telephone
	Email                   = guestrequests202010.Email
	Address                 = guestrequests202010.Address
	CountryName             = guestrequests202010.CountryName
	CommentName             = guestrequests202010.CommentName
	Comment                 = guestrequests202010.Comment
	ListItem                = guestrequests202010.ListItem
	CancelPenalty           = guestrequests202010.CancelPenalty
	PenaltyDescription      = guestrequests202010.PenaltyDescription
	HotelReservationID      = guestrequests202010.HotelReservationID
	GlobalProfile           = guestrequests202010.GlobalProfile
	CompanyInfo             = guestrequests202010.CompanyInfo
	CompanyName             = guestrequests202010.CompanyName
	AddressInfo             = guestrequests202010.AddressInfo
	TelephoneInfo           = guestrequests202010.TelephoneInfo
	BasicPropertyInfo       = guestrequests202010.BasicPropertyInfo
	NotifReportRQ           = guestrequests202010.NotifReportRQ
	Warning                 = guestrequests202010.Warning
	UniqueIDType            = guestrequests202010.UniqueIDType
	NotifReservation        = guestrequests202010.NotifReservation
	UniqueID                = guestrequests202010.UniqueID
	NotifReportRS           = guestrequests202010.NotifReportRS
)

const (
	ResStatusRequested = guestrequests202010.ResStatusRequested
	ResStatusReserved  = guestrequests202010.ResStatusReserved
	ResStatusCancelled = guestrequests202010.ResStatusCancelled
	ResStatusModify    = guestrequests202010.ResStatusModify

	GenderUnknown = guestrequests202010.GenderUnknown
	GenderMale    = guestrequests202010.GenderMale
	GenderFemale  = guestrequests202010.GenderFemale

	PhoneTechTypeVoice  = guestrequests202010.PhoneTechTypeVoice
	PhoneTechTypeFax    = guestrequests202010.PhoneTechTypeFax
	PhoneTechTypeMobile = guestrequests202010.PhoneTechTypeMobile

	CommentNameIncludedServices = guestrequests202010.CommentNameIncludedServices
	CommentNameCustomerComment  = guestrequests202010.CommentNameCustomerComment

	UniqueIDTypeReservation  = guestrequests202010.UniqueIDTypeReservation
	UniqueIDTypeCancellation = guestrequests202010.UniqueIDTypeCancellation
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

// HotelReservation is a quote, newsletter or catalog request, a booking
// reservation, a modification or a cancellation, as ResStatus says. A
// newsletter or catalog request is Requested with RoomStayReservation false.
type HotelReservation struct {
	UniqueID            UniqueID       `xml:"UniqueID"`
	RoomStays           *[]RoomStay    `xml:"RoomStays>RoomStay"`
	ResGuest            *ResGuest      `xml:"ResGuests>ResGuest"`
	ResGlobalInfo       *ResGlobalInfo `xml:"ResGlobalInfo"`
	CreateDateTime      DateTime       `xml:"CreateDateTime,attr"`
	LastModifyDateTime  *DateTime      `xml:"LastModifyDateTime,attr,omitempty"`
	ResStatus           ResStatus      `xml:"ResStatus,attr"`
	RoomStayReservation *bool          `xml:"RoomStayReservation,attr,omitempty"`
}

// RoomStay is a desired stay. In a quote request, the RoomStays with the same
// RoomStayGroupID form an alternative (AlpineBits 2022-10 §4.2.2).
type RoomStay struct {
	RoomStayGroupID string        `xml:"RoomStayGroupID,attr,omitempty"`
	RoomType        *RoomType     `xml:"RoomTypes>RoomType"`
	RatePlan        *RatePlan     `xml:"RatePlans>RatePlan"`
	GuestCounts     *[]GuestCount `xml:"GuestCounts>GuestCount"`
	TimeSpan        TimeSpan      `xml:"TimeSpan"`
	Guarantee       *Guarantee    `xml:"Guarantee"`
	Total           *Total        `xml:"Total"`
}

type ResGlobalInfo struct {
	Comments            *[]Comment            `xml:"Comments>Comment"`
	SpecialRequests     *[]SpecialRequest     `xml:"SpecialRequests>SpecialRequest"`
	CancelPenalty       *CancelPenalty        `xml:"CancelPenalties>CancelPenalty"`
	HotelReservationIDs *[]HotelReservationID `xml:"HotelReservationIDs>HotelReservationID"`
	GlobalProfile       *GlobalProfile        `xml:"Profiles>ProfileInfo>Profile"`
	BasicPropertyInfo   BasicPropertyInfo     `xml:"BasicPropertyInfo"`
}

// CodeContext tells who defines the RequestCode of a SpecialRequest.
type CodeContext string

const (
	CodeContextAlpineBits CodeContext = "ALPINEBITS"
	CodeContextHotel      CodeContext = "HOTEL"
)

// RequestCodePets is the only RequestCode CodeContext ALPINEBITS defines.
const RequestCodePets = "PETS"

// SpecialRequest tells whether the customer wants something: a Text, even an
// empty one, means they do and gives details; no Text means they do not.
// Without SpecialRequest it is unknown (AlpineBits 2022-10 §4.2.2).
type SpecialRequest struct {
	RequestCode string      `xml:"RequestCode,attr"`
	CodeContext CodeContext `xml:"CodeContext,attr"`
	// Text has TextFormat PlainText and no Language.
	Text *alpinebits.Text `xml:"Text"`
}
