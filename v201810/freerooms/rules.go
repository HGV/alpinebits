package freerooms

import (
	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rules"
)

// Validate checks the rules of AlpineBits 2018-10 §4.1 that the XSD cannot
// express.
func Validate(rq HotelAvailNotifRQ, caps []alpinebits.Capability) alpinebits.Errors {
	c := rules.New(caps)
	c.Hotel("", rq.AvailStatusMessages.HotelCode)
	if !rq.IsCompleteSet() {
		c.Need("", CapabilityAcceptDeltas, "deltas")
	}
	if rq.IsReset() {
		return c.Errors()
	}

	threshold := c.Has(CapabilityAcceptBookingThreshold)
	msgs := rq.AvailStatusMessages.AvailStatusMessages
	var rooms, categories bool
	var ranged []int
	for i, m := range msgs {
		at := alpinebits.At("AvailStatusMessage", i)
		s := m.StatusApplicationControl
		isRoom := s != nil && s.isRoom()
		switch {
		case m.BookingLimit == nil:
			c.Required(at, "BookingLimit required")
		case isRoom && *m.BookingLimit > 1:
			c.Invalid(at, "BookingLimit must be 0 or 1 for a room")
		}
		if m.BookingThreshold == nil {
			if threshold {
				c.Required(at, "BookingThreshold required")
			}
		} else {
			c.Need(at, CapabilityAcceptBookingThreshold, "BookingThreshold")
			if m.BookingLimit != nil && *m.BookingThreshold > *m.BookingLimit {
				c.Invalid(at, "BookingThreshold exceeds BookingLimit")
			}
		}

		if s == nil {
			c.Required(at, "StatusApplicationControl required")
			continue
		}
		sat := at.Field("StatusApplicationControl")
		hasCode := !rules.Blank(s.InvTypeCode)
		if !hasCode {
			c.Required(sat, "InvTypeCode required")
		}
		if s.End.Before(s.Start) {
			c.InvalidDateRange(sat, "End before Start")
		} else if hasCode {
			ranged = append(ranged, i)
		}
		if isRoom {
			rooms = true
		} else {
			categories = true
		}
	}

	switch {
	case rooms && categories:
		c.Invalid("", "rooms and categories must not be mixed")
	case rooms && !c.Has(CapabilityAcceptRooms):
		c.Invalid("", "rooms not supported")
	case categories && !c.Has(CapabilityAcceptCategories):
		c.Invalid("", "categories not supported")
	}
	rules.Overlaps(c, "", "AvailStatusMessage", msgs, ranged, func(m AvailStatusMessage) [2]string {
		return [2]string{m.StatusApplicationControl.InvTypeCode, m.StatusApplicationControl.InvCode}
	})
	return c.Errors()
}

// HotelInventory lists a hotel's room categories and rooms.
type HotelInventory struct {
	Categories map[string]struct{}            // by InvTypeCode
	Rooms      map[string]map[string]struct{} // by InvTypeCode, then InvCode
}

// CheckInventory reports AvailStatusMessages for categories or rooms not in
// hotel; codes match case sensitively (AlpineBits 2018-10 §4.1.1). A nil map
// skips that check. Handlers call it after Validate has passed.
func CheckInventory(rq HotelAvailNotifRQ, hotel HotelInventory) alpinebits.Errors {
	var errs alpinebits.Errors
	for i, m := range rq.AvailStatusMessages.AvailStatusMessages {
		s := m.StatusApplicationControl
		if s == nil {
			continue
		}
		at := alpinebits.At("AvailStatusMessage", i).Field("StatusApplicationControl")
		switch {
		case s.isRoom() && hotel.Rooms != nil:
			if _, ok := hotel.Rooms[s.InvTypeCode][s.InvCode]; !ok {
				errs.Addf(alpinebits.ErrCodeInvalidValue, at, "unknown room %s/%s", s.InvTypeCode, s.InvCode)
			}
		case !s.isRoom() && hotel.Categories != nil:
			if _, ok := hotel.Categories[s.InvTypeCode]; !ok {
				errs.Addf(alpinebits.ErrCodeInvalidValue, at, "unknown category %s", s.InvTypeCode)
			}
		}
	}
	return errs
}
