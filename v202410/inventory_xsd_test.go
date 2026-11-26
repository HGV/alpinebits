package v202410

import (
	"encoding/xml"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/HGV/alpinebits/v202410/inventory"
)

// The 2024-10 RQ samples are the 2018-10 ones, reformatted; the RS samples are
// those of 2022-10.
const (
	inventorySamples   = "../v201810/inventory/testdata/"
	inventoryRSSamples = "../v202210/inventory/testdata/"
)

func readSample[T any](t *testing.T, path string) T {
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

func roundTrip[T any](t *testing.T, v T) {
	t.Helper()
	out, err := xml.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := Version.Validate(out); err != nil {
		t.Fatalf("XSD: %v\n%s", err, out)
	}
	var back T
	if err := xml.Unmarshal(out, &back); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back, v) {
		t.Errorf("round trip changed the message:\n got %+v\nwant %+v\n%s", back, v, out)
	}
}

func inventoryWithABR(t *testing.T) inventory.HotelDescriptiveContentNotifRQ {
	rq := readSample[inventory.HotelDescriptiveContentNotifRQ](t, inventorySamples+"Inventory-OTA_HotelDescriptiveContentNotifRQ.xml")
	*rq.GuestRooms()[0].Amenities = []inventory.Amenity{{RoomAmenityCode: "26"}, {RoomAmenityCode: "7.ABR"}}
	return rq
}

// inventory.Validate relies on these being rejected by the XSD.
func TestXSDEnforcesInventoryRules(t *testing.T) {
	b, err := xml.Marshal(inventoryWithABR(t))
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	if err := Version.Validate(b); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	tests := []struct{ name, old, new string }{
		{"lowercase ABR", `"7.ABR"`, `"7.abr"`},
		{"code too long", `"7.ABR"`, `"1234"`},
		{"list only", `"7.ABR"`, `"ABR"`},
		// §4.3.1 describes AdditionalCodes, but the 2024-10 XSD lacks them.
		{"TPA_Extensions", `</FacilityInfo>`, `</FacilityInfo><TPA_Extensions><AdditionalCodes><AdditionalCode Type="CIN">IT1</AdditionalCode></AdditionalCodes></TPA_Extensions>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if n := strings.Count(doc, tt.old); n != 1 {
				t.Fatalf("%q occurs %d times", tt.old, n)
			}
			if err := Version.Validate([]byte(strings.Replace(doc, tt.old, tt.new, 1))); err == nil {
				t.Error("XSD accepted it")
			}
		})
	}
}

func TestInventoryTypesRoundTrip(t *testing.T) {
	t.Run("ABR amenities", func(t *testing.T) {
		roundTrip(t, inventoryWithABR(t))
	})
	for _, name := range []string{"-with-roomtype", "-delete-all"} {
		t.Run("RQ"+name, func(t *testing.T) {
			roundTrip(t, readSample[inventory.HotelDescriptiveContentNotifRQ](t, inventorySamples+"Inventory-OTA_HotelDescriptiveContentNotifRQ"+name+".xml"))
		})
	}
	for _, name := range []string{"success", "advisory", "warning", "error"} {
		t.Run("RS "+name, func(t *testing.T) {
			roundTrip(t, readSample[inventory.HotelDescriptiveContentNotifRS](t, inventoryRSSamples+"Inventory-Push-OTA_HotelDescriptiveContentNotifRS-"+name+".xml"))
		})
	}
}
