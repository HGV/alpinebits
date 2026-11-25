package freerooms

import (
	"github.com/HGV/alpinebits"
	freerooms202410 "github.com/HGV/alpinebits/v202410/freerooms"
)

// Validate checks the rules of AlpineBits 2026-10 §4.1 that the XSD cannot
// express. They are those of 2024-10.
func Validate(rq HotelInvCountNotifRQ, caps []alpinebits.Capability) alpinebits.Errors {
	return freerooms202410.Validate(rq, caps)
}

// CheckInventory reports Inventory elements for categories or rooms not in
// hotel, and categories whose Counts sum to more than their number of rooms.
func CheckInventory(rq HotelInvCountNotifRQ, hotel HotelInventory) alpinebits.Errors {
	return freerooms202410.CheckInventory(rq, hotel)
}
