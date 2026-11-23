package alpinebits

import (
	"context"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"testing"
)

func ping(json string) string {
	return `<OTA_PingRQ xmlns="http://www.opentravel.org/OTA/2003/05" Version="8.000"><EchoData>` +
		json + `</EchoData></OTA_PingRQ>`
}

func negotiatedJSON(t *testing.T, body string) string {
	t.Helper()
	var rs PingRS
	if err := xml.Unmarshal([]byte(body), &rs); err != nil {
		t.Fatalf("decode %q: %v", body, err)
	}
	return rs.Warnings.Message
}

func handshakeRouter(t *testing.T) (*Router, *[]Capability) {
	t.Helper()
	v18 := NewVersion("2018-10", nopValidator{})
	v24 := NewVersion("2024-10", nopValidator{})
	got := new([]Capability)
	h := func(_ context.Context, r Request[testRQ]) (testRS, error) {
		*got = r.Capabilities
		return testRS{Response: OK()}, nil
	}
	r := NewRouter(WithClientStore(NewInMemoryClientStore()))
	r.Version(v18,
		Handle(NewAction[testRQ, testRS](v18, "OTA_A:A", "action_A", nil), h, WithCapabilities("rooms", "deltas")),
		Handle(NewAction[testRQ, testRS](v18, "OTA_B:B", "action_B", nil), h),
	)
	r.Version(v24, Handle(NewAction[testRQ, testRS](v24, "OTA_A:A", "action_A", nil), h, WithCapabilities("rooms", "deltas")))
	return r, got
}

func TestHandshakeIntersectsAllVersions(t *testing.T) {
	r, _ := handshakeRouter(t)
	client := `{"versions":[
		{"version":"2099-10","actions":[{"action":"action_A"}]},
		{"version":"2024-10","actions":[{"action":"action_OTA_Ping"},{"action":"action_A","supports":["deltas","other"]}]},
		{"version":"2018-10","actions":[{"action":"action_OTA_Ping"},{"action":"action_A"}]}]}`
	code, body := post(t, r, "2024-10", "c1", "OTA_Ping:Handshaking", ping(client))
	if code != http.StatusOK {
		t.Fatalf("status %d: %s", code, body)
	}
	want := `{"versions":[` +
		`{"version":"2024-10","actions":[{"action":"action_A","supports":["deltas"]}]},` +
		`{"version":"2018-10","actions":[{"action":"action_A"},{"action":"action_OTA_Ping"}]}]}`
	if got := negotiatedJSON(t, body); got != want {
		t.Errorf("negotiated\n got %s\nwant %s", got, want)
	}
}

func TestHandshakeNegotiationIsEnforced(t *testing.T) {
	r, got := handshakeRouter(t)

	post(t, r, "2024-10", "c1", "OTA_A:A", `<TestRQ/>`)
	if want := []Capability{"rooms", "deltas"}; !reflect.DeepEqual(*got, want) {
		t.Errorf("caps before handshake %v, want %v", *got, want)
	}

	client := `{"versions":[{"version":"2024-10","actions":[{"action":"action_A","supports":["deltas"]}]}]}`
	if code, body := post(t, r, "2024-10", "c1", "OTA_Ping:Handshaking", ping(client)); code != http.StatusOK {
		t.Fatalf("handshake: %d %s", code, body)
	}

	tests := []struct {
		name, version, action string
		wantStatus            int
	}{
		{"negotiated action", "2024-10", "OTA_A:A", http.StatusOK},
		{"version the client did not announce", "2018-10", "OTA_A:A", http.StatusBadRequest},
		{"repeated handshake", "2018-10", "OTA_Ping:Handshaking", http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := `<TestRQ/>`
			if tt.action == "OTA_Ping:Handshaking" {
				doc = ping(client)
			}
			code, body := post(t, r, tt.version, "c1", tt.action, doc)
			if code != tt.wantStatus || code != http.StatusOK && body != "ERROR:action not negotiated" {
				t.Errorf("got %d %s, want %d", code, body, tt.wantStatus)
			}
		})
	}
	if want := []Capability{"deltas"}; !reflect.DeepEqual(*got, want) {
		t.Errorf("caps after handshake %v, want %v", *got, want)
	}
}

