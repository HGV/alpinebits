package common

import (
	"fmt"

	"github.com/HGV/x/timex"
)

var (
	ErrMissingHotelCode                    = newMissingAttributeError("HotelCode")
	ErrDeltasNotSupported                  = newError("deltas not supported")
	ErrMissingInvTypeCode                  = newMissingAttributeError("InvTypeCode")
	ErrMissingInvCode                      = newMissingAttributeError("InvCode")
	ErrOutOfOrderNotSupported              = newError("out of order not supported")
	ErrOutOfMarketNotSupported             = newError("out of market not supported")
	ErrClosingSeasonsNotSupported          = newError("closing seasons not supported")
	ErrUnexpectedInvCounts                 = newUnexpectedElementError("InvCounts")
	ErrAvailabilitiesOverlapClosingSeasons = newError("availabilities overlap closing seasons")
	ErrMissingCode                         = newMissingAttributeError("Code")
	ErrChildOccupancyNotSupported          = newError("child occupancy not supported")
	ErrMaxChildOccGreaterThanMaxOcc        = newError("child occupancy must be ≤ max occupancy")
	ErrStdOccLowerThanMinOcc               = newError("standard occupancy must be ≥ min occupancy")
	ErrMaxOccLowerThanStdOcc               = newError("max occupancy must be ≥ standard occupancy")
	ErrMissingMultimediaDescriptions       = newMissingElementError("MultimediaDescriptions")
	ErrMissingLongName                     = newMissingElementError("MultimediaDescription with attribute InfoCode = 25 (Long name)")
	ErrDuplicateLanguage                   = newError("duplicate language found for element Description")
	ErrRoomsNotSupported                   = newError("rooms not supported")
	ErrMissingRoomID                       = newMissingAttributeError("RoomID")
	ErrMissingID                           = newMissingAttributeError("UniqueID.ID")
	ErrMissingRoomStay                     = newMissingElementError("RoomStay")
	ErrDuplicateAlternativeRoomStay        = newError("at most one alternative room stay is allowed")
	ErrUnexpectedAlternativeRoomStay       = newError("alternative room stay is not allowed")
	ErrMissingRoomType                     = newMissingElementError("RoomType")
	ErrUnexpectedRoomType                  = newUnexpectedElementError("RoomType")
	ErrMissingRoomTypeCode                 = newMissingAttributeError("RoomTypeCode")
	ErrMissingRatePlan                     = newMissingElementError("RatePlan")
	ErrUnexpectedRatePlan                  = newUnexpectedElementError("RatePlan")
	ErrMissingRatePlanCode                 = newMissingAttributeError("RatePlanCode")
	ErrInvalidPercent                      = newError("percent must be ≤ 100")
	ErrMissingMealsIncluded                = newMissingElementError("MealsIncluded")
	ErrMissingGuestCount                   = newMissingElementError("GuestCount")
	ErrUnexpectedGuestCounts               = newUnexpectedElementError("GuestCounts")
	ErrDuplicateAdultGuestCount            = newError("duplicate element GuestCount for adults")
	ErrMissingStart                        = newMissingAttributeError("Start")
	ErrMissingEnd                          = newMissingAttributeError("End")
	ErrMissingTotal                        = newMissingElementError("Total")
	ErrUnexpectedTotal                     = newUnexpectedElementError("Total")
	ErrStartAfterEnd                       = newError("start must be ≤ end")
	ErrMissingDuration                     = newMissingAttributeError("Duration")
	ErrUnexpectedStartDateWindow           = newUnexpectedElementError("StartDateWindow")
	ErrUnexpectedDuration                  = newUnexpectedAttributeError("Duration")
	ErrMissingTimeSpan                     = newMissingElementError("TimeSpan")
	ErrMissingStartDateWindow              = newMissingElementError("StartDateWindow")
	ErrEarliestDateAfterLatestDate         = newError("earliest date must be ≤ latest date")
	ErrDurationOutOfRange                  = newError("duration exceeds the allowed date range")
	ErrInvalidNamePrefix                   = newError("invalid value for attribute NamePrefix")
	ErrMissingGivenName                    = newMissingAttributeError("GivenName")
	ErrMissingSurname                      = newMissingAttributeError("Surname")
	ErrInvalidNameTitle                    = newError("invalid value for attribute NameTitle")
	ErrInvalidAddressLine                  = newError("invalid value for attribute AddressLine")
	ErrInvalidCityName                     = newError("invalid value for attribute CityName")
	ErrInvalidPostalCode                   = newError("invalid value for attribute PostalCode")
	ErrInvalidCountryNameCode              = newError("invalid value for attribute CountryName.Code")
	ErrInvalidListItem                     = newError("invalid value for element ListItem")
	ErrInvalidCommentText                  = newError("invalid value for element Comment.Text")
	ErrInvalidPenaltyDescriptionText       = newError("invalid value for attribute element PenaltyDescription.Text")
	ErrInvalidResIDValue                   = newError("invalid value for attribute ResIDValue")
	ErrInvalidResIDSource                  = newError("invalid value for attribute ResIDSource")
	ErrInvalidResIDSourceContext           = newError("invalid value for attribute ResIDSourceContext")
	ErrInvalidCompanyNameCode              = newError("invalid value for attribute CompanyName.Code")
	ErrInvalidCompanyNameValue             = newError("invalid value for element CompanyName")
	ErrInvalidEmail                        = newError("invalid value for element Email")
	ErrMissingCurrencyCode                 = newMissingAttributeError("CurrencyCode")
	ErrRatePlanJoinNotSupported            = newError("rate plan join not supported")
	ErrMissingOfferRule                    = newMissingElementError("OfferRule")
	ErrOfferRuleBookingOffsetNotSupported  = newError("offer rule booking offset not supported")
	ErrOfferRuleDOWLOSNotSupported         = newError("offer rule days of week and lengths of stay not supported")
	ErrStayThroughNotAllowedInOfferRule    = newError("invalid value for attribute MinMaxMessageType inside element OfferRule")
	ErrMissingAdultOccupancy               = newMissingElementError("Occupancy with attribute AgeQualifyingCode = 10")
	ErrInvalidMinOccupancy                 = newError("min occupancy must be ≤ 99")
	ErrInvalidMaxOccupancy                 = newError("max occupancy must be ≤ 99")
	ErrDuplicateChildOccupancy             = newError("duplicate element Occupancy with attribute AgeQualifyingCode = 8")
	ErrDuplicateFreeNightOffer             = newError("duplicate free night offer")
	ErrDuplicateFamilyOffer                = newError("duplicate family offer")
	ErrFreeNightOfferNotSupported          = newError("free night offer not supported")
	ErrMissingNightsRequired               = newMissingAttributeError("NightsRequired")
	ErrMissingNightsDiscounted             = newMissingAttributeError("NightsDiscounted")
	ErrInvalidDiscountPattern              = newError("invalid value for attribute DiscountPattern")
	ErrFamilyOfferNotSupported             = newError("free night offer not supported")
	ErrInvalidGuestAgeQualifyngCode        = newError("invalid value for attribute Guest.AgeQualifyingCode")
	ErrRoomTypeBookingRulesNotSupported    = newError("room type booking rules not supported")
	ErrArrivalDOWNotSupported              = newError("arrival days of week not supported")
	ErrDepartureDOWNotSupported            = newError("departure days of week not supported")
	ErrMissingStaticRate                   = newMissingElementError("static Rate")
	ErrInvalidRateTimeUnit                 = newError("invalid value for attribute RateTimeUnit")
	ErrMissingBaseByGuestAmt               = newMissingElementError("BaseByGuestAmt")
	ErrMissingNumberOfGuests               = newMissingAttributeError("NumberOfGuests")
	ErrMissingAgeQualifyingCode            = newMissingAttributeError("AgeQualifyingCode")
	ErrMissingAmountAfterTax               = newMissingAttributeError("AmountAfterTax")
	ErrMissingAmount                       = newMissingAttributeError("Amount")
	ErrDuplicateAdditionalGuestAmountAdult = newError("duplicate element AdditionalGuestAmount with attribute AgeQualifyingCode = 10")
	ErrChildrenNotAllowed                  = newError("children not allowed")
	ErrMissingMinAge                       = newMissingAttributeError("MinAge")
	ErrMinAgeGreaterThanOrEqualsThanMaxAge = newError("attribute MinAge must be < attribute MaxAge")
	ErrSupplementsNotSupported             = newError("supplements not supported")
	ErrMissingAddToBasicRateIndicator      = newMissingAttributeError("AddToBasicRateIndicator")
	ErrMissingMandatoryIndicator           = newMissingAttributeError("MandatoryIndicator")
	ErrMissingChargeTypeCode               = newMissingAttributeError("ChargeTypeCode")
	ErrInvalidDOWString                    = newError("invalid value for attribute InvCode with attribute InvType = ALPINEBITSDOW")
	ErrUnexpectedOffers                    = newUnexpectedElementError("Offers")
	ErrUnexpectedDescription               = newUnexpectedElementError("Description")
	ErrUnexpectedBookingRules              = newUnexpectedElementError("BookingRules")
	ErrUnexpectedRates                     = newUnexpectedElementError("Rates")
	ErrUnexpectedSupplements               = newUnexpectedElementError("Supplements")
	ErrUnexpectedGuest                     = newUnexpectedElementError("Guest")
	ErrUnexpectedNightsRequired            = newUnexpectedAttributeError("NightsRequired")
	ErrUnexpectedNightsDiscounted          = newUnexpectedAttributeError("NightsDiscounted")
	ErrUnexpectedDiscountPattern           = newUnexpectedAttributeError("DiscountPattern")
	ErrUnexpectedInvTypeCode               = newUnexpectedAttributeError("InvTypeCode")
	ErrUnexpectedStart                     = newUnexpectedAttributeError("Start")
	ErrUnexpectedEnd                       = newUnexpectedAttributeError("End")
	ErrUnexpectedNumberOfGuests            = newUnexpectedAttributeError("NumberOfGuests")
	ErrUnexpectedAgeQualifyingCode         = newUnexpectedAttributeError("AgeQualifyingCode")
	ErrUnexpectedAmountAfterTax            = newUnexpectedAttributeError("AmountAfterTax")
	ErrUnexpectedBaseByGuestAmt            = newError("static rates can contain only one element BaseByGuestAmt")
	ErrUnexpectedAdditionalGuestAmounts    = newUnexpectedElementError("AdditionalGuestAmounts")
	ErrUnexpectedRateTimeUnit              = newUnexpectedAttributeError("RateTimeUnit")
	ErrUnexpectedUnitMultiplier            = newUnexpectedAttributeError("UnitMultiplier")
	ErrUnexpectedMealsIncluded             = newUnexpectedElementError("MealsIncluded")
	ErrUnexpectedType                      = newUnexpectedAttributeError("Type")
	ErrUnexpectedAmount                    = newUnexpectedAttributeError("Amount")
	ErrUnexpectedAddToBasicRateIndicator   = newUnexpectedAttributeError("AddToBasicRateIndicator")
	ErrUnexpectedMandatoryIndicator        = newUnexpectedAttributeError("MandatoryIndicator")
	ErrUnexpectedChargeTypeCode            = newUnexpectedAttributeError("ChargeTypeCode")
)

