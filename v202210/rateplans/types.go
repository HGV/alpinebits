// Package rateplans holds the AlpineBits 2022-10 RatePlans messages
// (OTA_HotelRatePlanNotif, §4.5). Their structure is unchanged from 2020-10, so
// the types are aliases.
package rateplans

import rateplans202010 "github.com/HGV/alpinebits/v202010/rateplans"

type (
	HotelRatePlanNotifRQ  = rateplans202010.HotelRatePlanNotifRQ
	HotelRatePlanNotifRS  = rateplans202010.HotelRatePlanNotifRS
	UniqueIDType          = rateplans202010.UniqueIDType
	Instance              = rateplans202010.Instance
	UniqueID              = rateplans202010.UniqueID
	RatePlans             = rateplans202010.RatePlans
	RatePlanNotifType     = rateplans202010.RatePlanNotifType
	RatePlanType          = rateplans202010.RatePlanType
	RatePlan              = rateplans202010.RatePlan
	CodeContext           = rateplans202010.CodeContext
	BookingRule           = rateplans202010.BookingRule
	TimeUnit              = rateplans202010.TimeUnit
	MinMaxMessageType     = rateplans202010.MinMaxMessageType
	LengthOfStay          = rateplans202010.LengthOfStay
	DOWRestrictions       = rateplans202010.DOWRestrictions
	DaysOfWeek            = rateplans202010.DaysOfWeek
	Status                = rateplans202010.Status
	Restriction           = rateplans202010.Restriction
	RestrictionStatus     = rateplans202010.RestrictionStatus
	Rate                  = rateplans202010.Rate
	BaseByGuestAmtType    = rateplans202010.BaseByGuestAmtType
	BaseByGuestAmt        = rateplans202010.BaseByGuestAmt
	AdditionalGuestAmount = rateplans202010.AdditionalGuestAmount
	MealsIncluded         = rateplans202010.MealsIncluded
	InvType               = rateplans202010.InvType
	ChargeTypeCode        = rateplans202010.ChargeTypeCode
	Supplement            = rateplans202010.Supplement
	PrerequisiteInvType   = rateplans202010.PrerequisiteInvType
	PrerequisiteInventory = rateplans202010.PrerequisiteInventory
	Offer                 = rateplans202010.Offer
	OfferRules            = rateplans202010.OfferRules
	OfferRule             = rateplans202010.OfferRule
	AgeQualifyingCode     = rateplans202010.AgeQualifyingCode
	Occupancy             = rateplans202010.Occupancy
	Guest                 = rateplans202010.Guest
	Discount              = rateplans202010.Discount
	DescriptionName       = rateplans202010.DescriptionName
	Description           = rateplans202010.Description
)

const (
	UniqueIDTypeReference = rateplans202010.UniqueIDTypeReference
	InstanceCompleteSet   = rateplans202010.InstanceCompleteSet

	RatePlanNotifTypeOverlay = rateplans202010.RatePlanNotifTypeOverlay
	RatePlanNotifTypeNew     = rateplans202010.RatePlanNotifTypeNew
	RatePlanNotifTypeRemove  = rateplans202010.RatePlanNotifTypeRemove

	RatePlanType12      = rateplans202010.RatePlanType12
	CodeContextRoomType = rateplans202010.CodeContextRoomType
	TimeUnitDay         = rateplans202010.TimeUnitDay

	MinMaxMessageTypeSetMinLOS         = rateplans202010.MinMaxMessageTypeSetMinLOS
	MinMaxMessageTypeSetForwardMinStay = rateplans202010.MinMaxMessageTypeSetForwardMinStay
	MinMaxMessageTypeSetMaxLOS         = rateplans202010.MinMaxMessageTypeSetMaxLOS
	MinMaxMessageTypeSetForwardMaxStay = rateplans202010.MinMaxMessageTypeSetForwardMaxStay

	StatusOpen        = rateplans202010.StatusOpen
	StatusClose       = rateplans202010.StatusClose
	RestrictionMaster = rateplans202010.RestrictionMaster

	BaseByGuestAmtType7  = rateplans202010.BaseByGuestAmtType7
	BaseByGuestAmtType25 = rateplans202010.BaseByGuestAmtType25

	InvTypeExtra           = rateplans202010.InvTypeExtra
	InvTypeAlpineBitsExtra = rateplans202010.InvTypeAlpineBitsExtra

	ChargeTypeCode1  = rateplans202010.ChargeTypeCode1
	ChargeTypeCode12 = rateplans202010.ChargeTypeCode12
	ChargeTypeCode18 = rateplans202010.ChargeTypeCode18
	ChargeTypeCode19 = rateplans202010.ChargeTypeCode19
	ChargeTypeCode20 = rateplans202010.ChargeTypeCode20
	ChargeTypeCode21 = rateplans202010.ChargeTypeCode21
	ChargeTypeCode24 = rateplans202010.ChargeTypeCode24

	PrerequisiteInvTypeAlpineBitsDOW = rateplans202010.PrerequisiteInvTypeAlpineBitsDOW
	PrerequisiteInvTypeRoomType      = rateplans202010.PrerequisiteInvTypeRoomType

	AgeQualifyingCode8  = rateplans202010.AgeQualifyingCode8
	AgeQualifyingCode10 = rateplans202010.AgeQualifyingCode10

	DescriptionNameTitle       = rateplans202010.DescriptionNameTitle
	DescriptionNameIntro       = rateplans202010.DescriptionNameIntro
	DescriptionNameDescription = rateplans202010.DescriptionNameDescription
	DescriptionNameGallery     = rateplans202010.DescriptionNameGallery
	DescriptionNameCodelist    = rateplans202010.DescriptionNameCodelist
)
