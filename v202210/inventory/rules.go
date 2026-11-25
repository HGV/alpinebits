package inventory

import (
	"github.com/HGV/alpinebits"
	inventory202010 "github.com/HGV/alpinebits/v202010/inventory"
)

// Validate checks the rules of AlpineBits 2022-10 §4.4.1 that the XSD cannot
// express; they are those of 2020-10, see inventory202010.Validate.
func Validate(rq HotelDescriptiveContentNotifRQ, caps []alpinebits.Capability) alpinebits.Errors {
	return inventory202010.Validate(rq, caps)
}

// Partition separates the GuestRooms defining room categories, the first with
// each Code, from those listing the rooms of a category.
func Partition(rooms []GuestRoom) (categories, individualRooms []GuestRoom) {
	return inventory202010.Partition(rooms)
}

// Renames returns the room categories to rename, old Code to new Code, given
// the Codes on record; see inventory202010.Renames.
func Renames(rq HotelDescriptiveContentNotifRQ, categories map[string]struct{}) map[string]string {
	return inventory202010.Renames(rq, categories)
}
