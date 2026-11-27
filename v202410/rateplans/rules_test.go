package rateplans

import (
	"encoding/xml"
	"os"
	"strings"
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
	plan := func(rq *HotelRatePlanNotifRQ) *RatePlan { return &rq.RatePlans.RatePlans[0] }
	tests := []struct {
		name   string
		change func(rq *HotelRatePlanNotifRQ)
	}{
		{"CompleteSet", func(rq *HotelRatePlanNotifRQ) {
			rq.UniqueID = &UniqueID{Type: UniqueIDTypeReference, ID: "1", Instance: InstanceCompleteSet}
			rq.RatePlans.RatePlans = append(rq.RatePlans.RatePlans, RatePlan{RatePlanCode: " "})
		}},
		{"second RatePlan", func(rq *HotelRatePlanNotifRQ) {
			rq.RatePlans.RatePlans = append(rq.RatePlans.RatePlans,
				RatePlan{RatePlanNotifType: RatePlanNotifTypeRemove, RatePlanCode: plan(rq).RatePlanCode, RatePlanID: "1"})
		}},
		{"RatePlan attributes", func(rq *HotelRatePlanNotifRQ) {
			p := plan(rq)
			p.RatePlanCode, p.CurrencyCode, p.RatePlanCategory = " ", "", "x"
		}},
		{"no RatePlanNotifType", func(rq *HotelRatePlanNotifRQ) { plan(rq).RatePlanNotifType = "" }},
		{"Remove with content", func(rq *HotelRatePlanNotifRQ) { plan(rq).RatePlanNotifType = RatePlanNotifTypeRemove }},
		{"Overlay", func(rq *HotelRatePlanNotifRQ) { plan(rq).RatePlanNotifType = RatePlanNotifTypeOverlay }},
		{"BookingRules", func(rq *HotelRatePlanNotifRQ) {
			p := plan(rq)
			*p.BookingRules = append(*p.BookingRules,
				BookingRule{Start: date("2014-04-17"), End: date("2014-05-01"),
					LengthsOfStay: &[]LengthOfStay{
						{Time: 0, TimeUnit: TimeUnitDay, MinMaxMessageType: MinMaxMessageTypeSetMinLOS},
						{Time: 1, TimeUnit: TimeUnitDay, MinMaxMessageType: MinMaxMessageTypeSetMinLOS},
					},
					RestrictionStatus: &RestrictionStatus{Status: StatusClose}},
				BookingRule{Code: "double", Start: date("2014-03-02"), End: date("2014-03-01")})
		}},
		{"no static Rate", func(rq *HotelRatePlanNotifRQ) { p := plan(rq); *p.Rates = (*p.Rates)[1:] }},
		{"static Rate", func(rq *HotelRatePlanNotifRQ) {
			p := plan(rq)
			r := &(*p.Rates)[0]
			r.MealsIncluded, r.RateTimeUnit = nil, ""
			second := *r
			second.BaseByGuestAmts = nil
			r.BaseByGuestAmts = &[]BaseByGuestAmt{{Type: BaseByGuestAmtType7, NumberOfGuests: 2}}
			*p.Rates = append(*p.Rates, second)
		}},
		{"date depending Rate", func(rq *HotelRatePlanNotifRQ) {
			r := &(*plan(rq).Rates)[1]
			r.InvTypeCode, r.End = " ", nil
			two := BaseByGuestAmt{NumberOfGuests: 2, AgeQualifyingCode: AgeQualifyingCode10, AmountAfterTax: 96}
			usd := two
			usd.CurrencyCode = "USD"
			r.BaseByGuestAmts = &[]BaseByGuestAmt{{}, two, usd}
			r.AdditionalGuestAmounts = &[]AdditionalGuestAmount{
				{AgeQualifyingCode: 8}, {AgeQualifyingCode: 7, Amount: new(alpinebits.Decimal(1.5555))},
				{AgeQualifyingCode: AgeQualifyingCode10, Amount: new(alpinebits.Decimal(1)), MinAge: 5, MaxAge: 3},
				{AgeQualifyingCode: AgeQualifyingCode10, Amount: new(alpinebits.Decimal(1))},
			}
		}},
		{"overlapping Rates", func(rq *HotelRatePlanNotifRQ) {
			p := plan(rq)
			r := (*p.Rates)[1]
			r.Start, r.End, r.BaseByGuestAmts = date("2014-03-08"), date("2014-03-10"), nil
			*p.Rates = append(*p.Rates, r)
		}},
		{"static Supplement", func(rq *HotelRatePlanNotifRQ) {
			s := &(*plan(rq).Supplements)[0]
			s.AddToBasicRateIndicator, s.MandatoryIndicator, s.ChargeTypeCode, s.Amount = false, nil, "", new(alpinebits.Decimal(1.0))
			s.PrerequisiteInventory = &PrerequisiteInventory{InvType: PrerequisiteInvTypeAlpineBitsDOW, InvCode: "1100002"}
		}},
		{"date depending Supplement", func(rq *HotelRatePlanNotifRQ) {
			s := &(*plan(rq).Supplements)[1]
			s.InvCode, s.End, s.ChargeTypeCode = " ", nil, ChargeTypeCode1
			s.PrerequisiteInventory = &PrerequisiteInventory{InvType: PrerequisiteInvTypeAlpineBitsDOW, InvCode: "1111111"}
		}},
		{"static Supplements per InvCode", func(rq *HotelRatePlanNotifRQ) {
			p := plan(rq)
			static := (*p.Supplements)[0]
			static.InvCode = "X"
			other := static
			other.PrerequisiteInventory = &PrerequisiteInventory{InvType: PrerequisiteInvTypeRoomType, InvCode: " "}
			other.ChargeTypeCode = ChargeTypeCode12
			p.Supplements = &[]Supplement{(*p.Supplements)[1], static, other}
		}},
		{"no Offer and Description", func(rq *HotelRatePlanNotifRQ) {
			p := plan(rq)
			p.Offers, p.Descriptions = nil, nil
		}},
		{"first Offer", func(rq *HotelRatePlanNotifRQ) {
			o := &(*plan(rq).Offers)[0]
			o.Discount = &Discount{Percent: 100}
			o.OfferRules.OfferRule.Occupancies = []Occupancy{
				{AgeQualifyingCode: AgeQualifyingCode10, MinAge: 19}, {AgeQualifyingCode: AgeQualifyingCode10},
				{AgeQualifyingCode: AgeQualifyingCode8, MinAge: 5, MaxAge: 3, MaxOccupancy: 100}, {AgeQualifyingCode: AgeQualifyingCode8},
			}
		}},
		{"no adult Occupancy", func(rq *HotelRatePlanNotifRQ) {
			(*plan(rq).Offers)[0].OfferRules.OfferRule.Occupancies = []Occupancy{{AgeQualifyingCode: AgeQualifyingCode8}}
		}},
		{"Offers without OfferRule and Discount", func(rq *HotelRatePlanNotifRQ) {
			p := plan(rq)
			(*p.Offers)[0].OfferRules = nil
			*p.Offers = append(*p.Offers, Offer{})
		}},
		{"family offers", func(rq *HotelRatePlanNotifRQ) {
			p := plan(rq)
			(*p.Offers)[0].OfferRules.OfferRule.Occupancies = []Occupancy{{AgeQualifyingCode: AgeQualifyingCode10}, {AgeQualifyingCode: AgeQualifyingCode8}}
			family := Offer{Discount: &Discount{Percent: 50, NightsRequired: 1}, Guests: &[]Guest{{AgeQualifyingCode: AgeQualifyingCode10}}}
			*p.Offers = append(*p.Offers, family, family)
		}},
		{"free nights offers", func(rq *HotelRatePlanNotifRQ) {
			p := plan(rq)
			(*p.Rates)[0].UnitMultiplier = 7
			*p.Offers = append(*p.Offers,
				Offer{Discount: &Discount{Percent: 100, NightsRequired: 2, NightsDiscounted: 3}},
				Offer{Discount: &Discount{Percent: 100, NightsRequired: 2}},
				Offer{Discount: &Discount{Percent: 100, NightsRequired: 3, NightsDiscounted: 1, DiscountPattern: "111"}})
		}},
		{"Descriptions", func(rq *HotelRatePlanNotifRQ) {
			p := plan(rq)
			p.Descriptions = append(p.Descriptions, p.Descriptions[0], Description{Name: DescriptionNameIntro, Texts: []alpinebits.Text{{TextFormat: alpinebits.TextFormatHTML, Value: "x"}}})
		}},
		{"currency", func(rq *HotelRatePlanNotifRQ) { plan(rq).CurrencyCode = "EUX" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := sample(t)
			tt.change(&rq)
			var want []string
			for _, e := range rateplans202210.Validate(rq, allCaps) {
				// 2024-10 names Full where 2022-10 names New.
				want = append(want, strings.ReplaceAll(e.Message, "RatePlanNotifType New", "RatePlanNotifType Full"))
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
}
