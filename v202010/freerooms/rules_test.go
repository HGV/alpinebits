package freerooms

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
	"github.com/HGV/x/timex"
)

var allCaps = []alpinebits.Capability{
	CapAcceptRooms, CapAcceptCategories, CapAcceptDeltas,
	CapAcceptOutOfOrder, CapAcceptOutOfMarket, CapAcceptClosingSeasons,
}

// The official AlpineBits 2020-10 samples. completeSet holds three periods of
// category DOUBLE, the second without InvCounts; closingSeasons is completeSet
// preceded by a closing season; delta updates two periods of DOUBLE; reset is
// a CompleteSet with one empty Inventory.
const (
	completeSet    = "FreeRooms-OTA_HotelInvCountNotifRQ.xml"
	closingSeasons = "FreeRooms-OTA_HotelInvCountNotifRQ-closing_seasons.xml"
	delta          = "FreeRooms-OTA_HotelInvCountNotifRQ-delta.xml"
	reset          = "FreeRooms-OTA_HotelInvCountNotifRQ-empty.xml"
)

func sample(t *testing.T, name string) HotelInvCountNotifRQ {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	var rq HotelInvCountNotifRQ
	if err := xml.Unmarshal(b, &rq); err != nil {
		t.Fatal(err)
	}
	return rq
}

func date(s string) timex.Date {
	d, err := timex.ParseDate(s)
	if err != nil {
		panic(err)
	}
	return d
}

func counts(c ...InvCount) *[]InvCount {
	return &c
}

func TestSamples(t *testing.T) {
	for _, name := range []string{completeSet, closingSeasons, delta, reset} {
		t.Run(name, func(t *testing.T) {
			rq := sample(t, name)
			if rq.HotelCode() != "123" || rq.Inventories.HotelName != "Frangart Inn" {
				t.Errorf("hotel %q %q", rq.HotelCode(), rq.Inventories.HotelName)
			}
			rulestest.Expect(t, Validate(rq, allCaps), nil)
		})
	}

	rq := sample(t, closingSeasons)
	if !rq.IsCompleteSet() || rq.IsReset() || rq.UniqueID.Type != UniqueIDTypeReference || rq.UniqueID.Instance != InstanceCompleteSet {
		t.Errorf("UniqueID %+v", rq.UniqueID)
	}
	invs := rq.Inventories.Inventories
	if len(invs) != 4 || !invs[0].IsClosingSeason() || invs[1].IsClosingSeason() || invs[1].IsRoom() {
		t.Fatalf("Inventories %+v", invs)
	}
	want := timex.DateRange{Start: date("2020-08-01"), End: date("2020-08-10")}
	if got := invs[1].DateRange(); got != want {
		t.Errorf("DateRange %v, want %v", got, want)
	}
	if n := invs[1].Count(CountTypeBookable); n != 3 {
		t.Errorf("bookable %d, want 3", n)
	}
	// A missing InvCount or InvCounts counts 0 (§4.1.1).
	if n := invs[1].Count(CountTypeOutOfOrder) + invs[2].Count(CountTypeBookable); n != 0 {
		t.Errorf("missing counts %d, want 0", n)
	}
	if !sample(t, reset).IsReset() || sample(t, delta).IsCompleteSet() {
		t.Error("reset or delta not recognized")
	}
}

func TestValidateCapabilities(t *testing.T) {
	tests := []struct {
		name   string
		sample string
		change func(rq *HotelInvCountNotifRQ)
		want   []string
	}{
		{"closing seasons, categories, out of order, out of market", closingSeasons, func(rq *HotelInvCountNotifRQ) {
			*rq.Inventories.Inventories[1].InvCounts = append(*rq.Inventories.Inventories[1].InvCounts,
				InvCount{CountType: CountTypeOutOfOrder, Count: 1}, InvCount{CountType: CountTypeOutOfMarket, Count: 1})
			rq.Inventories.Inventories[3].InvCounts = counts(InvCount{CountType: CountTypeOutOfOrder})
		}, []string{
			"Inventory[0]: closing seasons not supported",
			"Inventory[1]: categories not supported",
			"Inventory[1].InvCount[1]: CountType 6 (out of order) not supported",
			"Inventory[1].InvCount[2]: CountType 9 (out of market) not supported",
		}},
		{"deltas, rooms", delta, func(rq *HotelInvCountNotifRQ) {
			for _, inv := range rq.Inventories.Inventories {
				inv.StatusApplicationControl.InvCode = "101"
			}
		}, []string{
			"deltas not supported",
			"Inventory[0]: rooms not supported",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := sample(t, tt.sample)
			tt.change(&rq)
			rulestest.Expect(t, Validate(rq, allCaps), nil)
			rulestest.Expect(t, Validate(rq, nil), tt.want)
		})
	}
}

