// Package rateplans holds the AlpineBits 2018-10 RatePlans messages
// (OTA_HotelRatePlanNotif, §4.5).
package rateplans

import (
	"encoding/xml"

	"github.com/HGV/alpinebits"
	"github.com/HGV/x/timex"
)

type HotelRatePlanNotifRQ struct {
	XMLName   xml.Name  `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelRatePlanNotifRQ"`
	Version   string    `xml:"Version,attr"`
	UniqueID  *UniqueID `xml:"UniqueID,omitempty"`
	RatePlans RatePlans `xml:"RatePlans"`
}

// IsCompleteSet reports whether the message lists all rate plans of the hotel;
// the server expires the ones it does not contain (AlpineBits 2018-10 §4.5.3).
func (r HotelRatePlanNotifRQ) IsCompleteSet() bool {
	return r.UniqueID != nil
}

// IsReset reports whether r deletes all rate plans of the hotel: a CompleteSet
// with a single empty RatePlan.
func (r HotelRatePlanNotifRQ) IsReset() bool {
	plans := r.RatePlans.RatePlans
	return r.IsCompleteSet() && len(plans) == 1 && plans[0].RatePlanCode == "" && plans[0].isEmpty()
}

type UniqueIDType string

const (
	UniqueIDTypeReference UniqueIDType = "16"
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

// HotelCode implements alpinebits.HotelCoded.
func (r HotelRatePlanNotifRQ) HotelCode() string {
	return r.RatePlans.HotelCode
}

type RatePlans struct {
	HotelCode string     `xml:"HotelCode,attr,omitempty"`
	HotelName string     `xml:"HotelName,attr,omitempty"`
	RatePlans []RatePlan `xml:"RatePlan"`
}

// RatePlanNotifType represents the type of rate plan notification.
type RatePlanNotifType string

const (
	RatePlanNotifTypeOverlay RatePlanNotifType = "Overlay"
	RatePlanNotifTypeNew     RatePlanNotifType = "New"
	RatePlanNotifTypeRemove  RatePlanNotifType = "Remove"
)

// RatePlanType represents the type of rate plan.
type RatePlanType string

const (
	RatePlanType12 RatePlanType = "12"
)

type RatePlan struct {
	RatePlanNotifType RatePlanNotifType `xml:"RatePlanNotifType,attr,omitempty"`
	CurrencyCode      string            `xml:"CurrencyCode,attr,omitempty"`
	RatePlanCode      string            `xml:"RatePlanCode,attr,omitempty"`
	RatePlanType      RatePlanType      `xml:"RatePlanType,attr,omitempty"`
	RatePlanCategory  string            `xml:"RatePlanCategory,attr,omitempty"`
	RatePlanID        string            `xml:"RatePlanID,attr,omitempty"`
	RatePlanQualifier *bool             `xml:"RatePlanQualifier,attr,omitempty"`
	BookingRules      *[]BookingRule    `xml:"BookingRules>BookingRule"`
	Rates             *[]Rate           `xml:"Rates>Rate"`
	Supplements       *[]Supplement     `xml:"Supplements>Supplement"`
	Offers            *[]Offer          `xml:"Offers>Offer"`
	Descriptions      []Description     `xml:"Description,omitempty"`
}

func (p RatePlan) isEmpty() bool {
	return p.BookingRules == nil && p.Rates == nil && p.Supplements == nil && p.Offers == nil && len(p.Descriptions) == 0
}

// CodeContext represents the context for booking rule codes.
type CodeContext string

const (
	CodeContextRoomType CodeContext = "ROOMTYPE"
)

type BookingRule struct {
	CodeContext       CodeContext        `xml:"CodeContext,attr,omitempty"`
	Code              string             `xml:"Code,attr,omitempty"`
	Start             *timex.Date        `xml:"Start,attr,omitempty"`
	End               *timex.Date        `xml:"End,attr,omitempty"`
	LengthsOfStay     *[]LengthOfStay    `xml:"LengthsOfStay>LengthOfStay"`
	DOWRestrictions   *DOWRestrictions   `xml:"DOW_Restrictions,omitempty"`
	RestrictionStatus *RestrictionStatus `xml:"RestrictionStatus,omitempty"`
}

// DateRange returns the period from Start to End.
func (r BookingRule) DateRange() timex.DateRange {
	return dateRange(r.Start, r.End)
}

func dateRange(start, end *timex.Date) timex.DateRange {
	var r timex.DateRange
	if start != nil {
		r.Start = *start
	}
	if end != nil {
		r.End = *end
	}
	return r
}

// TimeUnit represents the unit of time.
type TimeUnit string

const (
	TimeUnitDay TimeUnit = "Day"
)

// MinMaxMessageType represents the type of min/max length of stay.
type MinMaxMessageType string

const (
	MinMaxMessageTypeSetMinLOS         MinMaxMessageType = "SetMinLOS"
	MinMaxMessageTypeSetForwardMinStay MinMaxMessageType = "SetForwardMinStay"
	MinMaxMessageTypeSetMaxLOS         MinMaxMessageType = "SetMaxLOS"
	MinMaxMessageTypeSetForwardMaxStay MinMaxMessageType = "SetForwardMaxStay"
)

type LengthOfStay struct {
	Time              int               `xml:"Time,attr"`
	TimeUnit          TimeUnit          `xml:"TimeUnit,attr"`
	MinMaxMessageType MinMaxMessageType `xml:"MinMaxMessageType,attr"`
}

type DOWRestrictions struct {
	ArrivalDaysOfWeek   *DaysOfWeek `xml:"ArrivalDaysOfWeek,omitempty"`
	DepartureDaysOfWeek *DaysOfWeek `xml:"DepartureDaysOfWeek,omitempty"`
}

// DaysOfWeek restricts the days that are false; nil days are not restricted.
type DaysOfWeek struct {
	Mon  *bool `xml:"Mon,attr,omitempty"`
	Tue  *bool `xml:"Tue,attr,omitempty"`
	Weds *bool `xml:"Weds,attr,omitempty"`
	Thur *bool `xml:"Thur,attr,omitempty"`
	Fri  *bool `xml:"Fri,attr,omitempty"`
	Sat  *bool `xml:"Sat,attr,omitempty"`
	Sun  *bool `xml:"Sun,attr,omitempty"`
}

// Status represents the restriction status.
type Status string

const (
	StatusOpen  Status = "Open"
	StatusClose Status = "Close"
)

// Restriction represents the restriction type.
type Restriction string

const (
	RestrictionMaster Restriction = "Master"
)

type RestrictionStatus struct {
	Restriction Restriction `xml:"Restriction,attr,omitempty"`
	Status      Status      `xml:"Status,attr,omitempty"`
}

type Rate struct {
	Start                  *timex.Date              `xml:"Start,attr,omitempty"`
	End                    *timex.Date              `xml:"End,attr,omitempty"`
	RateTimeUnit           TimeUnit                 `xml:"RateTimeUnit,attr,omitempty"`
	UnitMultiplier         int                      `xml:"UnitMultiplier,attr,omitempty"`
	InvTypeCode            string                   `xml:"InvTypeCode,attr,omitempty"`
	BaseByGuestAmts        *[]BaseByGuestAmt        `xml:"BaseByGuestAmts>BaseByGuestAmt"`
	AdditionalGuestAmounts *[]AdditionalGuestAmount `xml:"AdditionalGuestAmounts>AdditionalGuestAmount"`
	MealsIncluded          *MealsIncluded           `xml:"MealsIncluded,omitempty"`
}

// DateRange returns the period from Start to End.
func (r Rate) DateRange() timex.DateRange {
	return dateRange(r.Start, r.End)
}

// BaseByGuestAmtType represents the type of base by guest amount.
type BaseByGuestAmtType string

const (
	BaseByGuestAmtType7  BaseByGuestAmtType = "7"
	BaseByGuestAmtType25 BaseByGuestAmtType = "25"
)

type BaseByGuestAmt struct {
	NumberOfGuests    int                `xml:"NumberOfGuests,attr,omitempty"`
	AmountAfterTax    alpinebits.Decimal `xml:"AmountAfterTax,attr,omitempty"`
	CurrencyCode      string             `xml:"CurrencyCode,attr,omitempty"`
	Type              BaseByGuestAmtType `xml:"Type,attr,omitempty"`
	AgeQualifyingCode AgeQualifyingCode  `xml:"AgeQualifyingCode,attr,omitempty"`
}

type AdditionalGuestAmount struct {
	Amount            *alpinebits.Decimal `xml:"Amount,attr,omitempty"`
	AgeQualifyingCode AgeQualifyingCode   `xml:"AgeQualifyingCode,attr,omitempty"`
	MinAge            int                 `xml:"MinAge,attr,omitempty"`
	MaxAge            int                 `xml:"MaxAge,attr,omitempty"`
}

type MealsIncluded struct {
	Breakfast         bool   `xml:"Breakfast,attr,omitempty"`
	Lunch             bool   `xml:"Lunch,attr,omitempty"`
	Dinner            bool   `xml:"Dinner,attr,omitempty"`
	MealPlanCodes     string `xml:"MealPlanCodes,attr,omitempty"`
	MealPlanIndicator bool   `xml:"MealPlanIndicator,attr,omitempty"`
}

// InvType represents the inventory type for supplements.
type InvType string

const (
	InvTypeExtra           InvType = "EXTRA"
	InvTypeAlpineBitsExtra InvType = "ALPINEBITSEXTRA"
)

// ChargeTypeCode represents the charge type code.
type ChargeTypeCode string

const (
	ChargeTypeCode1  ChargeTypeCode = "1"
	ChargeTypeCode12 ChargeTypeCode = "12"
	ChargeTypeCode18 ChargeTypeCode = "18"
	ChargeTypeCode19 ChargeTypeCode = "19"
	ChargeTypeCode20 ChargeTypeCode = "20"
	ChargeTypeCode21 ChargeTypeCode = "21"
	ChargeTypeCode24 ChargeTypeCode = "24"
)

type Supplement struct {
	AddToBasicRateIndicator bool                   `xml:"AddToBasicRateIndicator,attr,omitempty"`
	ChargeTypeCode          ChargeTypeCode         `xml:"ChargeTypeCode,attr,omitempty"`
	Amount                  *alpinebits.Decimal    `xml:"Amount,attr,omitempty"` // nil: not available
	InvType                 InvType                `xml:"InvType,attr"`
	InvCode                 string                 `xml:"InvCode,attr"`
	MandatoryIndicator      *bool                  `xml:"MandatoryIndicator,attr,omitempty"`
	Start                   *timex.Date            `xml:"Start,attr,omitempty"`
	End                     *timex.Date            `xml:"End,attr,omitempty"`
	PrerequisiteInventory   *PrerequisiteInventory `xml:"PrerequisiteInventory,omitempty"`
	Descriptions            []Description          `xml:"Description,omitempty"`
}

// DateRange returns the period from Start to End.
func (s Supplement) DateRange() timex.DateRange {
	return dateRange(s.Start, s.End)
}

// PrerequisiteInvType represents the inventory type for prerequisites.
type PrerequisiteInvType string

const (
	PrerequisiteInvTypeAlpineBitsDOW PrerequisiteInvType = "ALPINEBITSDOW"
	PrerequisiteInvTypeRoomType      PrerequisiteInvType = "ROOMTYPE"
)

type PrerequisiteInventory struct {
	InvCode string              `xml:"InvCode,attr"`
	InvType PrerequisiteInvType `xml:"InvType,attr"`
}

type Offer struct {
	OfferCode    string        `xml:"OfferCode,attr,omitempty"`
	OfferRules   *OfferRules   `xml:"OfferRules,omitempty"`
	Discount     *Discount     `xml:"Discount,omitempty"`
	Guests       *[]Guest      `xml:"Guests>Guest"`
	Descriptions []Description `xml:"Description,omitempty"`
}

type OfferRules struct {
	OfferRule *OfferRule `xml:"OfferRule,omitempty"`
}

type OfferRule struct {
	MinAdvancedBookingOffset string           `xml:"MinAdvancedBookingOffset,attr,omitempty"`
	MaxAdvancedBookingOffset string           `xml:"MaxAdvancedBookingOffset,attr,omitempty"`
	LengthsOfStay            *[]LengthOfStay  `xml:"LengthsOfStay>LengthOfStay"`
	DOWRestrictions          *DOWRestrictions `xml:"DOW_Restrictions,omitempty"`
	Occupancies              []Occupancy      `xml:"Occupancy,omitempty"`
}

// AgeQualifyingCode is from the OTA list AQC. Most of its attributes are
// xs:integer, where "010" is a valid 10, so it is an int.
type AgeQualifyingCode int

const (
	AgeQualifyingCode8  AgeQualifyingCode = 8
	AgeQualifyingCode10 AgeQualifyingCode = 10
)

type Occupancy struct {
	AgeQualifyingCode AgeQualifyingCode `xml:"AgeQualifyingCode,attr"`
	MinAge            int               `xml:"MinAge,attr,omitempty"`
	MaxAge            int               `xml:"MaxAge,attr,omitempty"`
	MinOccupancy      int               `xml:"MinOccupancy,attr,omitempty"`
	MaxOccupancy      int               `xml:"MaxOccupancy,attr,omitempty"`
}

type Guest struct {
	AgeQualifyingCode       AgeQualifyingCode `xml:"AgeQualifyingCode,attr"`
	MaxAge                  int               `xml:"MaxAge,attr"`
	MinCount                int               `xml:"MinCount,attr"`
	FirstQualifyingPosition int               `xml:"FirstQualifyingPosition,attr"`
	LastQualifyingPosition  int               `xml:"LastQualifyingPosition,attr"`
}

type Discount struct {
	NightsRequired   int     `xml:"NightsRequired,attr,omitempty"`
	NightsDiscounted int     `xml:"NightsDiscounted,attr,omitempty"`
	DiscountPattern  string  `xml:"DiscountPattern,attr,omitempty"`
	Percent          float64 `xml:"Percent,attr,omitempty"`
}

// DescriptionName represents the name of description sections.
type DescriptionName string

const (
	DescriptionNameTitle       DescriptionName = "title"
	DescriptionNameIntro       DescriptionName = "intro"
	DescriptionNameDescription DescriptionName = "description"
	DescriptionNameGallery     DescriptionName = "gallery"
	DescriptionNameCodelist    DescriptionName = "codelist"
)

type Description struct {
	Name      DescriptionName   `xml:"Name,attr"`
	ListItems []string          `xml:"ListItem,omitempty"`
	Images    []string          `xml:"Image,omitempty"`
	Texts     []alpinebits.Text `xml:"Text,omitempty"`
	URLs      []string          `xml:"URL,omitempty"`
}

type HotelRatePlanNotifRS struct {
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelRatePlanNotifRS"`
	alpinebits.Response
}
