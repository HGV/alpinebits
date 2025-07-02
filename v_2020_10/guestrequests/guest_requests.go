package guestrequests

import (
	"encoding/xml"
	"time"

	"github.com/HGV/alpinebits/duration"
	"github.com/HGV/alpinebits/v_2020_10/common"
	"github.com/HGV/alpinebits/v_2020_10/rateplans"
	"github.com/HGV/alpinebits/version"
	"github.com/HGV/x/timex"
)

type ReadRQ struct {
	XMLName          xml.Name         `xml:"http://www.opentravel.org/OTA/2003/05 OTA_ReadRQ"`
	Version          string           `xml:"Version,attr"`
	HotelReadRequest HotelReadRequest `xml:"ReadRequests>HotelReadRequest"`
}

var _ version.HotelCodeProvider = (*ReadRQ)(nil)

func (r ReadRQ) HotelCode() string {
	return r.HotelReadRequest.HotelCode
}

type HotelReadRequest struct {
	HotelCode         string             `xml:"HotelCode,attr"`
	SelectionCriteria *SelectionCriteria `xml:"SelectionCriteria,omitempty"`
}

type SelectionCriteria struct {
	Start time.Time `xml:"Start,attr"`
}

type ResRetrieveRS struct {
	common.Response

	XMLName           xml.Name            `xml:"http://www.opentravel.org/OTA/2003/05 OTA_ResRetrieveRS"`
	Version           string              `xml:"Version,attr"`
	HotelReservations *[]HotelReservation `xml:"ReservationsList>HotelReservation"`
}

type ResStatus string

const (
	ResStatusRequested ResStatus = "Requested"
	ResStatusReserved  ResStatus = "Reserved"
	ResStatusModify    ResStatus = "Modify"
	ResStatusCancelled ResStatus = "Cancelled"
)

func (s ResStatus) IsReservation() bool {
	return s == ResStatusReserved || s == ResStatusModify
}

type UniqueIDType int

const (
	UniqueIDTypeReservation  UniqueIDType = 14
	UniqueIDTypeCancellation UniqueIDType = 15
)

type UniqueID struct {
	Type UniqueIDType `xml:"Type,attr"`
	ID   string       `xml:"ID,attr"`
}

type HotelReservation struct {
	CreateDateTime time.Time     `xml:"CreateDateTime,attr"`
	ResStatus      ResStatus     `xml:"ResStatus,attr"`
	UniqueID       UniqueID      `xml:"UniqueID"`
	RoomStays      []RoomStay    `xml:"RoomStays>RoomStay"`
	Customer       Customer      `xml:"ResGuests>ResGuest>Profiles>ProfileInfo>Profile>Customer"`
	ResGlobalInfo  ResGlobalInfo `xml:"ResGlobalInfo"`
}

type RoomStay struct {
	RoomType    *ResRoomType `xml:"RoomTypes>RoomType"`
	RatePlan    *ResRatePlan `xml:"RatePlans>RatePlan"`
	GuestCounts []GuestCount `xml:"GuestCounts>GuestCount"`
	TimeSpan    TimeSpan     `xml:"TimeSpan"`
	Total       *Total       `xml:"Total"`
}

func (r RoomStay) isPrimaryStay() bool {
	return !r.isAlternativeStay()
}

func (r RoomStay) isAlternativeStay() bool {
	return r.RoomType == nil &&
		r.RatePlan == nil &&
		len(r.GuestCounts) == 0 &&
		r.Total == nil
}

type ResRoomType struct {
	RoomTypeCode           string `xml:"RoomTypeCode,attr,omitempty"`
	RoomClassificationCode int    `xml:"RoomClassificationCode,attr,omitempty"`
	RoomType               *int   `xml:"RoomType,attr,omitempty"`
}

type ResRatePlan struct {
	RatePlanCode  string                   `xml:"RatePlanCode,attr,omitempty"`
	Commission    *Commission              `xml:"Commission"`
	MealsIncluded *rateplans.MealsIncluded `xml:"MealsIncluded"`
}

type Commission struct {
	Percent                 *int                     `xml:"Percent,attr"`
	CommissionPayableAmount *CommissionPayableAmount `xml:"CommissionPayableAmount"`
}

type CommissionPayableAmount struct {
	Amount       string `xml:"Amount,attr"`
	CurrencyCode string `xml:"CurrencyCode,attr"`
}

type GuestCount struct {
	Count int  `xml:"Count,attr"`
	Age   *int `xml:"Age,attr"`
}

type TimeSpan struct {
	Start           timex.Date       `xml:"Start,attr,omitempty"`
	End             timex.Date       `xml:"End,attr,omitempty"`
	Duration        *duration.Nights `xml:"Duration,attr"`
	StartDateWindow *StartDateWindow `xml:"StartDateWindow"`
}

