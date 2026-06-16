package alpinebits

import (
	"encoding/xml"
	"fmt"
	"strconv"
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
)

// ErrType represents OTA error types (EWT list).
type ErrType int

const (
	ErrTypeBusinessRule     ErrType = 3
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

// NewWarning creates a Warning with a message.
func NewWarning(msg string) Warning {
	return Warning{Message: msg}
}

// NewError creates an Error with code and message.
func NewError(code ErrCode, msg string) Error {
	return Error{Code: code, Message: msg}
}

// BusinessError creates a business rule error.
func BusinessError(code ErrCode, msg string) Error {
	return Error{Type: ErrTypeBusinessRule, Code: code, Message: msg}
}

// ApplicationError creates an application error.
func ApplicationError(code ErrCode, msg string) Error {
	return Error{Type: ErrTypeApplicationError, Code: code, Message: msg}
}

// BusinessErrorf creates a business rule error with a formatted message.
func BusinessErrorf(code ErrCode, format string, args ...any) Error {
	return Error{Type: ErrTypeBusinessRule, Code: code, Message: fmt.Sprintf(format, args...)}
}

// ApplicationErrorf creates an application error with a formatted message.
func ApplicationErrorf(code ErrCode, format string, args ...any) Error {
	return Error{Type: ErrTypeApplicationError, Code: code, Message: fmt.Sprintf(format, args...)}
}
