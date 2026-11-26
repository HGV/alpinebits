package rateplans

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
)

var allCaps = []alpinebits.Capability{
	CapAcceptArrivalDOW, CapAcceptDepartureDOW,
	CapAcceptRatePlanBookingRule, CapAcceptRatePlanRoomTypeBookingRule, CapAcceptRatePlanMixedBookingRule,
	CapAcceptSupplements, CapAcceptFreeNightsOffers, CapAcceptFamilyOffers, CapAcceptOverlay,
	CapAcceptRatePlanJoin, CapAcceptOfferRuleBookingOffset, CapAcceptOfferRuleDOWLOS,
}

// sample is the official RatePlans sample, which 2020-10 kept from 2018-10.
func sample(t *testing.T) HotelRatePlanNotifRQ {
	t.Helper()
	b, err := os.ReadFile("../../v201810/rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRQ.xml")
	if err != nil {
		t.Fatal(err)
	}
	var rq HotelRatePlanNotifRQ
	if err := xml.Unmarshal(b, &rq); err != nil {
		t.Fatal(err)
	}
	return rq
}

func TestValidateSample(t *testing.T) {
	rulestest.Expect(t, Validate(sample(t), allCaps), nil)
}

func TestValidateKeeps201810Rules(t *testing.T) {
	rq := sample(t)
	rq.RatePlans.RatePlans[0].RatePlanCode = " "
	rulestest.Expect(t, Validate(rq, allCaps), []string{"RatePlan[0]: RatePlanCode required"})
}

func TestValidateCurrency(t *testing.T) {
	tests := []struct {
		name   string
		change func(p *RatePlan)
		want   []string
	}{
		{"other ISO 4217 currency", func(p *RatePlan) { p.CurrencyCode = "CHF" }, nil},
		{"not ISO 4217", func(p *RatePlan) { p.CurrencyCode = "EUX" }, []string{
			"RatePlan[0]: CurrencyCode EUX is not an ISO 4217 code",
		}},
		{"amounts with too many decimal digits", func(p *RatePlan) { p.CurrencyCode = "JPY" }, []string{
			"RatePlan[0].Rate[1].AdditionalGuestAmount[0]: Amount has more than 0 decimal digits for JPY",
			"RatePlan[0].Rate[1].AdditionalGuestAmount[2]: Amount has more than 0 decimal digits for JPY",
			"RatePlan[0].Rate[1].AdditionalGuestAmount[4]: Amount has more than 0 decimal digits for JPY",
		}},
		{"BaseByGuestAmt", func(p *RatePlan) {
			b := &(*(*p.Rates)[1].BaseByGuestAmts)[0]
			b.CurrencyCode, b.AmountAfterTax = "CHF", 106.001
		}, []string{
			"RatePlan[0].Rate[1].BaseByGuestAmt[0]: CurrencyCode CHF differs from RatePlan CurrencyCode EUR",
			"RatePlan[0].Rate[1].BaseByGuestAmt[0]: AmountAfterTax has more than 2 decimal digits for EUR",
		}},
		{"Supplement", func(p *RatePlan) { (*p.Supplements)[1].Amount = new(alpinebits.Decimal(20.005)) }, []string{
			"RatePlan[0].Supplement[1]: Amount has more than 2 decimal digits for EUR",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := sample(t)
			tt.change(&rq.RatePlans.RatePlans[0])
			rulestest.Expect(t, Validate(rq, allCaps), tt.want)
		})
	}
}

func TestValidateOneRatePlanPerCode(t *testing.T) {
	remove := func(code string) RatePlan {
		return RatePlan{RatePlanNotifType: RatePlanNotifTypeRemove, CurrencyCode: "EUR", RatePlanCode: code}
	}
	rq := HotelRatePlanNotifRQ{RatePlans: RatePlans{HotelCode: "123", RatePlans: []RatePlan{remove("A"), remove("B"), remove("A")}}}
	rulestest.Expect(t, Validate(rq, allCaps), []string{"RatePlan[2]: duplicate RatePlanCode A"})
}
