package v201810

import (
	"encoding/xml"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v201810/rateplans"
)

const ratePlanNotif = `<OTA_HotelRatePlanNotifRQ ` + ns + ` Version="1.000">
<RatePlans HotelCode="123">
<RatePlan RatePlanNotifType="New" CurrencyCode="EUR" RatePlanCode="BB">
<BookingRules><BookingRule Start="2026-01-01" End="2026-01-31">
<LengthsOfStay><LengthOfStay Time="2" TimeUnit="Day" MinMaxMessageType="SetMinLOS"/></LengthsOfStay>
<RestrictionStatus Status="Open"/>
</BookingRule></BookingRules>
<Rates><Rate RateTimeUnit="Day" UnitMultiplier="1">
<BaseByGuestAmts><BaseByGuestAmt Type="7" NumberOfGuests="1" AmountAfterTax="100" CurrencyCode="EUR"/></BaseByGuestAmts>
<MealsIncluded MealPlanIndicator="true" MealPlanCodes="3"/>
</Rate></Rates>
<Supplements><Supplement InvType="EXTRA" InvCode="X" ChargeTypeCode="1" Amount="10">
<PrerequisiteInventory InvCode="DZ" InvType="ROOMTYPE"/>
</Supplement></Supplements>
<Offers><Offer><OfferRules><OfferRule>
<LengthsOfStay><LengthOfStay Time="7" TimeUnit="Day" MinMaxMessageType="SetMinLOS"/></LengthsOfStay>
<Occupancy AgeQualifyingCode="10"/>
</OfferRule></OfferRules></Offer></Offers>
</RatePlan>
</RatePlans>
</OTA_HotelRatePlanNotifRQ>`

// rateplans.Validate relies on these being rejected by the XSD.
func TestXSDEnforcesRatePlanRules(t *testing.T) {
	if err := Version.Validate([]byte(ratePlanNotif)); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	offer := `<Offer><OfferRules><OfferRule>
<LengthsOfStay><LengthOfStay Time="7" TimeUnit="Day" MinMaxMessageType="SetMinLOS"/></LengthsOfStay>
<Occupancy AgeQualifyingCode="10"/>
</OfferRule></OfferRules></Offer>`
	tests := []struct{ name, old, new string }{
		{"no RatePlan", ratePlanNotif[strings.Index(ratePlanNotif, "<RatePlan "):strings.Index(ratePlanNotif, "</RatePlans>")], ""},
		{"BookingRule TimeUnit", `Time="2" TimeUnit="Day"`, `Time="2" TimeUnit="Week"`},
		{"BookingRule negative Time", `Time="2"`, `Time="-1"`},
		{"RestrictionStatus Status", `Status="Open"`, `Status="Half"`},
		{"AmountAfterTax 0", `AmountAfterTax="100"`, `AmountAfterTax="0"`},
		{"BaseByGuestAmt CurrencyCode", `AmountAfterTax="100" CurrencyCode="EUR"`, `AmountAfterTax="100" CurrencyCode="USD"`},
		{"BaseByGuestAmt Type", `Type="7"`, `Type="8"`},
		{"MealPlanCodes", `MealPlanCodes="3"`, `MealPlanCodes="2"`},
		{"Supplement InvType missing", `InvType="EXTRA" `, ``},
		{"Supplement InvType", `InvType="EXTRA"`, `InvType="FOO"`},
		{"ChargeTypeCode", `ChargeTypeCode="1"`, `ChargeTypeCode="2"`},
		{"PrerequisiteInventory InvType", `InvType="ROOMTYPE"`, `InvType="FOO"`},
		{"four Offers", offer, strings.Repeat(offer, 4)},
		{"OfferRule TimeUnit", `Time="7" TimeUnit="Day"`, `Time="7" TimeUnit="Week"`},
		{"Occupancy AgeQualifyingCode", `AgeQualifyingCode="10"`, `AgeQualifyingCode="7"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if n := strings.Count(ratePlanNotif, tt.old); n != 1 {
				t.Fatalf("%q occurs %d times", tt.old, n)
			}
			doc := strings.Replace(ratePlanNotif, tt.old, tt.new, 1)
			if err := Version.Validate([]byte(doc)); err == nil {
				t.Error("XSD accepted it")
			}
		})
	}
}

// Values whose absence means something (a free supplement, a restricted day,
// an alternative rate plan) must survive marshalling and pass the XSD.
func TestRatePlanTypesRoundTrip(t *testing.T) {
	b, err := os.ReadFile("rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRQ.xml")
	if err != nil {
		t.Fatal(err)
	}
	var rq rateplans.HotelRatePlanNotifRQ
	if err := xml.Unmarshal(b, &rq); err != nil {
		t.Fatal(err)
	}
	p := &rq.RatePlans.RatePlans[0]
	p.RatePlanID, p.RatePlanQualifier = "M", new(false)
	(*p.BookingRules)[0].DOWRestrictions.ArrivalDaysOfWeek.Mon = new(false)
	(*p.Supplements)[1].Amount = new(alpinebits.Decimal(0.0))
	(*p.Offers)[0].OfferRules.OfferRule.MinAdvancedBookingOffset = "P1D"
	*p.Offers = append(*p.Offers,
		rateplans.Offer{Discount: &rateplans.Discount{Percent: 100, NightsRequired: 4, NightsDiscounted: 1}},
		rateplans.Offer{Discount: &rateplans.Discount{Percent: 100}, Guests: &[]rateplans.Guest{
			{AgeQualifyingCode: rateplans.AgeQualifyingCode8, MaxAge: 5, MinCount: 0, FirstQualifyingPosition: 1, LastQualifyingPosition: 1},
		}},
	)

	out, err := xml.Marshal(rq)
	if err != nil {
		t.Fatal(err)
	}
	if err := Version.Validate(out); err != nil {
		t.Fatalf("XSD: %v\n%s", err, out)
	}
	var back rateplans.HotelRatePlanNotifRQ
	if err := xml.Unmarshal(out, &back); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back, rq) {
		t.Errorf("round trip changed the message:\n got %+v\nwant %+v", back, rq)
	}
}

func TestRatePlanRemoveMarshalsValidXML(t *testing.T) {
	rq := rateplans.HotelRatePlanNotifRQ{Version: "1.000", RatePlans: rateplans.RatePlans{HotelCode: "123", RatePlans: []rateplans.RatePlan{
		{RatePlanNotifType: rateplans.RatePlanNotifTypeRemove, CurrencyCode: "EUR", RatePlanCode: "A"},
	}}}
	out, err := xml.Marshal(rq)
	if err != nil {
		t.Fatal(err)
	}
	if err := Version.Validate(out); err != nil {
		t.Errorf("XSD: %v\n%s", err, out)
	}
}
