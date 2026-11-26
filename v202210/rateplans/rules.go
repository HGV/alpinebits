package rateplans

import (
	"slices"
	"strings"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rules"
	rateplans202010 "github.com/HGV/alpinebits/v202010/rateplans"
	"github.com/HGV/x/timex"
)

// Validate checks the rules of AlpineBits 2022-10 §4.5 that the XSD cannot
// express.
func Validate(rq HotelRatePlanNotifRQ, caps []alpinebits.Capability) alpinebits.Errors {
	c := checker{rules.New(caps)}
	c.Hotel("", rq.RatePlans.HotelCode)
	if rq.IsReset() {
		return c.Errors()
	}
	plans := rq.RatePlans.RatePlans
	if rq.IsCompleteSet() {
		for i, plan := range plans {
			c.completeSetRatePlan(alpinebits.At("RatePlan", i), plan)
		}
	} else {
		for i, plan := range plans {
			c.ratePlan(alpinebits.At("RatePlan", i), plan)
		}
		if len(plans) > 1 && slices.ContainsFunc(plans, func(p RatePlan) bool { return complete(p.RatePlanNotifType) }) {
			c.Invalid("", "RatePlans with RatePlanNotifType New must be sent one per message")
		}
	}
	c.codesAndCurrencies(plans)
	return c.Errors()
}

// The checks are those of 2020-10 apart from the ones marked 2022-10;
// TestRulesMatch202010 keeps them in step.
type checker struct{ *rules.Checker }

// complete reports whether notif sends a whole rate plan.
func complete(notif RatePlanNotifType) bool {
	return notif == RatePlanNotifTypeNew
}

// dateRange reports whether start and end form a range usable for the overlap
// checks.
func (c *checker) dateRange(at alpinebits.Path, start, end *timex.Date) bool {
	switch {
	case start == nil || end == nil:
		c.Required(at, "Start and End required")
	case end.Before(*start):
		c.InvalidDateRange(at, "End before Start")
	default:
		return true
	}
	return false
}

func isEmptyPlan(p RatePlan) bool {
	return p.BookingRules == nil && p.Rates == nil && p.Supplements == nil && p.Offers == nil && len(p.Descriptions) == 0
}

func (c *checker) completeSetRatePlan(at alpinebits.Path, plan RatePlan) {
	if rules.Blank(plan.RatePlanCode) {
		c.Required(at, "RatePlanCode required")
	}
	if plan.RatePlanNotifType != "" || !isEmptyPlan(plan) {
		c.Invalid(at, "RatePlan of a CompleteSet must have no RatePlanNotifType and no content")
	}
}

func (c *checker) ratePlan(at alpinebits.Path, plan RatePlan) {
	notif := plan.RatePlanNotifType
	if rules.Blank(plan.RatePlanCode) {
		c.Required(at, "RatePlanCode required")
	}
	if plan.CurrencyCode == "" {
		c.Required(at, "CurrencyCode required")
	}
	if plan.RatePlanCategory != "" && plan.RatePlanType != RatePlanType12 {
		c.Invalid(at, "RatePlanCategory requires RatePlanType 12")
	}
	if plan.RatePlanID != "" || plan.RatePlanQualifier != nil {
		c.Need(at, CapabilityAcceptRatePlanJoin, "joining rate plans")
		if plan.RatePlanID == "" || plan.RatePlanQualifier == nil {
			c.Required(at, "RatePlanID and RatePlanQualifier must be given together")
		}
		if !complete(notif) {
			c.Invalid(at, "RatePlanID and RatePlanQualifier only allowed with RatePlanNotifType New")
		}
	}
	switch notif {
	case "":
		c.Required(at, "RatePlanNotifType required")
		return
	case RatePlanNotifTypeRemove:
		if !isEmptyPlan(plan) {
			c.Invalid(at, "RatePlan to Remove must be empty")
		}
		return
	case RatePlanNotifTypeOverlay:
		c.Need(at, CapabilityAcceptOverlay, "RatePlanNotifType Overlay")
	}

	c.bookingRules(at, plan)
	c.rates(at, plan)
	c.supplements(at, plan)
	if complete(notif) {
		c.offers(at, plan)
		if len(plan.Descriptions) == 0 {
			c.Required(at, "Description required")
		}
	} else if plan.Offers != nil {
		c.Invalid(at, "Offers only allowed with RatePlanNotifType New")
	}
	c.descriptions(at, plan.Descriptions)
}

