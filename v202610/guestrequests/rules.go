package guestrequests

import (
	"strings"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rules"
	guestrequests202410 "github.com/HGV/alpinebits/v202410/guestrequests"
)

// ValidateReadRQ checks the rules of AlpineBits 2026-10 §4.2.1 that the XSD
// cannot express. They are those of 2024-10.
func ValidateReadRQ(rq ReadRQ, caps []alpinebits.Capability) alpinebits.Errors {
	return guestrequests202410.ValidateReadRQ(rq, caps)
}

// ValidateNotifReportRQ checks the rules of AlpineBits 2026-10 §4.2.3 that the
// XSD cannot express. Whether the IDs are known is up to the handler.
func ValidateNotifReportRQ(rq NotifReportRQ, _ []alpinebits.Capability) alpinebits.Errors {
	c := rules.New(nil)
	for i, w := range rules.Items(rq.Warnings) {
		if rules.Blank(w.RecordID) {
			c.Required(alpinebits.At("Warning", i), "RecordID required")
		}
	}
	for i, res := range rules.Items(rq.HotelReservations) {
		at := alpinebits.At("HotelReservation", i)
		if rules.Blank(res.UniqueID.ID) {
			c.Required(at.Field("UniqueID"), "ID required")
		}
		// 2026-10: the hotel of the acknowledged reservation.
		if b := res.BasicPropertyInfo; b != nil {
			c.Hotel(at.Field("BasicPropertyInfo"), b.HotelCode)
		}
	}
	return c.Errors()
}

// ValidateResRetrieveRS checks the rules of AlpineBits 2026-10 §4.2.2 that the
// XSD cannot express. The router runs no rules on responses: handlers call it
// on the response they build.
func ValidateResRetrieveRS(rs ResRetrieveRS) alpinebits.Errors {
	c := checker{rules.New(nil)}
	for i, res := range rs.HotelReservations {
		c.hotelReservation(alpinebits.At("HotelReservation", i), res)
	}
	return c.Errors()
}

// ValidateHotelReservation checks res like ValidateResRetrieveRS, so a handler
// can leave out the reservations that fail instead of failing the response.
func ValidateHotelReservation(res HotelReservation) alpinebits.Errors {
	c := checker{rules.New(nil)}
	c.hotelReservation("", res)
	return c.Errors()
}

// The checks are those of 2024-10 apart from the ones marked 2026-10;
// TestRulesMatch202410 keeps them in step.
type checker struct{ *rules.Checker }

// isReservation reports whether s is a booking reservation or modification,
// which §4.2.2 treats alike.
func isReservation(s ResStatus) bool {
	return s == ResStatusReserved || s == ResStatusModify
}

// A cancellation may repeat any part of the reservation or quote request it
// cancels, so its parts are checked as far as both allow.
func (c *checker) hotelReservation(at alpinebits.Path, res HotelReservation) {
	status := res.ResStatus
	uat := at.Field("UniqueID")
	want := UniqueIDTypeReservation
	if status == ResStatusCancelled {
		want = UniqueIDTypeCancellation
	}
	if res.UniqueID.Type != want {
		c.Invalid(uat, "Type must be %s for ResStatus %s", want, status)
	}
	if rules.Blank(res.UniqueID.ID) {
		c.Required(uat, "ID required")
	}
	// The XSD cannot catch it: Go writes the zero time as a valid xs:dateTime.
	if res.CreateDateTime.Time.IsZero() {
		c.Required(at, "CreateDateTime required")
	}
	// RoomStays are mandatory for reservations, optional for quote requests and
	// not allowed in newsletter and catalog requests.
	newsletter := res.RoomStayReservation != nil && !*res.RoomStayReservation
	switch {
	case newsletter && isReservation(status):
		c.Invalid(at, "RoomStayReservation false not allowed in a reservation")
	case newsletter && res.RoomStays != nil:
		c.Invalid(at, "RoomStays not allowed in a newsletter or catalog request")
	case isReservation(status) && res.RoomStays == nil:
		c.Required(at, "RoomStays required")
	}
	// 2026-10: ResGuests are optional.
	if status != ResStatusCancelled {
		if res.ResGlobalInfo == nil {
			c.Required(at, "ResGlobalInfo required")
		}
	}
	c.roomStays(at, status, rules.Items(res.RoomStays))
	// Services, which RoomStays refer to by ServiceRPH.
	rphs := c.services(at, rules.Items(res.Services))
	for j, stay := range rules.Items(res.RoomStays) {
		for k, ref := range rules.Items(stay.ServiceRPHs) {
			if !rphs[ref.RPH] {
				c.Invalid(at.At("RoomStay", j).At("ServiceRPH", k), "ServiceRPH %s refers to no Service", ref.RPH)
			}
		}
	}
	// 2026-10: guests, which RoomStays refer to by ResGuestRPH.
	guests := make(map[string]bool)
	for k, g := range rules.Items(res.ResGuests) {
		if g.ResGuestRPH == "" {
			continue
		}
		if guests[g.ResGuestRPH] {
			c.Invalid(at.At("ResGuest", k), "duplicate ResGuestRPH %s", g.ResGuestRPH)
		}
		guests[g.ResGuestRPH] = true
	}
	for j, stay := range rules.Items(res.RoomStays) {
		for _, ref := range stay.ResGuestRPHs {
			if !guests[ref] {
				c.Invalid(at.At("RoomStay", j), "ResGuestRPH %s refers to no ResGuest", ref)
			}
		}
	}
	if g := res.ResGlobalInfo; g != nil {
		c.resGlobalInfo(at.Field("ResGlobalInfo"), status, *g)
	}
}

