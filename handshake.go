package alpinebits

import (
	"encoding/json"
	"encoding/xml"
)

// PingRQ is the OTA_Ping request.
type PingRQ struct {
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRQ"`
	Version  string   `xml:"Version,attr"`
	EchoData EchoData `xml:"EchoData"`
}

// EchoData contains capability negotiation JSON.
type EchoData struct {
	Message string `xml:",innerxml"`
}

// PingRS is the OTA_Ping response.
type PingRS struct {
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRS"`
	Version  string   `xml:"Version,attr"`
	Success  Success  `xml:"Success"`
	Warnings Warning  `xml:"Warnings>Warning"`
	EchoData EchoData `xml:"EchoData"`
}

// ActionCapabilities maps action names to their negotiated capabilities.
type ActionCapabilities map[string][]Capability

// ClientCapabilities parses and returns the client's declared capabilities from EchoData.
func (r *PingRQ) ClientCapabilities() ActionCapabilities {
	caps, err := parseCapabilities(r.EchoData.Message)
	if err != nil {
		return nil
	}
	return caps
}

// BuildResponse creates a PingRS with the negotiated capabilities.
func (r *PingRQ) BuildResponse(negotiated ActionCapabilities) PingRS {
	echoData := EchoData{
		Message: encodeCapabilities(negotiated),
	}
	return PingRS{
		Version: r.Version,
		Success: Success{},
		Warnings: Warning{
			Type:    ErrTypeAdvisory,
			Status:  StatusHandshake,
			Message: echoData.Message,
		},
		EchoData: r.EchoData,
	}
}

// capabilitiesJSON represents the JSON structure in EchoData.
type capabilitiesJSON struct {
	Versions []versionCapabilities `json:"versions"`
}

type versionCapabilities struct {
	Version string             `json:"version"`
	Actions []actionCapability `json:"actions"`
}

type actionCapability struct {
	Action   string   `json:"action"`
	Supports []string `json:"supports,omitempty"`
}

// parseCapabilities parses ActionCapabilities from EchoData JSON.
func parseCapabilities(echoJSON string) (ActionCapabilities, error) {
	if echoJSON == "" {
		return make(ActionCapabilities), nil
	}

	var caps capabilitiesJSON
	if err := json.Unmarshal([]byte(echoJSON), &caps); err != nil {
		return nil, err
	}

	result := make(ActionCapabilities)
	for _, ver := range caps.Versions {
		for _, act := range ver.Actions {
			var capabilities []Capability
			for _, sup := range act.Supports {
				capabilities = append(capabilities, Capability(sup))
			}
			result[act.Action] = capabilities
		}
	}

	return result, nil
}

// encodeCapabilities encodes ActionCapabilities to EchoData JSON.
func encodeCapabilities(caps ActionCapabilities) string {
	if len(caps) == 0 {
		return ""
	}

	var actions []actionCapability
	for action, supports := range caps {
		var supportsStrs []string
		for _, cap := range supports {
			supportsStrs = append(supportsStrs, string(cap))
		}
		actions = append(actions, actionCapability{
			Action:   action,
			Supports: supportsStrs,
		})
	}

	// Note: Version should come from the request context
	capsJSON := capabilitiesJSON{
		Versions: []versionCapabilities{{
			Version: "2018-10", // TODO: Get from context
			Actions: actions,
		}},
	}

	data, _ := json.Marshal(capsJSON)
	return string(data)
}

// handshakeAction is the internal action for OTA_Ping.
var handshakeAction = NewAction[PingRQ, PingRS](
	"OTA_Ping:Handshaking",
	"action_OTA_Ping",
)

// Intersect calculates the intersection of server and client capabilities.
// Returns only actions supported by both, with capabilities both support.
func Intersect(server, client ActionCapabilities) ActionCapabilities {
	result := make(ActionCapabilities)
	for action, serverCaps := range server {
		if clientCaps, ok := client[action]; ok {
			intersection := intersectCaps(serverCaps, clientCaps)
			if len(intersection) > 0 || len(serverCaps) == 0 {
				// Include action if caps intersect, or if server has no required caps
				result[action] = intersection
			}
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
