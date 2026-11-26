package v202210

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202210/guestrequests"
)

const (
	guestRequestsSamples = "guestrequests/testdata/GuestRequests-"
	// The empty and error responses are those of 2018-10.
	guestRequestsSamples2018 = "../v201810/guestrequests/testdata/GuestRequests-"
)

func withSpecialRequests(t *testing.T) guestrequests.ResRetrieveRS {
	rs := readSample[guestrequests.ResRetrieveRS](t, guestRequestsSamples+"OTA_ResRetrieveRS-reservation.xml")
	rs.HotelReservations[0].ResGlobalInfo.SpecialRequests = &[]guestrequests.SpecialRequest{
		{CodeContext: guestrequests.CodeContextAlpineBits, RequestCode: guestrequests.RequestCodePets, Text: &alpinebits.Text{TextFormat: alpinebits.TextFormatPlainText, Value: "two dogs"}},
		{CodeContext: guestrequests.CodeContextHotel, RequestCode: "VEGETARIAN_GUEST"},
	}
	return rs
}

// guestrequests.ValidateResRetrieveRS relies on these being rejected by the
// XSD.
func TestXSDEnforcesSpecialRequestRules(t *testing.T) {
	b, err := xml.Marshal(withSpecialRequests(t))
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	if err := Version.Validate(b); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	tests := []struct{ name, old, new string }{
		{"CodeContext", `CodeContext="HOTEL"`, `CodeContext="PORTAL"`},
		{"no RequestCode", `RequestCode="VEGETARIAN_GUEST" `, ``},
		{"2020-10 Name", `RequestCode="PETS" CodeContext="ALPINEBITS"`, `Name="PETS"`},
		{"Text with Language", `TextFormat="PlainText"`, `TextFormat="PlainText" Language="en"`},
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

// The new attributes and SpecialRequests must survive marshalling and pass
// the XSD.
func TestGuestRequestsTypesRoundTrip(t *testing.T) {
	for _, name := range []string{"OTA_ReadRQ", "OTA_ReadRQ-selection-criteria"} {
		t.Run(name, func(t *testing.T) {
			roundTrip(t, readSample[guestrequests.ReadRQ](t, guestRequestsSamples+name+".xml"))
		})
	}
	for _, path := range []string{
		guestRequestsSamples + "OTA_ResRetrieveRS-reservation.xml",
		guestRequestsSamples + "OTA_ResRetrieveRS-cancellation.xml",
		guestRequestsSamples + "OTA_ResRetrieveRS-newsletter.xml",
		guestRequestsSamples + "OTA_ResRetrieveRS-quote-request.xml",
		guestRequestsSamples2018 + "OTA_ResRetrieveRS-reservation-empty.xml",
		guestRequestsSamples2018 + "OTA_ResRetrieveRS-error.xml",
	} {
		t.Run(path, func(t *testing.T) {
			roundTrip(t, readSample[guestrequests.ResRetrieveRS](t, path))
		})
	}

	t.Run("SpecialRequests and LastModifyDateTime", func(t *testing.T) {
		rs := withSpecialRequests(t)
		rs.HotelReservations[0].LastModifyDateTime = &guestrequests.DateTime{Time: time.Date(2022, 3, 22, 9, 0, 0, 0, time.UTC)}
		if errs := guestrequests.ValidateResRetrieveRS(rs); errs != nil {
			t.Fatalf("rules: %v", errs)
		}
		roundTrip(t, rs)
	})
	t.Run("alternatives", func(t *testing.T) {
		rs := readSample[guestrequests.ResRetrieveRS](t, guestRequestsSamples+"OTA_ResRetrieveRS-quote-request.xml")
		stays := *rs.HotelReservations[0].RoomStays
		stays[0].TimeSpan.Duration = new(alpinebits.Nights(5))
		other := stays[0]
		stays[0].RoomStayGroupID, other.RoomStayGroupID = "1", "2"
		*rs.HotelReservations[0].RoomStays = append(stays, other)
		if errs := guestrequests.ValidateResRetrieveRS(rs); errs != nil {
			t.Fatalf("rules: %v", errs)
		}
		roundTrip(t, rs)
	})
}

func FuzzGuestRequests(f *testing.F) {
	paths, err := filepath.Glob(guestRequestsSamples + "OTA_ResRetrieveRS*.xml")
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
		if Version.Validate(doc) != nil || xml.Unmarshal(doc, &root) != nil || root.XMLName.Local != "OTA_ResRetrieveRS" {
			return
		}
		var rs guestrequests.ResRetrieveRS
		if err := xml.Unmarshal(doc, &rs); err != nil {
			t.Fatalf("valid document does not unmarshal: %v", err)
		}
		guestrequests.ValidateResRetrieveRS(rs)
		roundTrip(t, rs)
	})
}
