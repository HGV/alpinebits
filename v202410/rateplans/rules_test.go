package rateplans

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
	rateplans202210 "github.com/HGV/alpinebits/v202210/rateplans"
	"github.com/HGV/x/timex"
)

var allCaps = []alpinebits.Capability{
	CapabilityAcceptArrivalDOW, CapabilityAcceptDepartureDOW,
	CapabilityAcceptRatePlanBookingRule, CapabilityAcceptRatePlanRoomTypeBookingRule, CapabilityAcceptRatePlanMixedBookingRule,
	CapabilityAcceptSupplements, CapabilityAcceptFreeNightsOffers, CapabilityAcceptFamilyOffers, CapabilityAcceptOverlay,
	CapabilityAcceptRatePlanJoin, CapabilityAcceptOfferRuleBookingOffset, CapabilityAcceptOfferRuleDOWLOS, CapabilityAcceptFull,
}

func date(s string) *timex.Date {
	d, err := timex.ParseDate(s)
	if err != nil {
		panic(err)
	}
	return &d
}

// sample is the official 2018-10 sample: one New rate plan with a static and
// a date depending rate and supplement.
func sample(t *testing.T) HotelRatePlanNotifRQ {
	t.Helper()
	b, err := os.ReadFile("../../v201810/rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRQ.xml")
	if err != nil {
		t.Fatal(err)
	}
	var rq HotelRatePlanNotifRQ
	if err := xml.Unmarshal(b, &rq); err != nil {
		t.Fatal(err)
	}
	return rq
}

// The checks 2024-10 copied from 2022-10 must report what 2022-10 reports.
func TestRulesMatch202210(t *testing.T) {
	tests := []struct {
		name   string
		change func(p *RatePlan)
	}{
		{"Remove with content", func(p *RatePlan) { p.RatePlanNotifType = RatePlanNotifTypeRemove }},
		{"BookingRules", func(p *RatePlan) {
			*p.BookingRules = append(*p.BookingRules,
				BookingRule{Start: date("2014-04-17"), End: date("2014-05-01"),
					LengthsOfStay:     &[]LengthOfStay{{Time: 0, TimeUnit: TimeUnitDay, MinMaxMessageType: MinMaxMessageTypeSetMinLOS}},
					RestrictionStatus: &RestrictionStatus{Status: StatusClose}},
				BookingRule{Code: "double", Start: date("2014-03-02"), End: date("2014-03-01")})
		}},
		{"no static Rate", func(p *RatePlan) { *p.Rates = (*p.Rates)[1:] }},
		{"static Rate", func(p *RatePlan) {
			r := &(*p.Rates)[0]
			r.MealsIncluded, r.RateTimeUnit = nil, ""
			*p.Rates = append(*p.Rates, *r)
		}},
		{"date depending Rate", func(p *RatePlan) {
			r := &(*p.Rates)[1]
			r.InvTypeCode, r.End = " ", nil
			r.AdditionalGuestAmounts = &[]AdditionalGuestAmount{{AgeQualifyingCode: 8}, {AgeQualifyingCode: 7, Amount: new(alpinebits.Decimal(1.5555))}}
		}},
		{"overlapping Rates", func(p *RatePlan) {
			r := (*p.Rates)[1]
			r.Start, r.End = date("2014-03-08"), date("2014-03-10")
			*p.Rates = append(*p.Rates, r)
		}},
		{"static Supplement", func(p *RatePlan) {
			s := &(*p.Supplements)[0]
			s.AddToBasicRateIndicator, s.MandatoryIndicator, s.ChargeTypeCode, s.Amount = false, nil, "", new(alpinebits.Decimal(1.0))
			s.PrerequisiteInventory = &PrerequisiteInventory{InvType: PrerequisiteInvTypeAlpineBitsDOW, InvCode: "1100002"}
		}},
		{"date depending Supplement", func(p *RatePlan) {
			s := &(*p.Supplements)[1]
			s.End, s.ChargeTypeCode = nil, ChargeTypeCode1
			s.PrerequisiteInventory = &PrerequisiteInventory{InvType: PrerequisiteInvTypeAlpineBitsDOW, InvCode: "1111111"}
		}},
		{"static Supplements per InvCode", func(p *RatePlan) {
			static := (*p.Supplements)[0]
			static.InvCode = "X"
			p.Supplements = &[]Supplement{(*p.Supplements)[1], static, static}
		}},
		{"no Offer", func(p *RatePlan) { p.Offers = nil }},
		{"first Offer", func(p *RatePlan) {
			o := &(*p.Offers)[0]
			o.Discount = &Discount{Percent: 100}
			o.OfferRules.OfferRule.Occupancies = []Occupancy{{AgeQualifyingCode: AgeQualifyingCode10, MinAge: 19}, {AgeQualifyingCode: AgeQualifyingCode10}}
		}},
		{"free nights offers", func(p *RatePlan) {
			(*p.Rates)[0].UnitMultiplier = 7
			*p.Offers = append(*p.Offers, Offer{Discount: &Discount{Percent: 100, NightsRequired: 2, NightsDiscounted: 3}})
		}},
		{"Descriptions", func(p *RatePlan) {
			p.Descriptions = append(p.Descriptions, p.Descriptions[0], Description{Name: DescriptionNameIntro, Texts: []alpinebits.Text{{TextFormat: alpinebits.TextFormatHTML, Value: "x"}}})
		}},
		{"currency", func(p *RatePlan) { p.CurrencyCode = "EUX" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := sample(t)
			tt.change(&rq.RatePlans.RatePlans[0])
			var want []string
			for _, e := range rateplans202210.Validate(rq, allCaps) {
				want = append(want, e.Message)
			}
			if len(want) == 0 {
				t.Fatal("2022-10 reports no errors")
			}
			rulestest.Expect(t, Validate(rq, allCaps), want)
		})
	}
}

func TestValidate2024Changes(t *testing.T) {
	full := sample(t)
	full.RatePlans.RatePlans[0].RatePlanNotifType = RatePlanNotifTypeFull
	rulestest.Expect(t, Validate(full, allCaps), nil)
	rulestest.Expect(t, Validate(full, allCaps[:12]), []string{"RatePlan[0]: RatePlanNotifType Full not supported"})
	rulestest.Expect(t, Validate(sample(t), allCaps[:12]), []string{"RatePlan[0]: RatePlanNotifType New not supported"})

	overlay := sample(t)
	overlay.RatePlans.RatePlans[0].RatePlanNotifType = RatePlanNotifTypeOverlay
	rulestest.Expect(t, Validate(overlay, allCaps), []string{
		"RatePlan[0].Rate[0]: static Rate only allowed with RatePlanNotifType Full",
		"RatePlan[0].Supplement[0]: static Supplement only allowed with RatePlanNotifType Full",
		"RatePlan[0]: Offers only allowed with RatePlanNotifType Full",
	})
}
