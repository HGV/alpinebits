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

func TestHandleRejectsTwoVersionsOfSameName(t *testing.T) {
	r := NewRouter().Handle(NewAction[testRQ, testRS](NewVersion("2018-10", nopValidator{}), "OTA_A:A", "action_A", nil), testHandler)
	defer func() {
		if msg, _ := recover().(string); msg != "alpinebits: two versions named 2018-10" {
			t.Errorf("got panic %q", msg)
		}
	}()
	r.Handle(NewAction[testRQ, testRS](NewVersion("2018-10", nopValidator{}), "OTA_B:B", "action_B", nil), testHandler)
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

	if _, err := c.Send(context.Background(), a, testRQ{}); err == nil || !strings.Contains(err.Error(), "2022-10") {
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
