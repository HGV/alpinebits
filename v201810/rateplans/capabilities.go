package rateplans

import "github.com/HGV/alpinebits"

// Capabilities of the OTA_HotelRatePlanNotif:RatePlans action.
const (
	CapabilityAcceptArrivalDOW                  alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_ArrivalDOW"
	CapabilityAcceptDepartureDOW                alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_DepartureDOW"
	CapabilityAcceptRatePlanBookingRule         alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_RatePlan_BookingRule"
	CapabilityAcceptRatePlanRoomTypeBookingRule alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_RatePlan_RoomType_BookingRule"
	CapabilityAcceptRatePlanMixedBookingRule    alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_RatePlan_mixed_BookingRule"
	CapabilityAcceptSupplements                 alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_Supplements"
	CapabilityAcceptFreeNightsOffers            alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_FreeNightsOffers"
	CapabilityAcceptFamilyOffers                alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_FamilyOffers"
	CapabilityAcceptOverlay                     alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_overlay"
	CapabilityAcceptRatePlanJoin                alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_RatePlanJoin"
	CapabilityAcceptOfferRuleBookingOffset      alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_OfferRule_BookingOffset"
	CapabilityAcceptOfferRuleDOWLOS             alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_OfferRule_DOWLOS"
)
