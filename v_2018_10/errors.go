package v_2018_10

import "fmt"

var (
	ErrMissingHotelCode                        = newMissingAttributeError("HotelCode")
	ErrDeltasNotSupported                      = newError("deltas not supported")
	ErrMissingInvTypeCode                      = newMissingAttributeError("InvTypeCode")
	ErrMissingInvCode                          = newMissingAttributeError("InvCode")
	ErrBookingThresholdNotSupported            = newError("room status free but not bookable (booking threshold) not supported")
	ErrBookingThresholdGreaterThanBookingLimit = newError("attribute BookingThreshold must be ≤ attribute BookingLimit")
)

func ErrInvalidBookingLimit(n int) *Error {
	return newErrorf("attribute BookingLimit must be 0 or 1, got %d", n)
}

func ErrInvCodeNotFound(invCode string) *Error {
	return newErrorf("inv code not found %s", invCode)
}

func ErrInvTypeCodeNotFound(invTypeCode string) *Error {
	return newErrorf("inv type code not found %s", invTypeCode)
}

func newMissingAttributeError(attribute string) *Error {
	return newErrorf("missing required attribute %s", attribute)
}

func newMissingElementError(element string) *Error {
	return newErrorf("missing required element %s", element)
}

func newErrorf(message string, a ...any) *Error {
	return newError(fmt.Sprintf(message, a...))
}

func newError(message string) *Error {
	return &Error{
		Type:  ErrorWarningTypeApplicationError,
		Value: message,
	}
}
