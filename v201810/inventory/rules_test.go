package inventory

import (
	"encoding/xml"
	"os"
	"reflect"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
)

var allCaps = []alpinebits.Capability{CapUseRooms, CapOccupancyChildren}

// The official AlpineBits 2018-10 Inventory/Basic (push) samples: category DZ
// with all three MultimediaDescriptions and five rooms; a pitch without rooms;
// the removal of all categories.
const (
	notifSample    = "Inventory-OTA_HotelDescriptiveContentNotifRQ.xml"
	roomTypeSample = "Inventory-OTA_HotelDescriptiveContentNotifRQ-with-roomtype.xml"
	deleteSample   = "Inventory-OTA_HotelDescriptiveContentNotifRQ-delete-all.xml"
)

func sample(t *testing.T, file string) HotelDescriptiveContentNotifRQ {
	t.Helper()
	b, err := os.ReadFile("testdata/" + file)
	if err != nil {
		t.Fatal(err)
	}
	var rq HotelDescriptiveContentNotifRQ
	if err := xml.Unmarshal(b, &rq); err != nil {
		t.Fatal(err)
	}
	return rq
}

func TestValidateSamples(t *testing.T) {
	for _, file := range []string{notifSample, roomTypeSample, deleteSample} {
		t.Run(file, func(t *testing.T) {
			rulestest.Expect(t, Validate(sample(t, file), allCaps), nil)
		})
	}
}

func TestValidateCapabilities(t *testing.T) {
	rq := sample(t, notifSample)
	rulestest.Expect(t, Validate(rq, allCaps), nil)
	rulestest.Expect(t, Validate(rq, nil), []string{
		"GuestRoom[0]: MaxChildOccupancy not supported",
	})
}

// Without CapUseRooms the server must discard the rooms silently (§4.4.1).
func TestValidateIgnoresRoomsWithoutUseRooms(t *testing.T) {
	rq := sample(t, notifSample)
	rooms := &rq.HotelDescriptiveContent.FacilityInfo.GuestRooms
	*rooms = append(*rooms, GuestRoom{Code: "DZ"})
	rulestest.Expect(t, Validate(rq, []alpinebits.Capability{CapOccupancyChildren}), nil)
	rulestest.Expect(t, Validate(rq, allCaps), []string{"GuestRoom[6].TypeRoom: RoomID required"})
}

func TestValidateMessage(t *testing.T) {
	pitch := sample(t, roomTypeSample).GuestRooms()[0]
	tests := []struct {
		name   string
		change func(c *HotelDescriptiveContent)
		want   []string
	}{
		{"HotelName only", func(c *HotelDescriptiveContent) { c.HotelCode = "" }, nil},
		{"neither HotelCode nor HotelName", func(c *HotelDescriptiveContent) { c.HotelCode, c.HotelName = " ", "" }, []string{
			"missing HotelCode or HotelName",
		}},
		{"no FacilityInfo", func(c *HotelDescriptiveContent) { c.FacilityInfo = nil }, []string{
			"FacilityInfo required",
		}},
		{"blank Code", func(c *HotelDescriptiveContent) {
			c.FacilityInfo.GuestRooms[0].MinOccupancy = 0
			c.FacilityInfo.GuestRooms[2].Code = " "
		}, []string{
			"GuestRoom[2]: Code required",
		}},
		{"RoomIDs are per category", func(c *HotelDescriptiveContent) {
			c.FacilityInfo.GuestRooms = append(c.FacilityInfo.GuestRooms, pitch, GuestRoom{Code: pitch.Code, TypeRoom: TypeRoom{RoomID: "101"}})
		}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := sample(t, notifSample)
			tt.change(&rq.HotelDescriptiveContent)
			rulestest.Expect(t, Validate(rq, allCaps), tt.want)
		})
	}
}

