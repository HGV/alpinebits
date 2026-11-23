package freerooms

import "github.com/HGV/alpinebits"

// A server must accept rooms, categories or both (AlpineBits 2018-10 §4.1).
const (
	CapAcceptRooms            alpinebits.Capability = "OTA_HotelAvailNotif_accept_rooms"
	CapAcceptCategories       alpinebits.Capability = "OTA_HotelAvailNotif_accept_categories"
	CapAcceptDeltas           alpinebits.Capability = "OTA_HotelAvailNotif_accept_deltas"
	CapAcceptBookingThreshold alpinebits.Capability = "OTA_HotelAvailNotif_accept_BookingThreshold"
)
