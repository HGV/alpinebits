package rateplans

import (
	"github.com/HGV/alpinebits"
	rateplans202210 "github.com/HGV/alpinebits/v202210/rateplans"
)

// Capabilities of the OTA_HotelRatePlanNotif:RatePlans action.
const (
	CapabilityAcceptArrivalDOW                  = rateplans202210.CapabilityAcceptArrivalDOW
	CapabilityAcceptDepartureDOW                = rateplans202210.CapabilityAcceptDepartureDOW
	CapabilityAcceptRatePlanBookingRule         = rateplans202210.CapabilityAcceptRatePlanBookingRule
	CapabilityAcceptRatePlanRoomTypeBookingRule = rateplans202210.CapabilityAcceptRatePlanRoomTypeBookingRule
	CapabilityAcceptRatePlanMixedBookingRule    = rateplans202210.CapabilityAcceptRatePlanMixedBookingRule
	CapabilityAcceptSupplements                 = rateplans202210.CapabilityAcceptSupplements
	CapabilityAcceptFreeNightsOffers            = rateplans202210.CapabilityAcceptFreeNightsOffers
	CapabilityAcceptFamilyOffers                = rateplans202210.CapabilityAcceptFamilyOffers
	CapabilityAcceptOverlay                     = rateplans202210.CapabilityAcceptOverlay
	CapabilityAcceptRatePlanJoin                = rateplans202210.CapabilityAcceptRatePlanJoin
	CapabilityAcceptOfferRuleBookingOffset      = rateplans202210.CapabilityAcceptOfferRuleBookingOffset
	CapabilityAcceptOfferRuleDOWLOS             = rateplans202210.CapabilityAcceptOfferRuleDOWLOS

	// CapabilityAcceptFull: the server accepts RatePlanNotifType Full.
	CapabilityAcceptFull alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_full"
)