func TestValidateCategory(t *testing.T) {
	text := func(f alpinebits.TextFormat, lang string) alpinebits.Text {
		return alpinebits.Text{TextFormat: f, Language: lang, Value: "text"}
	}
	tests := []struct {
		name   string
		change func(g *GuestRoom)
		want   []string
	}{
		{"missing attributes", func(g *GuestRoom) {
			g.MinOccupancy, g.MaxOccupancy, g.TypeRoom = 0, 0, TypeRoom{}
		}, []string{
			"GuestRoom[0]: MinOccupancy required",
			"GuestRoom[0]: MaxOccupancy required",
			"GuestRoom[0].TypeRoom: StandardOccupancy required",
			"GuestRoom[0].TypeRoom: RoomClassificationCode required",
		}},
		{"occupancies above MaxOccupancy", func(g *GuestRoom) {
			g.MinOccupancy, g.TypeRoom.StandardOccupancy, g.MaxChildOccupancy = 3, 3, 3
		}, []string{
			"GuestRoom[0]: MinOccupancy exceeds MaxOccupancy",
			"GuestRoom[0]: StandardOccupancy exceeds MaxOccupancy",
			"GuestRoom[0]: MaxChildOccupancy exceeds MaxOccupancy",
		}},
		{"StandardOccupancy below MinOccupancy", func(g *GuestRoom) {
			g.MinOccupancy, g.TypeRoom.StandardOccupancy = 2, 1
		}, []string{
			"GuestRoom[0]: StandardOccupancy below MinOccupancy",
		}},
		// RoomType should only be used with its RoomClassificationCode.
		{"RoomType with another RoomClassificationCode", func(g *GuestRoom) {
			g.TypeRoom.RoomClassificationCode, g.TypeRoom.RoomType = 13, RoomTypeRoom
		}, nil},
		{"Amenity without RoomAmenityCode", func(g *GuestRoom) {
			*g.Amenities = append(*g.Amenities, Amenity{})
		}, []string{
			"GuestRoom[0].Amenity[1]: RoomAmenityCode required",
		}},
		{"no long name", func(g *GuestRoom) {
			*g.MultimediaDescriptions = (*g.MultimediaDescriptions)[1:]
		}, []string{
			"GuestRoom[0]: MultimediaDescription with InfoCode 25 required",
		}},
		{"no MultimediaDescriptions", func(g *GuestRoom) { g.MultimediaDescriptions = nil }, []string{
			"GuestRoom[0]: MultimediaDescription with InfoCode 25 required",
		}},
		{"MultimediaDescriptions", func(g *GuestRoom) {
			mds := *g.MultimediaDescriptions
			mds[0].TextItem, mds[0].ImageItems = nil, mds[2].ImageItems
			*g.MultimediaDescriptions = append(mds,
				MultimediaDescription{TextItem: mds[1].TextItem},
				MultimediaDescription{InfoCode: InfoCodePictures},
				mds[1],
			)
		}, []string{
			"GuestRoom[0].MultimediaDescription[0]: TextItems required",
			"GuestRoom[0].MultimediaDescription[3]: InfoCode required",
			"GuestRoom[0].MultimediaDescription[4]: ImageItems required",
			"GuestRoom[0].MultimediaDescription[4]: duplicate InfoCode 23",
			"GuestRoom[0].MultimediaDescription[5]: duplicate InfoCode 1",
		}},
		{"Descriptions", func(g *GuestRoom) {
			ti := g.MultimediaDescription(InfoCodeDescription).TextItem
			ti.Descriptions = append(ti.Descriptions,
				text(alpinebits.TextFormatPlainText, "en"),
				text(alpinebits.TextFormatHTML, "fr"),
				text(alpinebits.TextFormatHTML, "de"),
			)
			img := &(*g.MultimediaDescription(InfoCodePictures).ImageItems)[0]
			img.Descriptions = append(img.Descriptions, text(alpinebits.TextFormatPlainText, "de"))
		}, []string{
			"GuestRoom[0].MultimediaDescription[1].Description[3]: duplicate Description for Language en and TextFormat PlainText",
			"GuestRoom[0].MultimediaDescription[1].Description[4]: HTML Description requires a PlainText Description in the same Language",
			"GuestRoom[0].MultimediaDescription[2].ImageItem[0].Description[3]: duplicate Description for Language de and TextFormat PlainText",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := sample(t, notifSample)
			tt.change(&rq.HotelDescriptiveContent.FacilityInfo.GuestRooms[0])
			rulestest.Expect(t, Validate(rq, allCaps), tt.want)
		})
	}
}

func TestValidateRooms(t *testing.T) {
	tests := []struct {
		name   string
		change func(rooms []GuestRoom)
		want   []string
	}{
		{"room with category data", func(rooms []GuestRoom) {
			rooms[1].MaxOccupancy = 2
			rooms[2].TypeRoom.StandardOccupancy = 2
			rooms[3].MultimediaDescriptions = rooms[0].MultimediaDescriptions
		}, []string{
			"GuestRoom[1]: room of category DZ allows only Code and RoomID",
			"GuestRoom[2]: room of category DZ allows only Code and RoomID",
			"GuestRoom[3]: room of category DZ allows only Code and RoomID",
		}},
		{"RoomID", func(rooms []GuestRoom) {
			rooms[1].TypeRoom.RoomID = " "
			rooms[3].TypeRoom.RoomID = rooms[2].TypeRoom.RoomID
		}, []string{
			"GuestRoom[1].TypeRoom: RoomID required",
			"GuestRoom[3].TypeRoom: duplicate RoomID 102",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := sample(t, notifSample)
			tt.change(rq.GuestRooms())
			rulestest.Expect(t, Validate(rq, allCaps), tt.want)
		})
	}
}

func TestPartition(t *testing.T) {
	room := func(code, id string) GuestRoom { return GuestRoom{Code: code, TypeRoom: TypeRoom{RoomID: id}} }
	categories, rooms := Partition([]GuestRoom{room("DZ", ""), room("EZ", ""), room("DZ", "101"), room("EZ", "201")})
	if want := []GuestRoom{room("DZ", ""), room("EZ", "")}; !reflect.DeepEqual(categories, want) {
		t.Errorf("categories %v, want %v", categories, want)
	}
	if want := []GuestRoom{room("DZ", "101"), room("EZ", "201")}; !reflect.DeepEqual(rooms, want) {
		t.Errorf("rooms %v, want %v", rooms, want)
	}
}

func TestMinFull(t *testing.T) {
	g := sample(t, notifSample).GuestRooms()[0]
	if got := g.MinFull(); got != 1 {
		t.Errorf("MinFull() = %d, want 1", got)
	}
	g.MaxChildOccupancy = 0
	if got := g.MinFull(); got != 2 {
		t.Errorf("without MaxChildOccupancy MinFull() = %d, want 2", got)
	}
}

func TestRenames(t *testing.T) {
	category := func(code, oldCode string) GuestRoom { return GuestRoom{Code: code, OldCode: oldCode} }
	rq := HotelDescriptiveContentNotifRQ{HotelDescriptiveContent: HotelDescriptiveContent{FacilityInfo: &FacilityInfo{GuestRooms: []GuestRoom{
		category("single", "EZ"),
		category("double", "DZ"), // already renamed
		category("suite", "SU"),  // SU unknown: a new category
		category("triple", "EZ"), // EZ is renamed already
		category("single", "X"),  // a room of single
	}}}}
	got := Renames(rq, map[string]struct{}{"EZ": {}, "DZ": {}, "double": {}})
	if want := map[string]string{"EZ": "single"}; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
