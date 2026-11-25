// Package guestrequests holds the AlpineBits 2018-10 GuestRequests messages
// (OTA_Read and OTA_NotifReport, §4.2).
package guestrequests

import (
	"encoding/xml"
	"time"

	"github.com/HGV/alpinebits"
	"github.com/HGV/x/timex"
)

// DateTime is an xs:dateTime, whose time zone offset is optional.
type DateTime struct {
	Time time.Time
	// NoOffset marks a local time of unknown zone: the wall clock of Time is
	// written without offset. Unmarshalling such a value gives it in UTC.
	NoOffset bool
}

const localDateTime = "2006-01-02T15:04:05.999999999"

func (d DateTime) MarshalText() ([]byte, error) {
	if d.NoOffset {
		return []byte(d.Time.Format(localDateTime)), nil
	}
	return d.Time.MarshalText()
}

func (d *DateTime) UnmarshalText(b []byte) error {
	t, err := time.Parse(time.RFC3339, string(b))
	if err == nil {
		*d = DateTime{Time: t}
		return nil
	}
	if t, err := time.Parse(localDateTime, string(b)); err == nil {
		*d = DateTime{Time: t, NoOffset: true}
		return nil
	}
	return err
}

type ReadRQ struct {
	XMLName     xml.Name    `xml:"http://www.opentravel.org/OTA/2003/05 OTA_ReadRQ"`
	Version     string      `xml:"Version,attr"`
	ReadRequest ReadRequest `xml:"ReadRequests>HotelReadRequest"`
}

// HotelCode implements alpinebits.HotelCoded.
func (r ReadRQ) HotelCode() string {
	return r.ReadRequest.HotelCode
}

// ReadRequest identifies the hotel by HotelCode; HotelName is optional.
type ReadRequest struct {
	HotelCode         string             `xml:"HotelCode,attr,omitempty"`
	HotelName         string             `xml:"HotelName,attr,omitempty"`
	SelectionCriteria *SelectionCriteria `xml:"SelectionCriteria"`
}

// SelectionCriteria asks for the guest requests collected after Start,
// whether acknowledged or not. Without it the server sends those not yet
// acknowledged (AlpineBits 2018-10 §4.2.1).
type SelectionCriteria struct {
	Start DateTime `xml:"Start,attr"`
}

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

// ResStatus represents the status of a reservation.
type ResStatus string

const (
	ResStatusRequested ResStatus = "Requested"
	ResStatusReserved  ResStatus = "Reserved"
	ResStatusCancelled ResStatus = "Cancelled"
	ResStatusModify    ResStatus = "Modify"
)

