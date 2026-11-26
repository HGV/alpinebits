package rateplans

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
	"github.com/HGV/x/timex"
)

var allCaps = []alpinebits.Capability{
	CapabilityAcceptArrivalDOW, CapabilityAcceptDepartureDOW,
	CapabilityAcceptRatePlanBookingRule, CapabilityAcceptRatePlanRoomTypeBookingRule, CapabilityAcceptRatePlanMixedBookingRule,
	CapabilityAcceptSupplements, CapabilityAcceptFreeNightsOffers, CapabilityAcceptFamilyOffers, CapabilityAcceptOverlay,
	CapabilityAcceptRatePlanJoin, CapabilityAcceptOfferRuleBookingOffset, CapabilityAcceptOfferRuleDOWLOS,
}

func date(s string) *timex.Date {
	d, err := timex.ParseDate(s)
	if err != nil {
		panic(err)
	}
	return &d
}

// sample is the official AlpineBits 2018-10 RatePlans sample: one New rate
// plan with a static and a date depending rate, a static and a date depending
// supplement, a generic booking rule with DOW restrictions and one Offer.
func sample(t *testing.T) HotelRatePlanNotifRQ {
	t.Helper()
	b, err := os.ReadFile("testdata/RatePlans-OTA_HotelRatePlanNotifRQ.xml")
	if err != nil {
		t.Fatal(err)
	}
	var rq HotelRatePlanNotifRQ
	if err := xml.Unmarshal(b, &rq); err != nil {
		t.Fatal(err)
	}
	return rq
}

func TestValidateSample(t *testing.T) {
	rulestest.Expect(t, Validate(sample(t), allCaps), nil)
}

func TestValidateCapabilities(t *testing.T) {
	rq := sample(t)
	p := &rq.RatePlans.RatePlans[0]
	p.RatePlanID, p.RatePlanQualifier = "M", new(true)
	*p.BookingRules = append(*p.BookingRules, BookingRule{Code: "double", CodeContext: CodeContextRoomType, Start: date("2014-05-01"), End: date("2014-05-31")})
	rule := (*p.Offers)[0].OfferRules.OfferRule
	rule.MinAdvancedBookingOffset = "P1D"
	rule.LengthsOfStay = &[]LengthOfStay{{Time: 7, TimeUnit: TimeUnitDay, MinMaxMessageType: MinMaxMessageTypeSetMinLOS}}
	*p.Offers = append(*p.Offers,
		Offer{Discount: &Discount{Percent: 100, NightsRequired: 4, NightsDiscounted: 1, DiscountPattern: "0001"}},
		Offer{Discount: &Discount{Percent: 100}, Guests: &[]Guest{{AgeQualifyingCode: AgeQualifyingCode8, MaxAge: 5, MinCount: 2, FirstQualifyingPosition: 1, LastQualifyingPosition: 1}}},
	)

	rulestest.Expect(t, Validate(rq, allCaps), nil)
	rulestest.Expect(t, Validate(rq, nil), []string{
		"RatePlan[0]: joining rate plans not supported",
		"RatePlan[0].BookingRule[0]: ArrivalDaysOfWeek not supported",
		"RatePlan[0].BookingRule[0]: DepartureDaysOfWeek not supported",
		"RatePlan[0]: BookingRule without Code not supported",
		"RatePlan[0]: BookingRule with Code not supported",
		"RatePlan[0]: BookingRules with and without Code in one RatePlan not supported",
		"RatePlan[0]: Supplements not supported",
		"RatePlan[0].Offer[0].OfferRule: MinAdvancedBookingOffset and MaxAdvancedBookingOffset not supported",
		"RatePlan[0].Offer[0].OfferRule: DOW and LOS restrictions in OfferRule not supported",
		"RatePlan[0].Offer[1]: free nights offers not supported",
		"RatePlan[0].Offer[2]: family offers not supported",
	})
}

