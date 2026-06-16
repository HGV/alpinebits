package v201810

import (
	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v_2018_10/freerooms"
	"github.com/HGV/alpinebits/v_2018_10/guestrequests"
	"github.com/HGV/alpinebits/v_2018_10/inventory"
	"github.com/HGV/alpinebits/v_2018_10/rateplans"
)

var ActionHotelAvailNotif = alpinebits.NewAction[freerooms.HotelAvailNotifRQ, freerooms.HotelAvailNotifRS](
	"OTA_HotelAvailNotif:FreeRooms",
	"action_OTA_HotelAvailNotif",
)

var ActionRead = alpinebits.NewAction[guestrequests.ReadRQ, guestrequests.ResRetrieveRS](
	"OTA_Read:GuestRequests",
	"action_OTA_Read_GuestRequests",
)

var ActionNotifReport = alpinebits.NewAction[guestrequests.NotifReportRQ, guestrequests.NotifReportRS](
	"OTA_NotifReport:GuestRequests",
	"action_OTA_NotifReport_GuestRequests",
)

var ActionHotelDescriptiveContentNotif = alpinebits.NewAction[inventory.HotelDescriptiveContentNotifRQ, inventory.HotelDescriptiveContentNotifRS](
	"OTA_HotelDescriptiveContentNotif:Inventory",
	"action_OTA_HotelDescriptiveContentNotif_Inventory",
)

var ActionHotelRatePlanNotif = alpinebits.NewAction[rateplans.HotelRatePlanNotifRQ, rateplans.HotelRatePlanNotifRS](
	"OTA_HotelRatePlanNotif:RatePlans",
	"action_OTA_HotelRatePlanNotif_RatePlans",
)
