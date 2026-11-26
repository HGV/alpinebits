package inventory

import "github.com/HGV/alpinebits"

// Capabilities of the OTA_HotelDescriptiveContentNotif:Inventory action.
const (
	// CapabilityUseRooms: the server needs the rooms of each category, so clients
	// must send them all. Without it, Validate ignores the rooms.
	CapabilityUseRooms alpinebits.Capability = "OTA_HotelDescriptiveContentNotif_Inventory_use_rooms"
	// CapabilityOccupancyChildren allows MaxChildOccupancy.
	CapabilityOccupancyChildren alpinebits.Capability = "OTA_HotelDescriptiveContentNotif_Inventory_occupancy_children"
)
