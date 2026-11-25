package guestrequests

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
	guestrequests201810 "github.com/HGV/alpinebits/v201810/guestrequests"
	"github.com/HGV/x/timex"
)

// 2020-10 kept the 2018-10 samples.
const samples = "../../v201810/guestrequests/testdata/"

const (
	reservationSample  = "GuestRequests-OTA_ResRetrieveRS-reservation.xml"
	quoteSample        = "GuestRequests-OTA_ResRetrieveRS-request-with-roomtype.xml"
	cancellationSample = "GuestRequests-OTA_ResRetrieveRS-cancellation.xml"
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
	b, err := os.ReadFile(samples + name)
	if err != nil {
		t.Fatal(err)
	}
	var v T
	if err := xml.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestValidateRequests(t *testing.T) {
	rq := load[ReadRQ](t, "GuestRequests-OTA_ReadRQ.xml")
	rulestest.Expect(t, ValidateReadRQ(rq, nil), nil)
	rq.ReadRequest.HotelCode, rq.ReadRequest.HotelName = " ", ""
	rulestest.Expect(t, ValidateReadRQ(rq, nil), []string{"HotelCode required"})

	ack := load[NotifReportRQ](t, "GuestRequests-OTA_ReadRQ-ack.xml")
	rulestest.Expect(t, ValidateNotifReportRQ(ack, nil), nil)
	(*ack.Warnings)[0].RecordID = " "
	rulestest.Expect(t, ValidateNotifReportRQ(ack, nil), []string{"Warning[0]: RecordID required"})
}

func TestValidateResRetrieveRSSamples(t *testing.T) {
	tests := []struct {
		file string
		want []string
	}{
		{reservationSample, nil},
		{cancellationSample, nil},
		{"GuestRequests-OTA_ResRetrieveRS-reservation-empty.xml", nil},
		{"GuestRequests-OTA_ResRetrieveRS-error.xml", nil},
		// The official sample omits the Duration §4.2.2 requires with a
		// StartDateWindow.
		{quoteSample, []string{
			"HotelReservation[0].RoomStay[0].TimeSpan: Duration required with StartDateWindow",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			rulestest.Expect(t, ValidateResRetrieveRS(load[ResRetrieveRS](t, tt.file)), tt.want)
		})
	}
}

func from201810(res guestrequests201810.HotelReservation) HotelReservation {
	out := HotelReservation{
		UniqueID: res.UniqueID, RoomStays: res.RoomStays, ResGuest: res.ResGuest,
		CreateDateTime: res.CreateDateTime, ResStatus: res.ResStatus,
	}
	if g := res.ResGlobalInfo; g != nil {
		out.ResGlobalInfo = &ResGlobalInfo{
			Comments: g.Comments, CancelPenalty: g.CancelPenalty, HotelReservationIDs: g.HotelReservationIDs,
			GlobalProfile: g.GlobalProfile, BasicPropertyInfo: g.BasicPropertyInfo,
		}
	}
	return out
}

// The checks 2020-10 copied from 2018-10 must report what 2018-10 reports.
func TestRulesMatch201810(t *testing.T) {
	stay := func(res *guestrequests201810.HotelReservation) *RoomStay { return &(*res.RoomStays)[0] }
	tests := []struct {
		name, file string
		change     func(res *guestrequests201810.HotelReservation)
	}{
		{"HotelReservation", reservationSample, func(res *guestrequests201810.HotelReservation) {
			res.UniqueID = UniqueID{Type: UniqueIDTypeCancellation, ID: " "}
			res.CreateDateTime = DateTime{}
			res.RoomStays, res.ResGuest, res.ResGlobalInfo = nil, nil, nil
		}},
		{"RoomType", reservationSample, func(res *guestrequests201810.HotelReservation) {
			stay(res).RoomType = &RoomType{RoomTypeCode: " ", RoomType: 1}
		}},
		{"RatePlan", reservationSample, func(res *guestrequests201810.HotelReservation) {
			p := stay(res).RatePlan
			p.RatePlanCode, p.MealsIncluded, p.Commission = " ", nil, &Commission{}
		}},
		{"Commission both ways", reservationSample, func(res *guestrequests201810.HotelReservation) {
			stay(res).RatePlan.Commission.CommissionPayableAmount = &CommissionPayableAmount{Amount: new(alpinebits.Decimal(10.0)), CurrencyCode: "EUR"}
		}},
		{"CommissionPayableAmount", reservationSample, func(res *guestrequests201810.HotelReservation) {
			stay(res).RatePlan.Commission = &Commission{CommissionPayableAmount: &CommissionPayableAmount{CurrencyCode: "EUR"}}
		}},
		{"reservation without RoomType, RatePlan, GuestCounts and Total", reservationSample, func(res *guestrequests201810.HotelReservation) {
			s := stay(res)
			s.RoomType, s.RatePlan, s.GuestCounts, s.Total = nil, nil, nil, nil
		}},
		{"two adult GuestCounts", reservationSample, func(res *guestrequests201810.HotelReservation) {
			*stay(res).GuestCounts = append(*stay(res).GuestCounts, GuestCount{Count: 1})
		}},
		{"window in a reservation", reservationSample, func(res *guestrequests201810.HotelReservation) {
			stay(res).TimeSpan.StartDateWindow = &StartDateWindow{EarliestDate: *date("2012-01-01"), LatestDate: *date("2012-01-05")}
		}},
		{"End before Start", reservationSample, func(res *guestrequests201810.HotelReservation) {
			stay(res).TimeSpan.End = date("2011-12-31")
		}},
		{"no End", reservationSample, func(res *guestrequests201810.HotelReservation) { stay(res).TimeSpan.End = nil }},
		{"two periods", reservationSample, func(res *guestrequests201810.HotelReservation) {
			other := *stay(res)
			other.TimeSpan = TimeSpan{Start: date("2012-01-01"), End: date("2012-01-05")}
			*res.RoomStays = append(*res.RoomStays, other)
		}},
		{"alternative period in a reservation", reservationSample, func(res *guestrequests201810.HotelReservation) {
			*res.RoomStays = append(*res.RoomStays, RoomStay{TimeSpan: stay(res).TimeSpan})
		}},
		{"PaymentCard", reservationSample, func(res *guestrequests201810.HotelReservation) {
			p := &stay(res).Guarantee.PaymentCard
			p.CardCode, p.CardNumber.EncryptedValue = "V", "x"
		}},
		{"incomplete CardNumber", reservationSample, func(res *guestrequests201810.HotelReservation) {
			stay(res).Guarantee.PaymentCard.CardNumber = CardNumber{EncryptionMethod: "RSA-PKCSV2.2"}
		}},
		{"empty CardNumber", reservationSample, func(res *guestrequests201810.HotelReservation) {
			stay(res).Guarantee.PaymentCard.CardNumber = CardNumber{}
		}},
		{"Comments", reservationSample, func(res *guestrequests201810.HotelReservation) {
			comments := *res.ResGlobalInfo.Comments
			services, customer := comments[0], comments[1]
			services.ListItems = append(services.ListItems, ListItem{ListItem: 1, Language: "de", Value: "Sauna"})
			*res.ResGlobalInfo.Comments = []Comment{
				services,
				{Name: CommentNameCustomerComment, ListItems: []ListItem{{ListItem: 1, Language: "de", Value: "Hunde?"}}},
				customer,
				{Name: CommentNameIncludedServices, Text: "Parkplatz"},
			}
		}},
		{"BasicPropertyInfo blank", reservationSample, func(res *guestrequests201810.HotelReservation) {
			res.ResGlobalInfo.BasicPropertyInfo = BasicPropertyInfo{HotelCode: " "}
		}},
		{"window with Start", quoteSample, func(res *guestrequests201810.HotelReservation) {
			ts := &stay(res).TimeSpan
			ts.Start = date("2017-10-03")
			ts.StartDateWindow.LatestDate = ts.StartDateWindow.EarliestDate
		}},
		{"Duration without StartDateWindow", quoteSample, func(res *guestrequests201810.HotelReservation) {
			stay(res).TimeSpan.Duration = new(alpinebits.Nights(7))
			stay(res).TimeSpan.StartDateWindow = nil
		}},
		{"alternative period not different", quoteSample, func(res *guestrequests201810.HotelReservation) {
			*res.RoomStays = append(*res.RoomStays, RoomStay{TimeSpan: stay(res).TimeSpan})
		}},
		{"Guarantee and CancelPenalties in a quote request", quoteSample, func(res *guestrequests201810.HotelReservation) {
			stay(res).Guarantee = &Guarantee{PaymentCard: PaymentCard{CardCode: "VI", ExpireDate: "1216", CardHolderName: "Otto", CardNumber: CardNumber{PlainText: "1111"}}}
			res.ResGlobalInfo.CancelPenalty = &CancelPenalty{PenaltyDescription: PenaltyDescription{Text: "free"}}
		}},
		{"cancellation with Type 14", cancellationSample, func(res *guestrequests201810.HotelReservation) {
			res.UniqueID.Type = UniqueIDTypeReservation
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := load[guestrequests201810.ResRetrieveRS](t, tt.file).HotelReservations[0]
			tt.change(&res)
			var want []string
			for _, e := range guestrequests201810.ValidateHotelReservation(res) {
				want = append(want, e.Message)
			}
			if len(want) == 0 {
				t.Fatal("2018-10 reports no errors")
			}
			rulestest.Expect(t, ValidateHotelReservation(from201810(res)), want)
		})
	}
}

// quoteRequest is the official quote request sample with the Duration it
// lacks.
func quoteRequest(t *testing.T) ResRetrieveRS {
	t.Helper()
	rs := load[ResRetrieveRS](t, quoteSample)
	(*rs.HotelReservations[0].RoomStays)[0].TimeSpan.Duration = new(alpinebits.Nights(7))
	return rs
}

func TestValidate2020Changes(t *testing.T) {
	stay := func(res *HotelReservation) *RoomStay { return &(*res.RoomStays)[0] }
	dates := TimeSpan{Start: date("2017-10-03"), End: date("2017-10-10")}
	tests := []struct {
		name   string
		rs     func(t *testing.T) ResRetrieveRS
		change func(res *HotelReservation)
		want   []string
	}{
		{"quote request without RoomType", quoteRequest, func(res *HotelReservation) { stay(res).RoomType = nil }, nil},
		{"cancellation without RoomType", func(t *testing.T) ResRetrieveRS { return load[ResRetrieveRS](t, cancellationSample) }, func(res *HotelReservation) {
			res.RoomStays = &[]RoomStay{{GuestCounts: &[]GuestCount{{Count: 1}}, TimeSpan: dates}}
		}, nil},
		{"children only", reservation, func(res *HotelReservation) {
			*stay(res).GuestCounts = (*stay(res).GuestCounts)[1:]
		}, nil},
		{"no guests", reservation, func(res *HotelReservation) { stay(res).GuestCounts = &[]GuestCount{} }, []string{
			"HotelReservation[0].RoomStay[0]: GuestCount with at least one guest required",
		}},
		{"other currency", reservation, func(res *HotelReservation) {
			stay(res).Total.CurrencyCode = "CHF"
			stay(res).RatePlan.Commission = &Commission{CommissionPayableAmount: &CommissionPayableAmount{Amount: new(alpinebits.Decimal(10.5)), CurrencyCode: "CHF"}}
		}, nil},
		{"Total", reservation, func(res *HotelReservation) {
			stay(res).Total.CurrencyCode = "EUX"
			*res.RoomStays = append(*res.RoomStays, *stay(res))
			(*res.RoomStays)[1].Total = &Total{AmountAfterTax: 299.999, CurrencyCode: "EUR"}
		}, []string{
			"HotelReservation[0].RoomStay[0].Total: CurrencyCode EUX is not an ISO 4217 code",
			"HotelReservation[0].RoomStay[1].Total: AmountAfterTax has more than 2 decimal digits for EUR",
		}},
		{"CommissionPayableAmount", reservation, func(res *HotelReservation) {
			stay(res).RatePlan.Commission = &Commission{CommissionPayableAmount: &CommissionPayableAmount{Amount: new(alpinebits.Decimal(10.5)), CurrencyCode: "JPY"}}
		}, []string{
			"HotelReservation[0].RoomStay[0].RatePlan.Commission.CommissionPayableAmount: Amount has more than 0 decimal digits for JPY",
		}},
		{"pets", reservation, func(res *HotelReservation) {
			res.ResGlobalInfo.SpecialRequest = &SpecialRequest{Name: SpecialRequestNamePets, Text: &alpinebits.Text{TextFormat: alpinebits.TextFormatPlainText, Value: "two dogs"}}
		}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := tt.rs(t)
			tt.change(&rs.HotelReservations[0])
			rulestest.Expect(t, ValidateResRetrieveRS(rs), tt.want)
		})
	}
}

func reservation(t *testing.T) ResRetrieveRS {
	return load[ResRetrieveRS](t, reservationSample)
}
