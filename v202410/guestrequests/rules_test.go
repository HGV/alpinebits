package guestrequests

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
	guestrequests202210 "github.com/HGV/alpinebits/v202210/guestrequests"
	"github.com/HGV/x/timex"
)

const (
	reservationSample = "testdata/GuestRequests-OTA_ResRetrieveRS-reservation.xml"
	// The other 2024-10 samples are those of 2022-10 and 2018-10.
	samples2022 = "../../v202210/guestrequests/testdata/GuestRequests-"
	samples2018 = "../../v201810/guestrequests/testdata/GuestRequests-"
)

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
		{reservationSample, nil},
		{samples2022 + "OTA_ResRetrieveRS-cancellation.xml", nil},
		{samples2022 + "OTA_ResRetrieveRS-newsletter.xml", nil},
		// The official sample omits the Duration §4.2.2 requires with a
		// StartDateWindow.
		{samples2022 + "OTA_ResRetrieveRS-quote-request.xml", []string{
			"HotelReservation[0].RoomStay[0].TimeSpan: Duration required with StartDateWindow",
		}},
		{samples2018 + "OTA_ResRetrieveRS-reservation-empty.xml", nil},
		{samples2018 + "OTA_ResRetrieveRS-error.xml", nil},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			rulestest.Expect(t, ValidateResRetrieveRS(load[ResRetrieveRS](t, tt.path)), tt.want)
		})
	}
}

func from202210(res guestrequests202210.HotelReservation) HotelReservation {
	out := HotelReservation{
		UniqueID: res.UniqueID, ResGuest: res.ResGuest, CreateDateTime: res.CreateDateTime,
		LastModifyDateTime: res.LastModifyDateTime, ResStatus: res.ResStatus, RoomStayReservation: res.RoomStayReservation,
	}
	if res.RoomStays != nil {
		stays := []RoomStay{}
		for _, s := range *res.RoomStays {
			stays = append(stays, RoomStay{
				RoomStayGroupID: s.RoomStayGroupID, RoomType: s.RoomType, RatePlan: s.RatePlan, GuestCounts: s.GuestCounts,
				TimeSpan: s.TimeSpan, Guarantee: s.Guarantee, Total: s.Total,
			})
		}
		out.RoomStays = &stays
	}
	if g := res.ResGlobalInfo; g != nil {
		out.ResGlobalInfo = &ResGlobalInfo{
			Comments: g.Comments, SpecialRequests: g.SpecialRequests, CancelPenalty: g.CancelPenalty,
			HotelReservationIDs: g.HotelReservationIDs, GlobalProfile: g.GlobalProfile, BasicPropertyInfo: g.BasicPropertyInfo,
		}
	}
	return out
}

