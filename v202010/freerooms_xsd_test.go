package v202010

import (
	"encoding/xml"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202010/freerooms"
)

var freeRoomsCaps = []alpinebits.Capability{
	freerooms.CapAcceptRooms, freerooms.CapAcceptCategories, freerooms.CapAcceptDeltas,
	freerooms.CapAcceptOutOfOrder, freerooms.CapAcceptOutOfMarket, freerooms.CapAcceptClosingSeasons,
}

const freeRoomsDoc = `<OTA_HotelInvCountNotifRQ xmlns="http://www.opentravel.org/OTA/2003/05" Version="4">
<UniqueID Type="16" ID="1" Instance="CompleteSet"/>
<Inventories HotelCode="123" HotelName="Frangart Inn">
<Inventory><StatusApplicationControl Start="2020-08-31" End="2020-09-30" AllInvCode="true"/></Inventory>
<Inventory><StatusApplicationControl Start="2020-08-01" End="2020-08-10" InvTypeCode="DOUBLE" InvCode="101"/>
<InvCounts><InvCount CountType="2" Count="1"/></InvCounts></Inventory>
</Inventories>
</OTA_HotelInvCountNotifRQ>`

func replaceOnce(t *testing.T, doc, old, new string) string {
	t.Helper()
	if n := strings.Count(doc, old); n != 1 {
		t.Fatalf("%q occurs %d times", old, n)
	}
	return strings.Replace(doc, old, new, 1)
}

// freerooms.Validate relies on these being rejected by the XSD.
func TestXSDEnforcesFreeRoomsRules(t *testing.T) {
	if err := Version.Validate([]byte(freeRoomsDoc)); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	inventories := freeRoomsDoc[strings.Index(freeRoomsDoc, "<Inventory>"):strings.Index(freeRoomsDoc, "</Inventories>")]
	count := `<InvCount CountType="2" Count="1"/>`
	control := `<StatusApplicationControl Start="2020-08-01" End="2020-08-10" InvTypeCode="DOUBLE" InvCode="101"/>`
	tests := []struct{ name, old, new string }{
		{"no Inventory", inventories, ""},
		{"two UniqueIDs", `<UniqueID Type="16" ID="1" Instance="CompleteSet"/>`, `<UniqueID Type="16" ID="1" Instance="CompleteSet"/><UniqueID Type="16" ID="1" Instance="CompleteSet"/>`},
		{"UniqueID Type", `Type="16"`, `Type="17"`},
		{"UniqueID without ID", ` ID="1"`, ``},
		{"UniqueID Instance", `Instance="CompleteSet"`, `Instance="Delta"`},
		{"empty HotelCode", `HotelCode="123"`, `HotelCode=""`},
		{"empty HotelName", `HotelName="Frangart Inn"`, `HotelName=""`},
		{"no Start", `Start="2020-08-01" `, ``},
		{"no End", `End="2020-08-10" `, ``},
		{"invalid date", `End="2020-08-10"`, `End="2020-08-32"`},
		{"empty InvTypeCode", `InvTypeCode="DOUBLE"`, `InvTypeCode=""`},
		{"InvTypeCode longer than 8", `InvTypeCode="DOUBLE"`, `InvTypeCode="DOUBLEROOM"`},
		{"empty InvCode", `InvCode="101"`, `InvCode=""`},
		{"AllInvCode", `AllInvCode="true"`, `AllInvCode="yes"`},
		{"two StatusApplicationControl", control, control + control},
		{"InvCounts before StatusApplicationControl", control + "\n<InvCounts>" + count + "</InvCounts>", "<InvCounts>" + count + "</InvCounts>" + control},
		{"empty InvCounts", count, ""},
		{"four InvCount", count, strings.Repeat(count, 4)},
		{"CountType", `CountType="2"`, `CountType="3"`},
		{"negative Count", `Count="1"`, `Count="-1"`},
		{"no Count", ` Count="1"`, ``},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Version.Validate([]byte(replaceOnce(t, freeRoomsDoc, tt.old, tt.new))); err == nil {
				t.Error("XSD accepted it")
			}
		})
	}
}

