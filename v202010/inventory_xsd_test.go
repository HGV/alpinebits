package v202010

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202010/inventory"
)

// 2020-10 did not change Inventory/Basic, so the 2018-10 samples apply.
const inventorySamples = "../v201810/inventory/testdata/"

var inventoryCaps = []alpinebits.Capability{inventory.CapUseRooms, inventory.CapOccupancyChildren}

func inventorySample(t *testing.T, file string) inventory.HotelDescriptiveContentNotifRQ {
	t.Helper()
	b, err := os.ReadFile(inventorySamples + file)
	if err != nil {
		t.Fatal(err)
	}
	var rq inventory.HotelDescriptiveContentNotifRQ
	if err := xml.Unmarshal(b, &rq); err != nil {
		t.Fatal(err)
	}
	return rq
}

// Optional parts must be omitted, not marshalled empty, and absent containers
// must stay absent.
func TestInventoryTypesRoundTrip(t *testing.T) {
	tests := []struct {
		file   string
		change func(rq *inventory.HotelDescriptiveContentNotifRQ)
	}{
		{"Inventory-OTA_HotelDescriptiveContentNotifRQ.xml", func(rq *inventory.HotelDescriptiveContentNotifRQ) {
			g := &rq.HotelDescriptiveContent.FacilityInfo.GuestRooms[0]
			g.OldCode = "double"
			g.TypeRoom.RoomType, g.TypeRoom.Size = inventory.RoomTypeRoom, 20
			text := g.MultimediaDescription(inventory.InfoCodeDescription).TextItem
			text.Descriptions = append(text.Descriptions, alpinebits.Text{TextFormat: alpinebits.TextFormatHTML, Language: "en", Value: "<p>Description</p>"})
			images := g.MultimediaDescription(inventory.InfoCodePictures).ImageItems
			*images = append(*images, inventory.ImageItem{Category: 17, ImageFormat: inventory.ImageFormat{URL: "https://example.com/map.png"}})
		}},
		{"Inventory-OTA_HotelDescriptiveContentNotifRQ-with-roomtype.xml", func(rq *inventory.HotelDescriptiveContentNotifRQ) {
			rq.HotelDescriptiveContent.HotelCode = ""
		}},
		{"Inventory-OTA_HotelDescriptiveContentNotifRQ-delete-all.xml", nil},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			rq := inventorySample(t, tt.file)
			if tt.change != nil {
				tt.change(&rq)
			}
			if errs := inventory.Validate(rq, inventoryCaps); errs != nil {
				t.Fatalf("rules: %v", errs)
			}
			roundTrip(t, rq)
		})
	}
	t.Run("RS", func(t *testing.T) {
		out, err := xml.Marshal(inventory.HotelDescriptiveContentNotifRS{Response: alpinebits.OK()})
		if err != nil {
			t.Fatal(err)
		}
		if err := Version.Validate(out); err != nil {
			t.Errorf("XSD: %v\n%s", err, out)
		}
	})
}
