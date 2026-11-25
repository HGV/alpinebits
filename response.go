package alpinebits

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// Success is an empty element indicating successful processing.
type Success struct{}

// Status represents AlpineBits warning status codes.
type Status string

const (
	StatusHandshake     Status = "ALPINEBITS_HANDSHAKE"
	StatusSendInventory Status = "ALPINEBITS_SEND_INVENTORY"
	StatusSendFreeRooms Status = "ALPINEBITS_SEND_FREEROOMS"
	StatusSendRatePlans Status = "ALPINEBITS_SEND_RATEPLANS"
	StatusSendHandshake Status = "ALPINEBITS_SEND_HANDSHAKE" // since 2024-10
)

// ErrType represents OTA error types (EWT list). The AlpineBits XSD allows
// only these two.
type ErrType int

const (
	ErrTypeAdvisory         ErrType = 11
	ErrTypeApplicationError ErrType = 13
)

func (t ErrType) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: strconv.Itoa(int(t))}, nil
}

// ErrCode represents OTA error codes (ERR list).
type ErrCode int

const (
	ErrCodeInvalidValue     ErrCode = 320
	ErrCodeRequiredField    ErrCode = 321
	ErrCodeInvalidFormat    ErrCode = 322
	ErrCodeInvalidDateRange ErrCode = 323
)

func (c ErrCode) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: strconv.Itoa(int(c))}, nil
}

// Warning represents a business warning in the response.
type Warning struct {
	Type    ErrType `xml:"Type,attr,omitempty"`
	Code    ErrCode `xml:"Code,attr,omitempty"`
	Status  Status  `xml:"Status,attr,omitempty"`
	Message string  `xml:",chardata"`
}

// Warnings is a slice of warnings with helper methods.
type Warnings []Warning

// HasStatus checks if any warning contains the given status.
func (w *Warnings) HasStatus(status Status) bool {
	if w == nil {
		return false
	}
	for _, warn := range *w {
		if warn.Status == status {
			return true
		}
	}
	return false
}

// Error represents a business error in the response.
type Error struct {
	Type    ErrType `xml:"Type,attr,omitempty"`
	Code    ErrCode `xml:"Code,attr,omitempty"`
	Status  Status  `xml:"Status,attr,omitempty"`
	Message string  `xml:",chardata"`
}

// Errors is a slice of errors with helper methods.
type Errors []Error

// HasStatus checks if any error contains the given status.
func (e *Errors) HasStatus(status Status) bool {
	if e == nil {
		return false
	}
	for _, err := range *e {
		if err.Status == status {
			return true
		}
	}
	return false
}

// OTAVersion is the OTA specification version used in responses.
const OTAVersion = "1.0"

// NewWarning creates an advisory warning (Type 11).
func NewWarning(msg string) Warning {
	return Warning{Type: ErrTypeAdvisory, Message: msg}
}

// Response is embedded by every response type. Handlers return OK; error
// responses are built by the router, as the XSD forbids Errors next to
// Success or Warnings.
type Response struct {
	Version  string    `xml:"Version,attr"`
	Success  *Success  `xml:"Success"`
	Warnings *Warnings `xml:"Warnings>Warning"`
	Errors   *Errors   `xml:"Errors>Error"`
}

// OK returns a successful Response.
func OK(warnings ...Warning) Response {
	r := Response{Version: OTAVersion, Success: &Success{}}
	if len(warnings) > 0 {
		w := Warnings(warnings)
		r.Warnings = &w
	}
	return r
}

// setErrors is required of response types by NewAction.
func (r *Response) setErrors(errs Errors) {
	*r = Response{Version: OTAVersion, Errors: &errs}
}

// Rejection is the error returned by Reject.
type Rejection struct {
	Errors Errors
}

func (r *Rejection) Error() string {
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = e.Message
	}
	return "alpinebits: request rejected: " + strings.Join(msgs, "; ")
}

// Reject makes the router answer with errs. It returns nil if errs is empty,
// so handlers can always end with `return rs, alpinebits.Reject(errs...)`.
func Reject(errs ...Error) error {
	if len(errs) == 0 {
		return nil
	}
	return &Rejection{Errors: errs}
}
