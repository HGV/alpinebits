package v202610

import (
	"context"
	"testing"

	"github.com/HGV/alpinebits"
	v201810 "github.com/HGV/alpinebits/v201810"
	v202010 "github.com/HGV/alpinebits/v202010"
	v202210 "github.com/HGV/alpinebits/v202210"
	v202410 "github.com/HGV/alpinebits/v202410"
	"github.com/HGV/alpinebits/v202610/rateplans"
)

// Stops compiling when a version replaces an alias with its own type.
func TestHandlersServeEarlierVersions(t *testing.T) {
	rp := func(context.Context, alpinebits.Request[rateplans.HotelRatePlanNotifRQ]) (rateplans.HotelRatePlanNotifRS, error) {
		return rateplans.HotelRatePlanNotifRS{}, nil
	}
	_ = alpinebits.NewRouter().
		Handle(v201810.ActionHotelRatePlanNotif, rp).
		Handle(v202010.ActionHotelRatePlanNotif, rp).
		Handle(v202210.ActionHotelRatePlanNotif, rp).
		Handle(v202410.ActionHotelRatePlanNotif, rp).
		Handle(ActionHotelRatePlanNotif, rp)
}