func (c *checker) roomStays(at alpinebits.Path, status ResStatus, stays []RoomStay) {
	for j, stay := range stays {
		sat := at.At("RoomStay", j)
		c.roomStay(sat, status, stay)
		// RoomStayGroupID groups alternatives in quote requests.
		if stay.RoomStayGroupID != "" && isReservation(status) {
			c.Invalid(sat, "RoomStayGroupID not allowed in a reservation")
		}
	}
}

func (c *checker) roomStay(at alpinebits.Path, status ResStatus, stay RoomStay) {
	reservation := isReservation(status)
	// RoomType is mandatory only for reservations.
	if t := stay.RoomType; t == nil {
		if reservation {
			c.Required(at, "RoomType required")
		}
	} else {
		tat := at.Field("RoomType")
		if reservation && rules.Blank(t.RoomTypeCode) {
			c.Required(tat, "RoomTypeCode required")
		}
		if t.RoomType != 0 && t.RoomClassificationCode == nil {
			c.Required(tat, "RoomType requires RoomClassificationCode")
		}
	}

	if p := stay.RatePlan; p == nil {
		if reservation {
			c.Required(at, "RatePlan required")
		}
	} else {
		pat := at.Field("RatePlan")
		if reservation && rules.Blank(p.RatePlanCode) {
			c.Required(pat, "RatePlanCode required")
		}
		if m := p.Commission; m != nil {
			c.commission(pat.Field("Commission"), *m)
		}
		if reservation && p.MealsIncluded == nil {
			c.Required(pat, "MealsIncluded required")
		}
	}

	if r := stay.RoomRate; r != nil {
		c.roomRate(at.Field("RoomRate"), *r)
		// 2026-10: both RatePlanCodes name the same rate plan.
		if p := stay.RatePlan; p != nil && !rules.Blank(p.RatePlanCode) && r.RatePlanCode != p.RatePlanCode {
			c.Invalid(at.Field("RoomRate"), "RatePlanCode %s differs from RatePlan RatePlanCode %s", r.RatePlanCode, p.RatePlanCode)
		}
	}
	c.guestCounts(at, stay.GuestCounts)
	c.timeSpan(at.Field("TimeSpan"), status, stay.TimeSpan)

	if g := stay.Guarantee; g != nil {
		if status == ResStatusRequested {
			c.Invalid(at, "Guarantee not allowed in a quote request")
		}
		c.paymentCard(at.Field("PaymentCard"), g.PaymentCard)
	}
	if t := stay.Total; t == nil {
		if reservation {
			c.Required(at, "Total required")
		}
	} else if t.CurrencyCode != "" {
		c.Currency(at.Field("Total"), t.CurrencyCode)
		c.Amount(at.Field("Total"), "AmountAfterTax", t.AmountAfterTax, t.CurrencyCode)
	}
}

func (c *checker) roomRate(at alpinebits.Path, r RoomRate) {
	if rules.Blank(r.RatePlanCode) || rules.Blank(r.RoomTypeCode) {
		c.Required(at, "RatePlanCode and RoomTypeCode required")
	}
	for k, rate := range r.Rates {
		rat := at.At("Rate", k)
		if rate.ExpireDateExclusiveInd != nil && rate.ExpireDate == nil {
			c.Invalid(rat, "ExpireDateExclusiveInd requires ExpireDate")
		}
		c.amount(rat.Field("Base"), "AmountAfterTax", rate.Base.AmountAfterTax, rate.Base.CurrencyCode)
	}
}

// amount checks an amount of an ISO 4217 currency, when both are given.
func (c *checker) amount(at alpinebits.Path, name string, amount *alpinebits.Decimal, code string) {
	if amount != nil && code != "" {
		c.Currency(at, code)
		c.Amount(at, name, *amount, code)
	}
}

