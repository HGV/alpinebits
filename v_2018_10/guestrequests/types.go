package guestrequests

import (
	"encoding/xml"
	"time"

	"github.com/HGV/alpinebits"
	"github.com/HGV/x/timex"
)

type ReadRQ struct {
	XMLName     xml.Name    `xml:"http://www.opentravel.org/OTA/2003/05 OTA_ReadRQ"`
	Version     string      `xml:"Version,attr"`
	ReadRequest ReadRequest `xml:"ReadRequests>HotelReadRequest"`
}

// HotelCode implements alpinebits.HotelCoded.
func (r ReadRQ) HotelCode() string {
	return r.ReadRequest.HotelCode
}

type ReadRequest struct {
	HotelCode         string             `xml:"HotelCode,attr"`
	HotelName         string             `xml:"HotelName,attr,omitempty"`
	SelectionCriteria *SelectionCriteria `xml:"SelectionCriteria"`
}

type SelectionCriteria struct {
	Start time.Time `xml:"Start,attr"`
}

type ResRetrieveRS struct {
	XMLName           xml.Name              `xml:"http://www.opentravel.org/OTA/2003/05 OTA_ResRetrieveRS"`
	Version           string                `xml:"Version,attr"`
	Success           *alpinebits.Success   `xml:"Success"`
	Warnings          *alpinebits.Warnings  `xml:"Warnings>Warning"`
	Errors            *alpinebits.Errors    `xml:"Errors>Error"`
	HotelReservations *[]HotelReservation   `xml:"ReservationsList>HotelReservation"`
}

// ResStatus represents the status of a reservation.
type ResStatus string

const (
	ResStatusRequested ResStatus = "Requested"
	ResStatusReserved  ResStatus = "Reserved"
	ResStatusCancelled ResStatus = "Cancelled"
	ResStatusModify    ResStatus = "Modify"
)

type HotelReservation struct {
	UniqueID       UniqueID       `xml:"UniqueID"`
	RoomStays      *[]RoomStay    `xml:"RoomStays>RoomStay"`
	ResGuest       *ResGuest      `xml:"ResGuests>ResGuest"`
	ResGlobalInfo  *ResGlobalInfo `xml:"ResGlobalInfo"`
	CreateDateTime time.Time      `xml:"CreateDateTime,attr"`
	ResStatus      ResStatus      `xml:"ResStatus,attr"`
}

type RoomStay struct {
	RoomType    *RoomType     `xml:"RoomTypes>RoomType"`
	RatePlan    *RatePlan     `xml:"RatePlans>RatePlan"`
	GuestCounts *[]GuestCount `xml:"GuestCounts>GuestCount"`
	TimeSpan    TimeSpan      `xml:"TimeSpan"`
	Guarantee   *Guarantee    `xml:"Guarantee"`
	Total       *Total        `xml:"Total"`
}

type RoomType struct {
	RoomTypeCode           string `xml:"RoomTypeCode,attr,omitempty"`
	RoomClassificationCode int    `xml:"RoomClassificationCode,attr,omitempty"`
	RoomType               int    `xml:"RoomType,attr,omitempty"`
}

type RatePlan struct {
	RatePlanCode  string         `xml:"RatePlanCode,attr,omitempty"`
	Commission    *Commission    `xml:"Commission"`
	MealsIncluded *MealsIncluded `xml:"MealsIncluded"`
}

type Commission struct {
	Percent                 int                      `xml:"Percent,attr,omitempty"`
	CommissionPayableAmount *CommissionPayableAmount `xml:"CommissionPayableAmount"`
}

type CommissionPayableAmount struct {
	Amount       float64 `xml:"Amount,attr,omitempty"`
	CurrencyCode string  `xml:"CurrencyCode,attr,omitempty"`
}

type MealsIncluded struct {
	MealPlanIndicator bool   `xml:"MealPlanIndicator,attr"`
	MealPlanCodes     string `xml:"MealPlanCodes,attr"`
}

type GuestCount struct {
	Count int `xml:"Count,attr"`
	Age   int `xml:"Age,attr,omitempty"`
}

type TimeSpan struct {
	Start           timex.Date        `xml:"Start,attr,omitempty"`
	End             timex.Date        `xml:"End,attr,omitempty"`
	Duration        alpinebits.Nights `xml:"Duration,attr,omitempty"`
	StartDateWindow *StartDateWindow  `xml:"StartDateWindow"`
}

