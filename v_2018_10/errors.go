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
	ErrMissingID                               = newMissingAttributeError("UniqueID.ID")
	ErrMissingRoomStay                         = newMissingElementError("RoomStay")
	ErrMissingRoomTypeCode                     = newMissingAttributeError("RoomTypeCode")
	ErrMissingRatePlanCode                     = newMissingAttributeError("RatePlanCode")
	ErrInvalidPercent                          = newError("percent must be ≤ 100")
	ErrMissingMealsIncluded                    = newMissingElementError("MealsIncluded")
	ErrMissingGuestCount                       = newMissingElementError("GuestCount")
	ErrDuplicateAdultGuestCount                = newError("duplicate element GuestCount for adults")
	ErrMissingStart                            = newMissingAttributeError("Start")
	ErrMissingEnd                              = newMissingAttributeError("End")
	ErrMissingTotal                            = newMissingElementError("Total")
	ErrStartAfterEnd                           = newError("start must be ≤ end")
	ErrMissingDuration                         = newMissingAttributeError("Duration")
	ErrMissingStartDateWindow                  = newMissingElementError("StartDateWindow")
	ErrEarliestDateAfterLatestDate             = newError("earliest date must be ≤ latest date")
	ErrDurationOutOfRange                      = newError("duration exceeds the allowed date range")
	ErrInvalidNamePrefix                       = newError("invalid value for attribute NamePrefix")
	ErrMissingGivenName                        = newMissingAttributeError("GivenName")
	ErrMissingSurname                          = newMissingAttributeError("Surname")
	ErrInvalidNameTitle                        = newError("invalid value for attribute NameTitle")
	ErrInvalidAddressLine                      = newError("invalid value for attribute AddressLine")
	ErrInvalidCityName                         = newError("invalid value for attribute CityName")
	ErrInvalidPostalCode                       = newError("invalid value for attribute PostalCode")
	ErrInvalidCountryNameCode                  = newError("invalid value for attribute CountryName.Code")
	ErrInvalidListItem                         = newError("invalid value for element ListItem")
	ErrInvalidCommentText                      = newError("invalid value for element Comment.Text")
	ErrInvalidPenaltyDescriptionText           = newError("invalid value for attribute element PenaltyDescription.Text")
	ErrInvalidResIDValue                       = newError("invalid value for attribute ResIDValue")
	ErrInvalidResIDSource                      = newError("invalid value for attribute ResIDSource")
	ErrInvalidResIDSourceContext               = newError("invalid value for attribute ResIDSourceContext")
	ErrInvalidCompanyNameCode                  = newError("invalid value for attribute CompanyName.Code")
	ErrInvalidCompanyNameValue                 = newError("invalid value for element CompanyName")
	ErrInvalidEmail                            = newError("invalid value for element Email")
	ErrMissingCurrencyCode                     = newMissingAttributeError("CurrencyCode")
	ErrRatePlanJoinNotSupported                = newError("rate plan join not supported")
	ErrMissingOfferRule                        = newMissingElementError("OfferRule")
	ErrStayThroughNotAllowedInOfferRule        = newError("invalid value for attribute MinMaxMessageType inside element OfferRule")
	ErrMissingAdultOccupancy                   = newMissingElementError("Occupancy with attribute AgeQualifyingCode = 10")
	ErrInvalidMinOccupancy                     = newError("min occupancy must be ≤ 99")
	ErrInvalidMaxOccupancy                     = newError("max occupancy must be ≤ 99")
	ErrDuplicateChildOccupancy                 = newError("duplicate element Occupancy with attribute AgeQualifyingCode = 8")
	ErrDuplicateFreeNightOffer                 = newError("duplicate free night offer")
	ErrDuplicateFamilyOffer                    = newError("duplicate family offer")
	ErrMissingNightsRequired                   = newMissingAttributeError("NightsRequired")
	ErrMissingNightsDiscounted                 = newMissingAttributeError("NightsDiscounted")
	ErrInvalidDiscountPattern                  = newError("invalid value for attribute DiscountPattern")
	ErrInvalidGuestAgeQualifyngCode            = newError("invalid value for attribute Guest.AgeQualifyingCode")
	ErrRoomTypeBookingRulesNotSupported        = newError("room type booking rules not supported")
	ErrMissingStaticRate                       = newMissingElementError("static Rate")
	ErrMissingBaseByGuestAmt                   = newMissingElementError("BaseByGuestAmt")
	ErrMissingNumberOfGuests                   = newMissingAttributeError("NumberOfGuests")
	ErrMissingAgeQualifyingCode                = newMissingAttributeError("AgeQualifyingCode")
	ErrMissingAmountAfterTax                   = newMissingAttributeError("AmountAfterTax")
	ErrMissingAmount                           = newMissingAttributeError("Amount")
	ErrDuplicateAdditionalGuestAmountAdult     = newError("duplicate element AdditionalGuestAmount with attribute AgeQualifyingCode = 10")
	ErrChildrenNotAllowed                      = newError("children not allowed")
	ErrMissingMinAge                           = newMissingAttributeError("MinAge")
	ErrMinAgeGreaterThanOrEqualsThanMaxAge     = newError("attribute MinAge must be < attribute MaxAge")
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

func ErrInvalidUniqueID(status string, uidType int) *Error {
	return newErrorf("invalid value for attributes ResStatus %s and Type %d", status, uidType)
}

func ErrMinStayArrivalGratherThanMaxStayArrival(min, max int) *Error {
	return newErrorf("min stay arrival must be ≤ max stay arrival, got %d and %d", min, max)
}

func ErrMinStayGratherThanMaxStay(min, max int) *Error {
	return newErrorf("min stay must be ≤ max stay, got %d and %d", min, max)
}

func ErrDuplicateBaseByGuestAmt(numberOfGuests int) *Error {
	return newErrorf("duplicate element BaseByGuestAmt with attribute NumberOfGuests %d", numberOfGuests)
}

func ErrMinAgeOutOfRange(childMinAge, ratePlanChildMinAge int) *Error {
	return newErrorf("child min age must be ≥ rate plan child min age, got %d and %d", childMinAge, ratePlanChildMinAge)
}

func ErrMaxAgeOutOfRange(childMaxAge, ratePlanAdultMinAge int) *Error {
	return newErrorf("child max age must be < rate plan adult min age, got %d and %d", childMaxAge, ratePlanAdultMinAge)
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
