package rateplans

import (
	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rules"
	rateplans201810 "github.com/HGV/alpinebits/v201810/rateplans"
)

// Validate checks the rules of AlpineBits 2020-10 §4.5 that the XSD cannot
// express: those of 2018-10, ISO 4217 currencies with their decimal digits,
// and at most one RatePlan per RatePlanCode.
func Validate(rq HotelRatePlanNotifRQ, caps []alpinebits.Capability) alpinebits.Errors {
	errs := rateplans201810.Validate(rq, caps)
	if rq.IsReset() {
		return errs
	}
	c := checker{rules.New(caps)}
	seen := make(map[string]bool)
	for i, plan := range rq.RatePlans.RatePlans {
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
	return append(errs, c.Errors()...)
}

type checker struct{ *rules.Checker }

// amounts checks that the amounts of plan are in its CurrencyCode (§4.5.1).
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
	return rateplans201810.CheckInventory(rq, categories)
}