func (c *checker) bookingRules(at alpinebits.Path, plan RatePlan) {
	var generic, specific bool
	var ranged []int
	bookingRules := rules.Items(plan.BookingRules)
	for j, r := range bookingRules {
		rat := at.At("BookingRule", j)
		// Empty BookingRule elements delete the booking rules (§4.5.3).
		if plan.RatePlanNotifType == RatePlanNotifTypeOverlay && r == (BookingRule{}) {
			continue
		}
		if r.Code == "" {
			generic = true
		} else {
			specific = true
			if r.CodeContext == "" {
				c.Required(rat, "CodeContext required with Code")
			}
		}
		if c.dateRange(rat, r.Start, r.End) {
			ranged = append(ranged, j)
		}
		c.lengthsOfStay(rat, r.LengthsOfStay)
		if d := r.DOWRestrictions; d != nil {
			if d.ArrivalDaysOfWeek != nil {
				c.Need(rat, CapabilityAcceptArrivalDOW, "ArrivalDaysOfWeek")
			}
			if d.DepartureDaysOfWeek != nil {
				c.Need(rat, CapabilityAcceptDepartureDOW, "DepartureDaysOfWeek")
			}
		}
		if s := r.RestrictionStatus; s != nil && (s.Restriction == "" || s.Status == "") {
			c.Required(rat, "RestrictionStatus requires Restriction and Status")
		}
	}
	if generic {
		c.Need(at, CapabilityAcceptRatePlanBookingRule, "BookingRule without Code")
	}
	if specific {
		c.Need(at, CapabilityAcceptRatePlanRoomTypeBookingRule, "BookingRule with Code")
	}
	if generic && specific {
		c.Need(at, CapabilityAcceptRatePlanMixedBookingRule, "BookingRules with and without Code in one RatePlan")
	}
	rules.Overlaps(c.Checker, at, "BookingRule", bookingRules, ranged, func(r BookingRule) string { return r.Code })
}

func (c *checker) lengthsOfStay(at alpinebits.Path, los *[]LengthOfStay) {
	seen := make(map[MinMaxMessageType]bool)
	for k, l := range rules.Items(los) {
		lat := at.At("LengthOfStay", k)
		if l.Time == 0 {
			c.Invalid(lat, "Time must be > 0")
		}
		if seen[l.MinMaxMessageType] {
			c.Invalid(lat, "duplicate MinMaxMessageType %s", l.MinMaxMessageType)
		}
		seen[l.MinMaxMessageType] = true
	}
}

// A static Rate holds what applies to the whole rate plan; it has no dates
// and no room category (§4.5.1).
func isStaticRate(r Rate) bool {
	return r.InvTypeCode == "" && r.Start == nil && r.End == nil
}

// isDeletion reports whether r is a date depending Rate without content.
func isDeletion(r Rate) bool {
	return !isStaticRate(r) && r.BaseByGuestAmts == nil && r.AdditionalGuestAmounts == nil && r.MealsIncluded == nil
}

func (c *checker) rates(at alpinebits.Path, plan RatePlan) {
	overlay := plan.RatePlanNotifType == RatePlanNotifTypeOverlay
	withStatic := complete(plan.RatePlanNotifType)
	rates := rules.Items(plan.Rates)
	if withStatic && !slices.ContainsFunc(rates, isStaticRate) {
		c.Required(at, "static Rate required")
	}
	var ranged []int
	for j, r := range rates {
		rat := at.At("Rate", j)
		if isStaticRate(r) {
			if j > 0 {
				c.Invalid(rat, "static Rate must be the first Rate")
			}
			if !withStatic {
				c.Invalid(rat, "static Rate only allowed with RatePlanNotifType New")
			}
			c.staticRate(rat, r)
			continue
		}
		if rules.Blank(r.InvTypeCode) {
			c.Required(rat, "InvTypeCode required")
		}
		dated := c.dateRange(rat, r.Start, r.End)
		// 2022-10: empty date depending Rates delete in an Overlay (§4.5.3).
		if overlay && isDeletion(r) {
			continue
		}
		if dated {
			ranged = append(ranged, j)
		}
		c.baseByGuestAmts(rat, r.BaseByGuestAmts)
		c.additionalGuestAmounts(rat, r.AdditionalGuestAmounts)
	}
	rules.Overlaps(c.Checker, at, "Rate", rates, ranged, func(r Rate) string { return r.InvTypeCode })
}

