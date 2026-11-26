package v202410

import (
	"context"
	"testing"

	"github.com/HGV/alpinebits"
	v201810 "github.com/HGV/alpinebits/v201810"
	v202010 "github.com/HGV/alpinebits/v202010"
	v202210 "github.com/HGV/alpinebits/v202210"
	"github.com/HGV/alpinebits/v202410/guestrequests"
	"github.com/HGV/alpinebits/v202410/inventory"
)

// Stops compiling when a version replaces an alias with its own type.
func TestHandlersServeEarlierVersions(t *testing.T) {
	inv := func(context.Context, alpinebits.Request[inventory.HotelDescriptiveContentNotifRQ]) (inventory.HotelDescriptiveContentNotifRS, error) {
		return inventory.HotelDescriptiveContentNotifRS{}, nil
	}
	ack := func(context.Context, alpinebits.Request[guestrequests.NotifReportRQ]) (guestrequests.NotifReportRS, error) {
		return guestrequests.NotifReportRS{}, nil
	}
	_ = alpinebits.NewRouter().
		Handle(v201810.ActionHotelDescriptiveContentNotif, inv).
		Handle(v202010.ActionHotelDescriptiveContentNotif, inv).
		Handle(v202210.ActionHotelDescriptiveContentNotif, inv).
		Handle(ActionHotelDescriptiveContentNotif, inv).
		Handle(v201810.ActionNotifReport, ack).
		Handle(v202010.ActionNotifReport, ack).
		Handle(v202210.ActionNotifReport, ack).
		Handle(ActionNotifReport, ack)
}
