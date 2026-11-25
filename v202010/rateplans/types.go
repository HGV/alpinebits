// Package rateplans holds the AlpineBits 2020-10 RatePlans messages
// (OTA_HotelRatePlanNotif, §4.5). Their structure is unchanged from 2018-10,
// so the types are aliases; CurrencyCode is no longer fixed to EUR.
package rateplans

import rateplans201810 "github.com/HGV/alpinebits/v201810/rateplans"

type (
	HotelRatePlanNotifRQ  = rateplans201810.HotelRatePlanNotifRQ
	HotelRatePlanNotifRS  = rateplans201810.HotelRatePlanNotifRS
	UniqueIDType          = rateplans201810.UniqueIDType
	Instance              = rateplans201810.Instance
	UniqueID              = rateplans201810.UniqueID
	RatePlans             = rateplans201810.RatePlans
	RatePlanNotifType     = rateplans201810.RatePlanNotifType
	RatePlanType          = rateplans201810.RatePlanType
	RatePlan              = rateplans201810.RatePlan
	CodeContext           = rateplans201810.CodeContext
	BookingRule           = rateplans201810.BookingRule
	TimeUnit              = rateplans201810.TimeUnit
	MinMaxMessageType     = rateplans201810.MinMaxMessageType
	LengthOfStay          = rateplans201810.LengthOfStay
	DOWRestrictions       = rateplans201810.DOWRestrictions
	DaysOfWeek            = rateplans201810.DaysOfWeek
	Status                = rateplans201810.Status
	Restriction           = rateplans201810.Restriction
	RestrictionStatus     = rateplans201810.RestrictionStatus
	Rate                  = rateplans201810.Rate
	BaseByGuestAmtType    = rateplans201810.BaseByGuestAmtType
	BaseByGuestAmt        = rateplans201810.BaseByGuestAmt
	AdditionalGuestAmount = rateplans201810.AdditionalGuestAmount
	MealsIncluded         = rateplans201810.MealsIncluded
	InvType               = rateplans201810.InvType
	ChargeTypeCode        = rateplans201810.ChargeTypeCode
	Supplement            = rateplans201810.Supplement
	PrerequisiteInvType   = rateplans201810.PrerequisiteInvType
	PrerequisiteInventory = rateplans201810.PrerequisiteInventory
	Offer                 = rateplans201810.Offer
	OfferRules            = rateplans201810.OfferRules
	OfferRule             = rateplans201810.OfferRule
	AgeQualifyingCode     = rateplans201810.AgeQualifyingCode
	Occupancy             = rateplans201810.Occupancy
	Guest                 = rateplans201810.Guest
	Discount              = rateplans201810.Discount
	DescriptionName       = rateplans201810.DescriptionName
	Description           = rateplans201810.Description
)

const (
	UniqueIDTypeReference = rateplans201810.UniqueIDTypeReference
	InstanceCompleteSet   = rateplans201810.InstanceCompleteSet

	RatePlanNotifTypeOverlay = rateplans201810.RatePlanNotifTypeOverlay
	RatePlanNotifTypeNew     = rateplans201810.RatePlanNotifTypeNew
	RatePlanNotifTypeRemove  = rateplans201810.RatePlanNotifTypeRemove

	RatePlanType12      = rateplans201810.RatePlanType12
	CodeContextRoomType = rateplans201810.CodeContextRoomType
	TimeUnitDay         = rateplans201810.TimeUnitDay

	MinMaxMessageTypeSetMinLOS         = rateplans201810.MinMaxMessageTypeSetMinLOS
	MinMaxMessageTypeSetForwardMinStay = rateplans201810.MinMaxMessageTypeSetForwardMinStay
	MinMaxMessageTypeSetMaxLOS         = rateplans201810.MinMaxMessageTypeSetMaxLOS
	MinMaxMessageTypeSetForwardMaxStay = rateplans201810.MinMaxMessageTypeSetForwardMaxStay

	StatusOpen        = rateplans201810.StatusOpen
	StatusClose       = rateplans201810.StatusClose
	RestrictionMaster = rateplans201810.RestrictionMaster

	BaseByGuestAmtType7  = rateplans201810.BaseByGuestAmtType7
	BaseByGuestAmtType25 = rateplans201810.BaseByGuestAmtType25

	InvTypeExtra           = rateplans201810.InvTypeExtra
	InvTypeAlpineBitsExtra = rateplans201810.InvTypeAlpineBitsExtra

	ChargeTypeCode1  = rateplans201810.ChargeTypeCode1
	ChargeTypeCode12 = rateplans201810.ChargeTypeCode12
	ChargeTypeCode18 = rateplans201810.ChargeTypeCode18
	ChargeTypeCode19 = rateplans201810.ChargeTypeCode19
	ChargeTypeCode20 = rateplans201810.ChargeTypeCode20
	ChargeTypeCode21 = rateplans201810.ChargeTypeCode21
	ChargeTypeCode24 = rateplans201810.ChargeTypeCode24

	PrerequisiteInvTypeAlpineBitsDOW = rateplans201810.PrerequisiteInvTypeAlpineBitsDOW
	PrerequisiteInvTypeRoomType      = rateplans201810.PrerequisiteInvTypeRoomType

	AgeQualifyingCode8  = rateplans201810.AgeQualifyingCode8
	AgeQualifyingCode10 = rateplans201810.AgeQualifyingCode10

	DescriptionNameTitle       = rateplans201810.DescriptionNameTitle
	DescriptionNameIntro       = rateplans201810.DescriptionNameIntro
	DescriptionNameDescription = rateplans201810.DescriptionNameDescription
	DescriptionNameGallery     = rateplans201810.DescriptionNameGallery
	DescriptionNameCodelist    = rateplans201810.DescriptionNameCodelist
)
