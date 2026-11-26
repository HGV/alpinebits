// Package guestrequests holds the AlpineBits 2026-10 GuestRequests messages
// (OTA_Read and OTA_NotifReport, §4.2). ResRetrieveRS and the types leading to
// the several ResGuests, the new PaymentCard, the global Profiles, the single
// Base and the Services TimeSpan are its own; the others are aliases of
// 2024-10.
package guestrequests

import (
	"encoding/xml"
	"strings"

	"github.com/HGV/alpinebits"
	guestrequests202410 "github.com/HGV/alpinebits/v202410/guestrequests"
	"github.com/HGV/x/timex"
)

type (
	DateTime                = guestrequests202410.DateTime
	ReadRQ                  = guestrequests202410.ReadRQ
	ReadRequest             = guestrequests202410.ReadRequest
	SelectionCriteria       = guestrequests202410.SelectionCriteria
	ResStatus               = guestrequests202410.ResStatus
	RoomType                = guestrequests202410.RoomType
	RatePlan                = guestrequests202410.RatePlan
	Commission              = guestrequests202410.Commission
	CommissionPayableAmount = guestrequests202410.CommissionPayableAmount
	MealsIncluded           = guestrequests202410.MealsIncluded
	GuestCount              = guestrequests202410.GuestCount
	TimeSpan                = guestrequests202410.TimeSpan
	StartDateWindow         = guestrequests202410.StartDateWindow
	Total                   = guestrequests202410.Total
	Gender                  = guestrequests202410.Gender
	Customer                = guestrequests202410.Customer
	PersonName              = guestrequests202410.PersonName
	PhoneTechType           = guestrequests202410.PhoneTechType
	Telephone               = guestrequests202410.Telephone
	Email                   = guestrequests202410.Email
	Address                 = guestrequests202410.Address
	CountryName             = guestrequests202410.CountryName
	CommentName             = guestrequests202410.CommentName
	Comment                 = guestrequests202410.Comment
	ListItem                = guestrequests202410.ListItem
	CancelPenalty           = guestrequests202410.CancelPenalty
	PenaltyDescription      = guestrequests202410.PenaltyDescription
	HotelReservationID      = guestrequests202410.HotelReservationID
	CompanyName             = guestrequests202410.CompanyName
	TelephoneInfo           = guestrequests202410.TelephoneInfo
	BasicPropertyInfo       = guestrequests202410.BasicPropertyInfo
	Warning                 = guestrequests202410.Warning
	UniqueIDType            = guestrequests202410.UniqueIDType
	UniqueID                = guestrequests202410.UniqueID
	NotifReportRS           = guestrequests202410.NotifReportRS
	CodeContext             = guestrequests202410.CodeContext
	SpecialRequest          = guestrequests202410.SpecialRequest
	ServiceRPH              = guestrequests202410.ServiceRPH
	ServiceCategoryCode     = guestrequests202410.ServiceCategoryCode
	ServicePricingType      = guestrequests202410.ServicePricingType
	BankAcct                = guestrequests202410.BankAcct
	BankAcctNumber          = guestrequests202410.BankAcctNumber
	Voucher                 = guestrequests202410.Voucher
	AmountPercent           = guestrequests202410.AmountPercent
)