// The XSD accepts these, so freerooms.Validate must reject them.
func TestXSDLeavesFreeRoomsRulesToValidate(t *testing.T) {
	if errs := validateFreeRooms(t, freeRoomsDoc); len(errs) > 0 {
		t.Fatalf("baseline: %v", errs)
	}
	room := `<Inventory><StatusApplicationControl Start="2020-08-01" End="2020-08-10" InvTypeCode="DOUBLE" InvCode="101"/>`
	tests := []struct{ name, old, new, want string }{
		{"HotelName only", ` HotelCode="123"`, ``, "HotelCode required"},
		{"blank InvTypeCode", `InvTypeCode="DOUBLE"`, `InvTypeCode=" "`, "Inventory[1]: InvTypeCode required"},
		{"InvCode without InvTypeCode", ` InvTypeCode="DOUBLE"`, ``, "Inventory[1]: InvTypeCode required"},
		{"End before Start", `End="2020-08-10"`, `End="2020-07-31"`, "Inventory[1]: End before Start"},
		{"duplicate CountType", `<InvCount CountType="2" Count="1"/>`, `<InvCount CountType="2" Count="0"/><InvCount CountType="2" Count="1"/>`,
			"Inventory[1].InvCount[1]: duplicate CountType 2"},
		{"room counted twice", `Count="1"`, `Count="2"`, "Inventory[1]: Counts of a room must sum to at most 1"},
		{"empty Inventory among others", `</Inventories>`, `<Inventory/></Inventories>`, "Inventory[2]: StatusApplicationControl required"},
		{"overlapping periods", `</Inventories>`, room + `</Inventory></Inventories>`, "Inventory[2]: overlaps Inventory[1]"},
		{"rooms and categories mixed", `</Inventories>`, strings.Replace(room, ` InvCode="101"`, ``, 1) + `</Inventory></Inventories>`,
			"rooms and categories must not be mixed"},
		{"closing season in a delta", `<UniqueID Type="16" ID="1" Instance="CompleteSet"/>`, ``,
			"Inventory[0]: closing seasons only allowed in a CompleteSet"},
		{"closing season after other Inventory", `</Inventories>`, `<Inventory><StatusApplicationControl Start="2020-10-01" End="2020-10-31" AllInvCode="true"/></Inventory></Inventories>`,
			"Inventory[2]: closing seasons must precede the other Inventory elements"},
		{"closing season with InvCounts", `AllInvCode="true"/>`, `AllInvCode="true"/><InvCounts><InvCount CountType="2" Count="0"/></InvCounts>`,
			"Inventory[0]: closing season must have no InvCounts"},
		{"closing season with InvTypeCode", `AllInvCode="true"`, `AllInvCode="true" InvTypeCode="DOUBLE"`,
			"Inventory[0]: closing season must have no InvTypeCode and no InvCode"},
		{"closing season overlapping availability", `Start="2020-08-31"`, `Start="2020-08-05"`,
			"Inventory[1]: available rooms overlap closing season Inventory[0]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validateFreeRooms(t, replaceOnce(t, freeRoomsDoc, tt.old, tt.new))
			if len(errs) != 1 || errs[0].Message != tt.want {
				t.Errorf("got %v, want %q", errs, tt.want)
			}
		})
	}
}

func validateFreeRooms(t *testing.T, doc string) alpinebits.Errors {
	t.Helper()
	if err := Version.Validate([]byte(doc)); err != nil {
		t.Fatalf("XSD: %v", err)
	}
	var rq freerooms.HotelInvCountNotifRQ
	if err := xml.Unmarshal([]byte(doc), &rq); err != nil {
		t.Fatal(err)
	}
	return freerooms.Validate(rq, freeRoomsCaps)
}

func freeRoomsSample[T any](t *testing.T, name string) T {
	t.Helper()
	b, err := os.ReadFile("freerooms/testdata/FreeRooms-" + name + ".xml")
	if err != nil {
		t.Fatal(err)
	}
	var v T
	if err := xml.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func roundTrip[T any](t *testing.T, v T) {
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
}

// Optional attributes, elements and zero counts must survive marshalling and
// pass the XSD.
func TestFreeRoomsTypesRoundTrip(t *testing.T) {
	sample := func(name string) freerooms.HotelInvCountNotifRQ {
		return freeRoomsSample[freerooms.HotelInvCountNotifRQ](t, "OTA_HotelInvCountNotifRQ"+name)
	}
	tests := []struct {
		name string
		rq   func() freerooms.HotelInvCountNotifRQ
	}{
		{"CompleteSet with all CountTypes", func() freerooms.HotelInvCountNotifRQ {
			rq := sample("")
			rq.UniqueID.Type = freerooms.UniqueIDTypePurgedMasterReference
			rq.Inventories.Inventories[0].InvCounts = &[]freerooms.InvCount{
				{CountType: freerooms.CountTypeBookable, Count: 0},
				{CountType: freerooms.CountTypeOutOfOrder, Count: 1},
				{CountType: freerooms.CountTypeOutOfMarket, Count: 2},
			}
			return rq
		}},
		{"closing seasons by HotelName", func() freerooms.HotelInvCountNotifRQ {
			rq := sample("-closing_seasons")
			rq.Inventories.HotelCode = ""
			return rq
		}},
		{"delta of rooms by HotelCode", func() freerooms.HotelInvCountNotifRQ {
			rq := sample("-delta")
			rq.Inventories.HotelName = ""
			for _, inv := range rq.Inventories.Inventories {
				inv.StatusApplicationControl.InvCode = "101"
			}
			return rq
		}},
		{"reset", func() freerooms.HotelInvCountNotifRQ { return sample("-empty") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roundTrip(t, tt.rq())
		})
	}

	for _, name := range []string{"success", "advisory", "warning", "error"} {
		t.Run("RS "+name, func(t *testing.T) {
			roundTrip(t, freeRoomsSample[freerooms.HotelInvCountNotifRS](t, "OTA_HotelInvCountNotifRS-"+name))
		})
	}
}
