package alpinebits

import (
	"encoding/json"
	"encoding/xml"
	"maps"
	"slices"
)

// PingRQ is the OTA_Ping request.
type PingRQ struct {
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRQ"`
	Version  string   `xml:"Version,attr"`
	EchoData EchoData `xml:"EchoData"`
}

// EchoData contains capability negotiation JSON.
type EchoData struct {
	Message string `xml:",chardata"`
}

// PingRS is the OTA_Ping response.
type PingRS struct {
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRS"`
	Version  string   `xml:"Version,attr"`
	Success  Success  `xml:"Success"`
	Warnings Warning  `xml:"Warnings>Warning"`
	EchoData EchoData `xml:"EchoData"`
}

// ActionCapabilities maps handshake action names (e.g.
// "action_OTA_HotelAvailNotif") to capabilities.
type ActionCapabilities map[string][]Capability

// ClientCapabilities returns what the client announced in EchoData, keyed by
// version, or nil if EchoData is not handshake JSON.
func (r *PingRQ) ClientCapabilities() map[string]ActionCapabilities {
	caps, err := parseCapabilities(r.EchoData.Message)
	if err != nil {
		return nil
	}
	return caps
}

// BuildResponse creates a PingRS announcing negotiated, keyed by version.
func (r *PingRQ) BuildResponse(negotiated map[string]ActionCapabilities) PingRS {
	return PingRS{
		Version: r.Version,
		Warnings: Warning{
			Type:    ErrTypeAdvisory,
			Status:  StatusHandshake,
			Message: encodeCapabilities(negotiated),
		},
		EchoData: r.EchoData,
	}
}

type handshakeJSON struct {
	Versions []handshakeVersionJSON `json:"versions"`
}

type handshakeVersionJSON struct {
	Version string                `json:"version"`
	Actions []handshakeActionJSON `json:"actions"`
}

type handshakeActionJSON struct {
	Action   string   `json:"action"`
	Supports []string `json:"supports,omitempty"`
}

func parseCapabilities(echoJSON string) (map[string]ActionCapabilities, error) {
	result := make(map[string]ActionCapabilities)
	if echoJSON == "" {
		return result, nil
	}

	var data handshakeJSON
	if err := json.Unmarshal([]byte(echoJSON), &data); err != nil {
		return nil, err
	}
	for _, v := range data.Versions {
		actions := make(ActionCapabilities, len(v.Actions))
		for _, a := range v.Actions {
			caps := make([]Capability, len(a.Supports))
			for i, s := range a.Supports {
				caps[i] = Capability(s)
			}
			actions[a.Action] = caps
		}
		result[v.Version] = actions
	}
	return result, nil
}

// A version without actions is encoded as an empty list (AlpineBits 2024-10
// §3.2).
func encodeCapabilities(caps map[string]ActionCapabilities) string {
	data := handshakeJSON{Versions: []handshakeVersionJSON{}}
	for _, version := range slices.Backward(slices.Sorted(maps.Keys(caps))) {
		v := handshakeVersionJSON{Version: version, Actions: []handshakeActionJSON{}}
		for _, action := range slices.Sorted(maps.Keys(caps[version])) {
			a := handshakeActionJSON{Action: action}
			for _, c := range caps[version][action] {
				a.Supports = append(a.Supports, string(c))
			}
			v.Actions = append(v.Actions, a)
		}
		data.Versions = append(data.Versions, v)
	}
	b, _ := json.Marshal(data)
	return string(b)
}

// handshakeAction is shared by all versions, so it has none.
var handshakeAction = Action[PingRQ, PingRS]{
	name:          "OTA_Ping:Handshaking",
	handshakeName: "action_OTA_Ping",
}

// Intersect returns the actions and capabilities both sides support. Actions
// without common capabilities are kept.
func Intersect(server, client ActionCapabilities) ActionCapabilities {
	result := make(ActionCapabilities)
	for action, serverCaps := range server {
		if clientCaps, ok := client[action]; ok {
			result[action] = intersectCaps(serverCaps, clientCaps)
		}
	}
	return result
}

func intersectCaps(a, b []Capability) []Capability {
	set := make(map[Capability]struct{}, len(b))
	for _, cap := range b {
		set[cap] = struct{}{}
	}

	var result []Capability
	for _, cap := range a {
		if _, ok := set[cap]; ok {
			result = append(result, cap)
		}
	}
	return result
}
