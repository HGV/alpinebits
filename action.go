package alpinebits

import "encoding/xml"

// action is the non-generic interface for internal storage.
type action interface {
	Name() string
	HandshakeName() string
	Unmarshal(b []byte) (any, error)
}

// Action is a typed action that carries its request/response types.
// This ensures handlers have the correct signature at compile time.
type Action[RQ, RS any] struct {
	name          string
	handshakeName string
}

// NewAction creates a new typed action.
func NewAction[RQ, RS any](name, handshakeName string) Action[RQ, RS] {
	return Action[RQ, RS]{name: name, handshakeName: handshakeName}
}

// Name returns the action name (e.g., "OTA_HotelAvailNotif:FreeRooms").
func (a Action[RQ, RS]) Name() string {
	return a.name
}

// HandshakeName returns the handshake capability name.
func (a Action[RQ, RS]) HandshakeName() string {
	return a.handshakeName
}

// Unmarshal parses XML into the request type.
func (a Action[RQ, RS]) Unmarshal(b []byte) (any, error) {
	var rq RQ
	if err := xml.Unmarshal(b, &rq); err != nil {
		return nil, err
	}
	return &rq, nil
}
