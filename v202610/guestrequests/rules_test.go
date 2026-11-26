package guestrequests

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
	guestrequests202410 "github.com/HGV/alpinebits/v202410/guestrequests"
	"github.com/HGV/x/timex"
)

const (
	samples = "testdata/GuestRequests-"
	// The empty and error responses are those of 2018-10.
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
	for _, name := range []string{"Acknowledgments-OTA_NotifReportRQ", "Acknowledgments-OTA_NotifReportRQ-warning"} {
		rulestest.Expect(t, ValidateNotifReportRQ(load[NotifReportRQ](t, samples+name+".xml"), nil), nil)
	}
	tests := []struct {
		path string
		want []string
	}{
		{samples + "OTA_ResRetrieveRS-reservation.xml", nil},
		{samples + "OTA_ResRetrieveRS-reservation-alternatives.xml", nil},
		{samples + "OTA_ResRetrieveRS-allguests-and-company.xml", nil},
		{samples + "OTA_ResRetrieveRS-cancellation.xml", nil},
		{samples + "OTA_ResRetrieveRS-newsletter.xml", nil},
		// The official samples omit the Duration §4.2.2 requires with a
		// StartDateWindow, and the MealsIncluded it requires in reservations.
		{samples + "OTA_ResRetrieveRS-quote-request.xml", []string{
			"HotelReservation[0].RoomStay[0].TimeSpan: Duration required with StartDateWindow",
		}},
		{samples + "OTA_ResRetrieveRS-reservation-services-matrix.xml", []string{
			"HotelReservation[0].RoomStay[0].RatePlan: MealsIncluded required",
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

func from202410(res guestrequests202410.HotelReservation) HotelReservation {
	out := HotelReservation{
		UniqueID: res.UniqueID, CreateDateTime: res.CreateDateTime, LastModifyDateTime: res.LastModifyDateTime,
		ResStatus: res.ResStatus, RoomStayReservation: res.RoomStayReservation,
	}
	if res.RoomStays != nil {
		stays := []RoomStay{}
		for _, s := range *res.RoomStays {
			stay := RoomStay{
				RoomStayGroupID: s.RoomStayGroupID, RoomType: s.RoomType, RatePlan: s.RatePlan, GuestCounts: s.GuestCounts,
				TimeSpan: s.TimeSpan, Total: s.Total, ServiceRPHs: s.ServiceRPHs,
			}
			if r := s.RoomRate; r != nil {
				rate := RoomRate{RatePlanCode: r.RatePlanCode, RoomTypeCode: r.RoomTypeCode}
				for _, x := range r.Rates {
					rate.Rates = append(rate.Rates, Rate{
						EffectiveDate: x.EffectiveDate, ExpireDate: x.ExpireDate, ExpireDateExclusiveInd: x.ExpireDateExclusiveInd,
						RateTimeUnit: x.RateTimeUnit, UnitMultiplier: x.UnitMultiplier, Base: x.Bases[0],
					})
				}
				stay.RoomRate = &rate
			}
			if g := s.Guarantee; g != nil {
				p := g.PaymentCard
				stay.Guarantee = &Guarantee{PaymentCard: PaymentCard{
					CardCode: p.CardCode, ExpireDate: p.ExpireDate, CardHolderName: p.CardHolderName,
					CardNumber: CardNumber{PlainText: p.CardNumber.PlainText, EncryptedValue: p.CardNumber.EncryptedValue, EncryptionMethod: p.CardNumber.EncryptionMethod},
				}}
			}
			stays = append(stays, stay)
		}
		out.RoomStays = &stays
	}
	if res.Services != nil {
		services := []Service{}
		for _, s := range *res.Services {
			d := s.ServiceDetails
			details := ServiceDetails{GuestCounts: d.GuestCounts, Comment: d.Comment, Total: d.Total, Description: d.Description}
			if ts := d.TimeSpan; ts != nil {
				start, _ := ts.Start.MarshalText()
				details.TimeSpan = &ServiceTimeSpan{Start: string(start)}
				if ts.End != nil {
					end, _ := ts.End.MarshalText()
					details.TimeSpan.End = string(end)
				}
			}
			services = append(services, Service{
				ID: s.ID, Type: s.Type, ServiceCategoryCode: s.ServiceCategoryCode, ServiceInventoryCode: s.ServiceInventoryCode,
				ServiceRPH: s.ServiceRPH, ServicePricingType: s.ServicePricingType, Inclusive: s.Inclusive, Quantity: s.Quantity,
				ServiceDetails: details,
			})
		}
		out.Services = &services
	}
	if res.ResGuest != nil {
		out.ResGuests = &[]ResGuest{{Customer: res.ResGuest.Customer}}
	}
	if g := res.ResGlobalInfo; g != nil {
		info := ResGlobalInfo{
			Comments: g.Comments, SpecialRequests: g.SpecialRequests, CancelPenalty: g.CancelPenalty,
			HotelReservationIDs: g.HotelReservationIDs, BasicPropertyInfo: g.BasicPropertyInfo,
		}
		if g.DepositPayments != nil {
			deposits := []GuaranteePayment{}
			for _, d := range *g.DepositPayments {
				a := d.AcceptedPayment
				pay := GuaranteePayment{Start: d.Start, GuaranteeCode: d.GuaranteeCode, AmountPercent: d.AmountPercent, AcceptedPayment: AcceptedPayment{
					GuaranteeTypeCode: a.GuaranteeTypeCode, GuaranteeID: a.GuaranteeID, PaymentTransactionTypeCode: a.PaymentTransactionTypeCode,
					BankAcct: a.BankAcct, Voucher: a.Voucher,
				}}
				if p := a.PaymentCard; p != nil {
					card := PaymentCard{CardCode: p.CardCode, ExpireDate: p.ExpireDate, CardHolderName: p.CardHolderName}
					if n := p.CardNumber; n != nil {
						card.CardNumber = CardNumber{Mask: n.Mask, Token: n.Token, TokenProviderID: n.TokenProviderID}
					}
					pay.AcceptedPayment.PaymentCard = &card
				}
				deposits = append(deposits, pay)
			}
			info.DepositPayments = &deposits
		}
		out.ResGlobalInfo = &info
	}
	return out
}

// The checks 2026-10 copied from 2024-10 must report what 2024-10 reports.
func TestRulesMatch202410(t *testing.T) {
	type res = guestrequests202410.HotelReservation
	stay := func(r *res) *guestrequests202410.RoomStay { return &(*r.RoomStays)[0] }
	const (
		reservation  = "../../v202410/guestrequests/testdata/GuestRequests-OTA_ResRetrieveRS-reservation.xml"
		quote        = "../../v202210/guestrequests/testdata/GuestRequests-OTA_ResRetrieveRS-quote-request.xml"
		cancellation = "../../v202210/guestrequests/testdata/GuestRequests-OTA_ResRetrieveRS-cancellation.xml"
		newsletter   = "../../v202210/guestrequests/testdata/GuestRequests-OTA_ResRetrieveRS-newsletter.xml"
	)
	tests := []struct {
		name, file string
		change     func(r *res)
	}{
		{"HotelReservation", reservation, func(r *res) {
			r.UniqueID = UniqueID{Type: UniqueIDTypeCancellation, ID: " "}
			r.CreateDateTime = DateTime{}
			r.RoomStays, r.ResGlobalInfo = nil, nil
		}},
		{"RoomStayReservation false in a reservation", reservation, func(r *res) { r.RoomStayReservation = new(false) }},
		{"RoomType and RatePlan", reservation, func(r *res) {
			s := stay(r)
			s.RoomType = &RoomType{RoomTypeCode: " ", RoomType: 1}
			s.RatePlan.RatePlanCode, s.RatePlan.MealsIncluded = " ", nil
			s.RatePlan.Commission = &Commission{CommissionPayableAmount: &CommissionPayableAmount{Amount: new(alpinebits.Decimal(10.5)), CurrencyCode: "JPY"}}
		}},
		{"Commissions", reservation, func(r *res) {
			var stays []guestrequests202410.RoomStay
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
		// 2026-10 words the missing card number differently.
		{"PaymentCards", reservation, func(r *res) {
			var stays []guestrequests202410.RoomStay
			for _, n := range []guestrequests202410.CardNumber{{PlainText: "1111", EncryptedValue: "x", EncryptionMethod: "m"}, {EncryptedValue: "x"}} {
				s := *stay(r)
				s.Guarantee = &guestrequests202410.Guarantee{PaymentCard: guestrequests202410.PaymentCard{CardCode: "V", CardNumber: n}}
				stays = append(stays, s)
			}
			*r.RoomStays = stays
		}},
		{"Total", reservation, func(r *res) { stay(r).Total = &Total{AmountAfterTax: 299.999, CurrencyCode: "EUX"} }},
		{"RoomStayGroupID in a reservation", reservation, func(r *res) { stay(r).RoomStayGroupID = "1" }},
		{"RoomRate", reservation, func(r *res) {
			rate := stay(r).RoomRate
			rate.RoomTypeCode = " "
			rate.Rates[1].ExpireDateExclusiveInd = new(false)
		}},
		{"Services", reservation, func(r *res) {
			s := &(*r.Services)[0]
			s.ID, s.ServiceInventoryCode, s.ServiceDetails.Total.CurrencyCode = " ", " ", "EUX"
			*stay(r).ServiceRPHs = append(*stay(r).ServiceRPHs, ServiceRPH{RPH: "9"})
		}},
		{"ResGlobalInfo", reservation, func(r *res) {
			g := r.ResGlobalInfo
			comments := *g.Comments
			comments[0].ListItems = append(comments[0].ListItems, comments[0].ListItems[0])
			*g.Comments = []Comment{comments[0], {Name: CommentNameIncludedServices}, {Name: CommentNameCustomerComment}, {Name: CommentNameAdditionalInfo}}
			g.SpecialRequests = &[]SpecialRequest{
				{CodeContext: CodeContextAlpineBits, RequestCode: "X"},
				{CodeContext: CodeContextHotel, RequestCode: " "},
				{CodeContext: CodeContextHotel, RequestCode: "Y"},
				{CodeContext: CodeContextHotel, RequestCode: "Y"},
			}
			g.BasicPropertyInfo = BasicPropertyInfo{HotelCode: " "}
			d := &(*g.DepositPayments)[0]
			d.GuaranteeCode, d.Start, d.AmountPercent = "", "", nil
			d.AcceptedPayment.GuaranteeTypeCode, d.AcceptedPayment.GuaranteeID = "", " "
			d.AcceptedPayment.Voucher = &Voucher{ElectronicIndicator: true, Identifier: "V"}
		}},
		{"newsletter with RoomStays", newsletter, func(r *res) { r.RoomStays = &[]guestrequests202410.RoomStay{} }},
		{"Guarantee and CancelPenalties in a quote request", quote, func(r *res) {
			stay(r).Guarantee = &guestrequests202410.Guarantee{PaymentCard: guestrequests202410.PaymentCard{
				CardCode: "VI", ExpireDate: "1216", CardHolderName: "Otto", CardNumber: guestrequests202410.CardNumber{PlainText: "1111"},
			}}
			r.ResGlobalInfo.CancelPenalty = &CancelPenalty{PenaltyDescription: PenaltyDescription{Text: "free"}}
		}},
		{"cancellation with Type 14", cancellation, func(r *res) { r.UniqueID.Type = UniqueIDTypeReservation }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := load[guestrequests202410.ResRetrieveRS](t, tt.file).HotelReservations[0]
			tt.change(&old)
			var want []string
			for _, e := range guestrequests202410.ValidateHotelReservation(old) {
				want = append(want, e.Message)
			}
			if len(want) == 0 {
				t.Fatal("2024-10 reports no errors")
			}
			rulestest.Expect(t, ValidateHotelReservation(from202410(old)), want)
		})
	}
}

func TestValidate2026Changes(t *testing.T) {
	stay := func(res *HotelReservation, i int) *RoomStay { return &(*res.RoomStays)[i] }
	service := func(res *HotelReservation, i int) *Service { return &(*res.Services)[i] }
	tests := []struct {
		name, file string
		change     func(res *HotelReservation)
		want       []string
	}{
		{"without ResGuests", "OTA_ResRetrieveRS-reservation.xml", func(res *HotelReservation) { res.ResGuests = nil }, nil},
		{"ResGuestRPHs", "OTA_ResRetrieveRS-allguests-and-company.xml", func(res *HotelReservation) {
			(*res.ResGuests)[1].ResGuestRPH = "1"
			stay(res, 1).ResGuestRPHs = RPHs{"4"}
		}, []string{
			"HotelReservation[0].ResGuest[1]: duplicate ResGuestRPH 1",
			"HotelReservation[0].RoomStay[0]: ResGuestRPH 2 refers to no ResGuest",
			"HotelReservation[0].RoomStay[1]: ResGuestRPH 4 refers to no ResGuest",
		}},
		{"RatePlanCode of RoomRate", "OTA_ResRetrieveRS-reservation.xml", func(res *HotelReservation) {
			stay(res, 0).RoomRate.RatePlanCode = "other"
		}, []string{
			"HotelReservation[0].RoomStay[0].RoomRate: RatePlanCode other differs from RatePlan RatePlanCode 123456-xyz",
		}},
		{"Services", "OTA_ResRetrieveRS-reservation.xml", func(res *HotelReservation) {
			service(res, 4).ServiceDetails.TimeSpan.End = "2027-01-03"
			service(res, 3).Quantity = 2
		}, []string{
			"HotelReservation[0].Service[3]: Quantity must be 1 for ServicePricingType Per stay",
			"HotelReservation[0].Service[4].TimeSpan: Start and End must both be dates, times or date-times",
		}},
		{"original booking ID", "OTA_ResRetrieveRS-reservation.xml", func(res *HotelReservation) {
			res.ResGlobalInfo.HotelReservationIDs = &[]HotelReservationID{{ResIDType: 14}}
		}, []string{
			"HotelReservation[0].ResGlobalInfo.HotelReservationID[0]: ResID_Value required with ResID_Type 14",
		}},
		{"PaymentCard by Token", "OTA_ResRetrieveRS-reservation.xml", func(res *HotelReservation) {
			stay(res, 0).Guarantee.PaymentCard = PaymentCard{CardNumber: CardNumber{Mask: "XXXX1234", Token: "abc123"}}
		}, nil},
		{"PaymentCard by Mask only", "OTA_ResRetrieveRS-reservation.xml", func(res *HotelReservation) {
			stay(res, 0).Guarantee.PaymentCard = PaymentCard{CardNumber: CardNumber{Mask: "XXXX1234"}}
		}, []string{
			"HotelReservation[0].RoomStay[0].PaymentCard.CardNumber: PlainText, EncryptedValue and EncryptionMethod, or Token required",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := load[ResRetrieveRS](t, samples+tt.file)
			tt.change(&rs.HotelReservations[0])
			rulestest.Expect(t, ValidateResRetrieveRS(rs), tt.want)
		})
	}
}

func TestValidateNotifReportRQ(t *testing.T) {
	rq := NotifReportRQ{Warnings: &[]Warning{{RecordID: " "}}, HotelReservations: &[]NotifReservation{
		{UniqueID: UniqueID{Type: UniqueIDTypeReservation, ID: "1"}, BasicPropertyInfo: &BasicPropertyInfo{HotelCode: "123"}},
		{UniqueID: UniqueID{Type: UniqueIDTypeReservation, ID: " "}, BasicPropertyInfo: &BasicPropertyInfo{HotelCode: " "}},
	}}
	rulestest.Expect(t, ValidateNotifReportRQ(rq, nil), []string{
		"Warning[0]: RecordID required",
		"HotelReservation[1].UniqueID: ID required",
		"HotelReservation[1].BasicPropertyInfo: HotelCode required",
	})
}
