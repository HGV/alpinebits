package v_2018_10

import (
	"encoding/xml"

	"github.com/HGV/alpinebits/duration"
	"github.com/HGV/alpinebits/version"
	"github.com/HGV/x/timex"
)

type HotelRatePlanNotifRQ struct {
	XMLName   xml.Name  `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelRatePlanNotifRQ"`
	UniqueID  *UniqueID `xml:"UniqueID,omitempty"`
	RatePlans RatePlans `xml:"RatePlans"`
}

type RatePlans struct {
	HotelCode string     `xml:"HotelCode,attr"`
	HotelName string     `xml:"HotelName,attr"`
	RatePlans []RatePlan `xml:"RatePlan"`
}

type RatePlanNotifType string

const (
	RatePlanNotifTypeNew     RatePlanNotifType = "New"
	RatePlanNotifTypeOverlay RatePlanNotifType = "Overlay"
	RatePlanNotifTypeRemove  RatePlanNotifType = "Remove"
)

type RatePlan struct {
	RatePlanNotifType RatePlanNotifType   `xml:"RatePlanNotifType,attr"`
	RatePlanType      int                 `xml:"RatePlanType,attr,omitempty"` //12 promotional
	CurrencyCode      string              `xml:"CurrencyCode,attr"`
	RatePlanCode      string              `xml:"RatePlanCode,attr"`
	RatePlanID        string              `xml:"RatePlanID,attr,omitempty"`
	RatePlanQualifier bool                `xml:"RatePlanQualifier,attr,omitempty"`
	BookingRules      []BookingRule       `xml:"BookingRules>BookingRule"`
	Rates             []Rate              `xml:"Rates>Rate"`
	Supplements       []Supplement        `xml:"Supplements>Supplement"`
	Offers            []Offer             `xml:"Offers>Offer"`
	Descriptions      RatePlanDescription `xml:"Description"`
}

type CodeContext string

const (
	CodeContextRoomType CodeContext = "ROOMTYPE"
)

type BookingRule struct {
	Start               timex.Date         `xml:"Start,attr"`
	End                 timex.Date         `xml:"End,attr"`
	Code                string             `xml:"Code,attr"`
	CodeContext         CodeContext        `xml:"CodeContext,attr"`
	LengthsOfStay       []LengthOfStay     `xml:"LengthsOfStay>LengthOfStay"`
	ArrivalDaysOfWeek   *DaysOfWeek        `xml:"DOW_Restrictions>ArrivalDaysOfWeek"`
	DepartureDaysOfWeek *DaysOfWeek        `xml:"DOW_Restrictions>DepartureDaysOfWeek"`
	RestrictionStatus   *RestrictionStatus `xml:"RestrictionStatus,omitempty"`
}

var _ version.DateRangeProvider = (*BookingRule)(nil)

func (b BookingRule) DateRange() timex.DateRange {
	return timex.DateRange{
		Start: b.Start,
		End:   b.End,
	}
}

type StayType string

const (
	StayTypeMinArrival StayType = "SetMinLOS"
	StayTypeMinThrough StayType = "SetForwardMinStay"
	StayTypeMaxArrival StayType = "SetMaxLOS"
	StayTypeMaxThrough StayType = "SetForwardMaxStay"
)

type TimeUnit string

const (
	TimeUnitDay TimeUnit = "Day"
)

type LengthOfStay struct {
	Time              int      `xml:"Time,attr"`
	TimeUnit          TimeUnit `xml:"TimeUnit,attr"`
	MinMaxMessageType StayType `xml:"MinMaxMessageType,attr"`
}

type DaysOfWeek struct {
	Mon  *bool `xml:"Mon,attr,omitempty"`
	Tue  *bool `xml:"Tue,attr,omitempty"`
	Weds *bool `xml:"Weds,attr,omitempty"`
	Thur *bool `xml:"Thur,attr,omitempty"`
	Fri  *bool `xml:"Fri,attr,omitempty"`
	Sat  *bool `xml:"Sat,attr,omitempty"`
	Sun  *bool `xml:"Sun,attr,omitempty"`
}

type Restriction string

const (
	RestrictionMaster Restriction = "Master"
)

// TODO: Rename
type Status2 string

const (
	Status2Open  Status2 = "Open"
	Status2Close Status2 = "Close"
)

type RestrictionStatus struct {
	Restriction Restriction `xml:"Restriction,attr"`
	Status      Status2     `xml:"Status,attr"`
}