func TestValidateMessage(t *testing.T) {
	inventory := func(start, end, invTypeCode, invCode string, c *[]InvCount) Inventory {
		return Inventory{
			StatusApplicationControl: &StatusApplicationControl{Start: date(start), End: date(end), InvTypeCode: invTypeCode, InvCode: invCode},
			InvCounts:                c,
		}
	}
	bookable := counts(InvCount{CountType: CountTypeBookable, Count: 1})
	purge := &UniqueID{Type: UniqueIDTypePurgedMasterReference, ID: "1", Instance: InstanceCompleteSet}
	tests := []struct {
		name      string
		uniqueID  *UniqueID
		hotelCode string
		hotelName string
		invs      []Inventory
		caps      []alpinebits.Capability
		want      []string
	}{
		{"reset by Type 35", purge, "123", "", []Inventory{{}}, nil, nil},
		{"HotelName only", purge, " ", "Frangart Inn", []Inventory{{}}, nil, []string{
			"HotelCode required",
		}},
		{"empty Inventory in a delta", nil, "123", "", []Inventory{{}}, nil, []string{
			"deltas not supported",
			"Inventory[0]: StatusApplicationControl required",
		}},
		{"empty Inventory among others", purge, "123", "", []Inventory{
			inventory("2020-08-01", "2020-08-10", "DOUBLE", "", bookable),
			{InvCounts: bookable},
		}, allCaps, []string{
			"Inventory[1]: StatusApplicationControl required",
		}},
		{"rooms and categories mixed", purge, "123", "", []Inventory{
			inventory("2020-08-01", "2020-08-10", "DOUBLE", "", bookable),
			inventory("2020-08-01", "2020-08-10", "DOUBLE", "101", bookable),
		}, allCaps, []string{
			"rooms and categories must not be mixed",
		}},
		{"rooms", purge, "123", "", []Inventory{
			inventory("2020-08-01", "2020-08-10", "DOUBLE", "101", bookable),
			inventory("2020-08-05", "2020-08-15", "DOUBLE", "102", nil),
			inventory("2020-08-01", "2020-08-10", "SINGLE", "101", bookable),
			inventory("2020-08-11", "2020-08-20", "DOUBLE", "101", counts(
				InvCount{CountType: CountTypeBookable, Count: 1},
				InvCount{CountType: CountTypeOutOfMarket, Count: 1},
			)),
			inventory("2020-08-10", "2020-08-10", "DOUBLE", "101", counts(InvCount{CountType: CountTypeOutOfOrder, Count: 2})),
		}, allCaps, []string{
			"Inventory[3]: Counts of a room must sum to at most 1",
			"Inventory[4]: Counts of a room must sum to at most 1",
			"Inventory[4]: overlaps Inventory[0]",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := HotelInvCountNotifRQ{UniqueID: tt.uniqueID, Inventories: Inventories{HotelCode: tt.hotelCode, HotelName: tt.hotelName, Inventories: tt.invs}}
			rulestest.Expect(t, Validate(rq, tt.caps), tt.want)
		})
	}
}

func TestValidateInventory(t *testing.T) {
	tests := []struct {
		name   string
		change func(invs []Inventory) []Inventory
		want   []string
	}{
		{"End before Start", func(invs []Inventory) []Inventory {
			invs[1].StatusApplicationControl.End = date("2020-08-10")
			return invs
		}, []string{
			"Inventory[1]: End before Start",
		}},
		{"no InvTypeCode", func(invs []Inventory) []Inventory {
			invs[0].StatusApplicationControl.InvTypeCode = ""
			invs[2].StatusApplicationControl.InvTypeCode = " "
			return invs
		}, []string{
			"Inventory[0]: InvTypeCode required",
			"Inventory[2]: InvTypeCode required",
		}},
		{"InvCode without InvTypeCode", func(invs []Inventory) []Inventory {
			for i := range invs {
				invs[i].StatusApplicationControl.InvTypeCode, invs[i].StatusApplicationControl.InvCode = "", "101"
				invs[i].InvCounts = nil
			}
			return invs
		}, []string{
			"Inventory[0]: InvTypeCode required",
			"Inventory[1]: InvTypeCode required",
			"Inventory[2]: InvTypeCode required",
		}},
		{"duplicate CountType", func(invs []Inventory) []Inventory {
			invs[0].InvCounts = counts(
				InvCount{CountType: CountTypeBookable, Count: 1},
				InvCount{CountType: CountTypeOutOfOrder, Count: 1},
				InvCount{CountType: CountTypeBookable, Count: 2},
			)
			return invs
		}, []string{
			"Inventory[0].InvCount[2]: duplicate CountType 2",
		}},
		{"overlaps", func(invs []Inventory) []Inventory {
			other := invs[0]
			other.StatusApplicationControl = &StatusApplicationControl{Start: date("2020-08-05"), End: date("2020-08-25"), InvTypeCode: "SINGLE"}
			invs[1].StatusApplicationControl.Start = date("2020-08-10")
			invs[2].StatusApplicationControl.Start = date("2020-08-20")
			return append(invs, other)
		}, []string{
			"Inventory[1]: overlaps Inventory[0]",
			"Inventory[2]: overlaps Inventory[1]",
		}},
		{"order does not matter", func(invs []Inventory) []Inventory {
			return []Inventory{invs[2], invs[0], invs[1]}
		}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := sample(t, completeSet)
			rq.Inventories.Inventories = tt.change(rq.Inventories.Inventories)
			rulestest.Expect(t, Validate(rq, allCaps), tt.want)
		})
	}
}

