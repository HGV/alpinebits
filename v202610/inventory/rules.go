package inventory

import (
	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rules"
)

// Validate checks the rules of AlpineBits 2026-10 §4.3.1 that the XSD cannot
// express. Rooms are only checked if CapUseRooms is negotiated: otherwise the
// server must discard them silently.
func Validate(rq HotelDescriptiveContentNotifRQ, caps []alpinebits.Capability) alpinebits.Errors {
	c := checker{rules.New(caps)}
	content := rq.HotelDescriptiveContent
	c.Hotel("", content.HotelCode)
	if content.FacilityInfo == nil {
		c.Required("", "FacilityInfo required")
		return c.Errors()
	}
	rooms := content.FacilityInfo.GuestRooms
	blankCode := false
	for i, g := range rooms {
		if rules.Blank(g.Code) {
			c.Required(alpinebits.At("GuestRoom", i), "Code required")
			blankCode = true
		}
	}
	// Categories and rooms are told apart by Code.
	if blankCode {
		return c.Errors()
	}

	useRooms := c.Has(CapUseRooms)
	roomIDs := make(map[string]map[string]bool)
	for i, g := range rooms {
		at := alpinebits.At("GuestRoom", i)
		ids, ok := roomIDs[g.Code]
		if !ok {
			roomIDs[g.Code] = make(map[string]bool)
			c.category(at, g)
		} else if useRooms {
			c.room(at, g, ids)
		}
	}
	return c.Errors()
}

// The checks are those of 2024-10 apart from the ones marked 2026-10;
// TestRulesMatch202410 keeps them in step.
type checker struct{ *rules.Checker }

func (c *checker) category(at alpinebits.Path, g GuestRoom) {
	if g.MinOccupancy == 0 {
		c.Required(at, "MinOccupancy required")
	}
	if g.MaxOccupancy == 0 {
		c.Required(at, "MaxOccupancy required")
	}
	if g.MaxChildOccupancy != 0 {
		c.Need(at, CapOccupancyChildren, "MaxChildOccupancy")
	}
	t := g.TypeRoom
	if t.StandardOccupancy == 0 {
		c.Required(at.Field("TypeRoom"), "StandardOccupancy required")
	}
	if t.RoomClassificationCode == 0 {
		c.Required(at.Field("TypeRoom"), "RoomClassificationCode required")
	}

	// The cost of a stay needs MinOccupancy <= StandardOccupancy <=
	// MaxOccupancy (§4.4.2); absent values are 0.
	minOcc, std, maxOcc := g.MinOccupancy, t.StandardOccupancy, g.MaxOccupancy
	if maxOcc != 0 && minOcc > maxOcc {
		c.Invalid(at, "MinOccupancy exceeds MaxOccupancy")
	}
	if std != 0 && std < minOcc {
		c.Invalid(at, "StandardOccupancy below MinOccupancy")
	}
	if maxOcc != 0 && std > maxOcc {
		c.Invalid(at, "StandardOccupancy exceeds MaxOccupancy")
	}
	if maxOcc != 0 && g.MaxChildOccupancy > maxOcc {
		c.Invalid(at, "MaxChildOccupancy exceeds MaxOccupancy")
	}
	// 2026-10: StandardOccupancy <= MaxAdultOccupancy <= MaxOccupancy.
	if adults := g.MaxAdultOccupancy; adults != 0 {
		if std != 0 && adults < std {
			c.Invalid(at, "MaxAdultOccupancy below StandardOccupancy")
		}
		if maxOcc != 0 && adults > maxOcc {
			c.Invalid(at, "MaxAdultOccupancy exceeds MaxOccupancy")
		}
	}

	for k, a := range rules.Items(g.Amenities) {
		if rules.Blank(a.RoomAmenityCode) {
			c.Required(at.At("Amenity", k), "RoomAmenityCode required")
		}
	}
	c.multimediaDescriptions(at, rules.Items(g.MultimediaDescriptions))
}

func (c *checker) multimediaDescriptions(at alpinebits.Path, mds []MultimediaDescription) {
	seen := make(map[InfoCode]bool)
	for k, md := range mds {
		mat := at.At("MultimediaDescription", k)
		switch md.InfoCode {
		case 0:
			c.Required(mat, "InfoCode required")
			continue
		case InfoCodeLongName, InfoCodeDescription:
			if md.TextItem == nil {
				c.Required(mat, "TextItems required")
			} else {
				c.Texts(mat, "Description", md.TextItem.Descriptions)
			}
		case InfoCodePictures:
			if md.ImageItems == nil {
				c.Required(mat, "ImageItems required")
			}
			for j, img := range rules.Items(md.ImageItems) {
				c.Texts(mat.At("ImageItem", j), "Description", img.Descriptions)
			}
		}
		if seen[md.InfoCode] {
			c.Invalid(mat, "duplicate InfoCode %d", md.InfoCode)
		}
		seen[md.InfoCode] = true
	}
	if !seen[InfoCodeLongName] {
		c.Required(at, "MultimediaDescription with InfoCode 25 required")
	}
}

// room checks a GuestRoom listing a room of a category, whose ids are the
// RoomIDs listed before.
func (c *checker) room(at alpinebits.Path, g GuestRoom, ids map[string]bool) {
	id := g.TypeRoom.RoomID
	if g != (GuestRoom{Code: g.Code, TypeRoom: TypeRoom{RoomID: id}}) {
		c.Invalid(at, "room of category %s allows only Code and RoomID", g.Code)
	}
	switch {
	case rules.Blank(id):
		c.Required(at.Field("TypeRoom"), "RoomID required")
	case ids[id]:
		c.Invalid(at.Field("TypeRoom"), "duplicate RoomID %s", id)
	}
	ids[id] = true
}

// Partition separates the GuestRooms defining room categories, the first with
// each Code, from those listing the rooms of a category.
func Partition(rooms []GuestRoom) (categories, individualRooms []GuestRoom) {
	seen := make(map[string]bool)
	for _, room := range rooms {
		if seen[room.Code] {
			individualRooms = append(individualRooms, room)
		} else {
			categories = append(categories, room)
			seen[room.Code] = true
		}
	}
	return
}

// Renames returns the room categories to rename, old Code to new Code, given
// the Codes on record: those whose Code is not on record but whose ID is
// (§4.3.1). Apply them before replacing the categories.
func Renames(rq HotelDescriptiveContentNotifRQ, categories map[string]struct{}) map[string]string {
	renames := make(map[string]string)
	headings, _ := Partition(rq.GuestRooms())
	for _, g := range headings {
		if _, ok := categories[g.Code]; ok || g.OldCode == "" {
			continue
		}
		if _, ok := categories[g.OldCode]; !ok {
			continue
		}
		if _, done := renames[g.OldCode]; !done {
			renames[g.OldCode] = g.Code
		}
	}
	return renames
}
