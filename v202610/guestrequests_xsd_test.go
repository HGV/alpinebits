package v202610

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HGV/alpinebits/v202610/guestrequests"
)

// The empty and error responses are those of 2018-10.
const guestRequestsSamples2018 = "../v201810/guestrequests/testdata/GuestRequests-"

// guestrequests.ValidateResRetrieveRS relies on these being rejected by the
// XSD.
func TestXSDEnforcesGuestRequestsRules(t *testing.T) {
	b, err := os.ReadFile("guestrequests/testdata/GuestRequests-OTA_ResRetrieveRS-reservation.xml")
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	tests := []struct{ name, old, new string }{
		{"GuaranteeTypeCode", `GuaranteeTypeCode="28"`, `GuaranteeTypeCode="99"`},
		{"PaymentTransactionTypeCode", `PaymentTransactionTypeCode="reserve"`, `PaymentTransactionTypeCode="capture"`},
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

func TestGuestRequestsSamplesRoundTrip(t *testing.T) {
	paths, err := filepath.Glob("guestrequests/testdata/*.xml")
	if err != nil {
		t.Fatal(err)
	}
	paths = append(paths, guestRequestsSamples2018+"OTA_ResRetrieveRS-reservation-empty.xml", guestRequestsSamples2018+"OTA_ResRetrieveRS-error.xml")
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			switch {
			case strings.Contains(path, "OTA_ReadRQ"):
				roundTripSample[guestrequests.ReadRQ](t, path)
			case strings.Contains(path, "OTA_ResRetrieveRS"):
				roundTripSample[guestrequests.ResRetrieveRS](t, path)
			case strings.Contains(path, "OTA_NotifReportRQ"):
				roundTripSample[guestrequests.NotifReportRQ](t, path)
			case strings.Contains(path, "OTA_NotifReportRS"):
				roundTripSample[guestrequests.NotifReportRS](t, path)
			default:
				t.Fatal("unknown message")
			}
		})
	}
}

func FuzzGuestRequests(f *testing.F) {
	paths, err := filepath.Glob("guestrequests/testdata/*.xml")
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
		case "OTA_ResRetrieveRS":
			var rs guestrequests.ResRetrieveRS
			if err := xml.Unmarshal(doc, &rs); err != nil {
				t.Fatalf("valid document does not unmarshal: %v", err)
			}
			guestrequests.ValidateResRetrieveRS(rs)
			roundTrip[guestrequests.ResRetrieveRS](t, doc)
		case "OTA_NotifReportRQ":
			var rq guestrequests.NotifReportRQ
			if err := xml.Unmarshal(doc, &rq); err != nil {
				t.Fatalf("valid document does not unmarshal: %v", err)
			}
			guestrequests.ValidateNotifReportRQ(rq, nil)
			roundTrip[guestrequests.NotifReportRQ](t, doc)
		}
	})
}