func TestValidateMessage(t *testing.T) {
	plan := func(notif RatePlanNotifType, code string) RatePlan {
		return RatePlan{RatePlanNotifType: notif, CurrencyCode: "EUR", RatePlanCode: code}
	}
	completeSet := &UniqueID{Type: "16", ID: "1", Instance: "CompleteSet"}
	tests := []struct {
		name     string
		uniqueID *UniqueID
		plans    []RatePlan
		caps     []alpinebits.Capability
		want     []string
	}{
		{"reset", completeSet, []RatePlan{{}}, nil, nil},
		{"CompleteSet lists codes only", completeSet, []RatePlan{{RatePlanCode: "A"}, plan(RatePlanNotifTypeNew, " ")}, nil, []string{
			"RatePlan[1]: RatePlanCode required",
			"RatePlan[1]: RatePlan of a CompleteSet must have no RatePlanNotifType and no content",
		}},
		{"Remove", nil, []RatePlan{plan(RatePlanNotifTypeRemove, "A"), plan(RatePlanNotifTypeRemove, "B")}, nil, nil},
		{"Overlay deletes booking rules", nil, []RatePlan{{
			RatePlanNotifType: RatePlanNotifTypeOverlay, CurrencyCode: "EUR", RatePlanCode: "A",
			BookingRules: &[]BookingRule{{}},
		}}, allCaps, nil},
		{"Overlay not negotiated", nil, []RatePlan{plan(RatePlanNotifTypeOverlay, "A")}, nil, []string{
			"RatePlan[0]: RatePlanNotifType Overlay not supported",
		}},
		{"missing attributes", nil, []RatePlan{{RatePlanCode: " ", RatePlanCategory: "X"}}, nil, []string{
			"RatePlan[0]: RatePlanCode required",
			"RatePlan[0]: CurrencyCode required",
			"RatePlan[0]: RatePlanCategory requires RatePlanType 12",
			"RatePlan[0]: RatePlanNotifType required",
		}},
		{"RatePlanID outside New", nil, []RatePlan{{RatePlanNotifType: RatePlanNotifTypeRemove, CurrencyCode: "EUR", RatePlanCode: "A", RatePlanID: "M"}}, allCaps, []string{
			"RatePlan[0]: RatePlanID and RatePlanQualifier must be given together",
			"RatePlan[0]: RatePlanID and RatePlanQualifier only allowed with RatePlanNotifType New",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := HotelRatePlanNotifRQ{UniqueID: tt.uniqueID, RatePlans: RatePlans{HotelCode: "123", RatePlans: tt.plans}}
			rulestest.Expect(t, Validate(rq, tt.caps), tt.want)
		})
	}

	t.Run("New one per message", func(t *testing.T) {
		rq := sample(t)
		second := rq.RatePlans.RatePlans[0]
		second.RatePlanCode = "other"
		rq.RatePlans.RatePlans = append(rq.RatePlans.RatePlans, second)
		rulestest.Expect(t, Validate(rq, allCaps), []string{"RatePlans with RatePlanNotifType New must be sent one per message"})
	})
}

