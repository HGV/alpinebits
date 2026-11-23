package v201810_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/v201810"
	"github.com/HGV/alpinebits/v201810/freerooms"
	"github.com/HGV/alpinebits/v201810/guestrequests"
	"github.com/HGV/alpinebits/v201810/inventory"
	"github.com/HGV/alpinebits/v201810/rateplans"
)

func addSamples(f *testing.F) {
	paths, err := filepath.Glob("*/testdata/*.xml")
	if err != nil || len(paths) == 0 {
		f.Fatal("no samples", err)
	}
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(b)
	}
}

// roundTrip checks that a valid doc unmarshals into T, passes check without
// panicking, and marshals to a valid document that unmarshals to the same T.
func roundTrip[T any](t *testing.T, doc []byte, check func(T)) {
	var v T
	if err := xml.Unmarshal(doc, &v); err != nil {
		t.Fatalf("valid document does not unmarshal: %v", err)
	}
	check(v)
	out, err := xml.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := v201810.Version.Validate(out); err != nil {
		t.Fatalf("marshals to an invalid document: %v\n%s", err, out)
	}
	var back T
	if err := xml.Unmarshal(out, &back); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back, v) {
		t.Fatalf("round trip changed the message:\n got %+v\nwant %+v\n%s", back, v, out)
	}
}

func FuzzMessages(f *testing.F) {
	addSamples(f)
	f.Fuzz(func(t *testing.T, doc []byte) {
		if v201810.Version.Validate(doc) != nil {
			return
		}
		var root struct{ XMLName xml.Name }
		if xml.Unmarshal(doc, &root) != nil {
			return // libxml2 accepts some XML encoding/xml does not, such as version="1.1"
		}
		switch root.XMLName.Local {
		case "OTA_HotelAvailNotifRQ":
			roundTrip(t, doc, func(rq freerooms.HotelAvailNotifRQ) {
				freerooms.Validate(rq, nil)
				freerooms.Validate(rq, freeRoomsCaps)
				freerooms.CheckInventory(rq, freerooms.HotelInventory{})
			})
		case "OTA_HotelAvailNotifRS":
			roundTrip(t, doc, func(freerooms.HotelAvailNotifRS) {})
		case "OTA_HotelRatePlanNotifRQ":
			roundTrip(t, doc, func(rq rateplans.HotelRatePlanNotifRQ) {
				rateplans.Validate(rq, nil)
				rateplans.Validate(rq, ratePlansCaps)
				rateplans.CheckInventory(rq, nil)
			})
		case "OTA_HotelRatePlanNotifRS":
			roundTrip(t, doc, func(rateplans.HotelRatePlanNotifRS) {})
		case "OTA_HotelDescriptiveContentNotifRQ":
			roundTrip(t, doc, func(rq inventory.HotelDescriptiveContentNotifRQ) {
				inventory.Validate(rq, nil)
				inventory.Validate(rq, inventoryCaps)
			})
		case "OTA_HotelDescriptiveContentNotifRS":
			roundTrip(t, doc, func(inventory.HotelDescriptiveContentNotifRS) {})
		case "OTA_ReadRQ":
			roundTrip(t, doc, func(rq guestrequests.ReadRQ) { guestrequests.ValidateReadRQ(rq, nil) })
		case "OTA_ResRetrieveRS":
			roundTrip(t, doc, func(rs guestrequests.ResRetrieveRS) { guestrequests.ValidateResRetrieveRS(rs) })
		case "OTA_NotifReportRQ":
			roundTrip(t, doc, func(rq guestrequests.NotifReportRQ) { guestrequests.ValidateNotifReportRQ(rq, nil) })
		case "OTA_NotifReportRS":
			roundTrip(t, doc, func(guestrequests.NotifReportRS) {})
		}
	})
}

// FuzzRouter checks that the router answers every request with an
// AlpineBits response or an ERROR string, but never fails itself.
func FuzzRouter(f *testing.F) {
	actions := map[string]string{
		"OTA_HotelAvailNotifRQ":              v201810.ActionHotelAvailNotif.Name(),
		"OTA_HotelRatePlanNotifRQ":           v201810.ActionHotelRatePlanNotif.Name(),
		"OTA_HotelDescriptiveContentNotifRQ": v201810.ActionHotelDescriptiveContentNotif.Name(),
		"OTA_ReadRQ":                         v201810.ActionRead.Name(),
		"OTA_NotifReportRQ":                  v201810.ActionNotifReport.Name(),
	}
	paths, err := filepath.Glob("*/testdata/*.xml")
	if err != nil {
		f.Fatal(err)
	}
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			f.Fatal(err)
		}
		var root struct{ XMLName xml.Name }
		if xml.Unmarshal(b, &root) == nil && actions[root.XMLName.Local] != "" {
			f.Add(actions[root.XMLName.Local], b)
		}
	}
	f.Add("OTA_Ping:Handshaking", []byte(`<OTA_PingRQ xmlns="http://www.opentravel.org/OTA/2003/05" Version="8.000"><EchoData>`+
		`{"versions":[{"version":"2018-10","actions":[{"action":"action_OTA_Ping"},`+
		`{"action":"action_OTA_HotelAvailNotif","supports":["OTA_HotelAvailNotif_accept_rooms"]}]}]}</EchoData></OTA_PingRQ>`))

	r := alpinebits.NewRouter(alpinebits.WithClientStore(alpinebits.NewInMemoryClientStore())).
		Handle(v201810.ActionHotelAvailNotif, ok[freerooms.HotelAvailNotifRQ, freerooms.HotelAvailNotifRS](
			freerooms.HotelAvailNotifRS{Response: alpinebits.OK()}), alpinebits.WithCapabilities(freeRoomsCaps...)).
		Handle(v201810.ActionHotelRatePlanNotif, ok[rateplans.HotelRatePlanNotifRQ](
			rateplans.HotelRatePlanNotifRS{Response: alpinebits.OK()}), alpinebits.WithCapabilities(ratePlansCaps...)).
		Handle(v201810.ActionHotelDescriptiveContentNotif, ok[inventory.HotelDescriptiveContentNotifRQ](
			inventory.HotelDescriptiveContentNotifRS{Response: alpinebits.OK()}), alpinebits.WithCapabilities(inventoryCaps...)).
		Handle(v201810.ActionRead, ok[guestrequests.ReadRQ](
			guestrequests.ResRetrieveRS{Response: alpinebits.OK()})).
		Handle(v201810.ActionNotifReport, ok[guestrequests.NotifReportRQ](
			guestrequests.NotifReportRS{Response: alpinebits.OK()}))
	f.Fuzz(func(t *testing.T, action string, doc []byte) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		w.WriteField(alpinebits.FormAction, action)
		w.WriteField(alpinebits.FormRequest, string(doc))
		w.Close()
		req := httptest.NewRequest(http.MethodPost, "/", &body)
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.Header.Set(alpinebits.HeaderVersion, "2018-10")
		req.SetBasicAuth(user, password)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK && rec.Code != http.StatusBadRequest {
			t.Fatalf("%s: status %d: %s", action, rec.Code, rec.Body)
		}
	})
}

func ok[RQ, RS any](rs RS) alpinebits.HandlerFunc[RQ, RS] {
	return func(context.Context, alpinebits.Request[RQ]) (RS, error) { return rs, nil }
}
