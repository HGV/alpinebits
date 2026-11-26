package guestrequests

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
	guestrequests202010 "github.com/HGV/alpinebits/v202010/guestrequests"
	"github.com/HGV/x/timex"
)

// The empty and error responses are those of 2018-10.
const samples2018 = "../../v201810/guestrequests/testdata/"

func date(s string) *timex.Date {
	d, err := timex.ParseDate(s)
	if err != nil {
		panic(err)
	}
	return &d
}

func load[T any](t *testing.T, path string) T {
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

func TestValidateSamples(t *testing.T) {
	tests := []struct {
		path string
		want []string
	}{
		{"testdata/GuestRequests-OTA_ResRetrieveRS-reservation.xml", nil},
		{"testdata/GuestRequests-OTA_ResRetrieveRS-cancellation.xml", nil},
		{"testdata/GuestRequests-OTA_ResRetrieveRS-newsletter.xml", nil},
		// The official sample omits the Duration §4.2.2 requires with a
		// StartDateWindow.
		{"testdata/GuestRequests-OTA_ResRetrieveRS-quote-request.xml", []string{
			"HotelReservation[0].RoomStay[0].TimeSpan: Duration required with StartDateWindow",
		}},
		{samples2018 + "GuestRequests-OTA_ResRetrieveRS-reservation-empty.xml", nil},
		{samples2018 + "GuestRequests-OTA_ResRetrieveRS-error.xml", nil},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			rulestest.Expect(t, ValidateResRetrieveRS(load[ResRetrieveRS](t, tt.path)), tt.want)
		})
	}
}

func from202010(res guestrequests202010.HotelReservation) HotelReservation {
	out := HotelReservation{UniqueID: res.UniqueID, ResGuest: res.ResGuest, CreateDateTime: res.CreateDateTime, ResStatus: res.ResStatus}
	if res.RoomStays != nil {
		stays := []RoomStay{}
		for _, s := range *res.RoomStays {
			stays = append(stays, RoomStay{
				RoomType: s.RoomType, RatePlan: s.RatePlan, GuestCounts: s.GuestCounts,
				TimeSpan: s.TimeSpan, Guarantee: s.Guarantee, Total: s.Total,
			})
		}
		out.RoomStays = &stays
	}
	if g := res.ResGlobalInfo; g != nil {
		info := ResGlobalInfo{
			Comments: g.Comments, CancelPenalty: g.CancelPenalty, HotelReservationIDs: g.HotelReservationIDs,
			GlobalProfile: g.GlobalProfile, BasicPropertyInfo: g.BasicPropertyInfo,
		}
		if r := g.SpecialRequest; r != nil {
			info.SpecialRequests = &[]SpecialRequest{{RequestCode: RequestCodePets, CodeContext: CodeContextAlpineBits, Text: r.Text}}
		}
		out.ResGlobalInfo = &info
	}
	return out
}

// The checks 2022-10 copied from 2020-10 must report what 2020-10 reports.
func TestRulesMatch202010(t *testing.T) {
	type res = guestrequests202010.HotelReservation
	stay := func(r *res) *guestrequests202010.RoomStay { return &(*r.RoomStays)[0] }
	const (
		reservation  = samples2018 + "GuestRequests-OTA_ResRetrieveRS-reservation.xml"
		quote        = samples2018 + "GuestRequests-OTA_ResRetrieveRS-request-with-roomtype.xml"
		cancellation = samples2018 + "GuestRequests-OTA_ResRetrieveRS-cancellation.xml"
	)
	tests := []struct {
		name, file string
		change     func(r *res)
	}{
		{"HotelReservation", reservation, func(r *res) {
			r.UniqueID = UniqueID{Type: UniqueIDTypeCancellation, ID: " "}
			r.CreateDateTime = DateTime{}
			r.RoomStays, r.ResGuest, r.ResGlobalInfo = nil, nil, nil
		}},
		{"RoomType", reservation, func(r *res) { stay(r).RoomType = &RoomType{RoomTypeCode: " ", RoomType: 1} }},
		{"RatePlan", reservation, func(r *res) {
			p := stay(r).RatePlan
			p.RatePlanCode, p.MealsIncluded, p.Commission = " ", nil, &Commission{}
		}},
		{"CommissionPayableAmount", reservation, func(r *res) {
			stay(r).RatePlan.Commission = &Commission{CommissionPayableAmount: &CommissionPayableAmount{Amount: new(alpinebits.Decimal(10.5)), CurrencyCode: "JPY"}}
		}},
		{"reservation without RoomType, RatePlan, GuestCounts and Total", reservation, func(r *res) {
			s := stay(r)
			s.RoomType, s.RatePlan, s.GuestCounts, s.Total = nil, nil, nil, nil
		}},
		{"GuestCounts", reservation, func(r *res) {
			*stay(r).GuestCounts = append(*stay(r).GuestCounts, GuestCount{Count: 1})
		}},
		{"no guests", reservation, func(r *res) { stay(r).GuestCounts = &[]GuestCount{} }},
		{"window in a reservation", reservation, func(r *res) {
			stay(r).TimeSpan.StartDateWindow = &StartDateWindow{EarliestDate: *date("2012-01-01"), LatestDate: *date("2012-01-05")}
		}},
		{"End before Start", reservation, func(r *res) { stay(r).TimeSpan.End = date("2011-12-31") }},
		{"PaymentCard", reservation, func(r *res) {
			p := &stay(r).Guarantee.PaymentCard
			p.CardCode, p.CardNumber.EncryptedValue = "V", "x"
		}},
		{"empty CardNumber", reservation, func(r *res) { stay(r).Guarantee.PaymentCard.CardNumber = CardNumber{} }},
		{"Total", reservation, func(r *res) { stay(r).Total = &Total{AmountAfterTax: 299.999, CurrencyCode: "EUX"} }},
		{"Comments", reservation, func(r *res) {
			comments := *r.ResGlobalInfo.Comments
			*r.ResGlobalInfo.Comments = []Comment{comments[0], comments[0], {Name: CommentNameCustomerComment}}
		}},
		{"BasicPropertyInfo blank", reservation, func(r *res) {
			r.ResGlobalInfo.BasicPropertyInfo = BasicPropertyInfo{HotelCode: " "}
		}},
		{"window with Start", quote, func(r *res) {
			ts := &stay(r).TimeSpan
			ts.Start = date("2017-10-03")
			ts.StartDateWindow.LatestDate = ts.StartDateWindow.EarliestDate
		}},
		{"Guarantee and CancelPenalties in a quote request", quote, func(r *res) {
			stay(r).Guarantee = &Guarantee{PaymentCard: PaymentCard{CardCode: "VI", ExpireDate: "1216", CardHolderName: "Otto", CardNumber: CardNumber{PlainText: "1111"}}}
			r.ResGlobalInfo.CancelPenalty = &CancelPenalty{PenaltyDescription: PenaltyDescription{Text: "free"}}
		}},
		{"cancellation with Type 14", cancellation, func(r *res) { r.UniqueID.Type = UniqueIDTypeReservation }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := load[guestrequests202010.ResRetrieveRS](t, tt.file).HotelReservations[0]
			tt.change(&old)
			var want []string
			for _, e := range guestrequests202010.ValidateHotelReservation(old) {
				want = append(want, e.Message)
			}
			if len(want) == 0 {
				t.Fatal("2020-10 reports no errors")
			}
			rulestest.Expect(t, ValidateHotelReservation(from202010(old)), want)
		})
	}
}

