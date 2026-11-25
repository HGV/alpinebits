package v202410

import (
	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202410/freerooms"
)

var ActionHotelInvCountNotif = alpinebits.NewAction[freerooms.HotelInvCountNotifRQ, freerooms.HotelInvCountNotifRS](
	Version,
	"OTA_HotelInvCountNotif:FreeRooms",
	"action_OTA_HotelInvCountNotif",
	freerooms.Validate,
)