// timeForm tells an xs:date, an xs:time and an xs:dateTime apart.
func timeForm(s string) string {
	switch {
	case strings.Contains(s, "T"):
		return "date-time"
	case len(s) > 2 && s[2] == ':':
		return "time"
	default:
		return "date"
	}
}

// services checks the Services and returns their ServiceRPHs.
func (c *checker) services(at alpinebits.Path, services []Service) map[string]bool {
	rphs := make(map[string]bool)
	for k, s := range services {
		sat := at.At("Service", k)
		if rules.Blank(s.ID) {
			c.Required(sat, "ID required")
		}
		if rules.Blank(s.ServiceInventoryCode) {
			c.Required(sat, "ServiceInventoryCode required")
		}
		// Several Services may share a ServiceRPH: the occurrences of one
		// service.
		if s.ServiceRPH != "" {
			rphs[s.ServiceRPH] = true
		}
		t := s.ServiceDetails.Total
		c.amount(sat.Field("Total"), "AmountAfterTax", &t.AmountAfterTax, t.CurrencyCode)
		// 2026-10: Start and End take the same form, and a stay counts once
		// (§4.2.2, calculation of AmountAfterTax).
		if ts := s.ServiceDetails.TimeSpan; ts != nil && ts.End != "" && timeForm(ts.Start) != timeForm(ts.End) {
			c.Invalid(sat.Field("TimeSpan"), "Start and End must both be dates, times or date-times")
		}
		if s.ServicePricingType == ServicePricingTypePerStay && s.Quantity != 1 {
			c.Invalid(sat, "Quantity must be 1 for ServicePricingType Per stay")
		}
	}
	return rphs
}

func (c *checker) commission(at alpinebits.Path, m Commission) {
	a := m.CommissionPayableAmount
	switch {
	case m.Percent == nil && a == nil:
		c.Required(at, "Percent or CommissionPayableAmount required")
	case m.Percent != nil && a != nil:
		c.Invalid(at, "Percent and CommissionPayableAmount must not be given together")
	case a != nil && (a.Amount == nil || a.CurrencyCode == ""):
		c.Required(at.Field("CommissionPayableAmount"), "Amount and CurrencyCode required")
	case a != nil:
		c.Currency(at.Field("CommissionPayableAmount"), a.CurrencyCode)
		c.Amount(at.Field("CommissionPayableAmount"), "Amount", *a.Amount, a.CurrencyCode)
	}
}

// Adults are counted by at most one GuestCount without Age, and the guests
// must not total zero (§4.2.2).
func (c *checker) guestCounts(at alpinebits.Path, counts *[]GuestCount) {
	if counts == nil {
		c.Required(at, "GuestCounts required")
		return
	}
	adults, guests := 0, 0
	for k, g := range *counts {
		guests += g.Count
		if g.Age != nil {
			continue
		}
		if adults++; adults > 1 {
			c.Invalid(at.At("GuestCount", k), "at most one adult GuestCount allowed")
		}
	}
	if guests == 0 {
		c.Required(at, "GuestCount with at least one guest required")
	}
}

func (c *checker) timeSpan(at alpinebits.Path, status ResStatus, ts TimeSpan) {
	if ts.Duration != nil || ts.StartDateWindow != nil {
		if !isReservation(status) {
			c.window(at, ts)
			return
		}
		c.Invalid(at, "Duration and StartDateWindow not allowed in a reservation")
	}
	switch {
	case ts.Start == nil || ts.End == nil:
		c.Required(at, "Start and End required")
	case !ts.End.After(*ts.Start):
		c.InvalidDateRange(at, "End must be after Start")
	}
}

func (c *checker) window(at alpinebits.Path, ts TimeSpan) {
	if ts.Start != nil || ts.End != nil {
		c.Invalid(at, "Start and End not allowed with Duration and StartDateWindow")
	}
	w := ts.StartDateWindow
	switch {
	case ts.Duration == nil:
		c.Required(at, "Duration required with StartDateWindow")
	case w == nil:
		c.Required(at, "StartDateWindow required with Duration")
	}
	if w != nil && !w.LatestDate.After(w.EarliestDate) {
		c.InvalidDateRange(at.Field("StartDateWindow"), "LatestDate must be after EarliestDate")
	}
}

