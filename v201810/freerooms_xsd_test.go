package v201810

import (
	"encoding/xml"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/HGV/alpinebits/v201810/freerooms"
	"github.com/HGV/x/timex"
)

const freeRoomsNotif = `<OTA_HotelAvailNotifRQ ` + ns + ` Version="1.002">
<UniqueID Type="16" ID="1" Instance="CompleteSet"/>
<AvailStatusMessages HotelCode="123" HotelName="Frangart Inn">
<AvailStatusMessage BookingLimit="1" BookingLimitMessageType="SetLimit" BookingThreshold="0">
<StatusApplicationControl Start="2010-08-01" End="2010-08-10" InvTypeCode="double" InvCode="101S"/>
</AvailStatusMessage>
</AvailStatusMessages>
</OTA_HotelAvailNotifRQ>`

func TestFreeRoomsXSD(t *testing.T) {
	if err := Version.Validate([]byte(freeRoomsNotif)); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	control := `<StatusApplicationControl Start="2010-08-01" End="2010-08-10" InvTypeCode="double" InvCode="101S"/>`
	message := freeRoomsNotif[strings.Index(freeRoomsNotif, "<AvailStatusMessage "):strings.Index(freeRoomsNotif, "</AvailStatusMessages>")]
	tests := []struct {
		name, old, new string
		valid          bool // freerooms.Validate checks what the XSD accepts
	}{
		{"UniqueID Type", `Type="16"`, `Type="17"`, false},
		{"UniqueID Instance", `Instance="CompleteSet"`, `Instance="Delta"`, false},
		{"UniqueID without ID", ` ID="1"`, ``, false},
		{"two AvailStatusMessages", `</AvailStatusMessages>`, `</AvailStatusMessages><AvailStatusMessages HotelCode="1"><AvailStatusMessage/></AvailStatusMessages>`, false},
		{"no AvailStatusMessage", message, ``, false},
		{"empty HotelCode", `HotelCode="123"`, `HotelCode=""`, false},
		{"empty HotelName", `HotelName="Frangart Inn"`, `HotelName=""`, false},
		{"BookingLimitMessageType", `"SetLimit"`, `"AdjustLimit"`, false},
		{"negative BookingLimit", `BookingLimit="1"`, `BookingLimit="-1"`, false},
		{"negative BookingThreshold", `BookingThreshold="0"`, `BookingThreshold="-1"`, false},
		{"two StatusApplicationControl", control, control + control, false},
		{"StatusApplicationControl without Start", `Start="2010-08-01" `, ``, false},
		{"StatusApplicationControl without End", `End="2010-08-10" `, ``, false},
		{"invalid date", `"2010-08-10"`, `"2010-02-30"`, false},
		{"InvTypeCode too long", `InvTypeCode="double"`, `InvTypeCode="doubleroom"`, false},
		{"empty InvTypeCode", `InvTypeCode="double"`, `InvTypeCode=""`, false},
		{"empty InvCode", `InvCode="101S"`, `InvCode=""`, false},

		{"neither HotelCode nor HotelName", ` HotelCode="123" HotelName="Frangart Inn"`, ``, true},
		{"blank HotelCode only", `HotelCode="123" HotelName="Frangart Inn"`, `HotelCode=" "`, true},
		{"no StatusApplicationControl", control, ``, true},
		{"no InvTypeCode", ` InvTypeCode="double"`, ``, true},
		{"no BookingLimit", `BookingLimit="1" `, ``, true},
		{"BookingLimit 2 for a room", `BookingLimit="1"`, `BookingLimit="2"`, true},
		{"BookingThreshold exceeds BookingLimit", `BookingThreshold="0"`, `BookingThreshold="5"`, true},
		{"End before Start", `End="2010-08-10"`, `End="2010-07-01"`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if n := strings.Count(freeRoomsNotif, tt.old); n != 1 {
				t.Fatalf("%q occurs %d times", tt.old, n)
			}
			doc := strings.Replace(freeRoomsNotif, tt.old, tt.new, 1)
			err := Version.Validate([]byte(doc))
			if tt.valid && err != nil {
				t.Errorf("XSD rejected it: %v", err)
			}
			if !tt.valid && err == nil {
				t.Error("XSD accepted it")
			}
		})
	}
}

func readFreeRoomsSample[T any](t *testing.T, name string) T {
	t.Helper()
	b, err := os.ReadFile("freerooms/testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	var v T
	if err := xml.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func freeRoomsRoundTrip[T any](t *testing.T, v T) T {
	t.Helper()
	out, err := xml.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := Version.Validate(out); err != nil {
		t.Fatalf("XSD: %v\n%s", err, out)
	}
	var back T
	if err := xml.Unmarshal(out, &back); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back, v) {
		t.Errorf("round trip changed the message:\n got %+v\nwant %+v\n%s", back, v, out)
	}
	return back
}

// Zero limits and absent optional attributes must survive marshalling and pass
// the XSD.
func TestFreeRoomsTypesRoundTrip(t *testing.T) {
	rq := readFreeRoomsSample[freerooms.HotelAvailNotifRQ](t, "FreeRooms-OTA_HotelAvailNotifRQ.xml")
	rq.UniqueID.Type = freerooms.UniqueIDTypePurgedMasterReference
	msgs := rq.AvailStatusMessages.AvailStatusMessages
	*msgs[0].BookingLimit = 0
	msgs[1].BookingLimitMessageType = ""
	msgs[1].BookingThreshold = nil
	start, end := timex.Date{Year: 2010, Month: 9, Day: 1}, timex.Date{Year: 2010, Month: 9, Day: 30}
	rq.AvailStatusMessages.AvailStatusMessages = append(msgs, freerooms.AvailStatusMessage{
		BookingLimit:             new(4),
		BookingLimitMessageType:  freerooms.BookingLimitMessageTypeSetLimit,
		BookingThreshold:         new(1),
		StatusApplicationControl: &freerooms.StatusApplicationControl{Start: start, End: end, InvTypeCode: "single"},
	})
	freeRoomsRoundTrip(t, rq)

	delta := readFreeRoomsSample[freerooms.HotelAvailNotifRQ](t, "FreeRooms-OTA_HotelAvailNotifRQ.xml")
	delta.UniqueID = nil
	delta.AvailStatusMessages.HotelCode = ""
	freeRoomsRoundTrip(t, delta)
}

// The reset's AvailStatusMessage must stay empty, without attributes.
func TestFreeRoomsResetRoundTrip(t *testing.T) {
	rq := readFreeRoomsSample[freerooms.HotelAvailNotifRQ](t, "FreeRooms-OTA_HotelAvailNotifRQ-empty.xml")
	if back := freeRoomsRoundTrip(t, rq); !back.IsReset() {
		t.Error("no reset after the round trip")
	}
}

func TestFreeRoomsResponsesRoundTrip(t *testing.T) {
	for _, name := range []string{"success", "advisory", "warning", "error"} {
		t.Run(name, func(t *testing.T) {
			rs := readFreeRoomsSample[freerooms.HotelAvailNotifRS](t, "FreeRooms-OTA_HotelAvailNotifRS-"+name+".xml")
			freeRoomsRoundTrip(t, rs)
		})
	}
}
