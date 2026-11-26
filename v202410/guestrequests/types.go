// Package guestrequests holds the AlpineBits 2024-10 GuestRequests messages
// (OTA_Read and OTA_NotifReport, §4.2). ResRetrieveRS and the types leading to
// the new RoomRates, Services and DepositPayments are its own; the others are
// aliases of 2022-10.
package guestrequests

import (
	"encoding/xml"

	"github.com/HGV/alpinebits"
	guestrequests202210 "github.com/HGV/alpinebits/v202210/guestrequests"
	"github.com/HGV/x/timex"
)

type (
	DateTime                = guestrequests202210.DateTime
	ReadRQ                  = guestrequests202210.ReadRQ
	ReadRequest             = guestrequests202210.ReadRequest
	SelectionCriteria       = guestrequests202210.SelectionCriteria
	ResStatus               = guestrequests202210.ResStatus
	RoomType                = guestrequests202210.RoomType
	RatePlan                = guestrequests202210.RatePlan
	Commission              = guestrequests202210.Commission
	CommissionPayableAmount = guestrequests202210.CommissionPayableAmount
	MealsIncluded           = guestrequests202210.MealsIncluded
	GuestCount              = guestrequests202210.GuestCount
	TimeSpan                = guestrequests202210.TimeSpan
	StartDateWindow         = guestrequests202210.StartDateWindow
	Guarantee               = guestrequests202210.Guarantee
	PaymentCard             = guestrequests202210.PaymentCard
	CardNumber              = guestrequests202210.CardNumber
	Total                   = guestrequests202210.Total
	ResGuest                = guestrequests202210.ResGuest
	Gender                  = guestrequests202210.Gender
	Customer                = guestrequests202210.Customer
	PersonName              = guestrequests202210.PersonName
	PhoneTechType           = guestrequests202210.PhoneTechType
	Telephone               = guestrequests202210.Telephone
	Email                   = guestrequests202210.Email
	Address                 = guestrequests202210.Address
	CountryName             = guestrequests202210.CountryName
	CommentName             = guestrequests202210.CommentName
	Comment                 = guestrequests202210.Comment
	ListItem                = guestrequests202210.ListItem
	CancelPenalty           = guestrequests202210.CancelPenalty
	PenaltyDescription      = guestrequests202210.PenaltyDescription
	HotelReservationID      = guestrequests202210.HotelReservationID
	GlobalProfile           = guestrequests202210.GlobalProfile
	CompanyInfo             = guestrequests202210.CompanyInfo
	CompanyName             = guestrequests202210.CompanyName
	AddressInfo             = guestrequests202210.AddressInfo
	TelephoneInfo           = guestrequests202210.TelephoneInfo
	BasicPropertyInfo       = guestrequests202210.BasicPropertyInfo
	NotifReportRQ           = guestrequests202210.NotifReportRQ
	Warning                 = guestrequests202210.Warning
	UniqueIDType            = guestrequests202210.UniqueIDType
	NotifReservation        = guestrequests202210.NotifReservation
	UniqueID                = guestrequests202210.UniqueID
	NotifReportRS           = guestrequests202210.NotifReportRS
	CodeContext             = guestrequests202210.CodeContext
	SpecialRequest          = guestrequests202210.SpecialRequest
)

const (
	ResStatusRequested = guestrequests202210.ResStatusRequested
	ResStatusReserved  = guestrequests202210.ResStatusReserved
	ResStatusCancelled = guestrequests202210.ResStatusCancelled
	ResStatusModify    = guestrequests202210.ResStatusModify

	GenderUnknown = guestrequests202210.GenderUnknown
	GenderMale    = guestrequests202210.GenderMale
	GenderFemale  = guestrequests202210.GenderFemale

	PhoneTechTypeVoice  = guestrequests202210.PhoneTechTypeVoice
	PhoneTechTypeFax    = guestrequests202210.PhoneTechTypeFax
	PhoneTechTypeMobile = guestrequests202210.PhoneTechTypeMobile

	CommentNameIncludedServices = guestrequests202210.CommentNameIncludedServices
	CommentNameCustomerComment  = guestrequests202210.CommentNameCustomerComment

	UniqueIDTypeReservation  = guestrequests202210.UniqueIDTypeReservation
	UniqueIDTypeCancellation = guestrequests202210.UniqueIDTypeCancellation

	CodeContextAlpineBits = guestrequests202210.CodeContextAlpineBits
	CodeContextHotel      = guestrequests202210.CodeContextHotel
	RequestCodePets       = guestrequests202210.RequestCodePets

	CommentNameAdditionalInfo CommentName = "additional info"
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
	ResGuest            *ResGuest      `xml:"ResGuests>ResGuest"`
	ResGlobalInfo       *ResGlobalInfo `xml:"ResGlobalInfo"`
	CreateDateTime      DateTime       `xml:"CreateDateTime,attr"`
	LastModifyDateTime  *DateTime      `xml:"LastModifyDateTime,attr,omitempty"`
	ResStatus           ResStatus      `xml:"ResStatus,attr"`
	RoomStayReservation *bool          `xml:"RoomStayReservation,attr,omitempty"`
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
	// ServiceRPHs refer to the Services booked with the stay.
	ServiceRPHs *[]ServiceRPH `xml:"ServiceRPHs>ServiceRPH"`
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
	Bases                  []Base      `xml:"Base"`
}

