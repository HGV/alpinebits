package freerooms

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rulestest"
	"github.com/HGV/x/timex"
)

var allCaps = []alpinebits.Capability{CapabilityAcceptRooms, CapabilityAcceptCategories, CapabilityAcceptDeltas, CapabilityAcceptBookingThreshold}

func date(s string) timex.Date {
	d, err := timex.ParseDate(s)
	if err != nil {
		panic(err)
	}
	return d
}

func load(t *testing.T, name string) HotelAvailNotifRQ {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	var rq HotelAvailNotifRQ
	if err := xml.Unmarshal(b, &rq); err != nil {
		t.Fatal(err)
	}
	return rq
}

// sample is the official AlpineBits 2018-10 FreeRooms sample: a CompleteSet
// with two periods of room 101S, each with BookingLimit 1 and BookingThreshold
// 0.
func sample(t *testing.T) HotelAvailNotifRQ {
	return load(t, "FreeRooms-OTA_HotelAvailNotifRQ.xml")
}

func TestValidateSamples(t *testing.T) {
	rulestest.Expect(t, Validate(sample(t), allCaps), nil)

	reset := load(t, "FreeRooms-OTA_HotelAvailNotifRQ-empty.xml")
	if !reset.IsReset() {
		t.Error("empty sample is no reset")
	}
	rulestest.Expect(t, Validate(reset, nil), nil)
}

