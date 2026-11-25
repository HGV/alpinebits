package v202210

import (
	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202210/freerooms"
	"github.com/HGV/alpinebits/v202210/inventory"
)

var ActionHotelInvCountNotif = alpinebits.NewAction[freerooms.HotelInvCountNotifRQ, freerooms.HotelInvCountNotifRS](
	Version,
	"OTA_HotelInvCountNotif:FreeRooms",
	"action_OTA_HotelInvCountNotif",
	freerooms.Validate,
)

var ActionHotelDescriptiveContentNotif = alpinebits.NewAction[inventory.HotelDescriptiveContentNotifRQ, inventory.HotelDescriptiveContentNotifRS](
	Version,
	"OTA_HotelDescriptiveContentNotif:Inventory",
	"action_OTA_HotelDescriptiveContentNotif_Inventory",
	inventory.Validate,
)
