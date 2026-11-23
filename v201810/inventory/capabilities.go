package inventory

import "github.com/HGV/alpinebits"

// Capabilities of the OTA_HotelDescriptiveContentNotif:Inventory action.
const (
	// CapUseRooms: the server needs the rooms of each category, so clients
	// must send them all. Without it, Validate ignores the rooms.
	CapUseRooms alpinebits.Capability = "OTA_HotelDescriptiveContentNotif_Inventory_use_rooms"
	// CapOccupancyChildren allows MaxChildOccupancy.
	CapOccupancyChildren alpinebits.Capability = "OTA_HotelDescriptiveContentNotif_Inventory_occupancy_children"
)
