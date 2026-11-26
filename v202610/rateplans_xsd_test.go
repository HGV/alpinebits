package v202610

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
	"github.com/HGV/alpinebits/v202610/rateplans"
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
	tests := []struct {
		name string
		want []string
	}{
		{"", nil},
		{"-overlay", nil},
		{"-pets", nil},
		{"-remove", nil},
		// The sample sends whole rate plans in a CompleteSet, whose RatePlans
		// §4.4.3 allows only a RatePlanCode.
		{"-ratePlanJoin", []string{
			"RatePlan[0]: RatePlan of a CompleteSet must have no RatePlanNotifType and no content",
			"RatePlan[1]: RatePlan of a CompleteSet must have no RatePlanNotifType and no content",
		}},
		// The sample gives a date-dependent Supplement a day-of-week
		// prerequisite and sends two New rate plans in one message, which
		// §4.4.1 and §4.4.3 forbid.
		{"-special-offer", []string{
			"RatePlan[1].Supplement[1].PrerequisiteInventory: InvType must be ROOMTYPE in a date depending Supplement",
			"RatePlans with RatePlanNotifType New must be sent one per message",
		}},
	}
	for _, tt := range tests {
		t.Run("RQ"+tt.name, func(t *testing.T) {
			path := "rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRQ" + tt.name + ".xml"
			roundTripSample[rateplans.HotelRatePlanNotifRQ](t, path)
			b, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var rq rateplans.HotelRatePlanNotifRQ
			if err := xml.Unmarshal(b, &rq); err != nil {
				t.Fatal(err)
			}
			rulestest.Expect(t, rateplans.Validate(rq, ratePlansCaps), tt.want)
		})
	}
	for _, name := range []string{"success", "advisory", "warning", "error"} {
		t.Run("RS "+name, func(t *testing.T) {
			roundTripSample[rateplans.HotelRatePlanNotifRS](t, "rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRS-"+name+".xml")
		})
	}
}

func FuzzRatePlans(f *testing.F) {
	paths, err := filepath.Glob("rateplans/testdata/*RQ*.xml")
	if err != nil || len(paths) == 0 {
		f.Fatal("no samples", err)
	}
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, doc []byte) {
		var root struct{ XMLName xml.Name }
		if Version.Validate(doc) != nil || xml.Unmarshal(doc, &root) != nil || root.XMLName.Local != "OTA_HotelRatePlanNotifRQ" {
			return
		}
		var rq rateplans.HotelRatePlanNotifRQ
		if err := xml.Unmarshal(doc, &rq); err != nil {
			t.Fatalf("valid document does not unmarshal: %v", err)
		}
		rateplans.Validate(rq, nil)
		rateplans.Validate(rq, ratePlansCaps)
		rateplans.CheckInventory(rq, nil)
		roundTrip[rateplans.HotelRatePlanNotifRQ](t, doc)
	})
}