func TestValidateClosingSeasons(t *testing.T) {
	closingSeason := func(start, end string) Inventory {
		return Inventory{StatusApplicationControl: &StatusApplicationControl{Start: date(start), End: date(end), AllInvCode: true}}
	}
	tests := []struct {
		name   string
		change func(rq *HotelInvCountNotifRQ)
		want   []string
	}{
		{"in a delta", func(rq *HotelInvCountNotifRQ) { rq.UniqueID = nil }, []string{
			"Inventory[0]: closing seasons only allowed in a CompleteSet",
		}},
		{"after other Inventory elements", func(rq *HotelInvCountNotifRQ) {
			invs := rq.Inventories.Inventories
			rq.Inventories.Inventories = append(invs[1:], invs[0], closingSeason("2020-10-01", "2020-10-31"))
		}, []string{
			"Inventory[3]: closing seasons must precede the other Inventory elements",
			"Inventory[4]: closing seasons must precede the other Inventory elements",
		}},
		{"content", func(rq *HotelInvCountNotifRQ) {
			s := rq.Inventories.Inventories[0].StatusApplicationControl
			s.InvTypeCode, s.InvCode = "DOUBLE", "101"
			rq.Inventories.Inventories[0].InvCounts = counts(InvCount{CountType: CountTypeOutOfOrder, Count: 1})
		}, []string{
			"Inventory[0]: closing season must have no InvTypeCode and no InvCode",
			"Inventory[0]: closing season must have no InvCounts",
		}},
		{"End before Start", func(rq *HotelInvCountNotifRQ) {
			rq.Inventories.Inventories[0].StatusApplicationControl.End = date("2020-08-30")
		}, []string{
			"Inventory[0]: End before Start",
		}},
		{"overlapping closing seasons", func(rq *HotelInvCountNotifRQ) {
			invs := rq.Inventories.Inventories
			rq.Inventories.Inventories = append([]Inventory{invs[0], closingSeason("2020-09-30", "2020-10-15")}, invs[1:]...)
		}, []string{
			"Inventory[1]: overlaps Inventory[0]",
		}},
		{"overlapping availability", func(rq *HotelInvCountNotifRQ) {
			invs := rq.Inventories.Inventories
			invs[0].StatusApplicationControl.Start = date("2020-08-05")
			invs[2].InvCounts = counts(InvCount{CountType: CountTypeOutOfMarket, Count: 1})
		}, []string{
			"Inventory[1]: available rooms overlap closing season Inventory[0]",
			"Inventory[2]: available rooms overlap closing season Inventory[0]",
			"Inventory[3]: available rooms overlap closing season Inventory[0]",
		}},
		{"overlapping no availability", func(rq *HotelInvCountNotifRQ) {
			invs := rq.Inventories.Inventories
			invs[0].StatusApplicationControl.Start = date("2020-08-05")
			invs[1].InvCounts = counts(InvCount{CountType: CountTypeBookable}, InvCount{CountType: CountTypeOutOfOrder, Count: 3})
			invs[3].InvCounts = nil
		}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := sample(t, closingSeasons)
			tt.change(&rq)
			rulestest.Expect(t, Validate(rq, allCaps), tt.want)
		})
	}
}

func TestCheckInventory(t *testing.T) {
	rq := sample(t, closingSeasons)
	rq.Inventories.Inventories[3].StatusApplicationControl.InvTypeCode = "SUITE"
	hotel := HotelInventory{Categories: map[string]int{"DOUBLE": 2}}
	rulestest.Expect(t, CheckInventory(rq, hotel), []string{
		"Inventory[1]: Counts exceed the 2 rooms of category DOUBLE",
		"Inventory[3]: unknown category SUITE",
	})
	rulestest.Expect(t, CheckInventory(rq, HotelInventory{}), nil)
	rulestest.Expect(t, CheckInventory(sample(t, reset), hotel), nil)

	rq = sample(t, delta)
	for _, inv := range rq.Inventories.Inventories {
		inv.StatusApplicationControl.InvCode = "101"
	}
	rq.Inventories.Inventories[1].StatusApplicationControl.InvTypeCode = "SINGLE"
	hotel = HotelInventory{Rooms: map[string]map[string]struct{}{"DOUBLE": {"101": {}}, "SINGLE": {"201": {}}}}
	rulestest.Expect(t, CheckInventory(rq, hotel), []string{
		"Inventory[1]: unknown room 101 in category SINGLE",
	})
}
