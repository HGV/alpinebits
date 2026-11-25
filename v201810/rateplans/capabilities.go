package rateplans

import "github.com/HGV/alpinebits"

// Capabilities of the OTA_HotelRatePlanNotif:RatePlans action.
const (
	CapAcceptArrivalDOW                  alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_ArrivalDOW"
	CapAcceptDepartureDOW                alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_DepartureDOW"
	CapAcceptRatePlanBookingRule         alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_RatePlan_BookingRule"
	CapAcceptRatePlanRoomTypeBookingRule alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_RatePlan_RoomType_BookingRule"
	CapAcceptRatePlanMixedBookingRule    alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_RatePlan_mixed_BookingRule"
	CapAcceptSupplements                 alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_Supplements"
	CapAcceptFreeNightsOffers            alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_FreeNightsOffers"
	CapAcceptFamilyOffers                alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_FamilyOffers"
	CapAcceptOverlay                     alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_overlay"
	CapAcceptRatePlanJoin                alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_RatePlanJoin"
	CapAcceptOfferRuleBookingOffset      alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_OfferRule_BookingOffset"
	CapAcceptOfferRuleDOWLOS             alpinebits.Capability = "OTA_HotelRatePlanNotif_accept_OfferRule_DOWLOS"
)