const (
	ResStatusRequested = guestrequests202410.ResStatusRequested
	ResStatusReserved  = guestrequests202410.ResStatusReserved
	ResStatusCancelled = guestrequests202410.ResStatusCancelled
	ResStatusModify    = guestrequests202410.ResStatusModify

	GenderUnknown = guestrequests202410.GenderUnknown
	GenderMale    = guestrequests202410.GenderMale
	GenderFemale  = guestrequests202410.GenderFemale

	PhoneTechTypeVoice  = guestrequests202410.PhoneTechTypeVoice
	PhoneTechTypeFax    = guestrequests202410.PhoneTechTypeFax
	PhoneTechTypeMobile = guestrequests202410.PhoneTechTypeMobile

	CommentNameIncludedServices = guestrequests202410.CommentNameIncludedServices
	CommentNameCustomerComment  = guestrequests202410.CommentNameCustomerComment
	CommentNameAdditionalInfo   = guestrequests202410.CommentNameAdditionalInfo

	UniqueIDTypeReservation  = guestrequests202410.UniqueIDTypeReservation
	UniqueIDTypeCancellation = guestrequests202410.UniqueIDTypeCancellation

	CodeContextAlpineBits = guestrequests202410.CodeContextAlpineBits
	CodeContextHotel      = guestrequests202410.CodeContextHotel
	RequestCodePets       = guestrequests202410.RequestCodePets

	ServiceCategoryCodeBoard      = guestrequests202410.ServiceCategoryCodeBoard
	ServiceCategoryCodeSupplement = guestrequests202410.ServiceCategoryCodeSupplement
	ServiceCategoryCodeSpa        = guestrequests202410.ServiceCategoryCodeSpa
	ServiceCategoryCodeFood       = guestrequests202410.ServiceCategoryCodeFood
	ServiceCategoryCodeBeverage   = guestrequests202410.ServiceCategoryCodeBeverage
	ServiceCategoryCodeActivity   = guestrequests202410.ServiceCategoryCodeActivity
	ServiceCategoryCodeTouristTax = guestrequests202410.ServiceCategoryCodeTouristTax
	ServiceCategoryCodeTransfer   = guestrequests202410.ServiceCategoryCodeTransfer
	ServiceCategoryCodeOther      = guestrequests202410.ServiceCategoryCodeOther

	ServicePricingTypePerNight          = guestrequests202410.ServicePricingTypePerNight
	ServicePricingTypePerPerson         = guestrequests202410.ServicePricingTypePerPerson
	ServicePricingTypePerPersonPerNight = guestrequests202410.ServicePricingTypePerPersonPerNight
	ServicePricingTypePerStay           = guestrequests202410.ServicePricingTypePerStay
	ServicePricingTypePerUse            = guestrequests202410.ServicePricingTypePerUse

	ServiceCategoryCodePackage  ServiceCategoryCode = "PACKAGE"
	ServiceCategoryCodeResource ServiceCategoryCode = "RESOURCE"
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
	Services            *[]Service     `xml:"Services>Service"`
	ResGuests           *[]ResGuest    `xml:"ResGuests>ResGuest"`
	ResGlobalInfo       *ResGlobalInfo `xml:"ResGlobalInfo"`
	CreateDateTime      DateTime       `xml:"CreateDateTime,attr"`
	LastModifyDateTime  *DateTime      `xml:"LastModifyDateTime,attr,omitempty"`
	ResStatus           ResStatus      `xml:"ResStatus,attr"`
	RoomStayReservation *bool          `xml:"RoomStayReservation,attr,omitempty"`
}

// ResGuest is a guest; RoomStays refer to it by ResGuestRPH.
type ResGuest struct {
	ResGuestRPH string   `xml:"ResGuestRPH,attr,omitempty"`
	Customer    Customer `xml:"Profiles>ProfileInfo>Profile>Customer"`
}

// RoomStay is a desired stay. In a quote request, the RoomStays with the same
// RoomStayGroupID form an alternative.
type RoomStay struct {
	RoomStayGroupID string        `xml:"RoomStayGroupID,attr,omitempty"`
	RoomType        *RoomType     `xml:"RoomTypes>RoomType"`
	RatePlan        *RatePlan     `xml:"RatePlans>RatePlan"`
	RoomRate        *RoomRate     `xml:"RoomRates>RoomRate"`
	GuestCounts     *[]GuestCount `xml:"GuestCounts>GuestCount"`
	TimeSpan        TimeSpan      `xml:"TimeSpan"`
	Guarantee       *Guarantee    `xml:"Guarantee"`
	Total           *Total        `xml:"Total"`
	// ResGuestRPHs and ServiceRPHs refer to the guests and Services of the
	// stay.
	ResGuestRPHs RPHs          `xml:"ResGuestRPHs,omitempty"`
	ServiceRPHs  *[]ServiceRPH `xml:"ServiceRPHs>ServiceRPH"`
}

