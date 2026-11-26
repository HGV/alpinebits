package inventory

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
	inventory202410 "github.com/HGV/alpinebits/v202410/inventory"
)

var allCaps = []alpinebits.Capability{CapUseRooms, CapOccupancyChildren}

// The official 2026-10 sample: categories DOUBLE, SUITE and SINGLE with ABR
// amenities and MaxAdultOccupancy, and their rooms.
const basicSample = "testdata/Inventory-Push-OTA_HotelDescriptiveContentNotifRQ-basic.xml"

// 2026-10 kept the 2018-10 delete-all sample.
const deleteSample = "../../v201810/inventory/testdata/Inventory-OTA_HotelDescriptiveContentNotifRQ-delete-all.xml"

func load[T any](t *testing.T, path string) T {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var v T
	if err := xml.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestValidateSamples(t *testing.T) {
	// The official sample gives SUITE a MaxAdultOccupancy below its
	// StandardOccupancy, which §4.3.1 forbids.
	rulestest.Expect(t, Validate(load[HotelDescriptiveContentNotifRQ](t, basicSample), allCaps), []string{
		"GuestRoom[1]: MaxAdultOccupancy below StandardOccupancy",
	})
	rulestest.Expect(t, Validate(load[HotelDescriptiveContentNotifRQ](t, deleteSample), allCaps), nil)
}

func from202410(rq inventory202410.HotelDescriptiveContentNotifRQ) HotelDescriptiveContentNotifRQ {
	c := rq.HotelDescriptiveContent
	out := HotelDescriptiveContentNotifRQ{Version: rq.Version, HotelDescriptiveContent: HotelDescriptiveContent{HotelCode: c.HotelCode, HotelName: c.HotelName}}
	if c.FacilityInfo == nil {
		return out
	}
	rooms := []GuestRoom{}
	for _, g := range c.FacilityInfo.GuestRooms {
		rooms = append(rooms, GuestRoom{
			Code: g.Code, MinOccupancy: g.MinOccupancy, MaxOccupancy: g.MaxOccupancy, MaxChildOccupancy: g.MaxChildOccupancy,
			OldCode: g.OldCode, TypeRoom: g.TypeRoom, Amenities: g.Amenities, MultimediaDescriptions: g.MultimediaDescriptions,
		})
	}
	out.HotelDescriptiveContent.FacilityInfo = &FacilityInfo{GuestRooms: rooms}
	return out
}

// The checks 2026-10 copied from 2024-10 must report what 2024-10 reports.
func TestRulesMatch202410(t *testing.T) {
	type rq = inventory202410.HotelDescriptiveContentNotifRQ
	category := func(r *rq) *inventory202410.GuestRoom { return &r.HotelDescriptiveContent.FacilityInfo.GuestRooms[0] }
	tests := []struct {
		name   string
		caps   []alpinebits.Capability
		change func(r *rq)
	}{
		{"no HotelCode, no FacilityInfo", allCaps, func(r *rq) {
			r.HotelDescriptiveContent.HotelCode, r.HotelDescriptiveContent.FacilityInfo = " ", nil
		}},
		{"blank Code", allCaps, func(r *rq) { r.HotelDescriptiveContent.FacilityInfo.GuestRooms[2].Code = " " }},
		{"MaxChildOccupancy not negotiated", nil, func(*rq) {}},
		{"missing attributes", allCaps, func(r *rq) {
			g := category(r)
			g.MinOccupancy, g.MaxOccupancy, g.TypeRoom = 0, 0, TypeRoom{}
		}},
		{"occupancies above MaxOccupancy", allCaps, func(r *rq) {
			g := category(r)
			g.MinOccupancy, g.TypeRoom.StandardOccupancy, g.MaxChildOccupancy = 3, 3, 3
		}},
		{"StandardOccupancy below MinOccupancy", allCaps, func(r *rq) {
			g := category(r)
			g.MinOccupancy, g.TypeRoom.StandardOccupancy = 2, 1
		}},
		{"blank RoomAmenityCode", allCaps, func(r *rq) {
			*category(r).Amenities = append(*category(r).Amenities, Amenity{RoomAmenityCode: " "})
		}},
		{"MultimediaDescriptions", allCaps, func(r *rq) {
			g := category(r)
			mds := *g.MultimediaDescriptions
			mds[0].TextItem, mds[0].ImageItems = nil, mds[2].ImageItems
			*g.MultimediaDescriptions = append(mds[1:], MultimediaDescription{InfoCode: InfoCodePictures}, mds[1])
		}},
		{"rooms with category data", allCaps, func(r *rq) {
			rooms := r.GuestRooms()
			rooms[1].MaxOccupancy = 2
			rooms[2].Amenities = &[]Amenity{{RoomAmenityCode: "26"}}
		}},
		{"RoomID", allCaps, func(r *rq) {
			rooms := r.GuestRooms()
			rooms[1].TypeRoom.RoomID = " "
			rooms[3].TypeRoom.RoomID = rooms[2].TypeRoom.RoomID
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := load[rq](t, "../../v201810/inventory/testdata/Inventory-OTA_HotelDescriptiveContentNotifRQ.xml")
			tt.change(&old)
			var want []string
			for _, e := range inventory202410.Validate(old, tt.caps) {
				want = append(want, e.Message)
			}
			if len(want) == 0 {
				t.Fatal("2024-10 reports no errors")
			}
			rulestest.Expect(t, Validate(from202410(old), tt.caps), want)
		})
	}
}

func TestValidate2026Changes(t *testing.T) {
	tests := []struct {
		name   string
		change func(c *HotelDescriptiveContent)
		want   []string
	}{
		{"MaxAdultOccupancy in range", func(c *HotelDescriptiveContent) { c.FacilityInfo.GuestRooms[1].MaxAdultOccupancy = 4 }, nil},
		{"MaxAdultOccupancy above MaxOccupancy", func(c *HotelDescriptiveContent) {
			c.FacilityInfo.GuestRooms[1].MaxAdultOccupancy = 5
		}, []string{"GuestRoom[1]: MaxAdultOccupancy exceeds MaxOccupancy"}},
		{"room with MaxAdultOccupancy", func(c *HotelDescriptiveContent) {
			c.FacilityInfo.GuestRooms[1].MaxAdultOccupancy = 4
			c.FacilityInfo.GuestRooms[3].MaxAdultOccupancy = 2
		}, []string{"GuestRoom[3]: room of category DOUBLE allows only Code and RoomID"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := load[HotelDescriptiveContentNotifRQ](t, basicSample)
			tt.change(&rq.HotelDescriptiveContent)
			rulestest.Expect(t, Validate(rq, allCaps), tt.want)
		})
	}
}