func ErrInvCodeNotFound(invCode string) *Error {
	return newErrorf("inv code not found %s", invCode)
}

func ErrInvTypeCodeNotFound(invTypeCode string) *Error {
	return newErrorf("inv type code not found %s", invTypeCode)
}

func ErrInvalidInvCounts(n int) *Error {
	return newErrorf("invalid value for element InvCounts, expected one element InvCount, got %d", n)
}

func ErrInvalidCount(n int) *Error {
	return newErrorf("inv count must be 1, got %d", n)
}

func ErrDateRangeOverlaps(range1, range2 timex.DateRange) *Error {
	return newErrorf("date range [%s - %s] overlaps with [%s - %s]", range1.Start, range1.End, range2.Start, range2.End)
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

func ErrRatePlanNotFound(code string) *Error {
	return newErrorf("rate plan not found %s", code)
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

func ErrMissingBaseByGuestAmtWithStdOccupancy(std int) *Error {
	return newErrorf("missing element BaseByGuestAmt with attribute NumberOfGuests equal to the standard occupancy %d", std)
}

func ErrMinAgeOutOfRange(childMinAge, ratePlanChildMinAge int) *Error {
	return newErrorf("child min age must be ≥ rate plan child min age, got %d and %d", childMinAge, ratePlanChildMinAge)
}

func ErrMaxAgeOutOfRange(childMaxAge, ratePlanAdultMinAge int) *Error {
	return newErrorf("child max age must be < rate plan adult min age, got %d and %d", childMaxAge, ratePlanAdultMinAge)
}

func ErrInvalidInvType(invType string) *Error {
	return newErrorf("invalid value for attribute InvType %s", invType)
}

func newMissingAttributeError(attribute string) *Error {
	return newErrorf("missing required attribute %s", attribute)
}

func newMissingElementError(element string) *Error {
	return newErrorf("missing required element %s", element)
}

func newUnexpectedAttributeError(attribute string) *Error {
	return newErrorf("unexpected attribute found %s", attribute)
}

func newUnexpectedElementError(element string) *Error {
	return newErrorf("unexpected element found %s", element)
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
