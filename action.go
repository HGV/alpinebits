package alpinebits

import "encoding/xml"

// action is the non-generic interface for internal storage.
type action interface {
	Name() string
	HandshakeName() string
	Version() *Version
	unmarshal(b []byte) (any, error)
}

// Action is a typed action that carries its request/response types.
// This ensures handlers have the correct signature at compile time.
type Action[RQ, RS any] struct {
	version       *Version
	name          string
	handshakeName string
	validate      func(RQ, []Capability) Errors
	reject        func(Errors) RS
}

// NewAction creates an action of version v. validate checks the rules the XSD
// cannot express and may be nil; the router runs it before the handler and
// answers with its errors. RS must embed Response.
func NewAction[RQ, RS any, P interface {
	*RS
	setErrors(Errors)
}](v *Version, name, handshakeName string, validate func(RQ, []Capability) Errors) Action[RQ, RS] {
	if v == nil {
		panic("alpinebits: NewAction(" + name + "): nil Version")
	}
	return Action[RQ, RS]{
		version:       v,
		name:          name,
		handshakeName: handshakeName,
		validate:      validate,
		reject: func(errs Errors) RS {
			var rs RS
			P(&rs).setErrors(errs)
			return rs
		},
	}
}

// Version returns the version the action belongs to.
func (a Action[RQ, RS]) Version() *Version {
	return a.version
}

// Name returns the action name (e.g., "OTA_HotelAvailNotif:FreeRooms").
func (a Action[RQ, RS]) Name() string {
	return a.name
}

// HandshakeName returns the handshake capability name.
func (a Action[RQ, RS]) HandshakeName() string {
	return a.handshakeName
}

func (a Action[RQ, RS]) unmarshal(b []byte) (any, error) {
	var rq RQ
	if err := xml.Unmarshal(b, &rq); err != nil {
		return nil, err
	}
	return &rq, nil
}
