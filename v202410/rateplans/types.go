// Package rateplans holds the AlpineBits 2024-10 RatePlans messages
// (OTA_HotelRatePlanNotif, §4.4). Their structure is unchanged from 2022-10, so
// the types are aliases; RatePlanNotifType Full replaces New.
package rateplans

import rateplans202210 "github.com/HGV/alpinebits/v202210/rateplans"

type (
	HotelRatePlanNotifRQ  = rateplans202210.HotelRatePlanNotifRQ
	HotelRatePlanNotifRS  = rateplans202210.HotelRatePlanNotifRS
	UniqueIDType          = rateplans202210.UniqueIDType
	Instance              = rateplans202210.Instance
	UniqueID              = rateplans202210.UniqueID
	RatePlans             = rateplans202210.RatePlans
	RatePlanNotifType     = rateplans202210.RatePlanNotifType
	RatePlanType          = rateplans202210.RatePlanType
	RatePlan              = rateplans202210.RatePlan
	CodeContext           = rateplans202210.CodeContext
	BookingRule           = rateplans202210.BookingRule
	TimeUnit              = rateplans202210.TimeUnit
	MinMaxMessageType     = rateplans202210.MinMaxMessageType
	LengthOfStay          = rateplans202210.LengthOfStay
	DOWRestrictions       = rateplans202210.DOWRestrictions
	DaysOfWeek            = rateplans202210.DaysOfWeek
	Status                = rateplans202210.Status
	Restriction           = rateplans202210.Restriction
	RestrictionStatus     = rateplans202210.RestrictionStatus
	Rate                  = rateplans202210.Rate
	BaseByGuestAmtType    = rateplans202210.BaseByGuestAmtType
	BaseByGuestAmt        = rateplans202210.BaseByGuestAmt
	AdditionalGuestAmount = rateplans202210.AdditionalGuestAmount
	MealsIncluded         = rateplans202210.MealsIncluded
	InvType               = rateplans202210.InvType
	ChargeTypeCode        = rateplans202210.ChargeTypeCode
	Supplement            = rateplans202210.Supplement
	PrerequisiteInvType   = rateplans202210.PrerequisiteInvType
	PrerequisiteInventory = rateplans202210.PrerequisiteInventory
	Offer                 = rateplans202210.Offer
	OfferRules            = rateplans202210.OfferRules
	OfferRule             = rateplans202210.OfferRule
	AgeQualifyingCode     = rateplans202210.AgeQualifyingCode
	Occupancy             = rateplans202210.Occupancy
	Guest                 = rateplans202210.Guest
	Discount              = rateplans202210.Discount
	DescriptionName       = rateplans202210.DescriptionName
	Description           = rateplans202210.Description
)

const (
	UniqueIDTypeReference = rateplans202210.UniqueIDTypeReference
	InstanceCompleteSet   = rateplans202210.InstanceCompleteSet

	RatePlanNotifTypeOverlay = rateplans202210.RatePlanNotifTypeOverlay
	RatePlanNotifTypeNew     = rateplans202210.RatePlanNotifTypeNew
	RatePlanNotifTypeRemove  = rateplans202210.RatePlanNotifTypeRemove

	RatePlanType12      = rateplans202210.RatePlanType12
	CodeContextRoomType = rateplans202210.CodeContextRoomType
	TimeUnitDay         = rateplans202210.TimeUnitDay

	MinMaxMessageTypeSetMinLOS         = rateplans202210.MinMaxMessageTypeSetMinLOS
	MinMaxMessageTypeSetForwardMinStay = rateplans202210.MinMaxMessageTypeSetForwardMinStay
	MinMaxMessageTypeSetMaxLOS         = rateplans202210.MinMaxMessageTypeSetMaxLOS
	MinMaxMessageTypeSetForwardMaxStay = rateplans202210.MinMaxMessageTypeSetForwardMaxStay

	StatusOpen        = rateplans202210.StatusOpen
	StatusClose       = rateplans202210.StatusClose
	RestrictionMaster = rateplans202210.RestrictionMaster

	BaseByGuestAmtType7  = rateplans202210.BaseByGuestAmtType7
	BaseByGuestAmtType25 = rateplans202210.BaseByGuestAmtType25

	InvTypeExtra           = rateplans202210.InvTypeExtra
	InvTypeAlpineBitsExtra = rateplans202210.InvTypeAlpineBitsExtra

	ChargeTypeCode1  = rateplans202210.ChargeTypeCode1
	ChargeTypeCode12 = rateplans202210.ChargeTypeCode12
	ChargeTypeCode18 = rateplans202210.ChargeTypeCode18
	ChargeTypeCode19 = rateplans202210.ChargeTypeCode19
	ChargeTypeCode20 = rateplans202210.ChargeTypeCode20
	ChargeTypeCode21 = rateplans202210.ChargeTypeCode21
	ChargeTypeCode24 = rateplans202210.ChargeTypeCode24

	PrerequisiteInvTypeAlpineBitsDOW = rateplans202210.PrerequisiteInvTypeAlpineBitsDOW
	PrerequisiteInvTypeRoomType      = rateplans202210.PrerequisiteInvTypeRoomType

	AgeQualifyingCode8  = rateplans202210.AgeQualifyingCode8
	AgeQualifyingCode10 = rateplans202210.AgeQualifyingCode10

	DescriptionNameTitle       = rateplans202210.DescriptionNameTitle
	DescriptionNameIntro       = rateplans202210.DescriptionNameIntro
	DescriptionNameDescription = rateplans202210.DescriptionNameDescription
	DescriptionNameGallery     = rateplans202210.DescriptionNameGallery
	DescriptionNameCodelist    = rateplans202210.DescriptionNameCodelist

	// RatePlanNotifTypeFull sends a whole rate plan, in place of New.
	RatePlanNotifTypeFull RatePlanNotifType = "Full"
)