type Base struct {
	AmountAfterTax *alpinebits.Decimal `xml:"AmountAfterTax,attr,omitempty"`
	CurrencyCode   string              `xml:"CurrencyCode,attr,omitempty"`
}

type ServiceRPH struct {
	RPH string `xml:"RPH,attr"`
}

// ServiceCategoryCode represents the kind of a Service.
type ServiceCategoryCode string

const (
	ServiceCategoryCodeBoard      ServiceCategoryCode = "BOARD"
	ServiceCategoryCodeSupplement ServiceCategoryCode = "SUPPLEMENT"
	ServiceCategoryCodeSpa        ServiceCategoryCode = "SPA"
	ServiceCategoryCodeFood       ServiceCategoryCode = "FOOD"
	ServiceCategoryCodeBeverage   ServiceCategoryCode = "BEVERAGE"
	ServiceCategoryCodeActivity   ServiceCategoryCode = "ACTIVITY"
	ServiceCategoryCodeTouristTax ServiceCategoryCode = "TOURISTTAX"
	ServiceCategoryCodeTransfer   ServiceCategoryCode = "TRANSFER"
	ServiceCategoryCodeOther      ServiceCategoryCode = "OTHER"
)

// ServicePricingType represents how the price of a Service is computed.
type ServicePricingType string

const (
	ServicePricingTypePerNight          ServicePricingType = "Per night"
	ServicePricingTypePerPerson         ServicePricingType = "Per person"
	ServicePricingTypePerPersonPerNight ServicePricingType = "Per person per night"
	ServicePricingTypePerStay           ServicePricingType = "Per stay"
	ServicePricingTypePerUse            ServicePricingType = "Per use"
)

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

type ServiceTimeSpan struct {
	Start DateTime  `xml:"Start,attr"`
	End   *DateTime `xml:"End,attr,omitempty"`
}

type ResGlobalInfo struct {
	Comments            *[]Comment            `xml:"Comments>Comment"`
	SpecialRequests     *[]SpecialRequest     `xml:"SpecialRequests>SpecialRequest"`
	DepositPayments     *[]GuaranteePayment   `xml:"DepositPayments>GuaranteePayment"`
	CancelPenalty       *CancelPenalty        `xml:"CancelPenalties>CancelPenalty"`
	HotelReservationIDs *[]HotelReservationID `xml:"HotelReservationIDs>HotelReservationID"`
	GlobalProfile       *GlobalProfile        `xml:"Profiles>ProfileInfo>Profile"`
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
	GuaranteeTypeCode          string              `xml:"GuaranteeTypeCode,attr,omitempty"`
	GuaranteeID                string              `xml:"GuaranteeID,attr,omitempty"`
	PaymentTransactionTypeCode string              `xml:"PaymentTransactionTypeCode,attr,omitempty"`
	PaymentCard                *DepositPaymentCard `xml:"PaymentCard"`
	BankAcct                   *BankAcct           `xml:"BankAcct"`
	Voucher                    *Voucher            `xml:"Voucher"`
}

// DepositPaymentCard identifies a card by a mask or a token, never by its
// number.
type DepositPaymentCard struct {
	CardCode       string               `xml:"CardCode,attr,omitempty"`
	ExpireDate     string               `xml:"ExpireDate,attr,omitempty"`
	CardHolderName string               `xml:"CardHolderName,omitempty"`
	CardNumber     *TokenizedCardNumber `xml:"CardNumber"`
}

type TokenizedCardNumber struct {
	Mask            string `xml:"Mask,attr,omitempty"`
	Token           string `xml:"Token,attr,omitempty"`
	TokenProviderID string `xml:"TokenProviderID,attr,omitempty"`
}

type BankAcct struct {
	BankAcctName   string         `xml:"BankAcctName,omitempty"`
	BankAcctNumber BankAcctNumber `xml:"BankAcctNumber"`
}

type BankAcctNumber struct {
	Mask      string `xml:"Mask,attr,omitempty"`
	PlainText string `xml:"PlainText"`
}

// Voucher has ElectronicIndicator true.
type Voucher struct {
	ElectronicIndicator bool   `xml:"ElectronicIndicator,attr"`
	Identifier          string `xml:"Identifier,attr"`
}

type AmountPercent struct {
	Amount       alpinebits.Decimal `xml:"Amount,attr"`
	CurrencyCode string             `xml:"CurrencyCode,attr,omitempty"`
}
