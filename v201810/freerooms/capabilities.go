package freerooms

import "github.com/HGV/alpinebits"

// A server must accept rooms, categories or both (AlpineBits 2018-10 §4.1).
const (
	CapabilityAcceptRooms            alpinebits.Capability = "OTA_HotelAvailNotif_accept_rooms"
	CapabilityAcceptCategories       alpinebits.Capability = "OTA_HotelAvailNotif_accept_categories"
	CapabilityAcceptDeltas           alpinebits.Capability = "OTA_HotelAvailNotif_accept_deltas"
	CapabilityAcceptBookingThreshold alpinebits.Capability = "OTA_HotelAvailNotif_accept_BookingThreshold"
)
