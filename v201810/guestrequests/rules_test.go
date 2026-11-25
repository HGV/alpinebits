package guestrequests

import (
	"encoding/xml"
	"os"
	"slices"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
	"github.com/HGV/x/timex"
)

func date(s string) *timex.Date {
	d, err := timex.ParseDate(s)
	if err != nil {
		panic(err)
	}
	return &d
}

func load[T any](t *testing.T, name string) T {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	var v T
	if err := xml.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestValidateReadRQ(t *testing.T) {
	tests := []struct {
		name                 string
		hotelCode, hotelName string
		want                 []string
	}{
		{"sample", "123", "Frangart Inn", nil},
		{"HotelName only", " ", "Frangart Inn", []string{"HotelCode required"}},
		{"HotelCode only", "123", "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := load[ReadRQ](t, "GuestRequests-OTA_ReadRQ.xml")
			rq.ReadRequest.HotelCode, rq.ReadRequest.HotelName = tt.hotelCode, tt.hotelName
			rulestest.Expect(t, ValidateReadRQ(rq, nil), tt.want)
		})
	}
}

func TestValidateNotifReportRQ(t *testing.T) {
	rq := load[NotifReportRQ](t, "GuestRequests-OTA_ReadRQ-ack.xml")
	rulestest.Expect(t, ValidateNotifReportRQ(rq, nil), nil)

	(*rq.Warnings)[0].RecordID = " "
	(*rq.HotelReservations)[1].UniqueID.ID = " "
	rulestest.Expect(t, ValidateNotifReportRQ(rq, nil), []string{
		"Warning[0]: RecordID required",
		"HotelReservation[1].UniqueID: ID required",
	})
}

func TestValidateResRetrieveRSSamples(t *testing.T) {
	tests := []struct {
		file string
		want []string
	}{
		{"GuestRequests-OTA_ResRetrieveRS-reservation.xml", nil},
		{"GuestRequests-OTA_ResRetrieveRS-cancellation.xml", nil},
		{"GuestRequests-OTA_ResRetrieveRS-reservation-empty.xml", nil},
		{"GuestRequests-OTA_ResRetrieveRS-error.xml", nil},
		// The official sample omits the Duration §4.2.2 requires with a
		// StartDateWindow.
		{"GuestRequests-OTA_ResRetrieveRS-request-with-roomtype.xml", []string{
			"HotelReservation[0].RoomStay[0].TimeSpan: Duration required with StartDateWindow",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			rulestest.Expect(t, ValidateResRetrieveRS(load[ResRetrieveRS](t, tt.file)), tt.want)
		})
	}
}

