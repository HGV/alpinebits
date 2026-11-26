package v201810

import (
	"encoding/xml"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v201810/inventory"
)

const inventoryNotif = `<OTA_HotelDescriptiveContentNotifRQ ` + ns + ` Version="8.000">
<HotelDescriptiveContents><HotelDescriptiveContent HotelCode="123">
<FacilityInfo><GuestRooms>
<GuestRoom Code="DZ" MinOccupancy="1" MaxOccupancy="2" MaxChildOccupancy="1" ID="double">
<TypeRoom StandardOccupancy="2" RoomClassificationCode="42" RoomType="1" Size="20"/>
<Amenities><Amenity RoomAmenityCode="26"/></Amenities>
<MultimediaDescriptions>
<MultimediaDescription InfoCode="25"><TextItems><TextItem>
<Description TextFormat="PlainText" Language="en">Double room</Description>
</TextItem></TextItems></MultimediaDescription>
<MultimediaDescription InfoCode="23"><ImageItems><ImageItem Category="6">
<ImageFormat CopyrightNotice="(c)"><URL>https://example.com/dz.jpg</URL></ImageFormat>
<Description TextFormat="PlainText" Language="de">Zimmer</Description>
</ImageItem></ImageItems></MultimediaDescription>
</MultimediaDescriptions>
</GuestRoom>
<GuestRoom Code="DZ"><TypeRoom RoomID="101"/></GuestRoom>
</GuestRooms></FacilityInfo>
</HotelDescriptiveContent></HotelDescriptiveContents>
</OTA_HotelDescriptiveContentNotifRQ>`

// inventory.Validate relies on these being rejected by the XSD.
func TestInventoryXSDEnforcesRules(t *testing.T) {
	if err := Version.Validate([]byte(inventoryNotif)); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	content := `<HotelDescriptiveContent HotelCode="123">`
	tests := []struct{ name, old, new string }{
		{"two HotelDescriptiveContent", `</HotelDescriptiveContent>`, `</HotelDescriptiveContent>` + content + `</HotelDescriptiveContent>`},
		{"empty HotelCode", `HotelCode="123"`, `HotelCode=""`},
		{"unknown attribute", `HotelCode="123"`, `HotelCode="123" AreaID="1"`},
		{"no GuestRooms", `<FacilityInfo><GuestRooms>`, `<FacilityInfo><GuestRooms></GuestRooms><GuestRooms>`},
		{"no Code", `<GuestRoom Code="DZ" MinOccupancy`, `<GuestRoom MinOccupancy`},
		{"Code too long", `<GuestRoom Code="DZ" MinOccupancy`, `<GuestRoom Code="DOUBLEROOM" MinOccupancy`},
		{"MinOccupancy 0", `MinOccupancy="1"`, `MinOccupancy="0"`},
		{"MaxChildOccupancy 0", `MaxChildOccupancy="1"`, `MaxChildOccupancy="0"`},
		{"StandardOccupancy 0", `StandardOccupancy="2"`, `StandardOccupancy="0"`},
		{"RoomType", `RoomType="1"`, `RoomType="10"`},
		{"empty RoomID", `RoomID="101"`, `RoomID=""`},
		{"two TypeRoom", `<TypeRoom RoomID="101"/>`, `<TypeRoom RoomID="101"/><TypeRoom RoomID="102"/>`},
		{"elements out of order", `<TypeRoom RoomID="101"/>`, `<Amenities><Amenity RoomAmenityCode="26"/></Amenities><TypeRoom RoomID="101"/>`},
		{"empty Amenities", `<Amenity RoomAmenityCode="26"/>`, ``},
		{"RoomAmenityCode 0", `RoomAmenityCode="26"`, `RoomAmenityCode="0"`},
		{"ABR RoomAmenityCode", `RoomAmenityCode="26"`, `RoomAmenityCode="7.ABR"`},
		{"InfoCode", `InfoCode="25"`, `InfoCode="2"`},
		{"TextItems and ImageItems", `</TextItems>`, `</TextItems><ImageItems><ImageItem Category="6"><ImageFormat><URL>https://example.com</URL></ImageFormat></ImageItem></ImageItems>`},
		{"two TextItem", `</TextItem>`, `</TextItem><TextItem><Description TextFormat="PlainText" Language="de">Doppelzimmer</Description></TextItem>`},
		{"TextItem without Description", `<Description TextFormat="PlainText" Language="en">Double room</Description>`, ``},
		{"TextFormat", `TextFormat="PlainText" Language="en"`, `TextFormat="Markdown" Language="en"`},
		{"no Language", `TextFormat="PlainText" Language="en"`, `TextFormat="PlainText"`},
		{"Language uppercase", `Language="en"`, `Language="EN"`},
		{"empty Description", `>Double room<`, `><`},
		{"empty ImageItems", `<ImageItems><ImageItem Category="6">`, `<ImageItems></ImageItems><ImageItems><ImageItem Category="6">`},
		{"no Category", `<ImageItem Category="6">`, `<ImageItem>`},
		{"no ImageFormat", `<ImageFormat CopyrightNotice="(c)"><URL>https://example.com/dz.jpg</URL></ImageFormat>`, ``},
		{"URL scheme", `https://example.com/dz.jpg`, `ftp://example.com/dz.jpg`},
		{"HTML image Description", `TextFormat="PlainText" Language="de"`, `TextFormat="HTML" Language="de"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if n := strings.Count(inventoryNotif, tt.old); n != 1 {
				t.Fatalf("%q occurs %d times", tt.old, n)
			}
			doc := strings.Replace(inventoryNotif, tt.old, tt.new, 1)
			if err := Version.Validate([]byte(doc)); err == nil {
				t.Error("XSD accepted it")
			}
		})
	}
}

// Optional parts must be omitted, not marshalled empty, and absent containers
// must stay absent.
func TestInventoryTypesRoundTrip(t *testing.T) {
	allCaps := []alpinebits.Capability{inventory.CapabilityUseRooms, inventory.CapabilityOccupancyChildren}
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
			rq.HotelDescriptiveContent.HotelName = ""
		}},
		{"Inventory-OTA_HotelDescriptiveContentNotifRQ-delete-all.xml", nil},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			b, err := os.ReadFile("inventory/testdata/" + tt.file)
			if err != nil {
				t.Fatal(err)
			}
			var rq inventory.HotelDescriptiveContentNotifRQ
			if err := xml.Unmarshal(b, &rq); err != nil {
				t.Fatal(err)
			}
			if tt.change != nil {
				tt.change(&rq)
			}
			if errs := inventory.Validate(rq, allCaps); errs != nil {
				t.Fatalf("rules: %v", errs)
			}

			out, err := xml.Marshal(rq)
			if err != nil {
				t.Fatal(err)
			}
			if err := Version.Validate(out); err != nil {
				t.Fatalf("XSD: %v\n%s", err, out)
			}
			var back inventory.HotelDescriptiveContentNotifRQ
			if err := xml.Unmarshal(out, &back); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(back, rq) {
				t.Errorf("round trip changed the message:\n got %+v\nwant %+v", back, rq)
			}
		})
	}
}
