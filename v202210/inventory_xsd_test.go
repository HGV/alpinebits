package v202210

import (
	"encoding/xml"
	"os"
	"reflect"
	"testing"

	"github.com/HGV/alpinebits/v202210/inventory"
)

// The 2022-10 RQ samples are the 2018-10 ones, reformatted; the RS samples are
// new.
const inventorySamples = "../v201810/inventory/testdata/"

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

func TestInventoryTypesRoundTrip(t *testing.T) {
	for _, name := range []string{"", "-with-roomtype", "-delete-all"} {
		t.Run("RQ"+name, func(t *testing.T) {
			roundTrip(t, readSample[inventory.HotelDescriptiveContentNotifRQ](t, inventorySamples+"Inventory-OTA_HotelDescriptiveContentNotifRQ"+name+".xml"))
		})
	}
	for _, name := range []string{"success", "advisory", "warning", "error"} {
		t.Run("RS "+name, func(t *testing.T) {
			roundTrip(t, readSample[inventory.HotelDescriptiveContentNotifRS](t, "inventory/testdata/Inventory-Push-OTA_HotelDescriptiveContentNotifRS-"+name+".xml"))
		})
	}
}
