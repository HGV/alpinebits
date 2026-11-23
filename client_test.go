package alpinebits

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendServerError(t *testing.T) {
	v := NewVersion("2018-10", nopValidator{})
	a := NewAction[testRQ, testRS](v, "OTA_Test:Test", "action_OTA_Test", nil)
	tests := []struct {
		status int
		body   string
		want   ServerError
	}{
		{200, "ERROR:unknown or missing action", ServerError{200, "unknown or missing action"}},
		{401, "ERROR:invalid credentials", ServerError{401, "invalid credentials"}},
		{500, "ERROR:internal server error", ServerError{500, "internal server error"}},
	}
	for _, tt := range tests {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(tt.status)
			io.WriteString(w, tt.body)
		}))
		_, err := Send(context.Background(), NewClient(srv.URL, "user", "secret", v), a, testRQ{})
		srv.Close()
		if se, ok := errors.AsType[*ServerError](err); !ok || *se != tt.want {
			t.Errorf("%d %q: got %v, want %+v", tt.status, tt.body, err, tt.want)
		}
	}
}