func (c *checker) staticRate(at alpinebits.Path, r Rate) {
	amts := rules.Items(r.BaseByGuestAmts)
	if len(amts) != 1 || amts[0].Type == "" {
		c.Required(at, "static Rate requires one BaseByGuestAmt with Type")
	} else if b := amts[0]; b != (BaseByGuestAmt{Type: b.Type}) {
		c.Invalid(at.At("BaseByGuestAmt", 0), "static BaseByGuestAmt allows only Type")
	}
	if m := r.MealsIncluded; m == nil || !m.MealPlanIndicator || m.MealPlanCodes == "" {
		c.Required(at, "MealsIncluded with MealPlanIndicator true and MealPlanCodes required")
	}
	if (r.RateTimeUnit == "") != (r.UnitMultiplier == 0) {
		c.Required(at, "RateTimeUnit and UnitMultiplier must be given together")
	}
}

func (c *checker) baseByGuestAmts(at alpinebits.Path, amts *[]BaseByGuestAmt) {
	if len(rules.Items(amts)) == 0 {
		c.Required(at, "BaseByGuestAmt required")
	}
	seen := make(map[int]bool)
	for k, b := range rules.Items(amts) {
		bat := at.At("BaseByGuestAmt", k)
		switch {
		case b.NumberOfGuests == 0:
			c.Required(bat, "NumberOfGuests required")
		case seen[b.NumberOfGuests]:
			c.Invalid(bat, "duplicate NumberOfGuests %d", b.NumberOfGuests)
		}
		seen[b.NumberOfGuests] = true
		// Mandatory (§4.5.1); the XSD only requires > 0 if present.
		if b.AmountAfterTax == 0 {
			c.Required(bat, "AmountAfterTax required")
		}
		if b.AgeQualifyingCode != AgeQualifyingCode10 {
			c.Invalid(bat, "AgeQualifyingCode must be 10")
		}
	}
}

func (c *checker) additionalGuestAmounts(at alpinebits.Path, amts *[]AdditionalGuestAmount) {
	adults := 0
	for k, a := range rules.Items(amts) {
		aat := at.At("AdditionalGuestAmount", k)
		if a.Amount == nil {
			c.Required(aat, "Amount required")
		}
		if a.MinAge != 0 && a.MaxAge != 0 && a.MaxAge <= a.MinAge {
			c.Invalid(aat, "MaxAge must exceed MinAge")
		}
		switch a.AgeQualifyingCode {
		case AgeQualifyingCode10:
			if adults++; adults > 1 {
				c.Invalid(aat, "at most one adult AdditionalGuestAmount allowed")
			}
			if a.MinAge != 0 || a.MaxAge != 0 {
				c.Invalid(aat, "adult AdditionalGuestAmount must have neither MinAge nor MaxAge")
			}
		case AgeQualifyingCode8:
			if a.MinAge == 0 && a.MaxAge == 0 {
				c.Required(aat, "child AdditionalGuestAmount requires MinAge or MaxAge")
			}
		default:
			c.Invalid(aat, "AgeQualifyingCode must be 8 or 10")
		}
	}
}

// A static Supplement describes the supplement; date depending ones give its
// price for a period (§4.5.1).
func isStaticSupplement(s Supplement) bool {
	return s.Start == nil && s.End == nil
}

func (c *checker) supplements(at alpinebits.Path, plan RatePlan) {
	supplements := rules.Items(plan.Supplements)
	if len(supplements) == 0 {
		return
	}
	c.Need(at, CapabilityAcceptSupplements, "Supplements")
	isComplete := complete(plan.RatePlanNotifType)
	static := make(map[string]int)
	var ranged []int
	for j, s := range supplements {
		sat := at.At("Supplement", j)
		pat := sat.Field("PrerequisiteInventory")
		p := s.PrerequisiteInventory
		if rules.Blank(s.InvCode) {
			c.Required(sat, "InvCode required")
		}
		if p != nil && rules.Blank(p.InvCode) {
			c.Required(pat, "InvCode required")
		}
		if !isStaticSupplement(s) {
			if c.dateRange(sat, s.Start, s.End) {
				ranged = append(ranged, j)
			}
			if s.AddToBasicRateIndicator || s.MandatoryIndicator != nil || s.ChargeTypeCode != "" || len(s.Descriptions) > 0 {
				c.Invalid(sat, "date depending Supplement allows only Amount, Start, End and PrerequisiteInventory")
			}
			if p != nil && p.InvType != PrerequisiteInvTypeRoomType {
				c.Invalid(pat, "InvType must be ROOMTYPE in a date depending Supplement")
			}
			continue
		}
		if !isComplete {
			c.Invalid(sat, "static Supplement only allowed with RatePlanNotifType New")
		}
		if static[s.InvCode]++; static[s.InvCode] > 1 {
			c.Invalid(sat, "duplicate static Supplement for InvCode %s", s.InvCode)
		}
		if !s.AddToBasicRateIndicator {
			c.Invalid(sat, "AddToBasicRateIndicator must be true")
		}
		if s.MandatoryIndicator == nil {
			c.Required(sat, "MandatoryIndicator required")
		}
		switch s.ChargeTypeCode {
		case "":
			c.Required(sat, "ChargeTypeCode required")
		case ChargeTypeCode12:
			// 2022-10: §4.5.1 dropped 12 (per stay); the XSD still has it.
			c.Invalid(sat, "ChargeTypeCode 12 not allowed")
		}
		if s.Amount != nil {
			c.Invalid(sat, "Amount not allowed in a static Supplement")
		}
		if p != nil {
			switch {
			case p.InvType != PrerequisiteInvTypeAlpineBitsDOW:
				c.Invalid(pat, "InvType must be ALPINEBITSDOW in a static Supplement")
			case !isDaysOfWeek(p.InvCode):
				c.Invalid(pat, "InvCode must be 7 binary digits")
			}
		}
		c.descriptions(sat, s.Descriptions)
	}
	if isComplete {
		reported := make(map[string]bool)
		for j, s := range supplements {
			if static[s.InvCode] == 0 && !rules.Blank(s.InvCode) && !reported[s.InvCode] {
				c.Required(at.At("Supplement", j), "static Supplement required for InvCode %s", s.InvCode)
				reported[s.InvCode] = true
			}
		}
	}
	rules.Overlaps(c.Checker, at, "Supplement", supplements, ranged, func(s Supplement) [2]string {
		if p := s.PrerequisiteInventory; p != nil {
			return [2]string{s.InvCode, p.InvCode}
		}
		return [2]string{s.InvCode}
	})
}