type StartDateWindow struct {
	EarliestDate timex.Date `xml:"EarliestDate,attr"`
	LatestDate   timex.Date `xml:"LatestDate,attr"`
}

type Total struct {
	AmountAfterTax string `xml:"AmountAfterTax,attr"`
	CurrencyCode   string `xml:"CurrencyCode,attr"`
}

type Gender string

const (
	GenderMale    Gender = "Male"
	GenderFemale  Gender = "Female"
	GenderUnknown Gender = "Unknown"
)

type Customer struct {
	Gender     *Gender    `xml:"Gender,attr"`
	BirthDate  timex.Date `xml:"BirthDate,attr,omitempty"`
	Language   string     `xml:"Language,attr,omitempty"`
	PersonName PersonName `xml:"PersonName"`
	Phones     []Phone    `xml:"Telephone"`
	Email      *Email     `xml:"Email"`
	Address    *Address   `xml:"Address"`
}

type PersonName struct {
	NamePrefix *string `xml:"NamePrefix"`
	GivenName  string  `xml:"GivenName"`
	Surname    string  `xml:"Surname"`
	NameTitle  *string `xml:"NameTitle"`
}

type PhoneTechType string

const (
	PhoneTechTypeVoice  PhoneTechType = "1"
	PhoneTechTypeFax    PhoneTechType = "3"
	PhoneTechTypeMobile PhoneTechType = "5"
)

type Phone struct {
	PhoneTechType PhoneTechType `xml:"PhoneTechType,attr"`
	PhoneNumber   string        `xml:"PhoneNumber,attr"`
}

type Remark string

const (
	RemarkNewsletterYes Remark = "newsletter:yes"
	RemarkCatalogYes    Remark = "catalog:yes"
)

type Email struct {
	Remark Remark `xml:"Remark,attr,omitempty"`
	Value  string `xml:",innerxml"`
}

type Address struct {
	Language    string       `xml:"Language,attr,omitempty"`
	Remark      Remark       `xml:"Remark,attr,omitempty"`
	AddressLine *string      `xml:"AddressLine,omitempty"`
	CityName    *string      `xml:"CityName,omitempty"`
	PostalCode  *string      `xml:"PostalCode,omitempty"`
	StateProv   *StateProv   `xml:"StateProv,omitempty"`
	CountryName *CountryName `xml:"CountryName,omitempty"`
}

type StateProv struct {
	StateCode string `xml:"StateCode,attr"`
}

type CountryName struct {
	Code string `xml:"Code,attr"`
}

type ResGlobalInfo struct {
	Comments           *[]Comment          `xml:"Comments>Comment"`
	SpecialRequests    *[]SpecialRequest   `xml:"SpecialRequests>SpecialRequest"`
	CancelPenalty      *string             `xml:"CancelPenalties>CancelPenalty>PenaltyDescription>Text"`
	HotelReservationID *HotelReservationID `xml:"HotelReservationIDs>HotelReservationIDs"`
	Profile            *Profile            `xml:"Profiles>ProfileInfo>Profile"`
	BasicPropertyInfo  BasicPropertyInfo   `xml:"BasicPropertyInfo"`
}

type Comment struct {
	Name      string     `xml:"Name,attr"`
	ListItems []ListItem `xml:"ListItem,omitempty"`
	Text      *Text      `xml:"Text,omitempty"`
}

type ListItem struct {
	ListItem int    `xml:"ListItem,attr,omitempty"`
	Language string `xml:"Language,attr,omitempty"`
	Value    string `xml:",innerxml"`
}

type Text struct {
	Value string `xml:",innerxml"`
}

type SpecialRequest struct {
	Name string `xml:"Name,attr"`
	Text *Text  `xml:"Text"`
}

type HotelReservationID struct {
	ResIDType          int     `xml:"ResID_Type,attr"`
	ResIDValue         *string `xml:"ResID_Value,attr"`
	ResIDSource        *string `xml:"ResID_Source,attr"`
	ResIDSourceContext *string `xml:"ResID_SourceContext,attr"`
}

type ProfileType int

const (
	ProfileTypeTravelAgent = 4
)

type Profile struct {
	ProfileType ProfileType `xml:"ProfileType,attr"`
	CompanyInfo CompanyInfo `xml:"CompanyInfo"`
}

type CompanyInfo struct {
	CompanyName   CompanyName `xml:"CompanyName"`
	AddressInfo   *Address    `xml:"AddressInfo"`
	TelephoneInfo *Phone      `xml:"TelephoneInfo"`
	Email         *Email      `xml:"Email"`
}

type CompanyName struct {
	Code        string `xml:"Code,attr"`
	CodeContext string `xml:"CodeContext,attr"`
	Value       string `xml:",innerxml"`
}

type BasicPropertyInfo struct {
	HotelCode string `xml:"HotelCode,attr"`
	HotelName string `xml:"HotelName,attr,omitempty"`
}
