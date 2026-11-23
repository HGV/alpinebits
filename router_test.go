package alpinebits

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newRequest posts as the basic access authentication user client.
func newRequest(version, client, action, doc string) *http.Request {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	w.WriteField(FormAction, action)
	w.WriteField(FormRequest, doc)
	w.Close()
	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set(HeaderVersion, version)
	req.SetBasicAuth(client, "secret")
	return req
}

func serve(h http.Handler, req *http.Request) (int, string) {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func post(t *testing.T, h http.Handler, version, client, action, doc string) (int, string) {
	t.Helper()
	return serve(h, newRequest(version, client, action, doc))
}

func decodeRS(t *testing.T, body string) testRS {
	t.Helper()
	var rs testRS
	if err := xml.Unmarshal([]byte(body), &rs); err != nil {
		t.Fatalf("decode %q: %v", body, err)
	}
	return rs
}

func TestHandleRunsRulesAndRejections(t *testing.T) {
	v := NewVersion("2018-10", nopValidator{})
	a := NewAction[testRQ, testRS](v, "OTA_Test:Test", "action_OTA_Test",
		func(rq testRQ, _ []Capability) Errors {
			var errs Errors
			if rq.Bad {
				errs.Addf(ErrCodeInvalidValue, At("Item", 2), "bad")
			}
			return errs
		})

	var calls int
	var reply error
	r := NewRouter()
	r.Version(v, Handle(a, func(context.Context, Request[testRQ]) (testRS, error) {
		calls++
		return testRS{Response: OK()}, reply
	}))

	tests := []struct {
		name      string
		doc       string
		reply     error
		wantCalls int
		wantErrs  []string
	}{
		{"valid", `<TestRQ/>`, nil, 1, nil},
		{"rule fails, handler not called", `<TestRQ Bad="true"/>`, nil, 0, []string{"Item[2]: bad"}},
		{"handler rejects", `<TestRQ/>`, Reject(Error{Type: ErrTypeApplicationError, Message: "no room"}), 1, []string{"no room"}},
		{"handler rejects nothing", `<TestRQ/>`, Reject(), 1, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls, reply = 0, tt.reply
			code, body := post(t, r, "2018-10", "c", "OTA_Test:Test", tt.doc)
			if code != http.StatusOK {
				t.Fatalf("status %d: %s", code, body)
			}
			if calls != tt.wantCalls {
				t.Errorf("handler called %d times, want %d", calls, tt.wantCalls)
			}
			rs := decodeRS(t, body)
			if tt.wantErrs == nil {
				if rs.Success == nil || rs.Errors != nil {
					t.Errorf("want success, got %s", body)
				}
				return
			}
			if rs.Success != nil || rs.Warnings != nil || rs.Errors == nil {
				t.Fatalf("want only Errors, got %s", body)
			}
			var msgs []string
			for _, e := range *rs.Errors {
				msgs = append(msgs, e.Message)
			}
			if !equal(msgs, tt.wantErrs) {
				t.Errorf("errors %q, want %q", msgs, tt.wantErrs)
			}
		})
	}
}

func TestBodyLimit(t *testing.T) {
	v := NewVersion("2018-10", nopValidator{})
	a := NewAction[testRQ, testRS](v, "OTA_Test:Test", "action_OTA_Test", nil)
	handler := func(context.Context, Request[testRQ]) (testRS, error) { return testRS{Response: OK()}, nil }

	small := newRequest("2018-10", "c", "OTA_Test:Test", `<TestRQ/>`)
	limit := small.ContentLength
	r := NewRouter(WithMaxFormSize(limit))
	r.Version(v, Handle(a, handler))

	if code, body := serve(r, small); code != http.StatusOK {
		t.Errorf("at the limit: status %d: %s", code, body)
	}
	padded := `<TestRQ>` + strings.Repeat(" ", int(limit)) + `</TestRQ>`
	if code, body := post(t, r, "2018-10", "c", "OTA_Test:Test", padded); code != http.StatusRequestEntityTooLarge || body != "ERROR:request too large" {
		t.Errorf("over the limit: %d %q", code, body)
	}
}

