package v_2018_10

type Capability string

const (
	CapabilityHotelAvailNotifAcceptRooms                             Capability = "OTA_HotelAvailNotif_accept_rooms"
	CapabilityHotelAvailNotifAcceptCategories                        Capability = "OTA_HotelAvailNotif_accept_categories"
	CapabilityHotelAvailNotifAcceptDeltas                            Capability = "OTA_HotelAvailNotif_accept_deltas"
	CapabilityHotelAvailNotifAcceptBookingThreshold                  Capability = "OTA_HotelAvailNotif_accept_BookingThreshold"
	CapabilityHotelDescriptiveContentNotifInventoryUseRooms          Capability = "OTA_HotelDescriptiveContentNotif_Inventory_use_rooms"
	CapabilityHotelDescriptiveContentNotifInventoryOccupancyChildren Capability = "OTA_HotelDescriptiveContentNotif_Inventory_occupancy_children"
	CapabilityHotelRatePlanNotifAcceptArrivalDOW                     Capability = "OTA_HotelRatePlanNotif_accept_ArrivalDOW"
	CapabilityHotelRatePlanNotifAcceptDepartureDOW                   Capability = "OTA_HotelRatePlanNotif_accept_DepartureDOW"
	CapabilityHotelRatePlanNotifAcceptRatePlanBookingRule            Capability = "OTA_HotelRatePlanNotif_accept_RatePlan_BookingRule"
	CapabilityHotelRatePlanNotifAcceptRatePlanRoomTypeBookingRule    Capability = "OTA_HotelRatePlanNotif_accept_RatePlan_RoomType_BookingRule"
	CapabilityHotelRatePlanNotifAcceptRatePlanMixedBookingRule       Capability = "OTA_HotelRatePlanNotif_accept_RatePlan_mixed_BookingRule"
	CapabilityHotelRatePlanNotifAcceptSupplements                    Capability = "OTA_HotelRatePlanNotif_accept_Supplements"
	CapabilityHotelRatePlanNotifAcceptFreeNightsOffers               Capability = "OTA_HotelRatePlanNotif_accept_FreeNightsOffers"
	CapabilityHotelRatePlanNotifAcceptFamilyOffers                   Capability = "OTA_HotelRatePlanNotif_accept_FamilyOffers"
	CapabilityHotelRatePlanNotifAcceptOverlay                        Capability = "OTA_HotelRatePlanNotif_accept_overlay"
	CapabilityHotelRatePlanNotifAcceptRatePlanJoin                   Capability = "OTA_HotelRatePlanNotif_accept_RatePlanJoin"
	CapabilityHotelRatePlanNotifAcceptOfferRuleBookingOffset         Capability = "OTA_HotelRatePlanNotif_accept_OfferRule_BookingOffset"
	CapabilityHotelRatePlanNotifAcceptOfferRuleDOWLOS                Capability = "OTA_HotelRatePlanNotif_accept_OfferRule_DOWLOS"
)
