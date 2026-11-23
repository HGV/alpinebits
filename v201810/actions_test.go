package v201810

import (
	"bytes"
	"context"
	"encoding/xml"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v201810/freerooms"
	"github.com/HGV/alpinebits/v201810/guestrequests"
)

func post(t *testing.T, h http.Handler, action, doc string) (int, string) {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	w.WriteField(alpinebits.FormAction, action)
	w.WriteField(alpinebits.FormRequest, doc)
	w.Close()
	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set(alpinebits.HeaderVersion, "2018-10")
	req.SetBasicAuth("client", "secret")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

const ns = `xmlns="http://www.opentravel.org/OTA/2003/05"`

func freeRooms(uniqueID string, entries ...string) string {
	doc := `<OTA_HotelAvailNotifRQ ` + ns + ` Version="1.002">` + uniqueID + `<AvailStatusMessages HotelCode="123">`
	for _, e := range entries {
		doc += `<AvailStatusMessage BookingLimit="1" BookingLimitMessageType="SetLimit"><StatusApplicationControl ` + e + `/></AvailStatusMessage>`
	}
	return doc + `</AvailStatusMessages></OTA_HotelAvailNotifRQ>`
}

// Error responses must pass response validation, or the router answers 500.
func TestFreeRoomsAfterHandshake(t *testing.T) {
	r := alpinebits.NewRouter(alpinebits.WithClientStore(alpinebits.NewInMemoryClientStore())).
		Handle(ActionHotelAvailNotif,
			func(context.Context, alpinebits.Request[freerooms.HotelAvailNotifRQ]) (freerooms.HotelAvailNotifRS, error) {
				return freerooms.HotelAvailNotifRS{Response: alpinebits.OK()}, nil
			},
			alpinebits.WithCapabilities(freerooms.CapAcceptRooms, freerooms.CapAcceptDeltas)).
		Handle(ActionRead,
			func(context.Context, alpinebits.Request[guestrequests.ReadRQ]) (guestrequests.ResRetrieveRS, error) {
				return guestrequests.ResRetrieveRS{Response: alpinebits.OK()}, nil
			}).
		Handle(ActionNotifReport,
			func(context.Context, alpinebits.Request[guestrequests.NotifReportRQ]) (guestrequests.NotifReportRS, error) {
				return guestrequests.NotifReportRS{Response: alpinebits.OK()}, nil
			})

	client := `{"versions":[{"version":"2018-10","actions":[{"action":"action_OTA_Ping"},` +
		`{"action":"action_OTA_HotelAvailNotif","supports":["OTA_HotelAvailNotif_accept_rooms","OTA_HotelAvailNotif_accept_categories"]},` +
		`{"action":"action_OTA_Read"}]}]}`
	code, body := post(t, r, "OTA_Ping:Handshaking",
		`<OTA_PingRQ `+ns+` Version="8.000"><EchoData>`+client+`</EchoData></OTA_PingRQ>`)
	if code != http.StatusOK {
		t.Fatalf("handshake: %d %s", code, body)
	}
	var ping alpinebits.PingRS
	if err := xml.Unmarshal([]byte(body), &ping); err != nil {
		t.Fatal(err)
	}
	wantNegotiated := `{"versions":[{"version":"2018-10","actions":[` +
		`{"action":"action_OTA_HotelAvailNotif","supports":["OTA_HotelAvailNotif_accept_rooms"]},` +
		`{"action":"action_OTA_Ping"},{"action":"action_OTA_Read"}]}]}`
	if ping.Warnings.Message != wantNegotiated {
		t.Errorf("negotiated\n got %s\nwant %s", ping.Warnings.Message, wantNegotiated)
	}

	completeSet := `<UniqueID Type="16" ID="1" Instance="CompleteSet"/>`
	room101 := `Start="2026-01-01" End="2026-01-05" InvTypeCode="DZ" InvCode="101"`
	categoryDZ := `Start="2026-01-01" End="2026-01-05" InvTypeCode="DZ"`
	tests := []struct {
		name, action, doc string
		wantErrs          []string
	}{
		{"rooms", "OTA_HotelAvailNotif:FreeRooms", freeRooms(completeSet, room101), nil},
		{"delta, deltas not negotiated", "OTA_HotelAvailNotif:FreeRooms", freeRooms("", room101),
			[]string{"deltas not supported"}},
		{"rooms and categories mixed", "OTA_HotelAvailNotif:FreeRooms", freeRooms(completeSet, room101, categoryDZ),
			[]string{"rooms and categories must not be mixed"}},
		{"acknowledgement covered by action_OTA_Read", "OTA_NotifReport:GuestRequests",
			`<OTA_NotifReportRQ ` + ns + ` Version="7.000"><Success/><NotifDetails><HotelNotifReport><HotelReservations>` +
				`<HotelReservation><UniqueID Type="14" ID="abc"/></HotelReservation></HotelReservations></HotelNotifReport></NotifDetails></OTA_NotifReportRQ>`,
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body := post(t, r, tt.action, tt.doc)
			if code != http.StatusOK {
				t.Fatalf("status %d: %s", code, body)
			}
			var rs struct {
				alpinebits.Response
			}
			if err := xml.Unmarshal([]byte(body), &rs); err != nil {
				t.Fatal(err)
			}
			var got []string
			if rs.Errors != nil {
				for _, e := range *rs.Errors {
					got = append(got, e.Message)
				}
			}
			if !reflect.DeepEqual(got, tt.wantErrs) {
				t.Errorf("errors %q, want %q\n%s", got, tt.wantErrs, body)
			}
			if tt.wantErrs == nil && rs.Success == nil {
				t.Errorf("want Success: %s", body)
			}
		})
	}
}

func TestSchemaViolation(t *testing.T) {
	r := alpinebits.NewRouter().
		Handle(ActionHotelAvailNotif,
			func(context.Context, alpinebits.Request[freerooms.HotelAvailNotifRQ]) (freerooms.HotelAvailNotifRS, error) {
				t.Error("handler called")
				return freerooms.HotelAvailNotifRS{}, nil
			})
	code, body := post(t, r, "OTA_HotelAvailNotif:FreeRooms", `<OTA_HotelAvailNotifRQ `+ns+` Version="1.002"><Bogus/></OTA_HotelAvailNotifRQ>`)
	if code != http.StatusBadRequest || !strings.HasPrefix(body, "ERROR:XML validation error: ") || !strings.Contains(body, "Bogus") {
		t.Errorf("got %d %s", code, body)
	}
}
