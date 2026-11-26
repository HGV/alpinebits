package freerooms

import "github.com/HGV/alpinebits"

// A server must accept rooms, categories or both (AlpineBits 2020-10 §4.1).
const (
	CapabilityAcceptRooms          alpinebits.Capability = "OTA_HotelInvCountNotif_accept_rooms"
	CapabilityAcceptCategories     alpinebits.Capability = "OTA_HotelInvCountNotif_accept_categories"
	CapabilityAcceptDeltas         alpinebits.Capability = "OTA_HotelInvCountNotif_accept_deltas"
	CapabilityAcceptOutOfOrder     alpinebits.Capability = "OTA_HotelInvCountNotif_accept_out_of_order"
	CapabilityAcceptOutOfMarket    alpinebits.Capability = "OTA_HotelInvCountNotif_accept_out_of_market"
	CapabilityAcceptClosingSeasons alpinebits.Capability = "OTA_HotelInvCountNotif_accept_closing_seasons"
)
