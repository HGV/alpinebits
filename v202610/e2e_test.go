package v202610_test

import (
	"context"
	"encoding/xml"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v202410"
	"github.com/HGV/alpinebits/v202610"
	"github.com/HGV/alpinebits/v202610/freerooms"
)

var freeRoomsCaps = []alpinebits.Capability{
	freerooms.CapAcceptRooms, freerooms.CapAcceptCategories, freerooms.CapAcceptDeltas,
	freerooms.CapAcceptOutOfOrder, freerooms.CapAcceptOutOfMarket,
	freerooms.CapAcceptClosingSeasons, freerooms.CapAcceptCompleteSet,
}

// request is what the server saw of a request.
type request struct {
	contentType, action string
}

// serve serves 2024-10 and 2026-10 from one URL, as AlpineBits 2026-10 §2.1
// allows, and records the requests and what their handlers got.
func serve(t *testing.T) (url string, seen *[]request, got *[]freerooms.HotelInvCountNotifRQ) {
	t.Helper()
	seen, got = new([]request), new([]freerooms.HotelInvCountNotifRQ)
	h := func(_ context.Context, r alpinebits.Request[freerooms.HotelInvCountNotifRQ]) (freerooms.HotelInvCountNotifRS, error) {
		*got = append(*got, r.Data)
		return freerooms.HotelInvCountNotifRS{Response: alpinebits.OK()}, nil
	}
	r := alpinebits.NewRouter().
		Handle(v202410.ActionHotelInvCountNotif, h, alpinebits.WithCapabilities(freeRoomsCaps...))
	r.Handle(v202610.ActionHotelInvCountNotif, h, alpinebits.WithCapabilities(freeRoomsCaps...))
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mediaType, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
		*seen = append(*seen, request{mediaType, req.Header.Get(alpinebits.HeaderAction)})
		r.ServeHTTP(w, req)
	}))
	t.Cleanup(srv.Close)
	return srv.URL, seen, got
}

func TestClientUsesTransportOfItsVersion(t *testing.T) {
	url, seen, got := serve(t)
	b, err := os.ReadFile("freerooms/testdata/FreeRooms-OTA_HotelInvCountNotifRQ.xml")
	if err != nil {
		t.Fatal(err)
	}
	var rq freerooms.HotelInvCountNotifRQ
	if err := xml.Unmarshal(b, &rq); err != nil {
		t.Fatal(err)
	}

	for _, send := range []func(context.Context) (*freerooms.HotelInvCountNotifRS, error){
		func(ctx context.Context) (*freerooms.HotelInvCountNotifRS, error) {
			c := alpinebits.NewClient(url, "user", "secret", v202410.Version)
			return c.Send(ctx, v202410.ActionHotelInvCountNotif, rq)
		},
		func(ctx context.Context) (*freerooms.HotelInvCountNotifRS, error) {
			c := alpinebits.NewClient(url, "user", "secret", v202610.Version)
			return c.Send(ctx, v202610.ActionHotelInvCountNotif, rq)
		},
	} {
		rs, err := send(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		if rs.Success == nil {
			t.Errorf("no success: %+v", rs)
		}
	}

	want := []request{
		{"multipart/form-data", ""},
		{"application/xml", "OTA_HotelInvCountNotif:FreeRooms"},
	}
	if !reflect.DeepEqual(*seen, want) {
		t.Errorf("server saw %+v, want %+v", *seen, want)
	}
	if len(*got) != 2 || !reflect.DeepEqual((*got)[0], rq) || !reflect.DeepEqual((*got)[1], rq) {
		t.Errorf("handlers got %+v, want the sample twice", *got)
	}
}

// The official handshake announces 2026-10 and 2024-10; the server answers for
// both, as they are served from the same URL (AlpineBits 2026-10 §3.2).
func TestHandshakeSampleNegotiatesBothTransports(t *testing.T) {
	url, _, _ := serve(t)
	doc, err := os.ReadFile("testdata/Handshake-OTA_PingRQ.xml")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, strings.NewReader(string(doc)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set(alpinebits.HeaderVersion, "2026-10")
	req.Header.Set(alpinebits.HeaderAction, "OTA_Ping:Handshaking")
	req.SetBasicAuth("user", "secret")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d: %s", resp.StatusCode, body)
	}

	var rs alpinebits.PingRS
	if err := xml.Unmarshal(body, &rs); err != nil {
		t.Fatal(err)
	}
	want := `{"versions":[` +
		`{"version":"2026-10","actions":[{"action":"action_OTA_HotelInvCountNotif","supports":[` +
		`"OTA_HotelInvCountNotif_accept_rooms","OTA_HotelInvCountNotif_accept_categories","OTA_HotelInvCountNotif_accept_deltas",` +
		`"OTA_HotelInvCountNotif_accept_out_of_order","OTA_HotelInvCountNotif_accept_out_of_market",` +
		`"OTA_HotelInvCountNotif_accept_closing_seasons","OTA_HotelInvCountNotif_accept_complete_set"]}]},` +
		`{"version":"2024-10","actions":[]}]}`
	if got := rs.Warnings.Message; got != want {
		t.Errorf("negotiated\n got %s\nwant %s", got, want)
	}
}
