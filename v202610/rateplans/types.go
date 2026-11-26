// Package rateplans holds the AlpineBits 2026-10 RatePlans messages
// (OTA_HotelRatePlanNotif, §4.4). Their structure is unchanged from 2024-10, so
// the types are aliases; RatePlanNotifType is New again, and pets have their own
// values.
package rateplans

import rateplans202410 "github.com/HGV/alpinebits/v202410/rateplans"

type (
	HotelRatePlanNotifRQ  = rateplans202410.HotelRatePlanNotifRQ
	HotelRatePlanNotifRS  = rateplans202410.HotelRatePlanNotifRS
	UniqueIDType          = rateplans202410.UniqueIDType
	Instance              = rateplans202410.Instance
	UniqueID              = rateplans202410.UniqueID
	RatePlans             = rateplans202410.RatePlans
	RatePlanNotifType     = rateplans202410.RatePlanNotifType
	RatePlanType          = rateplans202410.RatePlanType
	RatePlan              = rateplans202410.RatePlan
	CodeContext           = rateplans202410.CodeContext
	BookingRule           = rateplans202410.BookingRule
	TimeUnit              = rateplans202410.TimeUnit
	MinMaxMessageType     = rateplans202410.MinMaxMessageType
	LengthOfStay          = rateplans202410.LengthOfStay
	DOWRestrictions       = rateplans202410.DOWRestrictions
	DaysOfWeek            = rateplans202410.DaysOfWeek
	Status                = rateplans202410.Status
	Restriction           = rateplans202410.Restriction
	RestrictionStatus     = rateplans202410.RestrictionStatus
	Rate                  = rateplans202410.Rate
	BaseByGuestAmtType    = rateplans202410.BaseByGuestAmtType
	BaseByGuestAmt        = rateplans202410.BaseByGuestAmt
	AdditionalGuestAmount = rateplans202410.AdditionalGuestAmount
	MealsIncluded         = rateplans202410.MealsIncluded
	InvType               = rateplans202410.InvType
	ChargeTypeCode        = rateplans202410.ChargeTypeCode
	Supplement            = rateplans202410.Supplement
	PrerequisiteInvType   = rateplans202410.PrerequisiteInvType
	PrerequisiteInventory = rateplans202410.PrerequisiteInventory
	Offer                 = rateplans202410.Offer
	OfferRules            = rateplans202410.OfferRules
	OfferRule             = rateplans202410.OfferRule
	AgeQualifyingCode     = rateplans202410.AgeQualifyingCode
	Occupancy             = rateplans202410.Occupancy
	Guest                 = rateplans202410.Guest
	Discount              = rateplans202410.Discount
	DescriptionName       = rateplans202410.DescriptionName
	Description           = rateplans202410.Description
)

const (
	UniqueIDTypeReference = rateplans202410.UniqueIDTypeReference
	InstanceCompleteSet   = rateplans202410.InstanceCompleteSet

	RatePlanNotifTypeOverlay = rateplans202410.RatePlanNotifTypeOverlay
	RatePlanNotifTypeNew     = rateplans202410.RatePlanNotifTypeNew
	RatePlanNotifTypeRemove  = rateplans202410.RatePlanNotifTypeRemove

	RatePlanType12      = rateplans202410.RatePlanType12
	CodeContextRoomType = rateplans202410.CodeContextRoomType
	TimeUnitDay         = rateplans202410.TimeUnitDay

	MinMaxMessageTypeSetMinLOS         = rateplans202410.MinMaxMessageTypeSetMinLOS
	MinMaxMessageTypeSetForwardMinStay = rateplans202410.MinMaxMessageTypeSetForwardMinStay
	MinMaxMessageTypeSetMaxLOS         = rateplans202410.MinMaxMessageTypeSetMaxLOS
	MinMaxMessageTypeSetForwardMaxStay = rateplans202410.MinMaxMessageTypeSetForwardMaxStay

	StatusOpen        = rateplans202410.StatusOpen
	StatusClose       = rateplans202410.StatusClose
	RestrictionMaster = rateplans202410.RestrictionMaster

	BaseByGuestAmtType7  = rateplans202410.BaseByGuestAmtType7
	BaseByGuestAmtType25 = rateplans202410.BaseByGuestAmtType25

	InvTypeExtra = rateplans202410.InvTypeExtra

	ChargeTypeCode1  = rateplans202410.ChargeTypeCode1
	ChargeTypeCode12 = rateplans202410.ChargeTypeCode12
	ChargeTypeCode18 = rateplans202410.ChargeTypeCode18
	ChargeTypeCode19 = rateplans202410.ChargeTypeCode19
	ChargeTypeCode20 = rateplans202410.ChargeTypeCode20
	ChargeTypeCode21 = rateplans202410.ChargeTypeCode21
	ChargeTypeCode24 = rateplans202410.ChargeTypeCode24

	PrerequisiteInvTypeAlpineBitsDOW = rateplans202410.PrerequisiteInvTypeAlpineBitsDOW
	PrerequisiteInvTypeRoomType      = rateplans202410.PrerequisiteInvTypeRoomType

	AgeQualifyingCode8  = rateplans202410.AgeQualifyingCode8
	AgeQualifyingCode10 = rateplans202410.AgeQualifyingCode10

	DescriptionNameTitle       = rateplans202410.DescriptionNameTitle
	DescriptionNameIntro       = rateplans202410.DescriptionNameIntro
	DescriptionNameDescription = rateplans202410.DescriptionNameDescription
	DescriptionNameGallery     = rateplans202410.DescriptionNameGallery
	DescriptionNameCodelist    = rateplans202410.DescriptionNameCodelist

	// CodeContextAlpineBitsPets marks the BookingRule allowing or forbidding pets;
	// its Code is ANY.
	CodeContextAlpineBitsPets CodeContext = "ALPINEBITS-PETS"

	// InvTypeAlpineBits replaces ALPINEBITSEXTRA; InvTypeAlpineBitsPets marks
	// pets' supplements.
	InvTypeAlpineBits     InvType = "ALPINEBITS"
	InvTypeAlpineBitsPets InvType = "ALPINEBITS-PETS"
)
