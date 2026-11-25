package inventory

import (
	"github.com/HGV/alpinebits"
	inventory202210 "github.com/HGV/alpinebits/v202210/inventory"
)

// Validate checks the rules of AlpineBits 2024-10 §4.3.1 that the XSD cannot
// express; they are those of 2022-10, see inventory202210.Validate.
func Validate(rq HotelDescriptiveContentNotifRQ, caps []alpinebits.Capability) alpinebits.Errors {
	return inventory202210.Validate(rq, caps)
}

// Partition separates the GuestRooms defining room categories, the first with
// each Code, from those listing the rooms of a category.
func Partition(rooms []GuestRoom) (categories, individualRooms []GuestRoom) {
	return inventory202210.Partition(rooms)
}

// Renames returns the room categories to rename, old Code to new Code, given
// the Codes on record; see inventory202210.Renames.
func Renames(rq HotelDescriptiveContentNotifRQ, categories map[string]struct{}) map[string]string {
	return inventory202210.Renames(rq, categories)
}
