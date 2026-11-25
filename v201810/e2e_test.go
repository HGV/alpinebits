package v201810_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"log/slog"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"slices"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v201810"
	"github.com/HGV/alpinebits/v201810/freerooms"
	"github.com/HGV/alpinebits/v201810/guestrequests"
	"github.com/HGV/alpinebits/v201810/inventory"
	"github.com/HGV/alpinebits/v201810/rateplans"
)

const user, password = "client", "secret"

var (
	freeRoomsCaps = []alpinebits.Capability{
		freerooms.CapAcceptRooms, freerooms.CapAcceptCategories,
		freerooms.CapAcceptDeltas, freerooms.CapAcceptBookingThreshold,
	}
	ratePlansCaps = []alpinebits.Capability{
		rateplans.CapAcceptArrivalDOW, rateplans.CapAcceptDepartureDOW,
		rateplans.CapAcceptRatePlanBookingRule, rateplans.CapAcceptRatePlanRoomTypeBookingRule,
		rateplans.CapAcceptRatePlanMixedBookingRule, rateplans.CapAcceptSupplements,
		rateplans.CapAcceptFreeNightsOffers, rateplans.CapAcceptFamilyOffers,
		rateplans.CapAcceptOverlay, rateplans.CapAcceptRatePlanJoin,
		rateplans.CapAcceptOfferRuleBookingOffset, rateplans.CapAcceptOfferRuleDOWLOS,
	}
	inventoryCaps = []alpinebits.Capability{inventory.CapUseRooms, inventory.CapOccupancyChildren}
)

func sample[T any](t *testing.T, path string) T {
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

// serve runs r behind basic access authentication, as AlpineBits §2 requires,
// and returns its URL.
func serve(t *testing.T, r *alpinebits.Router) string {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if u, p, ok := req.BasicAuth(); !ok || u != user || p != password {
			alpinebits.Unauthorized(w, "invalid or missing username/password")
			return
		}
		r.ServeHTTP(w, req)
	}))
	t.Cleanup(srv.Close)
	return srv.URL
}

func client(url string) *alpinebits.Client {
	return alpinebits.NewClient(url, user, password, v201810.Version)
}

// push sends rq through a router that handles a with caps and answers ok, and
// checks that the handler gets rq and the client gets ok.
func push[RQ, RS any](t *testing.T, a alpinebits.Action[RQ, RS], rq RQ, ok RS, caps ...alpinebits.Capability) {
	t.Helper()
	var got RQ
	r := alpinebits.NewRouter().
		Handle(a, func(_ context.Context, r alpinebits.Request[RQ]) (RS, error) {
			got = r.Data
			return ok, nil
		}, alpinebits.WithCapabilities(caps...))

	rs, err := client(serve(t, r)).Send(t.Context(), a, rq)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, rq) {
		t.Errorf("handler got\n%+v\nwant\n%+v", got, rq)
	}
	if gotXML, wantXML := marshal(t, *rs), marshal(t, ok); !bytes.Equal(gotXML, wantXML) {
		t.Errorf("client got %s, want %s", gotXML, wantXML)
	}
}

