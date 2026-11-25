package freerooms

import (
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
	"github.com/HGV/x/timex"
)

func TestValidate(t *testing.T) {
	day, err := timex.ParseDate("2022-08-01")
	if err != nil {
		t.Fatal(err)
	}
	completeSet := &UniqueID{Type: UniqueIDTypeReference, ID: "1", Instance: InstanceCompleteSet}
	category := Inventory{StatusApplicationControl: &StatusApplicationControl{Start: day, End: day, InvTypeCode: "DOUBLE"}}
	tests := []struct {
		name     string
		uniqueID *UniqueID
		invs     []Inventory
		caps     []alpinebits.Capability
		want     []string
	}{
		{"CompleteSet", completeSet, []Inventory{category}, []alpinebits.Capability{CapAcceptCategories, CapAcceptCompleteSet}, nil},
		{"CompleteSet not negotiated", completeSet, []Inventory{category}, []alpinebits.Capability{CapAcceptCategories}, []string{
			"complete sets not supported",
		}},
		{"reset not negotiated", completeSet, []Inventory{{}}, nil, []string{"complete sets not supported"}},
		{"2020-10 rules", nil, []Inventory{category}, []alpinebits.Capability{CapAcceptCategories}, []string{
			"deltas not supported",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := HotelInvCountNotifRQ{UniqueID: tt.uniqueID, Inventories: Inventories{HotelCode: "123", Inventories: tt.invs}}
			rulestest.Expect(t, Validate(rq, tt.caps), tt.want)
		})
	}
}
