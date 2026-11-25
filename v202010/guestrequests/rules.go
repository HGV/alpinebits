package guestrequests

import (
	"reflect"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rules"
	guestrequests201810 "github.com/HGV/alpinebits/v201810/guestrequests"
)

// ValidateReadRQ checks the rules of AlpineBits 2020-10 §4.2.1 that the XSD
// cannot express. They are those of 2018-10.
func ValidateReadRQ(rq ReadRQ, caps []alpinebits.Capability) alpinebits.Errors {
	return guestrequests201810.ValidateReadRQ(rq, caps)
}

// ValidateNotifReportRQ checks the rules of AlpineBits 2020-10 §4.2.3 that the
// XSD cannot express. They are those of 2018-10; see
// guestrequests201810.ValidateNotifReportRQ.
func ValidateNotifReportRQ(rq NotifReportRQ, caps []alpinebits.Capability) alpinebits.Errors {
	return guestrequests201810.ValidateNotifReportRQ(rq, caps)
}

// ValidateResRetrieveRS checks the rules of AlpineBits 2020-10 §4.2.2 that the
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

// The checks are those of 2018-10 apart from the ones marked 2020-10;
// TestRulesMatch201810 keeps them in step.
type checker struct{ *rules.Checker }

// isReservation reports whether s is a booking reservation or modification,
// which §4.2.2 treats alike.
func isReservation(s ResStatus) bool {
	return s == ResStatusReserved || s == ResStatusModify
}

func onlyTimeSpan(s RoomStay) bool {
	return s.RoomType == nil && s.RatePlan == nil && s.GuestCounts == nil && s.Guarantee == nil && s.Total == nil
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
	if status != ResStatusCancelled {
		if res.RoomStays == nil {
			c.Required(at, "RoomStays required")
		}
		if res.ResGuest == nil {
			c.Required(at, "ResGuests required")
		}
		if res.ResGlobalInfo == nil {
			c.Required(at, "ResGlobalInfo required")
		}
	}
	c.roomStays(at, status, rules.Items(res.RoomStays))
	if g := res.ResGlobalInfo; g != nil {
		c.resGlobalInfo(at.Field("ResGlobalInfo"), status, *g)
	}
}

func (c *checker) roomStays(at alpinebits.Path, status ResStatus, stays []RoomStay) {
	alternative := -1
	if n := len(stays); n > 1 && !isReservation(status) && onlyTimeSpan(stays[n-1]) {
		alternative = n - 1
	}
	for j, stay := range stays {
		sat := at.At("RoomStay", j)
		if j == alternative {
			c.timeSpan(sat.Field("TimeSpan"), status, stay.TimeSpan)
			if reflect.DeepEqual(stay.TimeSpan, stays[0].TimeSpan) {
				c.Invalid(sat, "alternative TimeSpan must differ from RoomStay[0]")
			}
			continue
		}
		c.roomStay(sat, status, stay)
		if j > 0 && !reflect.DeepEqual(stay.TimeSpan, stays[0].TimeSpan) {
			c.Invalid(sat, "TimeSpan differs from RoomStay[0]")
		}
	}
}

func (c *checker) roomStay(at alpinebits.Path, status ResStatus, stay RoomStay) {
	reservation := isReservation(status)
	// 2020-10: RoomType is mandatory only for reservations.
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
		// 2020-10: any ISO 4217 currency.
		c.Currency(at.Field("Total"), t.CurrencyCode)
		c.Amount(at.Field("Total"), "AmountAfterTax", t.AmountAfterTax, t.CurrencyCode)
	}
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
		// 2020-10: any ISO 4217 currency.
		c.Currency(at.Field("CommissionPayableAmount"), a.CurrencyCode)
		c.Amount(at.Field("CommissionPayableAmount"), "Amount", *a.Amount, a.CurrencyCode)
	}
}

// 2020-10: adults are counted by at most one GuestCount without Age, and the
// guests must not total zero (§4.2.2).
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
	// The XSD allows a single letter.
	if len(p.CardCode) != 2 {
		c.Invalid(at, "CardCode must be two letters")
	}
	n := p.CardNumber
	nat := at.Field("CardNumber")
	encrypted := n.EncryptedValue != "" || n.EncryptionMethod != ""
	switch {
	case n.PlainText == "" && !encrypted:
		c.Required(nat, "PlainText or EncryptedValue and EncryptionMethod required")
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
		}
		c.listItems(mat, m.ListItems)
	}
	if g.CancelPenalty != nil && status == ResStatusRequested {
		c.Invalid(at, "CancelPenalties not allowed in a quote request")
	}
	c.Hotel(at.Field("BasicPropertyInfo"), g.BasicPropertyInfo.HotelCode, g.BasicPropertyInfo.HotelName)
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
