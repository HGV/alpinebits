package alpinebits

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type nopValidator struct{}

func (nopValidator) Validate([]byte) error { return nil }

type testRQ struct {
	XMLName xml.Name `xml:"TestRQ"`
	Bad     bool     `xml:"Bad,attr"`
}

type testRS struct {
	XMLName xml.Name `xml:"TestRS"`
	Response
}

func testHandler(context.Context, Request[testRQ]) (testRS, error) { return testRS{}, nil }

func TestHandleRejectsActionOfOtherVersion(t *testing.T) {
	v1 := NewVersion("2020-10", nopValidator{})
	v2 := NewVersion("2022-10", nopValidator{})
	a := NewAction[testRQ, testRS](v1, "OTA_Test:Test", "action_OTA_Test", nil)

	NewRouter().Version(v1, Handle(a, testHandler)) // same version: no panic

	defer func() {
		msg, _ := recover().(string)
		if !strings.Contains(msg, "OTA_Test:Test") || !strings.Contains(msg, "2020-10") || !strings.Contains(msg, "2022-10") {
			t.Errorf("got panic %q, want one naming the action and both versions", msg)
		}
	}()
	NewRouter().Version(v2, Handle(a, testHandler))
}

func TestSendRejectsActionOfOtherVersion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Error("request was sent")
	}))
	defer srv.Close()

	v1 := NewVersion("2020-10", nopValidator{})
	v2 := NewVersion("2022-10", nopValidator{})
	a := NewAction[testRQ, testRS](v1, "OTA_Test:Test", "action_OTA_Test", nil)
	c := NewClient(srv.URL, "user", "pass", v2)

	if _, err := Send(context.Background(), c, a, testRQ{}); err == nil || !strings.Contains(err.Error(), "2022-10") {
		t.Errorf("got %v, want version mismatch error", err)
	}
}

func TestNewActionNilVersion(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewAction with nil Version did not panic")
		}
	}()
	NewAction[testRQ, testRS](nil, "OTA_Test:Test", "action_OTA_Test", nil)
}