func TestValidate2022Changes(t *testing.T) {
	const (
		reservation = "testdata/GuestRequests-OTA_ResRetrieveRS-reservation.xml"
		quote       = "testdata/GuestRequests-OTA_ResRetrieveRS-quote-request.xml"
		newsletter  = "testdata/GuestRequests-OTA_ResRetrieveRS-newsletter.xml"
	)
	stay := func(res *HotelReservation) *RoomStay { return &(*res.RoomStays)[0] }
	special := func(context CodeContext, code string) SpecialRequest {
		return SpecialRequest{CodeContext: context, RequestCode: code, Text: &alpinebits.Text{TextFormat: alpinebits.TextFormatPlainText}}
	}
	tests := []struct {
		name, file string
		change     func(res *HotelReservation)
		want       []string
	}{
		{"quote request without RoomStays", quote, func(res *HotelReservation) { res.RoomStays = nil }, nil},
		{"reservation without RoomStays", reservation, func(res *HotelReservation) { res.RoomStays = nil }, []string{
			"HotelReservation[0]: RoomStays required",
		}},
		{"newsletter with RoomStays", newsletter, func(res *HotelReservation) { res.RoomStays = &[]RoomStay{} }, []string{
			"HotelReservation[0]: RoomStays not allowed in a newsletter or catalog request",
		}},
		{"RoomStayReservation false in a reservation", reservation, func(res *HotelReservation) {
			res.RoomStayReservation = new(false)
		}, []string{
			"HotelReservation[0]: RoomStayReservation false not allowed in a reservation",
		}},
		{"stays of different periods", reservation, func(res *HotelReservation) {
			other := *stay(res)
			other.TimeSpan = TimeSpan{Start: date("2022-01-12"), End: date("2022-01-15")}
			*res.RoomStays = append(*res.RoomStays, other)
		}, nil},
		{"alternatives in a quote request", quote, func(res *HotelReservation) {
			other := *stay(res)
			stay(res).RoomStayGroupID, other.RoomStayGroupID = "1", "2"
			*res.RoomStays = append(*res.RoomStays, other)
		}, nil},
		{"RoomStayGroupID in a reservation", reservation, func(res *HotelReservation) { stay(res).RoomStayGroupID = "1" }, []string{
			"HotelReservation[0].RoomStay[0]: RoomStayGroupID not allowed in a reservation",
		}},
		{"SpecialRequests", reservation, func(res *HotelReservation) {
			res.ResGlobalInfo.SpecialRequests = &[]SpecialRequest{
				special(CodeContextAlpineBits, RequestCodePets),
				special(CodeContextHotel, "VEGETARIAN_GUEST"),
				special(CodeContextHotel, "PETS"),
			}
		}, nil},
		{"broken SpecialRequests", reservation, func(res *HotelReservation) {
			res.ResGlobalInfo.SpecialRequests = &[]SpecialRequest{
				special(CodeContextAlpineBits, "VEGETARIAN_GUEST"),
				special(CodeContextHotel, "X"),
				special(CodeContextHotel, "X"),
				special(CodeContextHotel, " "),
			}
		}, []string{
			"HotelReservation[0].ResGlobalInfo.SpecialRequest[0]: RequestCode VEGETARIAN_GUEST not defined in CodeContext ALPINEBITS",
			"HotelReservation[0].ResGlobalInfo.SpecialRequest[2]: duplicate RequestCode X in CodeContext HOTEL",
			"HotelReservation[0].ResGlobalInfo.SpecialRequest[3]: RequestCode required",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := load[ResRetrieveRS](t, tt.file)
			if tt.file == quote {
				stay(&rs.HotelReservations[0]).TimeSpan.Duration = new(alpinebits.Nights(5)) // missing in the sample
			}
			tt.change(&rs.HotelReservations[0])
			rulestest.Expect(t, ValidateResRetrieveRS(rs), tt.want)
		})
	}
}
