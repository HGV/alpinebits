package v202410

import (
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202410/rateplans"
)

var ratePlansCaps = []alpinebits.Capability{
	rateplans.CapAcceptArrivalDOW, rateplans.CapAcceptDepartureDOW,
	rateplans.CapAcceptRatePlanBookingRule, rateplans.CapAcceptRatePlanRoomTypeBookingRule,
	rateplans.CapAcceptRatePlanMixedBookingRule, rateplans.CapAcceptSupplements,
	rateplans.CapAcceptFreeNightsOffers, rateplans.CapAcceptFamilyOffers,
	rateplans.CapAcceptOverlay, rateplans.CapAcceptRatePlanJoin,
	rateplans.CapAcceptOfferRuleBookingOffset, rateplans.CapAcceptOfferRuleDOWLOS,
	rateplans.CapAcceptFull,
}

func TestRatePlanSamples(t *testing.T) {
	rq := readSample[rateplans.HotelRatePlanNotifRQ](t, "rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRQ.xml")
	if rq.RatePlans.RatePlans[0].RatePlanNotifType != rateplans.RatePlanNotifTypeFull {
		t.Fatal("the 2024-10 sample sends Full")
	}
	if errs := rateplans.Validate(rq, ratePlansCaps); errs != nil {
		t.Fatalf("rules: %v", errs)
	}
	roundTrip(t, rq)

	// The responses are those of 2022-10.
	for _, path := range []string{
		"../v202010/rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRS-success.xml",
		"../v202210/rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRS-advisory.xml",
		"../v202210/rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRS-warning.xml",
		"../v202010/rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRS-error.xml",
	} {
		t.Run(path, func(t *testing.T) {
			roundTrip(t, readSample[rateplans.HotelRatePlanNotifRS](t, path))
		})
	}
}
