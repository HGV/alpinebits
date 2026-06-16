package v202010

import (
	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v_2020_10/freerooms"
)

var ActionHotelInvCountNotif = alpinebits.NewAction[freerooms.HotelInvCountNotifRQ, freerooms.HotelInvCountNotifRS](
	"OTA_HotelInvCountNotif:FreeRooms",
	"action_OTA_HotelInvCountNotif",
)
