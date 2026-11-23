package v202010

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/HGV/alpinebits/v202010/freerooms"
)

func FuzzFreeRooms(f *testing.F) {
	paths, err := filepath.Glob("freerooms/testdata/*RQ*.xml")
	if err != nil || len(paths) == 0 {
		f.Fatal("no samples", err)
	}
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, doc []byte) {
		var root struct{ XMLName xml.Name }
		if Version.Validate(doc) != nil || xml.Unmarshal(doc, &root) != nil || root.XMLName.Local != "OTA_HotelInvCountNotifRQ" {
			return
		}
		var rq freerooms.HotelInvCountNotifRQ
		if err := xml.Unmarshal(doc, &rq); err != nil {
			t.Fatalf("valid document does not unmarshal: %v", err)
		}
		freerooms.Validate(rq, nil)
		freerooms.Validate(rq, freeRoomsCaps)
		freerooms.CheckInventory(rq, freerooms.HotelInventory{})
		roundTrip(t, rq)
	})
}
