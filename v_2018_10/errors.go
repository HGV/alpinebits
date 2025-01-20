package v_2018_10

import "fmt"

var (
	ErrMissingHotelCode                        = newMissingAttributeError("HotelCode")
	ErrDeltasNotSupported                      = newError("deltas not supported")
	ErrMissingInvTypeCode                      = newMissingAttributeError("InvTypeCode")
	ErrMissingInvCode                          = newMissingAttributeError("InvCode")
	ErrBookingThresholdNotSupported            = newError("room status free but not bookable (booking threshold) not supported")
	ErrBookingThresholdGreaterThanBookingLimit = newError("attribute BookingThreshold must be ≤ attribute BookingLimit")
	ErrMissingCode                             = newMissingAttributeError("Code")
	ErrChildOccupancyNotSupported              = newError("child occupancy not supported")
	ErrMaxChildOccGreaterThanMaxOcc            = newError("child occupancy must be ≤ max occupancy")
	ErrStdOccLowerThanMinOcc                   = newError("standard occupancy must be ≥ min occupancy")
	ErrMaxOccLowerThanStdOcc                   = newError("max occupancy must be ≥ standard occupancy")
	ErrMissingLongName                         = newMissingElementError("MultimediaDescription with attribute InfoCode = 25 (Long name)")
	ErrDuplicateLanguage                       = newError("duplicate language found for element Description")
	ErrMissingRoomID                           = newMissingAttributeError("RoomID")
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

func ErrInvalidRoomClassificationCode(roomClassificationCode int) *Error {
	return newErrorf("invalid value for attribute RoomClassificationCode %d", roomClassificationCode)
}

func ErrInvalidRoomType(roomType int) *Error {
	return newErrorf("invalid value for attribute RoomType %d", roomType)
}

func ErrInvalidRoomAmenityType(code int) *Error {
	return newErrorf("invalid value for attribute RoomAmenityCode %d", code)
}

func ErrInvalidPictureCategoryCode(code int) *Error {
	return newErrorf("invalid value for attribute Category %d", code)
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