func (c *checker) paymentCard(at alpinebits.Path, p PaymentCard) {
	// The XSD allows a single letter. 2026-10: CardCode is optional.
	if p.CardCode != "" && len(p.CardCode) != 2 {
		c.Invalid(at, "CardCode must be two letters")
	}
	n := p.CardNumber
	nat := at.Field("CardNumber")
	encrypted := n.EncryptedValue != "" || n.EncryptionMethod != ""
	switch {
	case n.PlainText == "" && !encrypted && n.Token == "":
		// 2026-10: or a Token.
		c.Required(nat, "PlainText, EncryptedValue and EncryptionMethod, or Token required")
	case n.PlainText != "" && encrypted:
		c.Invalid(nat, "PlainText not allowed with EncryptedValue and EncryptionMethod")
	case encrypted && (n.EncryptedValue == "" || n.EncryptionMethod == ""):
		c.Required(nat, "EncryptedValue and EncryptionMethod must be given together")
	}
}

func (c *checker) resGlobalInfo(at alpinebits.Path, status ResStatus, g ResGlobalInfo) {
	seen := make(map[CommentName]bool)
	for k, m := range rules.Items(g.Comments) {
		mat := at.At("Comment", k)
		if seen[m.Name] {
			c.Invalid(mat, "duplicate Name %s", m.Name)
		}
		seen[m.Name] = true
		switch m.Name {
		case CommentNameIncludedServices:
			if len(m.ListItems) == 0 {
				c.Required(mat, "included services require ListItem")
			}
		case CommentNameCustomerComment:
			if m.Text == "" {
				c.Required(mat, "customer comment requires Text")
			}
		case CommentNameAdditionalInfo:
			if m.Text == "" {
				c.Required(mat, "additional info requires Text")
			}
		}
		c.listItems(mat, m.ListItems)
	}
	// RequestCodes are distinct per CodeContext, and ALPINEBITS only defines
	// PETS.
	codes := make(map[[2]string]bool)
	for k, r := range rules.Items(g.SpecialRequests) {
		rat := at.At("SpecialRequest", k)
		switch {
		case rules.Blank(r.RequestCode):
			c.Required(rat, "RequestCode required")
		case r.CodeContext == CodeContextAlpineBits && r.RequestCode != RequestCodePets:
			c.Invalid(rat, "RequestCode %s not defined in CodeContext ALPINEBITS", r.RequestCode)
		}
		code := [2]string{string(r.CodeContext), r.RequestCode}
		if codes[code] {
			c.Invalid(rat, "duplicate RequestCode %s in CodeContext %s", r.RequestCode, r.CodeContext)
		}
		codes[code] = true
	}
	// Deposits, whose attributes §4.2.2 requires although the XSD does not.
	for k, p := range rules.Items(g.DepositPayments) {
		c.guaranteePayment(at.At("GuaranteePayment", k), p)
	}
	// 2026-10: the original booking ID.
	for k, id := range rules.Items(g.HotelReservationIDs) {
		if id.ResIDType == 14 && rules.Blank(id.ResIDValue) {
			c.Required(at.At("HotelReservationID", k), "ResID_Value required with ResID_Type 14")
		}
	}
	if g.CancelPenalty != nil && status == ResStatusRequested {
		c.Invalid(at, "CancelPenalties not allowed in a quote request")
	}
	c.Hotel(at.Field("BasicPropertyInfo"), g.BasicPropertyInfo.HotelCode)
}

func (c *checker) listItems(at alpinebits.Path, list []ListItem) {
	type key struct {
		item     int
		language string
	}
	seen := make(map[key]bool)
	for l, item := range list {
		k := key{item.ListItem, item.Language}
		if seen[k] {
			c.Invalid(at.At("ListItem", l), "duplicate ListItem %d for Language %s", item.ListItem, item.Language)
		}
		seen[k] = true
	}
}

func (c *checker) guaranteePayment(at alpinebits.Path, p GuaranteePayment) {
	if rules.Blank(p.GuaranteeCode) {
		c.Required(at, "GuaranteeCode required")
	}
	if rules.Blank(p.Start) {
		c.Required(at, "Start required")
	}
	if a := p.AmountPercent; a == nil || a.CurrencyCode == "" {
		c.Required(at, "AmountPercent with Amount and CurrencyCode required")
	} else {
		c.amount(at.Field("AmountPercent"), "Amount", &a.Amount, a.CurrencyCode)
	}
	a := p.AcceptedPayment
	aat := at.Field("AcceptedPayment")
	// 2026-10: the XSD enforces the values of GuaranteeTypeCode and
	// PaymentTransactionTypeCode, which adds refund.
	if a.GuaranteeTypeCode == "" {
		c.Required(aat, "GuaranteeTypeCode required")
	}
	if rules.Blank(a.GuaranteeID) {
		c.Required(aat, "GuaranteeID required")
	}
	n := 0
	for _, given := range []bool{a.PaymentCard != nil, a.BankAcct != nil, a.Voucher != nil} {
		if given {
			n++
		}
	}
	if n != 1 {
		c.Invalid(aat, "exactly one of PaymentCard, BankAcct and Voucher required")
	}
}
