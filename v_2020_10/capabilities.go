package v202010

import "github.com/HGV/alpinebits"

// Capabilities for AlpineBits 2020-10.
const (
	CapFreeRoomsRooms          alpinebits.Capability = "action_OTA_HotelInvCountNotif_accept_rooms"
	CapFreeRoomsDeltas         alpinebits.Capability = "action_OTA_HotelInvCountNotif_accept_deltas"
	CapFreeRoomsOutOfOrder     alpinebits.Capability = "action_OTA_HotelInvCountNotif_accept_out_of_order"
	CapFreeRoomsOutOfMarket    alpinebits.Capability = "action_OTA_HotelInvCountNotif_accept_out_of_market"
	CapFreeRoomsClosingSeasons alpinebits.Capability = "action_OTA_HotelInvCountNotif_accept_closing_seasons"
)