func TestValidateCapabilities(t *testing.T) {
	rq := sample(t)
	rq.UniqueID = nil
	rulestest.Expect(t, Validate(rq, allCaps), nil)
	rulestest.Expect(t, Validate(rq, nil), []string{
		"deltas not supported",
		"AvailStatusMessage[0]: BookingThreshold not supported",
		"rooms not supported",
	})

	for _, m := range rq.AvailStatusMessages.AvailStatusMessages {
		m.StatusApplicationControl.InvCode = ""
	}
	rulestest.Expect(t, Validate(rq, allCaps), nil)
	rulestest.Expect(t, Validate(rq, []alpinebits.Capability{CapabilityAcceptRooms, CapabilityAcceptDeltas, CapabilityAcceptBookingThreshold}), []string{
		"categories not supported",
	})
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name   string
		change func(rq *HotelAvailNotifRQ, msgs []AvailStatusMessage)
		want   []string
	}{
		{"HotelName only", func(rq *HotelAvailNotifRQ, _ []AvailStatusMessage) {
			rq.AvailStatusMessages.HotelCode = " "
		}, []string{
			"HotelCode required",
		}},
		{"HotelCode only", func(rq *HotelAvailNotifRQ, _ []AvailStatusMessage) {
			rq.AvailStatusMessages.HotelName = ""
		}, nil},
		{"purge", func(rq *HotelAvailNotifRQ, _ []AvailStatusMessage) {
			rq.UniqueID.Type = UniqueIDTypePurgedMasterReference
		}, nil},
		{"empty AvailStatusMessage outside a reset", func(rq *HotelAvailNotifRQ, msgs []AvailStatusMessage) {
			rq.AvailStatusMessages.AvailStatusMessages = append(msgs, AvailStatusMessage{})
		}, []string{
			"AvailStatusMessage[2]: BookingLimit required",
			"AvailStatusMessage[2]: BookingThreshold required",
			"AvailStatusMessage[2]: StatusApplicationControl required",
		}},
		{"reset as delta", func(rq *HotelAvailNotifRQ, _ []AvailStatusMessage) {
			rq.UniqueID = nil
			rq.AvailStatusMessages.AvailStatusMessages = []AvailStatusMessage{{}}
		}, []string{
			"AvailStatusMessage[0]: BookingLimit required",
			"AvailStatusMessage[0]: BookingThreshold required",
			"AvailStatusMessage[0]: StatusApplicationControl required",
		}},
		{"categories", func(_ *HotelAvailNotifRQ, msgs []AvailStatusMessage) {
			for _, m := range msgs {
				m.StatusApplicationControl.InvCode = ""
				*m.BookingLimit = 5
			}
		}, nil},
		{"rooms and categories mixed", func(_ *HotelAvailNotifRQ, msgs []AvailStatusMessage) {
			msgs[1].StatusApplicationControl.InvCode = ""
		}, []string{
			"rooms and categories must not be mixed",
		}},
		{"BookingLimit", func(_ *HotelAvailNotifRQ, msgs []AvailStatusMessage) {
			msgs[0].BookingLimit = nil
			*msgs[1].BookingLimit = 2
		}, []string{
			"AvailStatusMessage[0]: BookingLimit required",
			"AvailStatusMessage[1]: BookingLimit must be 0 or 1 for a room",
		}},
		{"BookingThreshold", func(_ *HotelAvailNotifRQ, msgs []AvailStatusMessage) {
			msgs[0].BookingThreshold = nil
			*msgs[1].BookingThreshold = 2
		}, []string{
			"AvailStatusMessage[0]: BookingThreshold required",
			"AvailStatusMessage[1]: BookingThreshold exceeds BookingLimit",
		}},
		{"StatusApplicationControl", func(_ *HotelAvailNotifRQ, msgs []AvailStatusMessage) {
			msgs[0].StatusApplicationControl.InvTypeCode = " "
			msgs[1].StatusApplicationControl.End = date("2010-08-20")
		}, []string{
			"AvailStatusMessage[0].StatusApplicationControl: InvTypeCode required",
			"AvailStatusMessage[1].StatusApplicationControl: End before Start",
		}},
		{"overlap on the shared last day", func(_ *HotelAvailNotifRQ, msgs []AvailStatusMessage) {
			msgs[1].StatusApplicationControl.Start = date("2010-08-10")
		}, []string{
			"AvailStatusMessage[1]: overlaps AvailStatusMessage[0]",
		}},
		{"same period for other rooms", func(rq *HotelAvailNotifRQ, msgs []AvailStatusMessage) {
			other := *msgs[0].StatusApplicationControl
			other.InvCode = "102S"
			sameCodeOtherCategory := *msgs[0].StatusApplicationControl
			sameCodeOtherCategory.InvTypeCode = "single"
			rq.AvailStatusMessages.AvailStatusMessages = append(msgs,
				AvailStatusMessage{BookingLimit: new(1), BookingThreshold: new(0), StatusApplicationControl: &other},
				AvailStatusMessage{BookingLimit: new(1), BookingThreshold: new(0), StatusApplicationControl: &sameCodeOtherCategory},
			)
		}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := sample(t)
			tt.change(&rq, rq.AvailStatusMessages.AvailStatusMessages)
			rulestest.Expect(t, Validate(rq, allCaps), tt.want)
		})
	}
}

func TestCheckInventory(t *testing.T) {
	rq := sample(t)
	msgs := rq.AvailStatusMessages.AvailStatusMessages
	msgs[1].StatusApplicationControl.InvCode = "101s"
	rq.AvailStatusMessages.AvailStatusMessages = append(msgs, AvailStatusMessage{
		BookingLimit:             new(3),
		StatusApplicationControl: &StatusApplicationControl{Start: date("2010-08-01"), End: date("2010-08-10"), InvTypeCode: "single"},
	})
	hotel := HotelInventory{
		Categories: map[string]struct{}{"double": {}},
		Rooms:      map[string]map[string]struct{}{"double": {"101S": {}}},
	}
	rulestest.Expect(t, CheckInventory(rq, hotel), []string{
		"AvailStatusMessage[1].StatusApplicationControl: unknown room double/101s",
		"AvailStatusMessage[2].StatusApplicationControl: unknown category single",
	})
	rulestest.Expect(t, CheckInventory(rq, HotelInventory{}), nil)
	rulestest.Expect(t, CheckInventory(load(t, "FreeRooms-OTA_HotelAvailNotifRQ-empty.xml"), hotel), nil)
}