func TestValidateReservation(t *testing.T) {
	stay := func(res *HotelReservation) *RoomStay { return &(*res.RoomStays)[0] }
	tests := []struct {
		name   string
		change func(res *HotelReservation)
		want   []string
	}{
		{"Modify", func(res *HotelReservation) { res.ResStatus = ResStatusModify }, nil},
		{"HotelReservation", func(res *HotelReservation) {
			res.UniqueID = UniqueID{Type: UniqueIDTypeCancellation, ID: " "}
			res.CreateDateTime = DateTime{}
			res.RoomStays, res.ResGuest, res.ResGlobalInfo = nil, nil, nil
		}, []string{
			"HotelReservation[0].UniqueID: Type must be 14 for ResStatus Reserved",
			"HotelReservation[0].UniqueID: ID required",
			"HotelReservation[0]: CreateDateTime required",
			"HotelReservation[0]: RoomStays required",
			"HotelReservation[0]: ResGuests required",
			"HotelReservation[0]: ResGlobalInfo required",
		}},
		{"RoomType", func(res *HotelReservation) {
			stay(res).RoomType = &RoomType{RoomTypeCode: " ", RoomType: 1}
		}, []string{
			"HotelReservation[0].RoomStay[0].RoomType: RoomTypeCode required",
			"HotelReservation[0].RoomStay[0].RoomType: RoomType requires RoomClassificationCode",
		}},
		{"RatePlan", func(res *HotelReservation) {
			p := stay(res).RatePlan
			p.RatePlanCode, p.MealsIncluded = " ", nil
			p.Commission = &Commission{}
		}, []string{
			"HotelReservation[0].RoomStay[0].RatePlan: RatePlanCode required",
			"HotelReservation[0].RoomStay[0].RatePlan.Commission: Percent or CommissionPayableAmount required",
			"HotelReservation[0].RoomStay[0].RatePlan: MealsIncluded required",
		}},
		{"Commission both ways", func(res *HotelReservation) {
			stay(res).RatePlan.Commission.CommissionPayableAmount = &CommissionPayableAmount{Amount: new(alpinebits.Decimal(10.0)), CurrencyCode: "EUR"}
		}, []string{
			"HotelReservation[0].RoomStay[0].RatePlan.Commission: Percent and CommissionPayableAmount must not be given together",
		}},
		{"CommissionPayableAmount", func(res *HotelReservation) {
			stay(res).RatePlan.Commission = &Commission{CommissionPayableAmount: &CommissionPayableAmount{CurrencyCode: "EUR"}}
		}, []string{
			"HotelReservation[0].RoomStay[0].RatePlan.Commission.CommissionPayableAmount: Amount and CurrencyCode required",
		}},
		{"no RoomType, RatePlan, GuestCounts and Total", func(res *HotelReservation) {
			s := stay(res)
			s.RoomType, s.RatePlan, s.GuestCounts, s.Total = nil, nil, nil, nil
		}, []string{
			"HotelReservation[0].RoomStay[0]: RoomType required",
			"HotelReservation[0].RoomStay[0]: RatePlan required",
			"HotelReservation[0].RoomStay[0]: GuestCounts required",
			"HotelReservation[0].RoomStay[0]: Total required",
		}},
		{"two adult GuestCounts", func(res *HotelReservation) {
			*stay(res).GuestCounts = append(*stay(res).GuestCounts, GuestCount{Count: 1})
		}, []string{
			"HotelReservation[0].RoomStay[0].GuestCount[3]: at most one adult GuestCount allowed",
		}},
		{"children only", func(res *HotelReservation) {
			*stay(res).GuestCounts = (*stay(res).GuestCounts)[1:]
		}, []string{
			"HotelReservation[0].RoomStay[0]: adult GuestCount required",
		}},
		{"window in a reservation", func(res *HotelReservation) {
			stay(res).TimeSpan.StartDateWindow = &StartDateWindow{EarliestDate: *date("2012-01-01"), LatestDate: *date("2012-01-05")}
		}, []string{
			"HotelReservation[0].RoomStay[0].TimeSpan: Duration and StartDateWindow not allowed in a reservation",
		}},
		{"End before Start", func(res *HotelReservation) {
			stay(res).TimeSpan.End = date("2011-12-31")
		}, []string{
			"HotelReservation[0].RoomStay[0].TimeSpan: End must be after Start",
		}},
		{"no End", func(res *HotelReservation) {
			stay(res).TimeSpan.End = nil
		}, []string{
			"HotelReservation[0].RoomStay[0].TimeSpan: Start and End required",
		}},
		{"two rooms", func(res *HotelReservation) {
			other := *stay(res)
			other.RoomType = &RoomType{RoomTypeCode: "double"}
			*res.RoomStays = append(*res.RoomStays, other)
		}, nil},
		{"two periods", func(res *HotelReservation) {
			other := *stay(res)
			other.TimeSpan = TimeSpan{Start: date("2012-01-01"), End: date("2012-01-05")}
			*res.RoomStays = append(*res.RoomStays, other)
		}, []string{
			"HotelReservation[0].RoomStay[1]: TimeSpan differs from RoomStay[0]",
		}},
		{"alternative period in a reservation", func(res *HotelReservation) {
			*res.RoomStays = append(*res.RoomStays, RoomStay{TimeSpan: stay(res).TimeSpan})
		}, []string{
			"HotelReservation[0].RoomStay[1]: RoomType required",
			"HotelReservation[0].RoomStay[1]: RatePlan required",
			"HotelReservation[0].RoomStay[1]: GuestCounts required",
			"HotelReservation[0].RoomStay[1]: Total required",
		}},
		{"PaymentCard", func(res *HotelReservation) {
			p := &stay(res).Guarantee.PaymentCard
			p.CardCode = "V"
			p.CardNumber.EncryptedValue = "x"
		}, []string{
			"HotelReservation[0].RoomStay[0].PaymentCard: CardCode must be two letters",
			"HotelReservation[0].RoomStay[0].PaymentCard.CardNumber: PlainText not allowed with EncryptedValue and EncryptionMethod",
		}},
		{"encrypted CardNumber", func(res *HotelReservation) {
			stay(res).Guarantee.PaymentCard.CardNumber = CardNumber{EncryptedValue: "x", EncryptionMethod: "RSA-PKCSV2.2"}
		}, nil},
		{"incomplete CardNumber", func(res *HotelReservation) {
			stay(res).Guarantee.PaymentCard.CardNumber = CardNumber{EncryptionMethod: "RSA-PKCSV2.2"}
		}, []string{
			"HotelReservation[0].RoomStay[0].PaymentCard.CardNumber: EncryptedValue and EncryptionMethod must be given together",
		}},
		{"empty CardNumber", func(res *HotelReservation) {
			stay(res).Guarantee.PaymentCard.CardNumber = CardNumber{}
		}, []string{
			"HotelReservation[0].RoomStay[0].PaymentCard.CardNumber: PlainText or EncryptedValue and EncryptionMethod required",
		}},
		{"Comments", func(res *HotelReservation) {
			comments := *res.ResGlobalInfo.Comments
			services, customer := comments[0], comments[1]
			services.ListItems = append(services.ListItems, ListItem{ListItem: 1, Language: "de", Value: "Sauna"}, ListItem{ListItem: 1, Language: "en", Value: "Parking"})
			*res.ResGlobalInfo.Comments = []Comment{
				services,
				{Name: CommentNameCustomerComment, ListItems: []ListItem{{ListItem: 1, Language: "de", Value: "Hunde?"}}},
				customer,
			}
		}, []string{
			"HotelReservation[0].ResGlobalInfo.Comment[0].ListItem[3]: duplicate ListItem 1 for Language de",
			"HotelReservation[0].ResGlobalInfo.Comment[1]: customer comment requires Text",
			"HotelReservation[0].ResGlobalInfo.Comment[2]: duplicate Name customer comment",
		}},
		{"included services as Text", func(res *HotelReservation) {
			(*res.ResGlobalInfo.Comments)[0] = Comment{Name: CommentNameIncludedServices, Text: "Parkplatz"}
		}, []string{
			"HotelReservation[0].ResGlobalInfo.Comment[0]: included services require ListItem",
		}},
		{"BasicPropertyInfo HotelName only", func(res *HotelReservation) {
			res.ResGlobalInfo.BasicPropertyInfo.HotelCode = " "
		}, []string{
			"HotelReservation[0].ResGlobalInfo.BasicPropertyInfo: HotelCode required",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := load[ResRetrieveRS](t, "GuestRequests-OTA_ResRetrieveRS-reservation.xml")
			tt.change(&rs.HotelReservations[0])
			rulestest.Expect(t, ValidateResRetrieveRS(rs), tt.want)
		})
	}
}

