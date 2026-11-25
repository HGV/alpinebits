package v202410

import (
	"context"
	"testing"

	"github.com/HGV/alpinebits"
	v201810 "github.com/HGV/alpinebits/v201810"
	v202010 "github.com/HGV/alpinebits/v202010"
	v202210 "github.com/HGV/alpinebits/v202210"
	"github.com/HGV/alpinebits/v202410/freerooms"
	"github.com/HGV/alpinebits/v202410/inventory"
)

// Stops compiling when a version replaces an alias with its own type.
func TestHandlersServeEarlierVersions(t *testing.T) {
	fr := func(context.Context, alpinebits.Request[freerooms.HotelInvCountNotifRQ]) (freerooms.HotelInvCountNotifRS, error) {
		return freerooms.HotelInvCountNotifRS{}, nil
	}
	inv := func(context.Context, alpinebits.Request[inventory.HotelDescriptiveContentNotifRQ]) (inventory.HotelDescriptiveContentNotifRS, error) {
		return inventory.HotelDescriptiveContentNotifRS{}, nil
	}
	_ = alpinebits.NewRouter().
		Handle(v202010.ActionHotelInvCountNotif, fr).
		Handle(v202210.ActionHotelInvCountNotif, fr).
		Handle(ActionHotelInvCountNotif, fr).
		Handle(v201810.ActionHotelDescriptiveContentNotif, inv).
		Handle(v202010.ActionHotelDescriptiveContentNotif, inv).
		Handle(v202210.ActionHotelDescriptiveContentNotif, inv).
		Handle(ActionHotelDescriptiveContentNotif, inv)
}
