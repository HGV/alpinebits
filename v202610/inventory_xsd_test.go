package v202610

import (
	"encoding/xml"
	"os"
	"strings"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202610/inventory"
)

// 2026-10 kept the 2018-10 delete-all sample and the 2022-10 RS samples.
const (
	inventoryBasicSample  = "inventory/testdata/Inventory-Push-OTA_HotelDescriptiveContentNotifRQ-basic.xml"
	inventoryDeleteSample = "../v201810/inventory/testdata/Inventory-OTA_HotelDescriptiveContentNotifRQ-delete-all.xml"
	inventoryRSSamples    = "../v202210/inventory/testdata/"
)

var inventoryCaps = []alpinebits.Capability{inventory.CapUseRooms, inventory.CapOccupancyChildren}

// inventoryWithCodes is the official sample with AdditionalCodes and a valid
// MaxAdultOccupancy for SUITE.
func inventoryWithCodes(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile(inventoryBasicSample)
	if err != nil {
		t.Fatal(err)
	}
	var rq inventory.HotelDescriptiveContentNotifRQ
	if err := xml.Unmarshal(b, &rq); err != nil {
		t.Fatal(err)
	}
	rq.GuestRooms()[1].MaxAdultOccupancy = 4
	rq.HotelDescriptiveContent.AdditionalCodes = &[]inventory.AdditionalCode{{Type: "CIN", Value: "IT12345678A1B23456"}, {Type: "CIR", Value: "021008-AGR-00001"}}
	if errs := inventory.Validate(rq, inventoryCaps); errs != nil {
		t.Fatalf("rules: %v", errs)
	}
	out, err := xml.Marshal(rq)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// inventory.Validate relies on these being rejected by the XSD.
func TestXSDEnforcesInventoryRules(t *testing.T) {
	doc := string(inventoryWithCodes(t))
	if err := Version.Validate([]byte(doc)); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	tests := []struct{ name, old, new string }{
		{"AdditionalCode without Type", ` Type="CIN"`, ``},
		{"empty AdditionalCode", `>IT12345678A1B23456<`, `><`},
		{"empty AdditionalCodes", `<AdditionalCodes>`, `<AdditionalCodes></AdditionalCodes><AdditionalCodes>`},
		{"MaxAdultOccupancy 0", `MaxAdultOccupancy="4"`, `MaxAdultOccupancy="0"`},
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

func TestInventorySamplesRoundTrip(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		roundTripSample[inventory.HotelDescriptiveContentNotifRQ](t, inventoryBasicSample)
	})
	t.Run("AdditionalCodes", func(t *testing.T) {
		roundTrip[inventory.HotelDescriptiveContentNotifRQ](t, inventoryWithCodes(t))
	})
	t.Run("delete-all", func(t *testing.T) {
		roundTripSample[inventory.HotelDescriptiveContentNotifRQ](t, inventoryDeleteSample)
	})
	for _, name := range []string{"success", "advisory", "warning", "error"} {
		t.Run("RS "+name, func(t *testing.T) {
			roundTripSample[inventory.HotelDescriptiveContentNotifRS](t, inventoryRSSamples+"Inventory-Push-OTA_HotelDescriptiveContentNotifRS-"+name+".xml")
		})
	}
}

func FuzzInventory(f *testing.F) {
	for _, p := range []string{inventoryBasicSample, inventoryDeleteSample} {
		b, err := os.ReadFile(p)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, doc []byte) {
		var root struct{ XMLName xml.Name }
		if Version.Validate(doc) != nil || xml.Unmarshal(doc, &root) != nil || root.XMLName.Local != "OTA_HotelDescriptiveContentNotifRQ" {
			return
		}
		var rq inventory.HotelDescriptiveContentNotifRQ
		if err := xml.Unmarshal(doc, &rq); err != nil {
			t.Fatalf("valid document does not unmarshal: %v", err)
		}
		inventory.Validate(rq, nil)
		inventory.Validate(rq, inventoryCaps)
		roundTrip[inventory.HotelDescriptiveContentNotifRQ](t, doc)
	})
}
