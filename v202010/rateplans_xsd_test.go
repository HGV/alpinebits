package v202010

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202010/rateplans"
)

var ratePlansCaps = []alpinebits.Capability{
	rateplans.CapAcceptArrivalDOW, rateplans.CapAcceptDepartureDOW,
	rateplans.CapAcceptRatePlanBookingRule, rateplans.CapAcceptRatePlanRoomTypeBookingRule,
	rateplans.CapAcceptRatePlanMixedBookingRule, rateplans.CapAcceptSupplements,
	rateplans.CapAcceptFreeNightsOffers, rateplans.CapAcceptFamilyOffers,
	rateplans.CapAcceptOverlay, rateplans.CapAcceptRatePlanJoin,
	rateplans.CapAcceptOfferRuleBookingOffset, rateplans.CapAcceptOfferRuleDOWLOS,
}

func readSample[T any](t *testing.T, path string) T {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var v T
	if err := xml.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestRatePlanTypesRoundTrip(t *testing.T) {
	// 2020-10 kept the 2018-10 sample, which is in EUR.
	rq := readSample[rateplans.HotelRatePlanNotifRQ](t, "../v201810/rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRQ.xml")
	p := &rq.RatePlans.RatePlans[0]
	p.CurrencyCode = "CHF"
	(*(*p.Rates)[1].BaseByGuestAmts)[0].CurrencyCode = "CHF"
	if errs := rateplans.Validate(rq, ratePlansCaps); errs != nil {
		t.Fatalf("rules: %v", errs)
	}
	roundTrip(t, rq)

	for _, name := range []string{"success", "advisory", "warning", "error"} {
		t.Run("RS "+name, func(t *testing.T) {
			roundTrip(t, readSample[rateplans.HotelRatePlanNotifRS](t, "rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRS-"+name+".xml"))
		})
	}
}