func isDaysOfWeek(s string) bool {
	return len(s) == 7 && strings.Trim(s, "01") == ""
}

func (c *checker) offers(at alpinebits.Path, plan RatePlan) {
	offers := rules.Items(plan.Offers)
	if len(offers) == 0 {
		c.Required(at, "Offer required")
		return
	}
	var freeNights, family bool
	for j, o := range offers {
		oat := at.At("Offer", j)
		if o.OfferRules != nil && o.OfferRules.OfferRule != nil {
			c.offerRuleCapabilities(oat.Field("OfferRule"), *o.OfferRules.OfferRule)
		}
		if j == 0 {
			c.firstOffer(oat, o)
			continue
		}
		d := o.Discount
		if d == nil {
			c.Required(oat, "Discount required")
			continue
		}
		if guests := rules.Items(o.Guests); len(guests) > 0 {
			c.Need(oat, CapabilityAcceptFamilyOffers, "family offers")
			if family {
				c.Invalid(oat, "at most one family offer allowed")
			}
			family = true
			if d.NightsRequired != 0 || d.NightsDiscounted != 0 || d.DiscountPattern != "" {
				c.Invalid(oat, "Discount of a family offer allows only Percent")
			}
			for k, g := range guests {
				if g.AgeQualifyingCode != AgeQualifyingCode8 {
					c.Invalid(oat.At("Guest", k), "AgeQualifyingCode must be 8")
				}
			}
			continue
		}
		c.Need(oat, CapabilityAcceptFreeNightsOffers, "free nights offers")
		if freeNights {
			c.Invalid(oat, "at most one free nights offer allowed")
		}
		freeNights = true
		switch {
		case d.NightsRequired == 0 || d.NightsDiscounted == 0:
			c.Required(oat, "NightsRequired and NightsDiscounted required")
		case d.NightsDiscounted > d.NightsRequired:
			c.Invalid(oat, "NightsDiscounted exceeds NightsRequired")
		case d.DiscountPattern != "":
			if want := strings.Repeat("0", d.NightsRequired-d.NightsDiscounted) + strings.Repeat("1", d.NightsDiscounted); d.DiscountPattern != want {
				c.Invalid(oat, "DiscountPattern must be %s", want)
			}
		}
		if rates := rules.Items(plan.Rates); len(rates) > 0 && isStaticRate(rates[0]) && rates[0].UnitMultiplier > 1 {
			c.Invalid(oat, "free nights offers require UnitMultiplier 1")
		}
	}
}

func (c *checker) offerRuleCapabilities(at alpinebits.Path, r OfferRule) {
	if r.MinAdvancedBookingOffset != "" || r.MaxAdvancedBookingOffset != "" {
		c.Need(at, CapabilityAcceptOfferRuleBookingOffset, "MinAdvancedBookingOffset and MaxAdvancedBookingOffset")
	}
	if len(rules.Items(r.LengthsOfStay)) > 0 || r.DOWRestrictions != nil {
		c.Need(at, CapabilityAcceptOfferRuleDOWLOS, "DOW and LOS restrictions in OfferRule")
	}
}

