package v202610

import (
	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202610/freerooms"
	"github.com/HGV/alpinebits/v202610/guestrequests"
	"github.com/HGV/alpinebits/v202610/inventory"
	"github.com/HGV/alpinebits/v202610/rateplans"
)

var ActionHotelInvCountNotif = alpinebits.NewAction[freerooms.HotelInvCountNotifRQ, freerooms.HotelInvCountNotifRS](
	Version,
	"OTA_HotelInvCountNotif:FreeRooms",
	"action_OTA_HotelInvCountNotif",
	freerooms.Validate,
)

var ActionRead = alpinebits.NewAction[guestrequests.ReadRQ, guestrequests.ResRetrieveRS](
	Version,
	"OTA_Read:GuestRequests",
	"action_OTA_Read",
	guestrequests.ValidateReadRQ,
)

// ActionNotifReport is covered by action_OTA_Read in the handshake
// (AlpineBits 2026-10 §3.3).
var ActionNotifReport = alpinebits.NewAction[guestrequests.NotifReportRQ, guestrequests.NotifReportRS](
	Version,
	"OTA_NotifReport:GuestRequests",
	"action_OTA_Read",
	guestrequests.ValidateNotifReportRQ,
)

var ActionHotelDescriptiveContentNotif = alpinebits.NewAction[inventory.HotelDescriptiveContentNotifRQ, inventory.HotelDescriptiveContentNotifRS](
	Version,
	"OTA_HotelDescriptiveContentNotif:Inventory",
	"action_OTA_HotelDescriptiveContentNotif_Inventory",
	inventory.Validate,
)

var ActionHotelRatePlanNotif = alpinebits.NewAction[rateplans.HotelRatePlanNotifRQ, rateplans.HotelRatePlanNotifRS](
	Version,
	"OTA_HotelRatePlanNotif:RatePlans",
	"action_OTA_HotelRatePlanNotif_RatePlans",
	rateplans.Validate,
)
