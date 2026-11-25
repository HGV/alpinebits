package v202010

import (
	"encoding/xml"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202010/guestrequests"
)

// 2020-10 kept the 2018-10 samples.
const guestRequestsSamples = "../v201810/guestrequests/testdata/GuestRequests-"

func pets(text *alpinebits.Text) *guestrequests.SpecialRequest {
	return &guestrequests.SpecialRequest{Name: guestrequests.SpecialRequestNamePets, Text: text}
}

// guestrequests.ValidateResRetrieveRS relies on these being rejected by the
// XSD.
func TestXSDEnforcesSpecialRequestRules(t *testing.T) {
	rs := readSample[guestrequests.ResRetrieveRS](t, guestRequestsSamples+"OTA_ResRetrieveRS-reservation.xml")
	rs.HotelReservations[0].ResGlobalInfo.SpecialRequest = pets(&alpinebits.Text{TextFormat: alpinebits.TextFormatPlainText, Value: "two dogs"})
	b, err := xml.Marshal(rs)
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	if err := Version.Validate(b); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	tests := []struct{ name, old, new string }{
		{"Name", `Name="PETS"`, `Name="CATS"`},
		{"no TextFormat", ` TextFormat="PlainText"`, ``},
		{"Language", `TextFormat="PlainText"`, `TextFormat="PlainText" Language="en"`},
		{"two SpecialRequest", `</SpecialRequest>`, `</SpecialRequest><SpecialRequest Name="PETS"></SpecialRequest>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Version.Validate([]byte(replaceOnce(t, doc, tt.old, tt.new))); err == nil {
				t.Error("XSD accepted it")
			}
		})
	}
}

// The new SpecialRequest and currencies other than EUR must survive
// marshalling and pass the XSD.
func TestGuestRequestsTypesRoundTrip(t *testing.T) {
	sample := func(name string) guestrequests.ResRetrieveRS {
		return readSample[guestrequests.ResRetrieveRS](t, guestRequestsSamples+"OTA_ResRetrieveRS-"+name+".xml")
	}
	tests := []struct {
		name   string
		change func(res *guestrequests.HotelReservation)
	}{
		{"pets described", func(res *guestrequests.HotelReservation) {
			res.ResGlobalInfo.SpecialRequest = pets(&alpinebits.Text{TextFormat: alpinebits.TextFormatPlainText, Value: "two dogs"})
		}},
		{"pets not described", func(res *guestrequests.HotelReservation) {
			res.ResGlobalInfo.SpecialRequest = pets(&alpinebits.Text{TextFormat: alpinebits.TextFormatPlainText})
		}},
		{"no pets", func(res *guestrequests.HotelReservation) {
			res.ResGlobalInfo.SpecialRequest = pets(nil)
		}},
		{"other currency", func(res *guestrequests.HotelReservation) {
			stay := &(*res.RoomStays)[0]
			stay.Total.CurrencyCode = "CHF"
			stay.RatePlan.Commission = &guestrequests.Commission{CommissionPayableAmount: &guestrequests.CommissionPayableAmount{
				Amount: new(alpinebits.Decimal(10.5)), CurrencyCode: "CHF",
			}}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := sample("reservation")
			tt.change(&rs.HotelReservations[0])
			if errs := guestrequests.ValidateResRetrieveRS(rs); errs != nil {
				t.Fatalf("rules: %v", errs)
			}
			roundTrip(t, rs)
		})
	}
	for _, name := range []string{"reservation", "request-with-roomtype", "cancellation", "reservation-empty", "error"} {
		t.Run(name, func(t *testing.T) {
			roundTrip(t, sample(name))
		})
	}
}