type Rate struct {
	RateTimeUnit           TimeUnit                `xml:"RateTimeUnit,attr,omitempty"`
	UnitMultiplier         int                     `xml:"UnitMultiplier,attr,omitempty"`
	InvTypeCode            string                  `xml:"InvTypeCode,attr,omitempty"`
	Start                  *timex.Date             `xml:"Start,attr,omitempty"`
	End                    *timex.Date             `xml:"End,attr,omitempty"`
	BaseByGuestAmts        []BaseByGuestAmt        `xml:"BaseByGuestAmts>BaseByGuestAmt"`
	AdditionalGuestAmounts []AdditionalGuestAmount `xml:"AdditionalGuestAmounts>AdditionalGuestAmount"`
	MealsIncluded          *MealsIncluded          `xml:"MealsIncluded,omitempty"`
}

var _ version.DateRangeProvider = (*Rate)(nil)

func (r Rate) DateRange() timex.DateRange {
	if r.Start == nil || r.End == nil {
		return timex.DateRange{}
	}
	return timex.DateRange{
		Start: *r.Start,
		End:   *r.End,
	}
}

func (r Rate) IsStaticRate() bool {
	return r.Start == nil &&
		r.End == nil &&
		len(r.BaseByGuestAmts) == 1 &&
		len(r.AdditionalGuestAmounts) == 0 &&
		r.MealsIncluded != nil
}

type RatePlanChargeType int

const (
	RatePlanChargeTypePerPerson RatePlanChargeType = 7
	RatePlanChargeTypePerRoom   RatePlanChargeType = 25
)

type AgeQualifyingCode int

const (
	AgeQualifyingCodeAdult AgeQualifyingCode = 10
	AgeQualifyingCodeChild AgeQualifyingCode = 8
)

type BaseByGuestAmt struct {
	Type              *RatePlanChargeType `xml:"Type,attr,omitempty"`
	NumberOfGuests    *int                `xml:"NumberOfGuests,attr,omitempty"`
	AgeQualifyingCode *AgeQualifyingCode  `xml:"AgeQualifyingCode,attr,omitempty"`
	AmountAfterTax    *string             `xml:"AmountAfterTax,attr,omitempty"`
}

type AdditionalGuestAmount struct {
	AgeQualifyingCode *AgeQualifyingCode `xml:"AgeQualifyingCode,attr"`
	MinAge            *int               `xml:"MinAge,attr,omitempty"`
	MaxAge            *int               `xml:"MaxAge,attr,omitempty"`
	Amount            *string            `xml:"Amount,attr"`
}

type MealPlan int

const (
	MealPlanAllInclusive    MealPlan = 1
	MealPlanBedAndBreakfast MealPlan = 3
	MealPlanFullBoard       MealPlan = 10
	MealPlanHalfBoard       MealPlan = 12
	MealPlanRoomOnly        MealPlan = 14
)

type MealsIncluded struct {
	MealPlanIndicator bool     `xml:"MealPlanIndicator,attr"`
	MealPlanCodes     MealPlan `xml:"MealPlanCodes,attr"`
}

type InvType string

const (
	InvTypeExtra InvType = "EXTRA"
)

type SupplementChargeType int

const (
	SupplementChargeTypePerPerson SupplementChargeType = 7
	SupplementChargeTypePerRoom   SupplementChargeType = 25
)

type Supplement struct {
	InvType                 InvType                `xml:"InvType,attr"`
	InvCode                 string                 `xml:"InvCode,attr"`
	AddToBasicRateIndicator *bool                  `xml:"AddToBasicRateIndicator,attr,omitempty"`
	MandatoryIndicator      *bool                  `xml:"MandatoryIndicator,attr,omitempty"`
	ChargeTypeCode          *SupplementChargeType  `xml:"ChargeTypeCode,attr,omitempty"`
	PrerequisiteInventory   *PrerequisiteInventory `xml:"PrerequisiteInventory,omitempty"`
	Descriptions            *RatePlanDescription   `xml:"Description,omitempty"`
	Start                   *timex.Date            `xml:"Start,attr,omitempty"`
	End                     *timex.Date            `xml:"End,attr,omitempty"`
	Amount                  *string                `xml:"Amount,attr,omitempty"`
}

var _ version.DateRangeProvider = (*Supplement)(nil)