// RPHs is a list of references, written separated by spaces.
type RPHs []string

func (r RPHs) MarshalText() ([]byte, error) {
	return []byte(strings.Join(r, " ")), nil
}

func (r *RPHs) UnmarshalText(b []byte) error {
	*r = strings.Fields(string(b))
	return nil
}

type Guarantee struct {
	PaymentCard PaymentCard `xml:"GuaranteesAccepted>GuaranteeAccepted>PaymentCard"`
}

// PaymentCard identifies a card by its number, in plain text or encrypted, or
// by a Mask or Token.
type PaymentCard struct {
	CardCode       string     `xml:"CardCode,attr,omitempty"`
	ExpireDate     string     `xml:"ExpireDate,attr,omitempty"`
	CardHolderName string     `xml:"CardHolderName,omitempty"`
	CardNumber     CardNumber `xml:"CardNumber"`
}

type CardNumber struct {
	PlainText        string `xml:"PlainText,omitempty"`
	EncryptionMethod string `xml:"EncryptionMethod,attr,omitempty"`
	EncryptedValue   string `xml:"EncryptedValue,attr,omitempty"`
	Mask             string `xml:"Mask,attr,omitempty"`
	Token            string `xml:"Token,attr,omitempty"`
	TokenProviderID  string `xml:"TokenProviderID,attr,omitempty"`
}

// RoomRate gives the daily rates of a stay in the room category RoomTypeCode
// with the rate plan RatePlanCode.
type RoomRate struct {
	RatePlanCode string `xml:"RatePlanCode,attr"`
	RoomTypeCode string `xml:"RoomTypeCode,attr"`
	Rates        []Rate `xml:"Rates>Rate"`
}

// Rate is the rate from EffectiveDate to ExpireDate.
type Rate struct {
	EffectiveDate          timex.Date  `xml:"EffectiveDate,attr"`
	ExpireDate             *timex.Date `xml:"ExpireDate,attr,omitempty"`
	ExpireDateExclusiveInd *bool       `xml:"ExpireDateExclusiveInd,attr,omitempty"`
	RateTimeUnit           string      `xml:"RateTimeUnit,attr,omitempty"`
	UnitMultiplier         int         `xml:"UnitMultiplier,attr,omitempty"`
	Base                   Base        `xml:"Base"`
}

type Base = guestrequests202410.Base

// Service is a service booked with the reservation; ServiceRPH links it to
// RoomStays.
type Service struct {
	ID                   string              `xml:"ID,attr"`
	Type                 string              `xml:"Type,attr"`
	ServiceCategoryCode  ServiceCategoryCode `xml:"ServiceCategoryCode,attr"`
	ServiceInventoryCode string              `xml:"ServiceInventoryCode,attr"`
	ServiceRPH           string              `xml:"ServiceRPH,attr,omitempty"`
	ServicePricingType   ServicePricingType  `xml:"ServicePricingType,attr"`
	Inclusive            bool                `xml:"Inclusive,attr"`
	Quantity             int                 `xml:"Quantity,attr"`
	ServiceDetails       ServiceDetails      `xml:"ServiceDetails"`
}

type ServiceDetails struct {
	GuestCounts *[]GuestCount    `xml:"GuestCounts>GuestCount"`
	TimeSpan    *ServiceTimeSpan `xml:"TimeSpan"`
	// Comment and Description have TextFormat PlainText and no Language.
	Comment     *alpinebits.Text `xml:"Comments>Comment>Text"`
	Total       Total            `xml:"Total"`
	Description *alpinebits.Text `xml:"ServiceDescription>Text"`
}

// ServiceTimeSpan gives Start and End each as an xs:date, xs:dateTime or
// xs:time.
type ServiceTimeSpan struct {
	Start string `xml:"Start,attr"`
	End   string `xml:"End,attr,omitempty"`
}