// The checks 2024-10 copied from 2022-10 must report what 2022-10 reports.
func TestRulesMatch202210(t *testing.T) {
	type res = guestrequests202210.HotelReservation
	stay := func(r *res) *guestrequests202210.RoomStay { return &(*r.RoomStays)[0] }
	const (
		reservation  = samples2022 + "OTA_ResRetrieveRS-reservation.xml"
		quote        = samples2022 + "OTA_ResRetrieveRS-quote-request.xml"
		cancellation = samples2022 + "OTA_ResRetrieveRS-cancellation.xml"
		newsletter   = samples2022 + "OTA_ResRetrieveRS-newsletter.xml"
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
		{"RoomStayReservation false in a reservation", reservation, func(r *res) { r.RoomStayReservation = new(false) }},
		{"RoomType and RatePlan", reservation, func(r *res) {
			s := stay(r)
			s.RoomType = &RoomType{RoomTypeCode: " ", RoomType: 1}
			s.RatePlan.RatePlanCode, s.RatePlan.MealsIncluded = " ", nil
			s.RatePlan.Commission = &Commission{CommissionPayableAmount: &CommissionPayableAmount{Amount: new(alpinebits.Decimal(10.5)), CurrencyCode: "JPY"}}
		}},
		{"Commissions", reservation, func(r *res) {
			var stays []guestrequests202210.RoomStay
			for _, m := range []Commission{{}, {Percent: new(10), CommissionPayableAmount: &CommissionPayableAmount{}}, {CommissionPayableAmount: &CommissionPayableAmount{}}} {
				s, p := *stay(r), *stay(r).RatePlan
				p.Commission = &m
				s.RatePlan = &p
				stays = append(stays, s)
			}
			*r.RoomStays = stays
		}},
		{"reservation without RoomType, RatePlan, GuestCounts and Total", reservation, func(r *res) {
			s := stay(r)
			s.RoomType, s.RatePlan, s.GuestCounts, s.Total = nil, nil, nil, nil
		}},
		{"GuestCounts", reservation, func(r *res) {
			other := *stay(r)
			other.GuestCounts = &[]GuestCount{}
			*stay(r).GuestCounts = append(*stay(r).GuestCounts, GuestCount{Count: 1})
			*r.RoomStays = append(*r.RoomStays, other)
		}},
		{"TimeSpan", reservation, func(r *res) {
			other := *stay(r)
			other.TimeSpan.End = nil
			stay(r).TimeSpan.End = date("2011-12-31")
			stay(r).TimeSpan.StartDateWindow = &StartDateWindow{EarliestDate: *date("2012-01-01"), LatestDate: *date("2012-01-05")}
			*r.RoomStays = append(*r.RoomStays, other)
		}},
		{"windows in a quote request", quote, func(r *res) {
			other := *stay(r)
			other.TimeSpan = TimeSpan{Duration: new(alpinebits.Nights(5))}
			ts := &stay(r).TimeSpan
			ts.Start = date("2022-10-03")
			ts.StartDateWindow.LatestDate = ts.StartDateWindow.EarliestDate
			*r.RoomStays = append(*r.RoomStays, other)
		}},
		{"PaymentCards", reservation, func(r *res) {
			var stays []guestrequests202210.RoomStay
			for _, n := range []CardNumber{{}, {PlainText: "1111", EncryptedValue: "x", EncryptionMethod: "m"}, {EncryptedValue: "x"}} {
				s := *stay(r)
				s.Guarantee = &Guarantee{PaymentCard: PaymentCard{CardCode: "V", CardNumber: n}}
				stays = append(stays, s)
			}
			*r.RoomStays = stays
		}},
		{"Total", reservation, func(r *res) { stay(r).Total = &Total{AmountAfterTax: 299.999, CurrencyCode: "EUX"} }},
		{"RoomStayGroupID in a reservation", reservation, func(r *res) { stay(r).RoomStayGroupID = "1" }},
		{"ResGlobalInfo", reservation, func(r *res) {
			g := r.ResGlobalInfo
			comments := *g.Comments
			comments[0].ListItems = append(comments[0].ListItems, comments[0].ListItems[0])
			*g.Comments = []Comment{comments[0], {Name: CommentNameIncludedServices}, {Name: CommentNameCustomerComment}}
			g.SpecialRequests = &[]SpecialRequest{
				{CodeContext: CodeContextAlpineBits, RequestCode: "X"},
				{CodeContext: CodeContextHotel, RequestCode: " "},
				{CodeContext: CodeContextHotel, RequestCode: "Y"},
				{CodeContext: CodeContextHotel, RequestCode: "Y"},
			}
			g.BasicPropertyInfo = BasicPropertyInfo{HotelCode: " "}
		}},
		{"newsletter with RoomStays", newsletter, func(r *res) { r.RoomStays = &[]guestrequests202210.RoomStay{} }},
		{"Guarantee and CancelPenalties in a quote request", quote, func(r *res) {
			stay(r).Guarantee = &Guarantee{PaymentCard: PaymentCard{CardCode: "VI", ExpireDate: "1216", CardHolderName: "Otto", CardNumber: CardNumber{PlainText: "1111"}}}
			r.ResGlobalInfo.CancelPenalty = &CancelPenalty{PenaltyDescription: PenaltyDescription{Text: "free"}}
		}},
		{"cancellation with Type 14", cancellation, func(r *res) { r.UniqueID.Type = UniqueIDTypeReservation }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := load[guestrequests202210.ResRetrieveRS](t, tt.file).HotelReservations[0]
			tt.change(&old)
			var want []string
			for _, e := range guestrequests202210.ValidateHotelReservation(old) {
				want = append(want, e.Message)
			}
			if len(want) == 0 {
				t.Fatal("2022-10 reports no errors")
			}
			rulestest.Expect(t, ValidateHotelReservation(from202210(old)), want)
		})
	}
}