// isReservation reports whether s is a booking reservation or modification,
// which AlpineBits 2018-10 §4.2.2 treats alike.
func (s ResStatus) isReservation() bool {
	return s == ResStatusReserved || s == ResStatusModify
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

// RoomStay is a desired stay. The last RoomStay of a quote request may hold
// only a TimeSpan: an alternative period for the others (AlpineBits 2018-10
// §4.2.2).
type RoomStay struct {
	RoomType    *RoomType     `xml:"RoomTypes>RoomType"`
	RatePlan    *RatePlan     `xml:"RatePlans>RatePlan"`
	GuestCounts *[]GuestCount `xml:"GuestCounts>GuestCount"`
	TimeSpan    TimeSpan      `xml:"TimeSpan"`
	Guarantee   *Guarantee    `xml:"Guarantee"`
	Total       *Total        `xml:"Total"`
}

func (s RoomStay) onlyTimeSpan() bool {
	return s.RoomType == nil && s.RatePlan == nil && s.GuestCounts == nil && s.Guarantee == nil && s.Total == nil
}

// RoomType is completely open in a quote request without attributes.
type RoomType struct {
	RoomTypeCode           string `xml:"RoomTypeCode,attr,omitempty"`
	RoomClassificationCode *int   `xml:"RoomClassificationCode,attr,omitempty"`
	RoomType               int    `xml:"RoomType,attr,omitempty"`
}

type RatePlan struct {
	RatePlanCode  string         `xml:"RatePlanCode,attr,omitempty"`
	Commission    *Commission    `xml:"Commission"`
	MealsIncluded *MealsIncluded `xml:"MealsIncluded"`
}

// Commission is given either by Percent or by CommissionPayableAmount.
type Commission struct {
	Percent                 *int                     `xml:"Percent,attr,omitempty"`
	CommissionPayableAmount *CommissionPayableAmount `xml:"CommissionPayableAmount"`
}

type CommissionPayableAmount struct {
	Amount       *alpinebits.Decimal `xml:"Amount,attr,omitempty"`
	CurrencyCode string              `xml:"CurrencyCode,attr,omitempty"`
}

// MealsIncluded requires MealPlanIndicator true.
type MealsIncluded struct {
	MealPlanIndicator bool   `xml:"MealPlanIndicator,attr"`
	MealPlanCodes     string `xml:"MealPlanCodes,attr"`
}

// GuestCount counts adults if Age is nil, children of that age otherwise.
type GuestCount struct {
	Count int  `xml:"Count,attr"`
	Age   *int `xml:"Age,attr,omitempty"`
}

// TimeSpan gives Start and End, or, in a quote request only, a window of
// possible arrival dates and a Duration.
type TimeSpan struct {
	Start           *timex.Date        `xml:"Start,attr,omitempty"`
	End             *timex.Date        `xml:"End,attr,omitempty"`
	Duration        *alpinebits.Nights `xml:"Duration,attr,omitempty"`
	StartDateWindow *StartDateWindow   `xml:"StartDateWindow"`
}

type StartDateWindow struct {
	EarliestDate timex.Date `xml:"EarliestDate,attr"`
	LatestDate   timex.Date `xml:"LatestDate,attr"`
}

type Guarantee struct {
	PaymentCard PaymentCard `xml:"GuaranteesAccepted>GuaranteeAccepted>PaymentCard"`
}

// PaymentCard requires CardHolderName in the XSD, although AlpineBits
// 2018-10 §4.2.2 calls it optional.
type PaymentCard struct {
	CardCode       string     `xml:"CardCode,attr"`
	ExpireDate     string     `xml:"ExpireDate,attr"`
	CardHolderName string     `xml:"CardHolderName"`
	CardNumber     CardNumber `xml:"CardNumber"`
}

// CardNumber holds either PlainText or EncryptedValue and EncryptionMethod.
type CardNumber struct {
	PlainText        string `xml:"PlainText,omitempty"`
	EncryptedValue   string `xml:"EncryptedValue,attr,omitempty"`
	EncryptionMethod string `xml:"EncryptionMethod,attr,omitempty"`
}

type Total struct {
	AmountAfterTax alpinebits.Decimal `xml:"AmountAfterTax,attr"`
	CurrencyCode   string             `xml:"CurrencyCode,attr"`
}

// ResGuest is the primary guest.
type ResGuest struct {
	Customer Customer `xml:"Profiles>ProfileInfo>Profile>Customer"`
}

// Gender represents the gender of a customer.
type Gender string

const (
	GenderUnknown Gender = "Unknown"
	GenderMale    Gender = "Male"
	GenderFemale  Gender = "Female"
)

type Customer struct {
	Gender     Gender      `xml:"Gender,attr,omitempty"`
	BirthDate  *timex.Date `xml:"BirthDate,attr,omitempty"`
	Language   string      `xml:"Language,attr,omitempty"`
	PersonName PersonName  `xml:"PersonName"`
	Telephones []Telephone `xml:"Telephone"`
	Email      *Email      `xml:"Email"`
	Address    *Address    `xml:"Address"`
}

type PersonName struct {
	NamePrefix string `xml:"NamePrefix,omitempty"`
	GivenName  string `xml:"GivenName"`
	Surname    string `xml:"Surname"`
	NameTitle  string `xml:"NameTitle,omitempty"`
}

// PhoneTechType represents the type of phone number.
type PhoneTechType string

const (
	PhoneTechTypeVoice  PhoneTechType = "1"
	PhoneTechTypeFax    PhoneTechType = "3"
	PhoneTechTypeMobile PhoneTechType = "5"
)

// Telephone requires PhoneTechType in the XSD, although AlpineBits 2018-10
// §4.2.2 calls it optional.
type Telephone struct {
	PhoneTechType PhoneTechType `xml:"PhoneTechType,attr"`
	PhoneNumber   string        `xml:"PhoneNumber,attr"`
}

// Email has Remark "newsletter:yes" or "newsletter:no"; none means unknown.
type Email struct {
	Value  string `xml:",chardata"`
	Remark string `xml:"Remark,attr,omitempty"`
}

// Address has Remark "catalog:yes" or "catalog:no"; none means unknown.
type Address struct {
	Remark      string       `xml:"Remark,attr,omitempty"`
	AddressLine string       `xml:"AddressLine,omitempty"`
	CityName    string       `xml:"CityName,omitempty"`
	PostalCode  string       `xml:"PostalCode,omitempty"`
	CountryName *CountryName `xml:"CountryName"`
}

type CountryName struct {
	Code string `xml:"Code,attr"`
}

type ResGlobalInfo struct {
	Comments            *[]Comment            `xml:"Comments>Comment"`
	CancelPenalty       *CancelPenalty        `xml:"CancelPenalties>CancelPenalty"`
	HotelReservationIDs *[]HotelReservationID `xml:"HotelReservationIDs>HotelReservationID"`
	GlobalProfile       *GlobalProfile        `xml:"Profiles>ProfileInfo>Profile"`
	BasicPropertyInfo   BasicPropertyInfo     `xml:"BasicPropertyInfo"`
}

// CommentName represents the type of comment.
type CommentName string

const (
	CommentNameIncludedServices CommentName = "included services"
	CommentNameCustomerComment  CommentName = "customer comment"
)

// Comment lists the included services as ListItems, or holds the customer
// comment as Text.
type Comment struct {
	Name      CommentName `xml:"Name,attr"`
	Text      string      `xml:"Text,omitempty"`
	ListItems []ListItem  `xml:"ListItem"`
}

type ListItem struct {
	Value    string `xml:",chardata"`
	ListItem int    `xml:"ListItem,attr"`
	Language string `xml:"Language,attr"`
}

type CancelPenalty struct {
	PenaltyDescription PenaltyDescription `xml:"PenaltyDescription"`
}

type PenaltyDescription struct {
	Text string `xml:"Text"`
}

type HotelReservationID struct {
	ResIDType          int    `xml:"ResID_Type,attr"`
	ResIDValue         string `xml:"ResID_Value,attr,omitempty"`
	ResIDSource        string `xml:"ResID_Source,attr,omitempty"`
	ResIDSourceContext string `xml:"ResID_SourceContext,attr,omitempty"`
}

// GlobalProfile describes the booking channel; ProfileType is "4".
type GlobalProfile struct {
	ProfileType string      `xml:"ProfileType,attr"`
	CompanyInfo CompanyInfo `xml:"CompanyInfo"`
}

type CompanyInfo struct {
	CompanyName   CompanyName    `xml:"CompanyName"`
	AddressInfo   *AddressInfo   `xml:"AddressInfo"`
	TelephoneInfo *TelephoneInfo `xml:"TelephoneInfo"`
	Email         string         `xml:"Email,omitempty"`
}

type CompanyName struct {
	Value       string `xml:",chardata"`
	Code        string `xml:"Code,attr"`
	CodeContext string `xml:"CodeContext,attr"`
}

type AddressInfo struct {
	AddressLine string      `xml:"AddressLine"`
	CityName    string      `xml:"CityName"`
	PostalCode  string      `xml:"PostalCode"`
	CountryName CountryName `xml:"CountryName"`
}

type TelephoneInfo struct {
	PhoneTechType PhoneTechType `xml:"PhoneTechType,attr"`
	PhoneNumber   string        `xml:"PhoneNumber,attr"`
}

// BasicPropertyInfo identifies the hotel by HotelCode; HotelName is optional.
type BasicPropertyInfo struct {
	HotelCode string `xml:"HotelCode,attr,omitempty"`
	HotelName string `xml:"HotelName,attr,omitempty"`
}

// NotifReportRQ acknowledges guest requests by HotelReservations and refuses
// them by Warnings (AlpineBits 2018-10 §4.2.3).
type NotifReportRQ struct {
	XMLName           xml.Name            `xml:"http://www.opentravel.org/OTA/2003/05 OTA_NotifReportRQ"`
	Version           string              `xml:"Version,attr"`
	Success           alpinebits.Success  `xml:"Success"`
	Warnings          *[]Warning          `xml:"Warnings>Warning"`
	HotelReservations *[]NotifReservation `xml:"NotifDetails>HotelNotifReport>HotelReservations>HotelReservation"`
}

// Warning refuses the guest request whose UniqueID has the ID RecordID. Type
// is from the OTA list EWT (3: biz rule), Code from ERR (450: unable to
// process).
type Warning struct {
	Type     alpinebits.ErrType `xml:"Type,attr"`
	Code     alpinebits.ErrCode `xml:"Code,attr"`
	RecordID string             `xml:"RecordID,attr"`
	Message  string             `xml:",chardata"`
}

type UniqueIDType string

const (
	UniqueIDTypeReservation  UniqueIDType = "14"
	UniqueIDTypeCancellation UniqueIDType = "15"
)

type NotifReservation struct {
	UniqueID UniqueID `xml:"UniqueID"`
}

type UniqueID struct {
	Type UniqueIDType `xml:"Type,attr"`
	ID   string       `xml:"ID,attr"`
}

// NotifReportRS answers OTA_NotifReport:GuestRequests. The 2018-10 XSD allows
// no Warnings here, although AlpineBits 2018-10 §4.2.4 allows advisory and
// warning outcomes: the router answers a response with Warnings with HTTP 500.
type NotifReportRS struct {
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_NotifReportRS"`
	alpinebits.Response
}