type ResGlobalInfo struct {
	Comments            *[]Comment            `xml:"Comments>Comment"`
	SpecialRequests     *[]SpecialRequest     `xml:"SpecialRequests>SpecialRequest"`
	DepositPayments     *[]GuaranteePayment   `xml:"DepositPayments>GuaranteePayment"`
	CancelPenalty       *CancelPenalty        `xml:"CancelPenalties>CancelPenalty"`
	HotelReservationIDs *[]HotelReservationID `xml:"HotelReservationIDs>HotelReservationID"`
	Profiles            *[]ProfileInfo        `xml:"Profiles>ProfileInfo"`
	BasicPropertyInfo   BasicPropertyInfo     `xml:"BasicPropertyInfo"`
}

// GuaranteePayment is a deposit paid or to be paid for the reservation.
type GuaranteePayment struct {
	Start           string          `xml:"Start,attr,omitempty"`
	GuaranteeCode   string          `xml:"GuaranteeCode,attr,omitempty"`
	AcceptedPayment AcceptedPayment `xml:"AcceptedPayments>AcceptedPayment"`
	AmountPercent   *AmountPercent  `xml:"AmountPercent"`
}

// AcceptedPayment holds one of PaymentCard, BankAcct and Voucher.
type AcceptedPayment struct {
	GuaranteeTypeCode          string       `xml:"GuaranteeTypeCode,attr,omitempty"`
	GuaranteeID                string       `xml:"GuaranteeID,attr,omitempty"`
	PaymentTransactionTypeCode string       `xml:"PaymentTransactionTypeCode,attr,omitempty"`
	PaymentCard                *PaymentCard `xml:"PaymentCard"`
	BankAcct                   *BankAcct    `xml:"BankAcct"`
	Voucher                    *Voucher     `xml:"Voucher"`
}

type ProfileInfo struct {
	Profile Profile `xml:"Profile"`
}

// ProfileType represents the role of a global Profile.
type ProfileType string

const (
	ProfileTypeCustomer       ProfileType = "1"
	ProfileTypeCompany        ProfileType = "3"
	ProfileTypeBookingChannel ProfileType = "4"
)

// Profile describes the customer, the company or the booking channel, as
// ProfileType says.
type Profile struct {
	ProfileType ProfileType  `xml:"ProfileType,attr"`
	Customer    *Customer    `xml:"Customer"`
	CompanyInfo *CompanyInfo `xml:"CompanyInfo"`
}

type CompanyInfo struct {
	CompanyName   CompanyName    `xml:"CompanyName"`
	AddressInfo   *AddressInfo   `xml:"AddressInfo"`
	TelephoneInfo *TelephoneInfo `xml:"TelephoneInfo"`
	Email         string         `xml:"Email,omitempty"`
}

type AddressInfo struct {
	Language    string      `xml:"Language,attr,omitempty"`
	AddressLine string      `xml:"AddressLine"`
	CityName    string      `xml:"CityName"`
	PostalCode  string      `xml:"PostalCode"`
	CountryName CountryName `xml:"CountryName"`
}

// NotifReportRQ acknowledges guest requests by HotelReservations and refuses
// them by Warnings.
type NotifReportRQ struct {
	XMLName           xml.Name            `xml:"http://www.opentravel.org/OTA/2003/05 OTA_NotifReportRQ"`
	Version           string              `xml:"Version,attr"`
	Success           alpinebits.Success  `xml:"Success"`
	Warnings          *[]Warning          `xml:"Warnings>Warning"`
	HotelReservations *[]NotifReservation `xml:"NotifDetails>HotelNotifReport>HotelReservations>HotelReservation"`
}

// NotifReservation acknowledges the guest request with UniqueID, optionally
// naming its hotel.
type NotifReservation struct {
	UniqueID          UniqueID           `xml:"UniqueID"`
	BasicPropertyInfo *BasicPropertyInfo `xml:"ResGlobalInfo>BasicPropertyInfo"`
}