func (c *checker) firstOffer(at alpinebits.Path, o Offer) {
	if o.Discount != nil || len(rules.Items(o.Guests)) > 0 {
		c.Invalid(at, "first Offer must have no Discount and no Guest")
	}
	if o.OfferRules == nil || o.OfferRules.OfferRule == nil {
		c.Required(at, "OfferRule required")
		return
	}
	rat := at.Field("OfferRule")
	c.lengthsOfStay(rat, o.OfferRules.OfferRule.LengthsOfStay)
	var adult, child *Occupancy
	for k, occ := range o.OfferRules.OfferRule.Occupancies {
		oat := rat.At("Occupancy", k)
		if occ.MinOccupancy > 99 || occ.MaxOccupancy > 99 {
			c.Invalid(oat, "MinOccupancy and MaxOccupancy must be at most 99")
		}
		switch occ.AgeQualifyingCode {
		case AgeQualifyingCode10:
			if adult != nil {
				c.Invalid(oat, "at most one adult Occupancy allowed")
			} else {
				adult = &occ
			}
			if occ.MinAge > 18 {
				c.Invalid(oat, "MinAge must be at most 18")
			}
		case AgeQualifyingCode8:
			if child != nil {
				c.Invalid(oat, "at most one child Occupancy allowed")
			} else {
				child = &occ
			}
			if occ.MinAge != 0 && occ.MaxAge != 0 && occ.MaxAge <= occ.MinAge {
				c.Invalid(oat, "MaxAge must exceed MinAge")
			}
		}
	}
	switch {
	case adult == nil:
		c.Required(rat, "adult Occupancy required")
	case child != nil && adult.MinAge == 0:
		c.Invalid(rat, "child Occupancy requires MinAge of the adult Occupancy")
	}
}

func (c *checker) descriptions(at alpinebits.Path, ds []Description) {
	seen := make(map[DescriptionName]bool)
	for k, d := range ds {
		dat := at.At("Description", k)
		if seen[d.Name] {
			c.Invalid(dat, "duplicate Name %s", d.Name)
		}
		seen[d.Name] = true
		switch d.Name {
		case DescriptionNameTitle, DescriptionNameIntro, DescriptionNameDescription:
			for k, t := range d.Texts {
				if t.Language == "" {
					c.Required(dat.At("Text", k), "Language required")
				}
			}
			c.Texts(dat, "Text", d.Texts)
		}
	}
}

// codesAndCurrencies requires at most one RatePlan per RatePlanCode and ISO
// 4217 currencies, with amounts in their decimal digits.
func (c *checker) codesAndCurrencies(plans []RatePlan) {
	seen := make(map[string]bool)
	for i, plan := range plans {
		at := alpinebits.At("RatePlan", i)
		if code := plan.RatePlanCode; !rules.Blank(code) {
			if seen[code] {
				c.Invalid(at, "duplicate RatePlanCode %s", code)
			}
			seen[code] = true
		}
		if plan.CurrencyCode != "" {
			c.Currency(at, plan.CurrencyCode)
			c.amounts(at, plan)
		}
	}
}

// amounts checks that the amounts of plan are in its CurrencyCode.
func (c *checker) amounts(at alpinebits.Path, plan RatePlan) {
	code := plan.CurrencyCode
	for j, r := range rules.Items(plan.Rates) {
		rat := at.At("Rate", j)
		for k, b := range rules.Items(r.BaseByGuestAmts) {
			bat := rat.At("BaseByGuestAmt", k)
			if b.CurrencyCode != "" && b.CurrencyCode != code {
				c.Invalid(bat, "CurrencyCode %s differs from RatePlan CurrencyCode %s", b.CurrencyCode, code)
			}
			c.Amount(bat, "AmountAfterTax", b.AmountAfterTax, code)
		}
		for k, a := range rules.Items(r.AdditionalGuestAmounts) {
			if a.Amount != nil {
				c.Amount(rat.At("AdditionalGuestAmount", k), "Amount", *a.Amount, code)
			}
		}
	}
	for j, s := range rules.Items(plan.Supplements) {
		if s.Amount != nil {
			c.Amount(at.At("Supplement", j), "Amount", *s.Amount, code)
		}
	}
}

// CheckInventory reports references to room categories (InvTypeCode) not in
// categories.
func CheckInventory(rq HotelRatePlanNotifRQ, categories map[string]struct{}) alpinebits.Errors {
	return rateplans202010.CheckInventory(rq, categories)
}