func TestHandlerGetsRequestContext(t *testing.T) {
	type key struct{}
	v := NewVersion("2018-10", nopValidator{})
	a := NewAction[testRQ, testRS](v, "OTA_Test:Test", "action_OTA_Test", nil)
	var got any
	r := NewRouter()
	r.Version(v, Handle(a, func(ctx context.Context, _ Request[testRQ]) (testRS, error) {
		got = ctx.Value(key{})
		return testRS{Response: OK()}, nil
	}))

	req := newRequest("2018-10", "c", "OTA_Test:Test", `<TestRQ/>`)
	if code, body := serve(r, req.WithContext(context.WithValue(req.Context(), key{}, "v"))); code != http.StatusOK {
		t.Fatalf("status %d: %s", code, body)
	}
	if got != "v" {
		t.Errorf("handler context value %v, want v", got)
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// failValidator rejects documents with an <Invalid> element.
type failValidator struct{}

func (failValidator) Validate(doc []byte) error {
	if bytes.Contains(doc, []byte("<Invalid")) {
		return errors.Join(errors.New("first violation"), errors.New("second violation"))
	}
	return nil
}

func TestServeHTTPAnswers(t *testing.T) {
	var log bytes.Buffer
	defer slog.SetDefault(slog.Default())
	slog.SetDefault(slog.New(slog.NewTextHandler(&log, nil)))

	v := NewVersion("2018-10", failValidator{})
	r := NewRouter()
	r.Version(v,
		Handle(NewAction[testRQ, testRS](v, "OTA_Test:Test", "action_OTA_Test", nil),
			func(context.Context, Request[testRQ]) (testRS, error) { return testRS{Response: OK()}, nil }),
		Handle(NewAction[testRQ, testRS](v, "OTA_Fail:Fail", "action_OTA_Fail", nil),
			func(context.Context, Request[testRQ]) (testRS, error) { return testRS{}, errors.New("database down") }),
	)
	gzipped := newRequest("2018-10", "c", "OTA_Test:Test", `<TestRQ/>`)
	gzipped.Header.Set("Content-Encoding", "gzip")

	tests := []struct {
		name       string
		req        *http.Request
		wantStatus int
		wantBody   string
	}{
		{"GET", httptest.NewRequest(http.MethodGet, "/", nil), 405, "ERROR:method must be POST"},
		{"compressed", gzipped, 400, "ERROR:unsupported GZIP compression"},
		{"no action", newRequest("2018-10", "c", "", `<TestRQ/>`), 400, "ERROR:unknown or missing action"},
		{"unknown action", newRequest("2018-10", "c", "OTA_Other:Other", `<TestRQ/>`), 400, "ERROR:unknown or missing action"},
		{"unsupported version", newRequest("2099-10", "c", "OTA_Test:Test", `<TestRQ/>`), 400,
			"ERROR:your current alpinebits version does not match one of the servers supported versions"},
		{"no version", newRequest("", "c", "OTA_Test:Test", `<TestRQ/>`), 400, "ERROR:no valid client protocol version provided"},
		{"no request", newRequest("2018-10", "c", "OTA_Test:Test", ""), 400, "ERROR:missing request"},
		{"schema violations", newRequest("2018-10", "c", "OTA_Test:Test", `<Invalid/>`), 400,
			"ERROR:XML validation error: first violation\nsecond violation"},
		{"unparsable", newRequest("2018-10", "c", "OTA_Test:Test", `<TestRQ Bad="maybe"/>`), 400,
			`ERROR:XML validation error: strconv.ParseBool: parsing "maybe": invalid syntax`},
		{"handler error", newRequest("2018-10", "c", "OTA_Fail:Fail", `<TestRQ/>`), 500, "ERROR:internal server error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if code, body := serve(r, tt.req); code != tt.wantStatus || body != tt.wantBody {
				t.Errorf("got %d %q\nwant %d %q", code, body, tt.wantStatus, tt.wantBody)
			}
		})
	}
	if !strings.Contains(log.String(), "database down") {
		t.Errorf("handler error not logged: %s", log.String())
	}
}

func TestUnauthorized(t *testing.T) {
	rec := httptest.NewRecorder()
	Unauthorized(rec, "invalid credentials")
	if rec.Code != http.StatusUnauthorized || rec.Body.String() != "ERROR:invalid credentials" {
		t.Errorf("got %d %q", rec.Code, rec.Body.String())
	}
}

func TestRequestUserAndClientID(t *testing.T) {
	v := NewVersion("2018-10", nopValidator{})
	a := NewAction[testRQ, testRS](v, "OTA_Test:Test", "action_OTA_Test", nil)
	var got Request[testRQ]
	r := NewRouter()
	r.Version(v, Handle(a, func(_ context.Context, r Request[testRQ]) (testRS, error) {
		got = r
		return testRS{Response: OK()}, nil
	}))
	req := newRequest("2018-10", "user", "OTA_Test:Test", `<TestRQ/>`)
	req.Header.Set(HeaderClientID, "pms 1.0")
	serve(r, req)
	if got.User != "user" || got.ClientID != "pms 1.0" {
		t.Errorf("User %q, ClientID %q", got.User, got.ClientID)
	}
}

func TestMethodNotAllowedSetsAllow(t *testing.T) {
	rec := httptest.NewRecorder()
	NewRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusMethodNotAllowed || rec.Header().Get("Allow") != http.MethodPost {
		t.Errorf("got %d, Allow %q", rec.Code, rec.Header().Get("Allow"))
	}
}
