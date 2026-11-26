package rateplans

import (
	"github.com/HGV/alpinebits"
	rateplans202210 "github.com/HGV/alpinebits/v202210/rateplans"
)

// Capabilities of the OTA_HotelRatePlanNotif:RatePlans action.
const (
	CapAcceptArrivalDOW                  = rateplans202210.CapAcceptArrivalDOW
	CapAcceptDepartureDOW                = rateplans202210.CapAcceptDepartureDOW
	CapAcceptRatePlanBookingRule         = rateplans202210.CapAcceptRatePlanBookingRule
	CapAcceptRatePlanRoomTypeBookingRule = rateplans202210.CapAcceptRatePlanRoomTypeBookingRule
	CapAcceptRatePlanMixedBookingRule    = rateplans202210.CapAcceptRatePlanMixedBookingRule
	CapAcceptSupplements                 = rateplans202210.CapAcceptSupplements
	CapAcceptFreeNightsOffers            = rateplans202210.CapAcceptFreeNightsOffers
	CapAcceptFamilyOffers                = rateplans202210.CapAcceptFamilyOffers
	CapAcceptOverlay                     = rateplans202210.CapAcceptOverlay
	CapAcceptRatePlanJoin                = rateplans202210.CapAcceptRatePlanJoin
	CapAcceptOfferRuleBookingOffset      = rateplans202210.CapAcceptOfferRuleBookingOffset
	CapAcceptOfferRuleDOWLOS             = rateplans202210.CapAcceptOfferRuleDOWLOS

	// CapAcceptFull: the server accepts RatePlanNotifType Full.
	CapAcceptFull alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_full"
)
