package v_2018_10

import (
	"encoding/xml"
	"fmt"

	"github.com/HGV/alpinebits/v_2018_10/freerooms"
	"github.com/HGV/alpinebits/v_2018_10/guestrequests"
	"github.com/HGV/alpinebits/v_2018_10/handshake"
	"github.com/HGV/alpinebits/v_2018_10/inventory"
	"github.com/HGV/alpinebits/v_2018_10/rateplans"
	"github.com/HGV/alpinebits/version"
)

type Action string

var _ version.Action = new(Action)

const (
	ActionPing                                  Action = "OTA_Ping:Handshaking"
	ActionHotelAvailNotif                       Action = "OTA_HotelAvailNotif:FreeRooms"
	ActionReadGuestRequests                     Action = "OTA_Read:GuestRequests"
	ActionNotifReportGuestRequests              Action = "OTA_NotifReport:GuestRequests"
	ActionHotelDescriptiveContentNotifInventory Action = "OTA_HotelDescriptiveContentNotif:Inventory"
	ActionHotelDescriptiveContentNotifInfo      Action = "OTA_HotelDescriptiveContentNotif:Info"
	ActionHotelRatePlanNotifRatePlans           Action = "OTA_HotelRatePlanNotif:RatePlans"
)

func (a Action) Unmarshal(b []byte) (any, error) {
	var v any

	switch a {
	case ActionPing:
		v = new(handshake.PingRQ)
	case ActionHotelAvailNotif:
		v = new(freerooms.HotelAvailNotifRQ)
	case ActionReadGuestRequests:
		v = new(guestrequests.ReadRQ)
	case ActionNotifReportGuestRequests:
		v = new(guestrequests.NotifReportRQ)
	case ActionHotelDescriptiveContentNotifInventory:
		v = new(inventory.HotelDescriptiveContentNotifRQ)
	case ActionHotelRatePlanNotifRatePlans:
		v = new(rateplans.HotelRatePlanNotifRQ)
	default:
		return nil, fmt.Errorf("unhandled action: %s", a)
	}

	if err := xml.Unmarshal(b, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func (a Action) HandshakeName() string {
	switch a {
	case ActionPing:
		return "action_OTA_Ping"
	case ActionHotelAvailNotif:
		return "action_OTA_HotelAvailNotif"
	case ActionReadGuestRequests:
		return "action_OTA_Read"
	case ActionNotifReportGuestRequests:
		return "action_OTA_HotelResNotif_GuestRequests"
	case ActionHotelDescriptiveContentNotifInventory:
		return "action_OTA_HotelDescriptiveContentNotif_Inventory"
	case ActionHotelDescriptiveContentNotifInfo:
		return "action_OTA_HotelDescriptiveContentNotif_Info"
	case ActionHotelRatePlanNotifRatePlans:
		return "action_OTA_HotelRatePlanNotif_RatePlans"
	default:
		return ""
	}
}

func (a Action) String() string {
	return string(a)
}
