package v202410

import (
	"context"
	"testing"

	"github.com/HGV/alpinebits"
	v202010 "github.com/HGV/alpinebits/v202010"
	v202210 "github.com/HGV/alpinebits/v202210"
	"github.com/HGV/alpinebits/v202410/freerooms"
)

// Stops compiling when a version replaces an alias with its own type.
func TestFreeRoomsHandlerServesVersionsSince202010(t *testing.T) {
	h := func(context.Context, alpinebits.Request[freerooms.HotelInvCountNotifRQ]) (freerooms.HotelInvCountNotifRS, error) {
		return freerooms.HotelInvCountNotifRS{}, nil
	}
	_ = alpinebits.NewRouter().
		Handle(v202010.ActionHotelInvCountNotif, h).
		Handle(v202210.ActionHotelInvCountNotif, h).
		Handle(ActionHotelInvCountNotif, h)
}
