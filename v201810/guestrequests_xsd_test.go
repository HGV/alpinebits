package v201810

import (
	"context"
	"encoding/xml"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v201810/guestrequests"
	"github.com/HGV/x/timex"
)

func guestRequestsSample(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile("guestrequests/testdata/GuestRequests-" + name + ".xml")
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func unmarshalGuestRequests[T any](t *testing.T, doc string) T {
	t.Helper()
	var v T
	if err := xml.Unmarshal([]byte(doc), &v); err != nil {
		t.Fatal(err)
	}
	return v
}

// The guestrequests rules rely on these being rejected by the XSD; some show
// why the types use pointers, omitempty or MarshalXML.
func TestGuestRequestsXSDEnforcesRules(t *testing.T) {
	read := guestRequestsSample(t, "OTA_ReadRQ")
	ack := guestRequestsSample(t, "OTA_ReadRQ-ack")
	reservation := guestRequestsSample(t, "OTA_ResRetrieveRS-reservation")
	empty := guestRequestsSample(t, "OTA_ResRetrieveRS-reservation-empty")
	notifReportRS := `<OTA_NotifReportRS ` + ns + ` Version="1.0"><Success/></OTA_NotifReportRS>`
	for _, doc := range []string{read, ack, reservation, empty, notifReportRS} {
		if err := Version.Validate([]byte(doc)); err != nil {
			t.Fatalf("baseline: %v\n%s", err, doc)
		}
	}
	acks := ack[strings.Index(ack, "<HotelReservation>") : strings.LastIndex(ack, "</HotelReservation>")+len("</HotelReservation>")]
	tests := []struct{ name, doc, old, new string }{
		{"HotelCode empty", read, `HotelCode="123"`, `HotelCode=""`},
		{"SelectionCriteria without Start", read, ` Start="2012-03-21T15:00:00+01:00"`, ``},
		{"Start without time", read, `Start="2012-03-21T15:00:00+01:00"`, `Start="2012-03-21"`},
		{"acknowledgement without Success", ack, `<Success/>`, ``},
		{"refusal without RecordID", ack, ` RecordID="f054bbd2f5ebab9"`, ``},
		{"empty HotelReservations", ack, acks, ``},
		{"acknowledged UniqueID Type", ack, `Type="14" ID="6b34fe24ac2ff810"`, `Type="16" ID="6b34fe24ac2ff810"`},
		{"no ReservationsList", empty, `<ReservationsList/>`, ``},
		{"ResStatus", reservation, `ResStatus="Reserved"`, `ResStatus="Booked"`},
		{"UniqueID Type", reservation, `Type="14"`, `Type="16"`},
		{"RoomTypeCode too long", reservation, `RoomTypeCode="bigsuite"`, `RoomTypeCode="bigsuite2"`},
		{"Commission CurrencyCode", reservation, `<Commission Percent="15"/>`, `<Commission><CommissionPayableAmount Amount="10" CurrencyCode="USD"/></Commission>`},
		{"MealPlanIndicator false", reservation, `MealPlanIndicator="true"`, `MealPlanIndicator="false"`},
		{"MealPlanCodes", reservation, `MealPlanCodes="1"`, `MealPlanCodes="2"`},
		{"GuestCount Count 0", reservation, `<GuestCount Count="2"/>`, `<GuestCount Count="0"/>`},
		{"TimeSpan Start empty", reservation, `Start="2012-01-01"`, `Start=""`},
		{"Duration in days", reservation, `End="2012-01-12"`, `End="2012-01-12" Duration="P11D"`},
		{"ExpireDate", reservation, `ExpireDate="1216"`, `ExpireDate="12/16"`},
		{"no CardHolderName", reservation, `<CardHolderName>Otto Mustermann</CardHolderName>`, ``},
		{"Total without CurrencyCode", reservation, `<Total AmountAfterTax="299" CurrencyCode="EUR"/>`, `<Total AmountAfterTax="299"/>`},
		{"ResGuest Profile with ProfileType", reservation, `<Profile>`, `<Profile ProfileType="4">`},
		{"BirthDate empty", reservation, `BirthDate="1980-01-01"`, `BirthDate=""`},
		{"Telephone without PhoneTechType", reservation, `<Telephone PhoneTechType="1" `, `<Telephone `},
		{"Comment with ListItem and Text", reservation, `Skipass</ListItem>`, `Skipass</ListItem><Text>x</Text>`},
		{"three Comments", reservation, `<Comments>`, `<Comments><Comment Name="customer comment"><Text>x</Text></Comment>`},
		{"ListItem without Language", reservation, `<ListItem ListItem="1" Language="de">`, `<ListItem ListItem="1">`},
		{"CompanyName without CodeContext", reservation, ` CodeContext="ABC"`, ``},
		{"ProfileType", reservation, `ProfileType="4"`, `ProfileType="5"`},
		{"NotifReportRS with Warnings", notifReportRS, `<Success/>`, `<Success/><Warnings><Warning Type="11">unknown ID</Warning></Warnings>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if n := strings.Count(tt.doc, tt.old); n != 1 {
				t.Fatalf("%q occurs %d times", tt.old, n)
			}
			doc := strings.Replace(tt.doc, tt.old, tt.new, 1)
			if err := Version.Validate([]byte(doc)); err == nil {
				t.Error("XSD accepted it")
			}
		})
	}
}

func guestRequestsRoundTrip[T any](t *testing.T, v T) {
	t.Helper()
	out, err := xml.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := Version.Validate(out); err != nil {
		t.Fatalf("XSD: %v\n%s", err, out)
	}
	back := unmarshalGuestRequests[T](t, string(out))
	if !reflect.DeepEqual(back, v) {
		t.Errorf("round trip changed the message:\n got %+v\nwant %+v\n%s", back, v, out)
	}
}

func guestRequestsDate(s string) *timex.Date {
	d, err := timex.ParseDate(s)
	if err != nil {
		panic(err)
	}
	return &d
}

// Values whose absence means something (an adult, a commission given as
// amount, an open quote request) must survive marshalling and pass the XSD.
func TestGuestRequestsTypesRoundTrip(t *testing.T) {
	t.Run("OTA_ReadRQ", func(t *testing.T) {
		rq := unmarshalGuestRequests[guestrequests.ReadRQ](t, guestRequestsSample(t, "OTA_ReadRQ"))
		guestRequestsRoundTrip(t, rq)

		rq.ReadRequest.HotelCode = ""
		rq.ReadRequest.SelectionCriteria.Start = guestrequests.DateTime{}
		if err := rq.ReadRequest.SelectionCriteria.Start.UnmarshalText([]byte("2012-03-21T15:00:00")); err != nil {
			t.Fatal(err)
		}
		guestRequestsRoundTrip(t, rq)
	})

	t.Run("OTA_NotifReportRQ", func(t *testing.T) {
		rq := unmarshalGuestRequests[guestrequests.NotifReportRQ](t, guestRequestsSample(t, "OTA_ReadRQ-ack"))
		if w := *rq.Warnings; len(w) != 1 || w[0].RecordID != "f054bbd2f5ebab9" || w[0].Type != 3 || w[0].Code != 450 {
			t.Errorf("Warnings %+v", w)
		}
		guestRequestsRoundTrip(t, rq)

		rq.HotelReservations = nil
		guestRequestsRoundTrip(t, rq)
		rq.Warnings = nil
		guestRequestsRoundTrip(t, rq)
	})

	t.Run("reservation", func(t *testing.T) {
		rs := unmarshalGuestRequests[guestrequests.ResRetrieveRS](t, guestRequestsSample(t, "OTA_ResRetrieveRS-reservation"))
		guestRequestsRoundTrip(t, rs)

		res := &rs.HotelReservations[0]
		stay := &(*res.RoomStays)[0]
		stay.RoomType.RoomClassificationCode, stay.RoomType.RoomType = new(0), 1
		stay.RatePlan.Commission = &guestrequests.Commission{CommissionPayableAmount: &guestrequests.CommissionPayableAmount{Amount: new(alpinebits.Decimal(0.0)), CurrencyCode: "EUR"}}
		*stay.GuestCounts = append(*stay.GuestCounts, guestrequests.GuestCount{Count: 1, Age: new(0)})
		stay.Guarantee.PaymentCard.CardNumber = guestrequests.CardNumber{EncryptedValue: "x", EncryptionMethod: "RSA-PKCSV2.2"}
		customer := &res.ResGuest.Customer
		customer.BirthDate, customer.Telephones, customer.Email, customer.Address = nil, nil, nil, nil
		(*res.ResGlobalInfo.Comments)[0].ListItems[0].ListItem = 0
		res.ResGlobalInfo.GlobalProfile.CompanyInfo.AddressInfo = nil
		res.ResGlobalInfo.BasicPropertyInfo.HotelName = ""
		res.ResStatus = guestrequests.ResStatusModify
		if errs := guestrequests.ValidateResRetrieveRS(rs); len(errs) > 0 {
			t.Fatalf("%+v", errs)
		}
		guestRequestsRoundTrip(t, rs)
	})

	t.Run("quote request", func(t *testing.T) {
		rs := unmarshalGuestRequests[guestrequests.ResRetrieveRS](t, guestRequestsSample(t, "OTA_ResRetrieveRS-request-with-roomtype"))
		if c := rs.HotelReservations[0].CreateDateTime; !c.NoOffset {
			t.Errorf("CreateDateTime %v has an offset", c)
		}
		guestRequestsRoundTrip(t, rs)

		res := &rs.HotelReservations[0]
		stay := &(*res.RoomStays)[0]
		stay.TimeSpan.Duration = new(alpinebits.Nights(5))
		stay.RoomType = &guestrequests.RoomType{}
		stay.RatePlan = &guestrequests.RatePlan{MealsIncluded: &guestrequests.MealsIncluded{MealPlanIndicator: true, MealPlanCodes: "3"}}
		stay.Total = &guestrequests.Total{AmountAfterTax: 0, CurrencyCode: "EUR"}
		*res.RoomStays = append(*res.RoomStays, guestrequests.RoomStay{TimeSpan: guestrequests.TimeSpan{
			Start: guestRequestsDate("2017-11-03"), End: guestRequestsDate("2017-11-08"),
		}})
		res.ResGlobalInfo.Comments = &[]guestrequests.Comment{{Name: guestrequests.CommentNameCustomerComment, Text: "Hunde?"}}
		if errs := guestrequests.ValidateResRetrieveRS(rs); len(errs) > 0 {
			t.Fatalf("%+v", errs)
		}
		guestRequestsRoundTrip(t, rs)
	})

	for _, name := range []string{"cancellation", "reservation-empty", "error"} {
		t.Run(name, func(t *testing.T) {
			guestRequestsRoundTrip(t, unmarshalGuestRequests[guestrequests.ResRetrieveRS](t, guestRequestsSample(t, "OTA_ResRetrieveRS-"+name)))
		})
	}
}

// The error response the router builds when the rules fail must pass the
// router's response validation, or it answers 500.
func TestGuestRequestsRouter(t *testing.T) {
	r := alpinebits.NewRouter().
		Handle(ActionRead,
			func(context.Context, alpinebits.Request[guestrequests.ReadRQ]) (guestrequests.ResRetrieveRS, error) {
				return guestrequests.ResRetrieveRS{Response: alpinebits.OK()}, nil
			})
	code, body := post(t, r, "OTA_Read:GuestRequests",
		`<OTA_ReadRQ `+ns+` Version="1.0"><ReadRequests><HotelReadRequest HotelCode=" "/></ReadRequests></OTA_ReadRQ>`)
	if code != http.StatusOK {
		t.Fatalf("status %d: %s", code, body)
	}
	rs := unmarshalGuestRequests[guestrequests.ResRetrieveRS](t, body)
	if rs.Errors == nil || len(*rs.Errors) != 1 || (*rs.Errors)[0].Message != "HotelCode required" || rs.HotelReservations != nil {
		t.Errorf("got %s", body)
	}
}
