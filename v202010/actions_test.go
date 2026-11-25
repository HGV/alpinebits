package v202010

import (
	"context"
	"testing"

	"github.com/HGV/alpinebits"
	v201810 "github.com/HGV/alpinebits/v201810"
	"github.com/HGV/alpinebits/v202010/guestrequests"
	"github.com/HGV/alpinebits/v202010/inventory"
	"github.com/HGV/alpinebits/v202010/rateplans"
)

// Stops compiling when 2020-10 replaces an alias with its own type.
func TestHandlersServeVersionsSince201810(t *testing.T) {
	inv := func(context.Context, alpinebits.Request[inventory.HotelDescriptiveContentNotifRQ]) (inventory.HotelDescriptiveContentNotifRS, error) {
		return inventory.HotelDescriptiveContentNotifRS{}, nil
	}
	rp := func(context.Context, alpinebits.Request[rateplans.HotelRatePlanNotifRQ]) (rateplans.HotelRatePlanNotifRS, error) {
		return rateplans.HotelRatePlanNotifRS{}, nil
	}
	ack := func(context.Context, alpinebits.Request[guestrequests.NotifReportRQ]) (guestrequests.NotifReportRS, error) {
		return guestrequests.NotifReportRS{}, nil
	}
	_ = alpinebits.NewRouter().
		Handle(v201810.ActionHotelDescriptiveContentNotif, inv).
		Handle(ActionHotelDescriptiveContentNotif, inv).
		Handle(v201810.ActionHotelRatePlanNotif, rp).
		Handle(ActionHotelRatePlanNotif, rp).
		Handle(v201810.ActionNotifReport, ack).
		Handle(ActionNotifReport, ack)
}
