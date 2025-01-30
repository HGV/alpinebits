package handshake

import (
	"encoding/json"
	"maps"
	"slices"
	"strings"

	"github.com/HGV/x"
	"github.com/juliangruber/go-intersect/v2"
)

type HandshakeData map[string]map[string][]string

var _ json.Marshaler = (*HandshakeData)(nil)
var _ json.Unmarshaler = (*HandshakeData)(nil)

type (
	handshakeData struct {
		Versions []version `json:"versions"`
	}
	version struct {
		Version string   `json:"version"`
		Actions []action `json:"actions,omitempty"`
	}
	action struct {
		Action       string   `json:"action"`
		Capabilities []string `json:"supports,omitempty"`
	}
)

func (h HandshakeData) MarshalJSON() ([]byte, error) {
	var versions []version
	for _, versionKey := range slices.SortedFunc(maps.Keys(h), compareVersionsDescending) {
		var actions []action
		for _, actionKey := range slices.Sorted(maps.Keys(h[versionKey])) {
			actions = append(actions, action{
				Action:       actionKey,
				Capabilities: h[versionKey][actionKey],
			})
		}
		versions = append(versions, version{
			Version: versionKey,
			Actions: actions,
		})
	}
	data := handshakeData{
		Versions: versions,
	}
	return json.Marshal(data)
}

func compareVersionsDescending(v1, v2 string) int {
	return strings.Compare(v2, v1)
}

func (h *HandshakeData) UnmarshalJSON(data []byte) error {
	var handshakeData handshakeData
	if err := json.Unmarshal(data, &handshakeData); err != nil {
		return err
	}

	*h = make(HandshakeData)
	for _, version := range handshakeData.Versions {
		actions := make(map[string][]string)
		for _, action := range version.Actions {
			actions[action.Action] = action.Capabilities
		}
		(*h)[version.Version] = actions
	}
	return nil
}

func (h HandshakeData) Intersect(other HandshakeData) HandshakeData {
	result := make(HandshakeData)
	for version, actions := range h {
		if otherActions, ok := other[version]; ok {
			result[version] = make(map[string][]string)
			for action, capabilities := range actions {
				if otherCapabilities, ok := otherActions[action]; ok {
					intersectedCapabilities := intersect.SimpleGeneric(
						capabilities,
						otherCapabilities,
					)
					result[version][action] = x.If(
						len(intersectedCapabilities) > 0,
						intersectedCapabilities,
						nil,
					)
				}
			}
		}
	}
	return result
}