func TestValidate2024Changes(t *testing.T) {
	stay := func(res *HotelReservation) *RoomStay { return &(*res.RoomStays)[0] }
	deposit := func(res *HotelReservation) *GuaranteePayment { return &(*res.ResGlobalInfo.DepositPayments)[0] }
	tests := []struct {
		name   string
		change func(res *HotelReservation)
		want   []string
	}{
		{"RoomRate", func(res *HotelReservation) {
			r := stay(res).RoomRate
			r.RoomTypeCode = " "
			r.Rates[0].Bases[0].AmountAfterTax = new(alpinebits.Decimal(100.555))
			r.Rates[1].ExpireDateExclusiveInd = new(false)
		}, []string{
			"HotelReservation[0].RoomStay[0].RoomRate: RatePlanCode and RoomTypeCode required",
			"HotelReservation[0].RoomStay[0].RoomRate.Rate[0].Base[0]: AmountAfterTax has more than 2 decimal digits for EUR",
			"HotelReservation[0].RoomStay[0].RoomRate.Rate[1]: ExpireDateExclusiveInd requires ExpireDate",
		}},
		{"Services", func(res *HotelReservation) {
			s := &(*res.Services)[0]
			s.ID, s.ServiceInventoryCode, s.ServiceDetails.Total.CurrencyCode = " ", " ", "EUX"
			*stay(res).ServiceRPHs = append(*stay(res).ServiceRPHs, ServiceRPH{RPH: "9"})
		}, []string{
			"HotelReservation[0].Service[0]: ID required",
			"HotelReservation[0].Service[0]: ServiceInventoryCode required",
			"HotelReservation[0].Service[0].Total: CurrencyCode EUX is not an ISO 4217 code",
			"HotelReservation[0].RoomStay[0].ServiceRPH[2]: ServiceRPH 9 refers to no Service",
		}},
		{"additional info as ListItem", func(res *HotelReservation) {
			*res.ResGlobalInfo.Comments = []Comment{{Name: CommentNameAdditionalInfo, ListItems: []ListItem{{ListItem: 1, Language: "en", Value: "x"}}}}
		}, []string{
			"HotelReservation[0].ResGlobalInfo.Comment[0]: additional info requires Text",
		}},
		{"DepositPayments", func(res *HotelReservation) {
			d := deposit(res)
			d.GuaranteeCode, d.Start, d.AmountPercent = "", "", nil
			a := &d.AcceptedPayment
			a.GuaranteeTypeCode, a.GuaranteeID, a.PaymentTransactionTypeCode = "99", " ", "refund"
			a.BankAcct = &BankAcct{BankAcctNumber: BankAcctNumber{PlainText: "123"}}
			(*res.ResGlobalInfo.DepositPayments)[1].AcceptedPayment.GuaranteeTypeCode = ""
		}, []string{
			"HotelReservation[0].ResGlobalInfo.GuaranteePayment[0]: GuaranteeCode required",
			"HotelReservation[0].ResGlobalInfo.GuaranteePayment[0]: Start required",
			"HotelReservation[0].ResGlobalInfo.GuaranteePayment[0]: AmountPercent with Amount and CurrencyCode required",
			"HotelReservation[0].ResGlobalInfo.GuaranteePayment[0].AcceptedPayment: GuaranteeTypeCode must be 3, 28 or 46",
			"HotelReservation[0].ResGlobalInfo.GuaranteePayment[0].AcceptedPayment: GuaranteeID required",
			"HotelReservation[0].ResGlobalInfo.GuaranteePayment[0].AcceptedPayment: PaymentTransactionTypeCode must be charge or reserve",
			"HotelReservation[0].ResGlobalInfo.GuaranteePayment[0].AcceptedPayment: exactly one of PaymentCard, BankAcct and Voucher required",
			"HotelReservation[0].ResGlobalInfo.GuaranteePayment[1].AcceptedPayment: GuaranteeTypeCode required",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := load[ResRetrieveRS](t, reservationSample)
			tt.change(&rs.HotelReservations[0])
			rulestest.Expect(t, ValidateResRetrieveRS(rs), tt.want)
		})
	}
}
