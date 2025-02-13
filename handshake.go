package alpinebits

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

func NewHandshakeDataFromRouter(r Router) HandshakeData {
	handshakeData := make(HandshakeData)
	for _, version := range r.versionRoutes {
		actions := make(map[string][]string)
		for _, action := range version.actionRoutes {
			if action.excludeFromHandshake {
				continue
			}
			actions[action.action.HandshakeName()] = action.capabilities
		}
		handshakeData[version.version.String()] = actions
	}
	return handshakeData
}

func (h HandshakeData) NegotiatedVersion() (string, map[string][]string) {
	sortedKeys := slices.SortedFunc(maps.Keys(h), compareVersionsDescending)
	if len(sortedKeys) > 0 {
		version := sortedKeys[0]
		return version, h[version]
	}
	return "", nil
}

type (
	handshakeData struct {
		Versions []handshakeVersion `json:"versions"`
	}
	handshakeVersion struct {
		Version string            `json:"version"`
		Actions []handshakeAction `json:"actions,omitempty"`
	}
	handshakeAction struct {
		Action       string   `json:"action"`
		Capabilities []string `json:"supports,omitempty"`
	}
)

func (h HandshakeData) MarshalJSON() ([]byte, error) {
	var versions []handshakeVersion
	for _, versionKey := range slices.SortedFunc(maps.Keys(h), compareVersionsDescending) {
		var actions []handshakeAction
		for _, actionKey := range slices.Sorted(maps.Keys(h[versionKey])) {
			actions = append(actions, handshakeAction{
				Action:       actionKey,
				Capabilities: h[versionKey][actionKey],
			})
		}
		versions = append(versions, handshakeVersion{
			Version: versionKey,
			Actions: actions,
		})
	}
	return json.Marshal(handshakeData{Versions: versions})
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