func marshal(t *testing.T, v any) []byte {
	t.Helper()
	b, err := xml.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestPushSamples(t *testing.T) {
	t.Run("FreeRooms", func(t *testing.T) {
		push(t, v201810.ActionHotelAvailNotif,
			sample[freerooms.HotelAvailNotifRQ](t, "freerooms/testdata/FreeRooms-OTA_HotelAvailNotifRQ.xml"),
			freerooms.HotelAvailNotifRS{Response: alpinebits.OK()}, freeRoomsCaps...)
	})
	t.Run("RatePlans", func(t *testing.T) {
		push(t, v201810.ActionHotelRatePlanNotif,
			sample[rateplans.HotelRatePlanNotifRQ](t, "rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRQ.xml"),
			rateplans.HotelRatePlanNotifRS{Response: alpinebits.OK()}, ratePlansCaps...)
	})
	for _, name := range []string{"", "-with-roomtype", "-delete-all"} {
		t.Run("Inventory"+name, func(t *testing.T) {
			push(t, v201810.ActionHotelDescriptiveContentNotif,
				sample[inventory.HotelDescriptiveContentNotifRQ](t, "inventory/testdata/Inventory-OTA_HotelDescriptiveContentNotifRQ"+name+".xml"),
				inventory.HotelDescriptiveContentNotifRS{Response: alpinebits.OK()}, inventoryCaps...)
		})
	}
	t.Run("GuestRequests acknowledgement", func(t *testing.T) {
		push(t, v201810.ActionNotifReport,
			sample[guestrequests.NotifReportRQ](t, "guestrequests/testdata/GuestRequests-OTA_ReadRQ-ack.xml"),
			guestrequests.NotifReportRS{Response: alpinebits.OK()})
	})
}

func TestPullGuestRequests(t *testing.T) {
	for _, name := range []string{"reservation", "cancellation", "reservation-empty"} {
		t.Run(name, func(t *testing.T) {
			rs := sample[guestrequests.ResRetrieveRS](t, "guestrequests/testdata/GuestRequests-OTA_ResRetrieveRS-"+name+".xml")
			push(t, v201810.ActionRead,
				sample[guestrequests.ReadRQ](t, "guestrequests/testdata/GuestRequests-OTA_ReadRQ.xml"), rs)
		})
	}
}

func TestHandshakeLimitsLaterRequests(t *testing.T) {
	okFreeRooms := func(context.Context, alpinebits.Request[freerooms.HotelAvailNotifRQ]) (freerooms.HotelAvailNotifRS, error) {
		return freerooms.HotelAvailNotifRS{Response: alpinebits.OK()}, nil
	}
	okRatePlans := func(context.Context, alpinebits.Request[rateplans.HotelRatePlanNotifRQ]) (rateplans.HotelRatePlanNotifRS, error) {
		return rateplans.HotelRatePlanNotifRS{Response: alpinebits.OK()}, nil
	}
	r := alpinebits.NewRouter(alpinebits.WithClientStore(alpinebits.NewInMemoryClientStore())).
		Handle(v201810.ActionHotelAvailNotif, okFreeRooms, alpinebits.WithCapabilities(freeRoomsCaps...)).
		Handle(v201810.ActionHotelRatePlanNotif, okRatePlans, alpinebits.WithCapabilities(ratePlansCaps...))
	c := client(serve(t, r))

	negotiated, err := c.Handshake(t.Context(), alpinebits.ActionCapabilities{
		"action_OTA_Ping":            nil,
		"action_OTA_HotelAvailNotif": {freerooms.CapAcceptCategories, freerooms.CapAcceptDeltas},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := alpinebits.ActionCapabilities{
		"action_OTA_Ping":            nil,
		"action_OTA_HotelAvailNotif": {freerooms.CapAcceptCategories, freerooms.CapAcceptDeltas},
	}
	if !maps.EqualFunc(negotiated, want, slices.Equal) {
		t.Errorf("negotiated %v, want %v", negotiated, want)
	}

	// The sample sends rooms with BookingThreshold, which were not negotiated.
	rq := sample[freerooms.HotelAvailNotifRQ](t, "freerooms/testdata/FreeRooms-OTA_HotelAvailNotifRQ.xml")
	rs, err := c.Send(t.Context(), v201810.ActionHotelAvailNotif, rq)
	if err != nil {
		t.Fatal(err)
	}
	if got := messages(rs.Errors); !reflect.DeepEqual(got, []string{
		"AvailStatusMessage[0]: BookingThreshold not supported",
		"rooms not supported",
	}) {
		t.Errorf("errors %q", got)
	}

	_, err = c.Send(t.Context(), v201810.ActionHotelRatePlanNotif,
		sample[rateplans.HotelRatePlanNotifRQ](t, "rateplans/testdata/RatePlans-OTA_HotelRatePlanNotifRQ.xml"))
	if se, ok := errors.AsType[*alpinebits.ServerError](err); !ok || *se != (alpinebits.ServerError{StatusCode: 400, Message: "action not negotiated"}) {
		t.Errorf("RatePlans after handshake: %v", err)
	}
}

func messages(errs *alpinebits.Errors) []string {
	if errs == nil {
		return nil
	}
	var msgs []string
	for _, e := range *errs {
		msgs = append(msgs, e.Message)
	}
	return msgs
}

func TestFailuresReachClient(t *testing.T) {
	rq := sample[freerooms.HotelAvailNotifRQ](t, "freerooms/testdata/FreeRooms-OTA_HotelAvailNotifRQ.xml")
	var handlerErr error
	r := alpinebits.NewRouter().
		Handle(v201810.ActionHotelAvailNotif,
			func(context.Context, alpinebits.Request[freerooms.HotelAvailNotifRQ]) (freerooms.HotelAvailNotifRS, error) {
				return freerooms.HotelAvailNotifRS{Response: alpinebits.OK()}, handlerErr
			}, alpinebits.WithCapabilities(freeRoomsCaps...))
	url := serve(t, r)
	c := client(url)

	t.Run("rules", func(t *testing.T) {
		rq := rq
		rq.AvailStatusMessages.HotelCode, rq.AvailStatusMessages.HotelName = "", ""
		rs, err := c.Send(t.Context(), v201810.ActionHotelAvailNotif, rq)
		if err != nil {
			t.Fatal(err)
		}
		if got := messages(rs.Errors); !reflect.DeepEqual(got, []string{"HotelCode required"}) {
			t.Errorf("errors %q", got)
		}
	})
	t.Run("rejected", func(t *testing.T) {
		handlerErr = alpinebits.Reject(alpinebits.Error{Type: alpinebits.ErrTypeApplicationError, Code: 392, Message: "unknown hotel"})
		defer func() { handlerErr = nil }()
		rs, err := c.Send(t.Context(), v201810.ActionHotelAvailNotif, rq)
		if err != nil {
			t.Fatal(err)
		}
		if got := messages(rs.Errors); !reflect.DeepEqual(got, []string{"unknown hotel"}) || rs.Success != nil {
			t.Errorf("errors %q, success %v", got, rs.Success)
		}
	})
	t.Run("handler fails", func(t *testing.T) {
		handlerErr = errors.New("database down")
		defer func() { handlerErr = nil }()
		defer slog.SetDefault(slog.Default())
		slog.SetDefault(slog.New(slog.DiscardHandler))
		_, err := c.Send(t.Context(), v201810.ActionHotelAvailNotif, rq)
		if se, ok := errors.AsType[*alpinebits.ServerError](err); !ok || se.StatusCode != 500 {
			t.Errorf("got %v, want status 500", err)
		}
	})
	t.Run("wrong password", func(t *testing.T) {
		c := alpinebits.NewClient(url, user, "wrong", v201810.Version)
		_, err := c.Send(t.Context(), v201810.ActionHotelAvailNotif, rq)
		want := alpinebits.ServerError{StatusCode: 401, Message: "invalid or missing username/password"}
		if se, ok := errors.AsType[*alpinebits.ServerError](err); !ok || *se != want {
			t.Errorf("got %v, want %+v", err, want)
		}
	})
}
