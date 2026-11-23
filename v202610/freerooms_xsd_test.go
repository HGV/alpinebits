package v202610

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/HGV/alpinebits/v202610/freerooms"
)

func roundTripSample[T any](t *testing.T, path string) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	roundTrip[T](t, b)
}

// roundTrip checks that a valid doc unmarshals into T and marshals to a valid
// document that unmarshals to the same T.
func roundTrip[T any](t *testing.T, doc []byte) {
	t.Helper()
	var v T
	if err := xml.Unmarshal(doc, &v); err != nil {
		t.Fatal(err)
	}
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

func TestFreeRoomsSamplesRoundTrip(t *testing.T) {
	for _, name := range []string{"", "-closing_seasons", "-delta", "-empty", "-out-of-order-and-out-of-market", "-rooms"} {
		t.Run("RQ"+name, func(t *testing.T) {
			roundTripSample[freerooms.HotelInvCountNotifRQ](t, "freerooms/testdata/FreeRooms-OTA_HotelInvCountNotifRQ"+name+".xml")
		})
	}
	for _, name := range []string{"success", "advisory", "warning", "error"} {
		t.Run("RS "+name, func(t *testing.T) {
			roundTripSample[freerooms.HotelInvCountNotifRS](t, "freerooms/testdata/FreeRooms-OTA_HotelInvCountNotifRS-"+name+".xml")
		})
	}
}

func FuzzFreeRooms(f *testing.F) {
	paths, err := filepath.Glob("freerooms/testdata/*.xml")
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
		if Version.Validate(doc) != nil || xml.Unmarshal(doc, &root) != nil {
			return
		}
		switch root.XMLName.Local {
		case "OTA_HotelInvCountNotifRQ":
			roundTrip[freerooms.HotelInvCountNotifRQ](t, doc)
		case "OTA_HotelInvCountNotifRS":
			roundTrip[freerooms.HotelInvCountNotifRS](t, doc)
		}
	})
}
