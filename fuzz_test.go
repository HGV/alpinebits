package alpinebits

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var xsDecimal = regexp.MustCompile(`^[+-]?(\d+(\.\d*)?|\.\d+)$`)

func FuzzDecimal(f *testing.F) {
	for _, s := range []string{"0", " 76.80 ", "-0.5", "1500000", "0.00001", ".5", "+3.", "1e5", "NaN"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		var d Decimal
		err := d.UnmarshalText([]byte(s))
		if trimmed := strings.TrimSpace(s); err != nil && len(trimmed) <= 20 && xsDecimal.MatchString(trimmed) {
			t.Fatalf("xs:decimal %q rejected: %v", s, err)
		}
		if err != nil || math.IsInf(float64(d), 0) || math.IsNaN(float64(d)) {
			return
		}
		out, err := d.MarshalText()
		if err != nil || !xsDecimal.Match(out) {
			t.Fatalf("%q marshals to %q, %v; want an xs:decimal", s, out, err)
		}
		var back Decimal
		if err := back.UnmarshalText(out); err != nil || back != d {
			t.Fatalf("%q: %q unmarshals to %v, %v; want %v", s, out, back, err, d)
		}
	})
}

func FuzzHandshakeJSON(f *testing.F) {
	f.Add(`{"versions":[{"version":"2024-10","actions":[{"action":"action_A","supports":["deltas"]}]},{"version":"2018-10"}]}`)
	f.Add(`{}`)
	f.Fuzz(func(t *testing.T, s string) {
		caps, err := parseCapabilities(s)
		if err != nil {
			return
		}
		back, err := parseCapabilities(encodeCapabilities(caps))
		if err != nil || !reflect.DeepEqual(back, caps) {
			t.Fatalf("%q: parsed %v, encoded and parsed again %v, %v", s, caps, back, err)
		}
	})
}

// FuzzServeHTTP checks that the router answers any request in either transport
// without failing itself.
func FuzzServeHTTP(f *testing.F) {
	f.Add("application/xml", "2026-10", "OTA_Test:Test", "", `<TestRQ/>`)
	f.Add("application/xml", "2026-10", "OTA_Ping:Handshaking", "",
		ping(`{"versions":[{"version":"2026-10"},{"version":"2024-10","actions":[{"action":"action_OTA_Test"}]}]}`))
	f.Add("multipart/form-data; boundary=b", "2024-10", "", "",
		"--b\r\nContent-Disposition: form-data; name=\"action\"\r\n\r\nOTA_Test:Test\r\n"+
			"--b\r\nContent-Disposition: form-data; name=\"request\"\r\n\r\n<TestRQ/>\r\n--b--\r\n")
	f.Add("application/xml", "", "OTA_Test:Test", "gzip", `<TestRQ Bad="maybe"/>`)

	ok := func(context.Context, Request[testRQ]) (testRS, error) { return testRS{Response: OK()}, nil }
	r := NewRouter(WithMaxFormSize(1<<16), WithClientStore(NewInMemoryClientStore()))
	for _, name := range []string{"2024-10", "2026-10"} {
		v := NewVersion(name, nopValidator{})
		r.Handle(NewAction[testRQ, testRS](v, "OTA_Test:Test", "action_OTA_Test",
			func(rq testRQ, _ []Capability) Errors {
				if rq.Bad {
					return Errors{{Type: ErrTypeApplicationError, Code: ErrCodeInvalidValue, Message: "bad"}}
				}
				return nil
			}), ok)
	}
	f.Fuzz(func(t *testing.T, contentType, version, action, encoding, body string) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", contentType)
		req.Header.Set(HeaderVersion, version)
		req.Header.Set(HeaderAction, action)
		req.Header.Set("Content-Encoding", encoding)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		switch rec.Code {
		case http.StatusOK, http.StatusBadRequest, http.StatusRequestEntityTooLarge:
		default:
			t.Fatalf("status %d: %s", rec.Code, rec.Body)
		}
	})
}