type StartDateWindow struct {
	EarliestDate timex.Date `xml:"EarliestDate,attr"`
	LatestDate   timex.Date `xml:"LatestDate,attr"`
}

type Guarantee struct {
	PaymentCard PaymentCard `xml:"GuaranteesAccepted>GuaranteeAccepted>PaymentCard"`
}

type PaymentCard struct {
	CardCode       string     `xml:"CardCode,attr"`
	ExpireDate     string     `xml:"ExpireDate,attr"`
	CardHolderName string     `xml:"CardHolderName"`
	CardNumber     CardNumber `xml:"CardNumber"`
}

type CardNumber struct {
	PlainText        string `xml:"PlainText,omitempty"`
	EncryptedValue   string `xml:"EncryptedValue,attr,omitempty"`
	EncryptionMethod string `xml:"EncryptionMethod,attr,omitempty"`
}

type Total struct {
	AmountAfterTax float64 `xml:"AmountAfterTax,attr"`
	CurrencyCode   string  `xml:"CurrencyCode,attr"`
}

type ResGuest struct {
	Profile Profile `xml:"Profiles>ProfileInfo>Profile"`
}

type Profile struct {
	ProfileType string       `xml:"ProfileType,attr,omitempty"`
	Customer    *Customer    `xml:"Customer"`
	CompanyInfo *CompanyInfo `xml:"CompanyInfo"`
}

// Gender represents the gender of a customer.
type Gender string

const (
	GenderUnknown Gender = "Unknown"
	GenderMale    Gender = "Male"
	GenderFemale  Gender = "Female"
)

type Customer struct {
	Gender     Gender     `xml:"Gender,attr,omitempty"`
	BirthDate  timex.Date `xml:"BirthDate,attr,omitempty"`
	Language   string     `xml:"Language,attr,omitempty"`
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

type Telephone struct {
	PhoneTechType PhoneTechType `xml:"PhoneTechType,attr"`
	PhoneNumber   string        `xml:"PhoneNumber,attr"`
}

type Email struct {
	Value  string `xml:",chardata"`
	Remark string `xml:"Remark,attr,omitempty"`
}

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
	Comments            *[]Comment           `xml:"Comments>Comment"`
	CancelPenalty       *CancelPenalty       `xml:"CancelPenalties>CancelPenalty"`
	HotelReservationIDs *[]HotelReservationID `xml:"HotelReservationIDs>HotelReservationID"`
	GlobalProfile       *GlobalProfile       `xml:"Profiles>ProfileInfo>Profile"`
	BasicPropertyInfo   BasicPropertyInfo    `xml:"BasicPropertyInfo"`
}

// CommentName represents the type of comment.
type CommentName string

const (
	CommentNameIncludedServices CommentName = "included services"
	CommentNameCustomerComment  CommentName = "customer comment"
)

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

type BasicPropertyInfo struct {
	HotelCode string `xml:"HotelCode,attr,omitempty"`
	HotelName string `xml:"HotelName,attr,omitempty"`
}

type NotifReportRQ struct {
	XMLName           xml.Name             `xml:"http://www.opentravel.org/OTA/2003/05 OTA_NotifReportRQ"`
	Version           string               `xml:"Version,attr"`
	Success           *alpinebits.Success  `xml:"Success"`
	Warnings          *alpinebits.Warnings `xml:"Warnings>Warning"`
	HotelReservations []NotifReservation   `xml:"NotifDetails>HotelNotifReport>HotelReservations>HotelReservation"`
}

type UniqueIDType int

const (
	UniqueIDTypeReservation  UniqueIDType = 14
	UniqueIDTypeCancellation UniqueIDType = 15
)

type NotifReservation struct {
	UniqueID UniqueID `xml:"UniqueID"`
}

type UniqueID struct {
	Type UniqueIDType `xml:"Type,attr"`
	ID   string       `xml:"ID,attr"`
}

type NotifReportRS struct {
	XMLName  xml.Name             `xml:"http://www.opentravel.org/OTA/2003/05 OTA_NotifReportRS"`
	Version  string               `xml:"Version,attr"`
	Success  *alpinebits.Success  `xml:"Success"`
	Warnings *alpinebits.Warnings `xml:"Warnings>Warning"`
	Errors   *alpinebits.Errors   `xml:"Errors>Error"`
}
