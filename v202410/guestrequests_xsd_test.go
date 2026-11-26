package v202410

import (
	"bytes"
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits/v202410/guestrequests"
)

// The other 2024-10 samples are those of 2022-10 and 2018-10.
var guestRequestsSamples = []string{
	"../v202210/guestrequests/testdata/GuestRequests-OTA_ResRetrieveRS-cancellation.xml",
	"../v202210/guestrequests/testdata/GuestRequests-OTA_ResRetrieveRS-newsletter.xml",
	"../v202210/guestrequests/testdata/GuestRequests-OTA_ResRetrieveRS-quote-request.xml",
	"../v201810/guestrequests/testdata/GuestRequests-OTA_ResRetrieveRS-reservation-empty.xml",
	"../v201810/guestrequests/testdata/GuestRequests-OTA_ResRetrieveRS-error.xml",
}

// guestRequestsReservation is the official sample with digits for its IBAN:
// the 2024-10 XSD allows only up to 19 digits in BankAcctNumber.
func guestRequestsReservation(tb testing.TB) []byte {
	tb.Helper()
	b, err := os.ReadFile("guestrequests/testdata/GuestRequests-OTA_ResRetrieveRS-reservation.xml")
	if err != nil {
		tb.Fatal(err)
	}
	return bytes.Replace(b, []byte("IT12 3456 7890 1234 5678 90"), []byte("1234567890"), 1)
}

func TestGuestRequestsTypesRoundTrip(t *testing.T) {
	t.Run("reservation", func(t *testing.T) {
		var rs guestrequests.ResRetrieveRS
		if err := xml.Unmarshal(guestRequestsReservation(t), &rs); err != nil {
			t.Fatal(err)
		}
		roundTrip(t, rs)
	})
	for _, path := range guestRequestsSamples {
		t.Run(path, func(t *testing.T) {
			roundTrip(t, readSample[guestrequests.ResRetrieveRS](t, path))
		})
	}
}

func FuzzGuestRequests(f *testing.F) {
	f.Add(guestRequestsReservation(f))
	for _, p := range guestRequestsSamples {
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