func (s Supplement) DateRange() timex.DateRange {
	if s.Start == nil || s.End == nil {
		return timex.DateRange{}
	}
	return timex.DateRange{
		Start: *s.Start,
		End:   *s.End,
	}
}

func (s Supplement) isStaticSupplement() bool {
	return (s.AddToBasicRateIndicator != nil && *s.AddToBasicRateIndicator) &&
		s.MandatoryIndicator != nil &&
		s.ChargeTypeCode != nil &&
		s.Start == nil &&
		s.End == nil &&
		s.Amount == nil
}

func (s Supplement) isDateDependingSupplement() bool {
	return !s.isStaticSupplement()
}

type PrerequisiteInventoryInvType string

const (
	PrerequisiteInventoryInvTypeAlpineBitsDOW PrerequisiteInventoryInvType = "ALPINEBITSDOW"
	PrerequisiteInventoryInvTypeRoomType      PrerequisiteInventoryInvType = "ROOMTYPE"
)

type PrerequisiteInventory struct {
	InvType PrerequisiteInventoryInvType `xml:"InvType,attr"`
	InvCode string                       `xml:"InvCode,attr"`
}

type Offer struct {
	OfferRule *OfferRule `xml:"OfferRules>OfferRule"`
	Discount  *Discount  `xml:"Discount,omitempty"`
	Guest     *Guest     `xml:"Guests>Guest"`
}

func (o Offer) IsFreeNightOffer() bool {
	return o.Discount != nil &&
		o.Discount.NightsRequired != 0 &&
		o.Discount.NightsDiscounted != 0
}

func (o Offer) IsFamilyOffer() bool {
	return o.Discount != nil &&
		o.Guest != nil
}

type OfferRule struct {
	MinAdvancedBookingOffset *duration.Days `xml:"MinAdvancedBookingOffset,attr,omitempty"`
	MaxAdvancedBookingOffset *duration.Days `xml:"MaxAdvancedBookingOffset,attr,omitempty"`
	LengthsOfStay            []LengthOfStay `xml:"LengthsOfStay>LengthOfStay"`
	ArrivalDaysOfWeek        *DaysOfWeek    `xml:"DOW_Restrictions>ArrivalDaysOfWeek"`
	DepartureDaysOfWeek      *DaysOfWeek    `xml:"DOW_Restrictions>DepartureDaysOfWeek"`
	Occupancies              []Occupancy    `xml:"Occupancy,omitempty"`
}

type Occupancy struct {
	AgeQualifyingCode AgeQualifyingCode `xml:"AgeQualifyingCode,attr"`
	MinAge            *int              `xml:"MinAge,attr,omitempty"`
	MaxAge            *int              `xml:"MaxAge,attr,omitempty"`
	MinOccupancy      *int              `xml:"MinOccupancy,attr,omitempty"`
	MaxOccupancy      *int              `xml:"MaxOccupancy,attr,omitempty"`
}

func (o Occupancy) isAdult() bool {
	return o.AgeQualifyingCode == AgeQualifyingCodeAdult
}

func (o Occupancy) isChild() bool {
	return o.AgeQualifyingCode == AgeQualifyingCodeChild
}

type Discount struct {
	Percent          int    `xml:"Percent,attr"`
	NightsRequired   int    `xml:"NightsRequired,attr,omitempty"`
	NightsDiscounted int    `xml:"NightsDiscounted,attr,omitempty"`
	DiscountPattern  string `xml:"DiscountPattern,attr,omitempty"`
}

type Guest struct {
	AgeQualifyingCode       AgeQualifyingCode `xml:"AgeQualifyingCode,attr"`
	MaxAge                  int               `xml:"MaxAge,attr"`
	MinCount                int               `xml:"MinCount,attr"`
	FirstQualifyingPosition int               `xml:"FirstQualifyingPosition,attr"`
	LastQualifyingPosition  int               `xml:"LastQualifyingPosition,attr"`
}

type RatePlanDescription struct {
	Titles       []Description
	Intros       []Description
	Descriptions []Description
	Codelist     []ListItem
	Gallery      []GalleryItem
}

// TODO: Custom Unmarshal

type GalleryItem struct {
	Image           URL
	Descriptions    []Description
	CopyrightNotice string
	Attribution     URL
}

type HotelRatePlanNotifRS struct {
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelRatePlanNotifRS"`
	Version string   `xml:"Version,attr"`
}
