package v202410

import (
	"context"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202010"
	"github.com/HGV/alpinebits/v202210"
	"github.com/HGV/alpinebits/v202410/freerooms"
)

// Stops compiling when a version replaces an alias with its own type.
func TestFreeRoomsHandlerServesVersionsSince202010(t *testing.T) {
	h := func(context.Context, alpinebits.Request[freerooms.HotelInvCountNotifRQ]) (freerooms.HotelInvCountNotifRS, error) {
		return freerooms.HotelInvCountNotifRS{}, nil
	}
	r := alpinebits.NewRouter()
	r.Version(v202010.Version, alpinebits.Handle(v202010.ActionHotelInvCountNotif, h))
	r.Version(v202210.Version, alpinebits.Handle(v202210.ActionHotelInvCountNotif, h))
	r.Version(Version, alpinebits.Handle(ActionHotelInvCountNotif, h))
}