func TestIntersectKeepsActionsWithoutCommonCapabilities(t *testing.T) {
	server := ActionCapabilities{"action_A": {"rooms", "deltas"}, "action_B": nil}
	client := ActionCapabilities{"action_A": nil, "action_C": nil}
	want := ActionCapabilities{"action_A": nil}
	if got := Intersect(server, client); !reflect.DeepEqual(got, want) {
		t.Errorf("Intersect = %v, want %v", got, want)
	}
}

func TestClientCapabilitiesDecodesXMLText(t *testing.T) {
	want := map[string]ActionCapabilities{"2018-10": {"action_A": {"rooms"}}}
	for _, echo := range []string{
		`{&quot;versions&quot;:[{&quot;version&quot;:&quot;2018-10&quot;,&quot;actions&quot;:[{&quot;action&quot;:&quot;action_A&quot;,&quot;supports&quot;:[&quot;rooms&quot;]}]}]}`,
		`<![CDATA[{"versions":[{"version":"2018-10","actions":[{"action":"action_A","supports":["rooms"]}]}]}]]>`,
	} {
		var rq PingRQ
		if err := xml.Unmarshal([]byte(ping(echo)), &rq); err != nil {
			t.Fatal(err)
		}
		if got := rq.ClientCapabilities(); !reflect.DeepEqual(got, want) {
			t.Errorf("ClientCapabilities(%s) = %v, want %v", echo, got, want)
		}
	}
}

func TestHandshakeWithUnknownVersionHeader(t *testing.T) {
	r, _ := handshakeRouter(t)
	client := `{"versions":[{"version":"2018-10","actions":[{"action":"action_OTA_Ping"},{"action":"action_A"}]}]}`
	want := `{"versions":[{"version":"2018-10","actions":[{"action":"action_A"},{"action":"action_OTA_Ping"}]}]}`
	for _, header := range []string{"2099-10", ""} {
		code, body := post(t, r, header, "c1", "OTA_Ping:Handshaking", ping(client))
		if code != http.StatusOK {
			t.Fatalf("header %q: status %d: %s", header, code, body)
		}
		if got := negotiatedJSON(t, body); got != want {
			t.Errorf("header %q: negotiated %s, want %s", header, got, want)
		}
	}
}

func TestNegotiationIsPerClient(t *testing.T) {
	r, got := handshakeRouter(t)
	client := `{"versions":[{"version":"2024-10","actions":[{"action":"action_A","supports":["deltas"]}]}]}`
	post(t, r, "2024-10", "c1", "OTA_Ping:Handshaking", ping(client))

	post(t, r, "2024-10", "c2", "OTA_A:A", `<TestRQ/>`)
	if want := []Capability{"rooms", "deltas"}; !reflect.DeepEqual(*got, want) {
		t.Errorf("client without handshake: caps %v, want %v", *got, want)
	}
	post(t, r, "2024-10", "c1", "OTA_A:A", `<TestRQ/>`)
	if want := []Capability{"deltas"}; !reflect.DeepEqual(*got, want) {
		t.Errorf("client after handshake: caps %v, want %v", *got, want)
	}
}

type failingStore struct{}

func (failingStore) Load(context.Context, string) (map[string]ActionCapabilities, bool, error) {
	return nil, false, errors.New("database down")
}

func (failingStore) Save(context.Context, string, map[string]ActionCapabilities) error {
	return errors.New("database down")
}

func TestClientStoreErrors(t *testing.T) {
	defer slog.SetDefault(slog.Default())
	slog.SetDefault(slog.New(slog.DiscardHandler))

	v := NewVersion("2018-10", nopValidator{})
	r := NewRouter(WithClientStore(failingStore{}))
	r.Version(v, Handle(NewAction[testRQ, testRS](v, "OTA_A:A", "action_A", nil), testHandler))

	client := `{"versions":[{"version":"2018-10","actions":[{"action":"action_A"}]}]}`
	for _, tt := range []struct{ action, doc string }{
		{"OTA_Ping:Handshaking", ping(client)},
		{"OTA_A:A", `<TestRQ/>`},
	} {
		if code, body := post(t, r, "2018-10", "c1", tt.action, tt.doc); code != http.StatusInternalServerError {
			t.Errorf("%s: got %d %s, want 500", tt.action, code, body)
		}
	}
}