func TestValidateHotelReservation(t *testing.T) {
	rs := load[ResRetrieveRS](t, "GuestRequests-OTA_ResRetrieveRS-reservation.xml")
	res := rs.HotelReservations[0]
	rulestest.Expect(t, ValidateHotelReservation(res), nil)
	res.UniqueID.ID = " "
	rulestest.Expect(t, ValidateHotelReservation(res), []string{"UniqueID: ID required"})
}

// quoteRequest is the official quote request sample with the Duration it
// lacks.
func quoteRequest(t *testing.T) ResRetrieveRS {
	t.Helper()
	rs := load[ResRetrieveRS](t, "GuestRequests-OTA_ResRetrieveRS-request-with-roomtype.xml")
	(*rs.HotelReservations[0].RoomStays)[0].TimeSpan.Duration = new(alpinebits.Nights(7))
	return rs
}

func TestValidateQuoteRequest(t *testing.T) {
	stay := func(res *HotelReservation) *RoomStay { return &(*res.RoomStays)[0] }
	dates := TimeSpan{Start: date("2017-10-03"), End: date("2017-10-10")}
	tests := []struct {
		name   string
		change func(res *HotelReservation)
		want   []string
	}{
		{"Start and End", func(res *HotelReservation) { stay(res).TimeSpan = dates }, nil},
		{"open RoomType, RatePlan without code", func(res *HotelReservation) {
			s := stay(res)
			s.RoomType = &RoomType{}
			s.RatePlan = &RatePlan{}
		}, nil},
		{"no RoomType", func(res *HotelReservation) { stay(res).RoomType = nil }, []string{
			"HotelReservation[0].RoomStay[0]: RoomType required",
		}},
		{"no RoomStays", func(res *HotelReservation) { res.RoomStays = nil }, []string{
			"HotelReservation[0]: RoomStays required",
		}},
		{"window with Start", func(res *HotelReservation) {
			ts := &stay(res).TimeSpan
			ts.Start = date("2017-10-03")
			ts.StartDateWindow.LatestDate = ts.StartDateWindow.EarliestDate
		}, []string{
			"HotelReservation[0].RoomStay[0].TimeSpan: Start and End not allowed with Duration and StartDateWindow",
			"HotelReservation[0].RoomStay[0].TimeSpan.StartDateWindow: LatestDate must be after EarliestDate",
		}},
		{"Duration without StartDateWindow", func(res *HotelReservation) {
			stay(res).TimeSpan.StartDateWindow = nil
		}, []string{
			"HotelReservation[0].RoomStay[0].TimeSpan: StartDateWindow required with Duration",
		}},
		{"alternative period", func(res *HotelReservation) {
			*res.RoomStays = append(*res.RoomStays, RoomStay{TimeSpan: dates})
		}, nil},
		{"alternative period not different", func(res *HotelReservation) {
			*res.RoomStays = append(*res.RoomStays, RoomStay{TimeSpan: stay(res).TimeSpan})
		}, []string{
			"HotelReservation[0].RoomStay[1]: alternative TimeSpan must differ from RoomStay[0]",
		}},
		{"two alternative periods", func(res *HotelReservation) {
			later := TimeSpan{Start: date("2017-11-03"), End: date("2017-11-10")}
			*res.RoomStays = append(*res.RoomStays, RoomStay{TimeSpan: dates}, RoomStay{TimeSpan: later})
		}, []string{
			"HotelReservation[0].RoomStay[1]: RoomType required",
			"HotelReservation[0].RoomStay[1]: GuestCounts required",
			"HotelReservation[0].RoomStay[1]: TimeSpan differs from RoomStay[0]",
		}},
		{"alternative period alone", func(res *HotelReservation) {
			*res.RoomStays = []RoomStay{{TimeSpan: dates}}
		}, []string{
			"HotelReservation[0].RoomStay[0]: RoomType required",
			"HotelReservation[0].RoomStay[0]: GuestCounts required",
		}},
		{"Guarantee and CancelPenalties", func(res *HotelReservation) {
			stay(res).Guarantee = &Guarantee{PaymentCard{CardCode: "VI", ExpireDate: "1216", CardHolderName: "Otto", CardNumber: CardNumber{PlainText: "1111"}}}
			res.ResGlobalInfo.CancelPenalty = &CancelPenalty{PenaltyDescription{Text: "free"}}
		}, []string{
			"HotelReservation[0].RoomStay[0]: Guarantee not allowed in a quote request",
			"HotelReservation[0].ResGlobalInfo: CancelPenalties not allowed in a quote request",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := quoteRequest(t)
			rulestest.Expect(t, ValidateResRetrieveRS(rs), nil)
			tt.change(&rs.HotelReservations[0])
			rulestest.Expect(t, ValidateResRetrieveRS(rs), tt.want)
		})
	}
}

