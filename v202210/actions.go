package v202210

import (
	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202210/freerooms"
)

var ActionHotelInvCountNotif = alpinebits.NewAction[freerooms.HotelInvCountNotifRQ, freerooms.HotelInvCountNotifRS](
	Version,
	"OTA_HotelInvCountNotif:FreeRooms",
	"action_OTA_HotelInvCountNotif",
	nil, // no rules yet
)