func TestValidateRatePlan(t *testing.T) {
	staticRate := func(p *RatePlan) *Rate { return &(*p.Rates)[0] }
	rate := func(p *RatePlan) *Rate { return &(*p.Rates)[1] }
	tests := []struct {
		name   string
		change func(p *RatePlan)
		want   []string
	}{
		{"Remove with content", func(p *RatePlan) { p.RatePlanNotifType = RatePlanNotifTypeRemove }, []string{
			"RatePlan[0]: RatePlan to Remove must be empty",
		}},
		{"Overlay with static data", func(p *RatePlan) { p.RatePlanNotifType = RatePlanNotifTypeOverlay }, []string{
			"RatePlan[0].Rate[0]: static Rate only allowed with RatePlanNotifType New",
			"RatePlan[0].Supplement[0]: static Supplement only allowed with RatePlanNotifType New",
			"RatePlan[0]: Offers only allowed with RatePlanNotifType New",
		}},
		{"BookingRules", func(p *RatePlan) {
			*p.BookingRules = append(*p.BookingRules,
				BookingRule{Start: date("2014-04-17"), End: date("2014-05-01"),
					LengthsOfStay: &[]LengthOfStay{
						{Time: 0, TimeUnit: TimeUnitDay, MinMaxMessageType: MinMaxMessageTypeSetMinLOS},
						{Time: 3, TimeUnit: TimeUnitDay, MinMaxMessageType: MinMaxMessageTypeSetMinLOS},
					},
					RestrictionStatus: &RestrictionStatus{Status: StatusClose}},
				BookingRule{Code: "double", Start: date("2014-03-02"), End: date("2014-03-01")},
			)
		}, []string{
			"RatePlan[0].BookingRule[1].LengthOfStay[0]: Time must be > 0",
			"RatePlan[0].BookingRule[1].LengthOfStay[1]: duplicate MinMaxMessageType SetMinLOS",
			"RatePlan[0].BookingRule[1]: RestrictionStatus requires Restriction and Status",
			"RatePlan[0].BookingRule[2]: CodeContext required with Code",
			"RatePlan[0].BookingRule[2]: End before Start",
			"RatePlan[0].BookingRule[1]: overlaps BookingRule[0]",
		}},
		{"no static Rate", func(p *RatePlan) { *p.Rates = (*p.Rates)[1:] }, []string{
			"RatePlan[0]: static Rate required",
		}},
		{"static Rate", func(p *RatePlan) {
			r := staticRate(p)
			(*r.BaseByGuestAmts)[0].NumberOfGuests = 2
			r.MealsIncluded = nil
			r.RateTimeUnit = ""
			*p.Rates = append(*p.Rates, *r)
		}, []string{
			"RatePlan[0].Rate[0].BaseByGuestAmt[0]: static BaseByGuestAmt allows only Type",
			"RatePlan[0].Rate[0]: MealsIncluded with MealPlanIndicator true and MealPlanCodes required",
			"RatePlan[0].Rate[0]: RateTimeUnit and UnitMultiplier must be given together",
			"RatePlan[0].Rate[2]: static Rate must be the first Rate",
			"RatePlan[0].Rate[2].BaseByGuestAmt[0]: static BaseByGuestAmt allows only Type",
			"RatePlan[0].Rate[2]: MealsIncluded with MealPlanIndicator true and MealPlanCodes required",
			"RatePlan[0].Rate[2]: RateTimeUnit and UnitMultiplier must be given together",
		}},
		{"date depending Rate", func(p *RatePlan) {
			r := rate(p)
			r.InvTypeCode = " "
			r.End = nil
			r.BaseByGuestAmts = &[]BaseByGuestAmt{{NumberOfGuests: 1, AgeQualifyingCode: 10, AmountAfterTax: 106}, {NumberOfGuests: 1}}
			r.AdditionalGuestAmounts = &[]AdditionalGuestAmount{
				{AgeQualifyingCode: 10, Amount: new(alpinebits.Decimal(76.8))},
				{AgeQualifyingCode: 10, MinAge: 3, Amount: new(alpinebits.Decimal(1.0))},
				{AgeQualifyingCode: 8},
				{AgeQualifyingCode: 8, MinAge: 6, MaxAge: 3, Amount: new(alpinebits.Decimal(0.0))},
				{AgeQualifyingCode: 7, Amount: new(alpinebits.Decimal(1.0))},
			}
		}, []string{
			"RatePlan[0].Rate[1]: InvTypeCode required",
			"RatePlan[0].Rate[1]: Start and End required",
			"RatePlan[0].Rate[1].BaseByGuestAmt[1]: duplicate NumberOfGuests 1",
			"RatePlan[0].Rate[1].BaseByGuestAmt[1]: AmountAfterTax required",
			"RatePlan[0].Rate[1].BaseByGuestAmt[1]: AgeQualifyingCode must be 10",
			"RatePlan[0].Rate[1].AdditionalGuestAmount[1]: at most one adult AdditionalGuestAmount allowed",
			"RatePlan[0].Rate[1].AdditionalGuestAmount[1]: adult AdditionalGuestAmount must have neither MinAge nor MaxAge",
			"RatePlan[0].Rate[1].AdditionalGuestAmount[2]: Amount required",
			"RatePlan[0].Rate[1].AdditionalGuestAmount[2]: child AdditionalGuestAmount requires MinAge or MaxAge",
			"RatePlan[0].Rate[1].AdditionalGuestAmount[3]: MaxAge must exceed MinAge",
			"RatePlan[0].Rate[1].AdditionalGuestAmount[4]: AgeQualifyingCode must be 8 or 10",
		}},
		{"overlapping Rates", func(p *RatePlan) {
			r := *rate(p)
			r.Start, r.End = date("2014-03-08"), date("2014-03-10")
			other := r
			other.InvTypeCode = "single"
			*p.Rates = append(*p.Rates, r, other)
		}, []string{
			"RatePlan[0].Rate[2]: overlaps Rate[1]",
		}},
		{"static Supplement", func(p *RatePlan) {
			s := &(*p.Supplements)[0]
			s.AddToBasicRateIndicator, s.MandatoryIndicator, s.ChargeTypeCode, s.Amount = false, nil, "", new(alpinebits.Decimal(1.0))
			s.PrerequisiteInventory = &PrerequisiteInventory{InvType: PrerequisiteInvTypeRoomType, InvCode: "double"}
		}, []string{
			"RatePlan[0].Supplement[0]: AddToBasicRateIndicator must be true",
			"RatePlan[0].Supplement[0]: MandatoryIndicator required",
			"RatePlan[0].Supplement[0]: ChargeTypeCode required",
			"RatePlan[0].Supplement[0]: Amount not allowed in a static Supplement",
			"RatePlan[0].Supplement[0].PrerequisiteInventory: InvType must be ALPINEBITSDOW in a static Supplement",
		}},
		{"static Supplement days of week", func(p *RatePlan) {
			(*p.Supplements)[0].PrerequisiteInventory = &PrerequisiteInventory{InvType: PrerequisiteInvTypeAlpineBitsDOW, InvCode: "1100002"}
		}, []string{
			"RatePlan[0].Supplement[0].PrerequisiteInventory: InvCode must be 7 binary digits",
		}},
		{"date depending Supplement", func(p *RatePlan) {
			s := &(*p.Supplements)[1]
			s.End = nil
			s.ChargeTypeCode = ChargeTypeCode1
			s.PrerequisiteInventory = &PrerequisiteInventory{InvType: PrerequisiteInvTypeAlpineBitsDOW, InvCode: "1111111"}
		}, []string{
			"RatePlan[0].Supplement[1]: Start and End required",
			"RatePlan[0].Supplement[1]: date depending Supplement allows only Amount, Start, End and PrerequisiteInventory",
			"RatePlan[0].Supplement[1].PrerequisiteInventory: InvType must be ROOMTYPE in a date depending Supplement",
		}},
		{"static Supplements per InvCode", func(p *RatePlan) {
			static := (*p.Supplements)[0]
			static.InvCode = "X"
			p.Supplements = &[]Supplement{(*p.Supplements)[1], static, static}
		}, []string{
			"RatePlan[0].Supplement[2]: duplicate static Supplement for InvCode X",
			"RatePlan[0].Supplement[0]: static Supplement required for InvCode 0x539",
		}},
		{"overlapping Supplements", func(p *RatePlan) {
			s := (*p.Supplements)[1]
			s.Start, s.End = date("2014-10-11"), date("2014-10-20")
			forDouble := s
			forDouble.PrerequisiteInventory = &PrerequisiteInventory{InvType: PrerequisiteInvTypeRoomType, InvCode: "double"}
			*p.Supplements = append(*p.Supplements, s, forDouble)
		}, []string{
			"RatePlan[0].Supplement[2]: overlaps Supplement[1]",
		}},
		{"no Offer", func(p *RatePlan) { p.Offers = nil }, []string{
			"RatePlan[0]: Offer required",
		}},
		{"first Offer", func(p *RatePlan) {
			o := &(*p.Offers)[0]
			o.Discount = &Discount{Percent: 100}
			o.OfferRules.OfferRule.Occupancies = []Occupancy{
				{AgeQualifyingCode: AgeQualifyingCode10, MinAge: 19},
				{AgeQualifyingCode: AgeQualifyingCode10},
				{AgeQualifyingCode: AgeQualifyingCode8, MinAge: 6, MaxAge: 3, MinOccupancy: 100},
				{AgeQualifyingCode: AgeQualifyingCode8},
			}
		}, []string{
			"RatePlan[0].Offer[0]: first Offer must have no Discount and no Guest",
			"RatePlan[0].Offer[0].OfferRule.Occupancy[0]: MinAge must be at most 18",
			"RatePlan[0].Offer[0].OfferRule.Occupancy[1]: at most one adult Occupancy allowed",
			"RatePlan[0].Offer[0].OfferRule.Occupancy[2]: MinOccupancy and MaxOccupancy must be at most 99",
			"RatePlan[0].Offer[0].OfferRule.Occupancy[2]: MaxAge must exceed MinAge",
			"RatePlan[0].Offer[0].OfferRule.Occupancy[3]: at most one child Occupancy allowed",
		}},
		{"children without adult MinAge", func(p *RatePlan) {
			(*p.Offers)[0].OfferRules.OfferRule.Occupancies[0].MinAge = 0
		}, []string{
			"RatePlan[0].Offer[0].OfferRule: child Occupancy requires MinAge of the adult Occupancy",
		}},
		{"no adult", func(p *RatePlan) {
			rule := (*p.Offers)[0].OfferRules.OfferRule
			rule.Occupancies = rule.Occupancies[1:]
		}, []string{
			"RatePlan[0].Offer[0].OfferRule: adult Occupancy required",
		}},
		{"free nights offers", func(p *RatePlan) {
			staticRate(p).UnitMultiplier = 7
			*p.Offers = append(*p.Offers,
				Offer{Discount: &Discount{Percent: 100, NightsRequired: 4, NightsDiscounted: 1, DiscountPattern: "0011"}},
				Offer{Discount: &Discount{Percent: 100, NightsRequired: 2, NightsDiscounted: 3}},
			)
		}, []string{
			"RatePlan[0].Offer[1]: DiscountPattern must be 0001",
			"RatePlan[0].Offer[1]: free nights offers require UnitMultiplier 1",
			"RatePlan[0].Offer[2]: at most one free nights offer allowed",
			"RatePlan[0].Offer[2]: NightsDiscounted exceeds NightsRequired",
			"RatePlan[0].Offer[2]: free nights offers require UnitMultiplier 1",
		}},
		{"family offers", func(p *RatePlan) {
			guest := Guest{AgeQualifyingCode: AgeQualifyingCode10, MaxAge: 5, MinCount: 1, FirstQualifyingPosition: 1, LastQualifyingPosition: 1}
			*p.Offers = append(*p.Offers,
				Offer{Discount: &Discount{Percent: 100, NightsRequired: 2}, Guests: &[]Guest{guest}},
				Offer{Discount: &Discount{Percent: 100}, Guests: &[]Guest{{AgeQualifyingCode: AgeQualifyingCode8}}},
			)
		}, []string{
			"RatePlan[0].Offer[1]: Discount of a family offer allows only Percent",
			"RatePlan[0].Offer[1].Guest[0]: AgeQualifyingCode must be 8",
			"RatePlan[0].Offer[2]: at most one family offer allowed",
		}},
		{"Offer without Discount", func(p *RatePlan) { *p.Offers = append(*p.Offers, Offer{}) }, []string{
			"RatePlan[0].Offer[1]: Discount required",
		}},
		{"Descriptions", func(p *RatePlan) {
			p.Descriptions = []Description{
				{Name: DescriptionNameTitle, Texts: []alpinebits.Text{
					{TextFormat: alpinebits.TextFormatPlainText, Language: "en"},
					{TextFormat: alpinebits.TextFormatPlainText, Language: "en"},
					{TextFormat: alpinebits.TextFormatHTML, Language: "de"},
					{TextFormat: alpinebits.TextFormatPlainText},
				}},
				{Name: DescriptionNameTitle},
			}
		}, []string{
			"RatePlan[0].Description[0].Text[3]: Language required",
			"RatePlan[0].Description[0].Text[1]: duplicate Text for Language en and TextFormat PlainText",
			"RatePlan[0].Description[0].Text[2]: HTML Text requires a PlainText Text in the same Language",
			"RatePlan[0].Description[1]: duplicate Name title",
		}},
		{"no Description", func(p *RatePlan) { p.Descriptions = nil }, []string{
			"RatePlan[0]: Description required",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := sample(t)
			tt.change(&rq.RatePlans.RatePlans[0])
			rulestest.Expect(t, Validate(rq, allCaps), tt.want)
		})
	}
}

func TestCheckInventory(t *testing.T) {
	rq := sample(t)
	p := &rq.RatePlans.RatePlans[0]
	*p.BookingRules = append(*p.BookingRules, BookingRule{Code: "suite", CodeContext: CodeContextRoomType})
	*p.Supplements = append(*p.Supplements, Supplement{InvCode: "0x539", PrerequisiteInventory: &PrerequisiteInventory{InvType: PrerequisiteInvTypeRoomType, InvCode: "single"}})

	rulestest.Expect(t, CheckInventory(rq, map[string]struct{}{"double": {}}), []string{
		"RatePlan[0].BookingRule[1]: unknown category suite",
		"RatePlan[0].Supplement[2].PrerequisiteInventory: unknown category single",
	})
}

func TestValidateHotel(t *testing.T) {
	rq := sample(t)
	rq.RatePlans.HotelName = ""
	rulestest.Expect(t, Validate(rq, allCaps), nil)
	rq.RatePlans.HotelCode = " "
	rulestest.Expect(t, Validate(rq, allCaps), []string{"HotelCode required"})
}