func TestValidateCancellation(t *testing.T) {
	reservation := load[ResRetrieveRS](t, "GuestRequests-OTA_ResRetrieveRS-reservation.xml").HotelReservations[0]
	request := quoteRequest(t).HotelReservations[0]
	tests := []struct {
		name   string
		change func(res *HotelReservation)
		want   []string
	}{
		{"Type 14", func(res *HotelReservation) { res.UniqueID.Type = UniqueIDTypeReservation }, []string{
			"HotelReservation[0].UniqueID: Type must be 15 for ResStatus Cancelled",
		}},
		{"repeats the reservation", func(res *HotelReservation) {
			res.RoomStays, res.ResGuest, res.ResGlobalInfo = reservation.RoomStays, reservation.ResGuest, reservation.ResGlobalInfo
		}, nil},
		{"repeats the quote request", func(res *HotelReservation) {
			stays := append(slices.Clone(*request.RoomStays), RoomStay{TimeSpan: TimeSpan{Start: date("2017-10-03"), End: date("2017-10-10")}})
			res.RoomStays, res.ResGuest, res.ResGlobalInfo = &stays, request.ResGuest, request.ResGlobalInfo
		}, nil},
		{"broken parts", func(res *HotelReservation) {
			res.RoomStays = &[]RoomStay{{GuestCounts: &[]GuestCount{{Count: 1}}, TimeSpan: TimeSpan{Start: date("2017-10-03")}}}
			res.ResGlobalInfo = &ResGlobalInfo{}
		}, []string{
			"HotelReservation[0].RoomStay[0]: RoomType required",
			"HotelReservation[0].RoomStay[0].TimeSpan: Start and End required",
			"HotelReservation[0].ResGlobalInfo.BasicPropertyInfo: HotelCode required",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := load[ResRetrieveRS](t, "GuestRequests-OTA_ResRetrieveRS-cancellation.xml")
			tt.change(&rs.HotelReservations[0])
			rulestest.Expect(t, ValidateResRetrieveRS(rs), tt.want)
		})
	}
}

func TestDateTime(t *testing.T) {
	tests := []struct {
		in       string
		noOffset bool
		utc      string
	}{
		{"2012-03-21T15:00:00+01:00", false, "2012-03-21T14:00:00Z"},
		{"2012-03-21T15:00:00.5Z", false, "2012-03-21T15:00:00.5Z"},
		{"2017-09-03T19:47:50", true, "2017-09-03T19:47:50Z"},
		{"2017-09-03T19:47:50.25", true, "2017-09-03T19:47:50.25Z"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			var d DateTime
			if err := d.UnmarshalText([]byte(tt.in)); err != nil {
				t.Fatal(err)
			}
			if d.NoOffset != tt.noOffset || d.Time.UTC().Format("2006-01-02T15:04:05.999999999Z07:00") != tt.utc {
				t.Errorf("got %v NoOffset %v", d.Time, d.NoOffset)
			}
			out, err := d.MarshalText()
			if err != nil {
				t.Fatal(err)
			}
			if string(out) != tt.in {
				t.Errorf("marshals as %s", out)
			}
		})
	}
	var d DateTime
	if err := d.UnmarshalText([]byte("2017-09-03")); err == nil {
		t.Error("accepted a date")
	}
}
