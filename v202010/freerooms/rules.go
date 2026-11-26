package freerooms

import (
	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rules"
)

// Validate checks the rules of AlpineBits 2020-10 §4.1 that the XSD cannot
// express.
func Validate(rq HotelInvCountNotifRQ, caps []alpinebits.Capability) alpinebits.Errors {
	c := checker{rules.New(caps)}
	c.Hotel("", rq.Inventories.HotelCode)
	if !rq.IsCompleteSet() {
		c.Need("", CapabilityAcceptDeltas, "deltas")
	}
	if rq.IsReset() {
		return c.Errors()
	}

	invs := rq.Inventories.Inventories
	var ranged, closing []int
	var rooms, categories, others bool
	for i, inv := range invs {
		at := alpinebits.At("Inventory", i)
		if inv.IsClosingSeason() {
			if c.closingSeason(at, inv, rq.IsCompleteSet(), others) {
				closing = append(closing, i)
			}
			continue
		}
		others = true
		if c.inventory(at, inv) {
			ranged = append(ranged, i)
		}
		if inv.StatusApplicationControl != nil {
			if inv.IsRoom() {
				rooms = true
			} else {
				categories = true
			}
		}
	}
	if rooms && categories {
		c.Invalid("", "rooms and categories must not be mixed")
	}

	// Only periods of the same room or category must not overlap (§4.1.1).
	rules.Overlaps(c.Checker, "", "Inventory", invs, ranged, func(inv Inventory) [2]string {
		return [2]string{inv.StatusApplicationControl.InvTypeCode, inv.StatusApplicationControl.InvCode}
	})
	rules.Overlaps(c.Checker, "", "Inventory", invs, closing, func(Inventory) struct{} { return struct{}{} })
	c.closedButAvailable(invs, closing, ranged)
	return c.Errors()
}

type checker struct{ *rules.Checker }

func (c *checker) dateRange(at alpinebits.Path, s *StatusApplicationControl) bool {
	if s.End.Before(s.Start) {
		c.InvalidDateRange(at, "End before Start")
		return false
	}
	return true
}

// inventory checks the Inventory of a room or category and reports whether
// it takes part in the overlap checks.
func (c *checker) inventory(at alpinebits.Path, inv Inventory) bool {
	ok := false
	if s := inv.StatusApplicationControl; s == nil {
		c.Required(at, "StatusApplicationControl required")
	} else {
		ok = c.dateRange(at, s)
		if rules.Blank(s.InvTypeCode) {
			c.Required(at, "InvTypeCode required")
			ok = false
		}
		if inv.IsRoom() {
			c.Need(at, CapabilityAcceptRooms, "rooms")
			// A room is a single room (§4.1.1).
			if inv.total() > 1 {
				c.Invalid(at, "Counts of a room must sum to at most 1")
			}
		} else {
			c.Need(at, CapabilityAcceptCategories, "categories")
		}
	}
	seen := make(map[CountType]bool)
	for k, n := range rules.Items(inv.InvCounts) {
		nat := at.At("InvCount", k)
		switch n.CountType {
		case CountTypeOutOfOrder:
			c.Need(nat, CapabilityAcceptOutOfOrder, "CountType 6 (out of order)")
		case CountTypeOutOfMarket:
			c.Need(nat, CapabilityAcceptOutOfMarket, "CountType 9 (out of market)")
		}
		if seen[n.CountType] {
			c.Invalid(nat, "duplicate CountType %s", n.CountType)
		}
		seen[n.CountType] = true
	}
	return ok
}

// closingSeason checks a closing season and reports whether it takes part in
// the overlap checks. others reports whether other Inventory elements precede
// it.
func (c *checker) closingSeason(at alpinebits.Path, inv Inventory, completeSet, others bool) bool {
	s := inv.StatusApplicationControl
	c.Need(at, CapabilityAcceptClosingSeasons, "closing seasons")
	if !completeSet {
		c.Invalid(at, "closing seasons only allowed in a CompleteSet")
	}
	if others {
		c.Invalid(at, "closing seasons must precede the other Inventory elements")
	}
	ok := c.dateRange(at, s)
	if s.InvTypeCode != "" || s.InvCode != "" {
		c.Invalid(at, "closing season must have no InvTypeCode and no InvCode")
	}
	if inv.InvCounts != nil {
		c.Invalid(at, "closing season must have no InvCounts")
	}
	return ok
}

// closedButAvailable reports rooms or categories available during a closing
// season. Closing seasons may overlap periods without free rooms (§4.1.1).
func (c *checker) closedButAvailable(invs []Inventory, closing, ranged []int) {
	for _, j := range ranged {
		if !invs[j].isAvailable() {
			continue
		}
		r := invs[j].DateRange()
		for _, i := range closing {
			cr := invs[i].DateRange()
			if !r.Start.After(cr.End) && !cr.Start.After(r.End) {
				c.InvalidDateRange(alpinebits.At("Inventory", j), "available rooms overlap closing season Inventory[%d]", i)
			}
		}
	}
}

// isAvailable reports whether i counts free rooms. §4.1.1 describes out of
// market rooms as "available but not bookable"; out of order rooms are not
// available.
func (i Inventory) isAvailable() bool {
	return i.Count(CountTypeBookable) > 0 || i.Count(CountTypeOutOfMarket) > 0
}

// HotelInventory lists a hotel's room categories and rooms.
type HotelInventory struct {
	Categories map[string]int                 // number of rooms by InvTypeCode
	Rooms      map[string]map[string]struct{} // by InvTypeCode, then InvCode
}

// CheckInventory reports Inventory elements for categories or rooms not in
// hotel, and categories whose Counts sum to more than their number of rooms
// (AlpineBits 2020-10 §4.1.1). A nil map skips its checks.
func CheckInventory(rq HotelInvCountNotifRQ, hotel HotelInventory) alpinebits.Errors {
	var errs alpinebits.Errors
	for i, inv := range rq.Inventories.Inventories {
		s := inv.StatusApplicationControl
		if s == nil || inv.IsClosingSeason() {
			continue
		}
		at := alpinebits.At("Inventory", i)
		switch {
		case inv.IsRoom() && hotel.Rooms != nil:
			if _, ok := hotel.Rooms[s.InvTypeCode][s.InvCode]; !ok {
				errs.Addf(alpinebits.ErrCodeInvalidValue, at, "unknown room %s in category %s", s.InvCode, s.InvTypeCode)
			}
		case !inv.IsRoom() && hotel.Categories != nil:
			n, ok := hotel.Categories[s.InvTypeCode]
			if !ok {
				errs.Addf(alpinebits.ErrCodeInvalidValue, at, "unknown category %s", s.InvTypeCode)
			} else if inv.total() > n {
				errs.Addf(alpinebits.ErrCodeInvalidValue, at, "Counts exceed the %d rooms of category %s", n, s.InvTypeCode)
			}
		}
	}
	return errs
}
