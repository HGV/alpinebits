package freerooms

import "github.com/HGV/alpinebits"

// A server must accept rooms, categories or both (AlpineBits 2020-10 §4.1).
const (
	CapAcceptRooms          alpinebits.Capability = "OTA_HotelInvCountNotif_accept_rooms"
	CapAcceptCategories     alpinebits.Capability = "OTA_HotelInvCountNotif_accept_categories"
	CapAcceptDeltas         alpinebits.Capability = "OTA_HotelInvCountNotif_accept_deltas"
	CapAcceptOutOfOrder     alpinebits.Capability = "OTA_HotelInvCountNotif_accept_out_of_order"
	CapAcceptOutOfMarket    alpinebits.Capability = "OTA_HotelInvCountNotif_accept_out_of_market"
	CapAcceptClosingSeasons alpinebits.Capability = "OTA_HotelInvCountNotif_accept_closing_seasons"
)
