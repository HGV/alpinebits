package freerooms

import (
	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rules"
	freerooms202010 "github.com/HGV/alpinebits/v202010/freerooms"
)

// Validate checks the rules of AlpineBits 2022-10 §4.1 that the XSD cannot
// express: those of 2020-10, and a CompleteSet requires CapabilityAcceptCompleteSet.
func Validate(rq HotelInvCountNotifRQ, caps []alpinebits.Capability) alpinebits.Errors {
	errs := freerooms202010.Validate(rq, caps)
	c := rules.New(caps)
	if rq.IsCompleteSet() {
		c.Need("", CapabilityAcceptCompleteSet, "complete sets")
	}
	return append(errs, c.Errors()...)
}

// CheckInventory reports Inventory elements for categories or rooms not in
// hotel, and categories whose Counts sum to more than their number of rooms.
func CheckInventory(rq HotelInvCountNotifRQ, hotel HotelInventory) alpinebits.Errors {
	return freerooms202010.CheckInventory(rq, hotel)
}
