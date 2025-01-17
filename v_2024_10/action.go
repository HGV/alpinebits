package v_2024_10

import (
	"encoding/xml"
	"fmt"

	"github.com/HGV/alpinebits/version"
)

type Action string

var _ version.Action = new(Action)

const (
	ActionPing                                  Action = "OTA_Ping:Handshaking"
	ActionHotelInvCountNotif                    Action = "OTA_HotelInvCountNotif:FreeRooms"
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
		v = new(PingRQ)
	default:
		return nil, fmt.Errorf("unhandled action: %s", a)
	}

	if err := xml.Unmarshal(b, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func (a Action) String() string {
	return string(a)
}
