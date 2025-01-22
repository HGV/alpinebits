package v_2018_10

import (
	"net/mail"
	"strings"

	"github.com/HGV/alpinebits/duration"
)

type ReadValidator struct{}

var _ Validatable[ReadRQ] = (*ReadValidator)(nil)

func (v ReadValidator) Validate(r ReadRQ) error {
	if err := validateHotelCode(r.HotelReadRequest.HotelCode); err != nil {
		return err
	}
	return nil
}

type ResRetrieveValidator struct {
	roomTypeCodes map[string]struct{}
	guestRequests []ResStatus // TODO: Naming
}

var _ Validatable[ResRetrieveRS] = (*ResRetrieveValidator)(nil)

type ResRetrieveValidatorFunc func(*ResRetrieveValidator)

func NewResRetrieveValidator(opts ...ResRetrieveValidatorFunc) ResRetrieveValidator {
	var v ResRetrieveValidator
	for _, opt := range opts {
		opt(&v)
	}
	return v
}

func WithRoomTypeCodes(mapping map[string]struct{}) ResRetrieveValidatorFunc {
	return func(v *ResRetrieveValidator) {
		v.roomTypeCodes = mapping
	}
}

func (v ResRetrieveValidator) Validate(r ResRetrieveRS) error {
	for _, res := range r.HotelReservations {
		if err := v.validateHotelReservation(res); err != nil {
			return err
		}
	}
	return nil
}

func (v ResRetrieveValidator) validateHotelReservation(h HotelReservation) error {
	if err := v.validateUniqueID(h.UniqueID, h.ResStatus); err != nil {
		return err
	}
	v.guestRequests = append(v.guestRequests, h.ResStatus)

	if err := v.validateRoomStays(h.RoomStays); err != nil {
		return err
	}

	if err := v.validateCustomer(h.Customer); err != nil {
		return err
	}

	if err := v.validateResGlobalInfo(h.ResGlobalInfo); err != nil {
		return err
	}

	return nil
}

func (v ResRetrieveValidator) validateUniqueID(uid UniqueID2, resStatus ResStatus) error {
	switch resStatus {
	case ResStatusRequested, ResStatusReserved, ResStatusModify:
		if uid.Type != UniqueIDType2Reservation {
			return ErrInvalidUniqueID(string(resStatus), int(uid.Type))
		}
	case ResStatusCancelled:
		if uid.Type != UniqueIDType2Cancellation {
			return ErrInvalidUniqueID(string(resStatus), int(uid.Type))
		}
	}

	if strings.TrimSpace(uid.ID) == "" {
		return ErrMissingID
	}

	return nil
}

func (v ResRetrieveValidator) validateRoomStays(roomStays []RoomStay) error {
	if len(roomStays) == 0 {
		return ErrMissingRoomStay
	}

	for _, roomStay := range roomStays {
		if err := v.validateRoomStay(roomStay); err != nil {
			return err
		}

		// validateAlternativeRoomStay...
	}

	return nil
}

func (v ResRetrieveValidator) validateRoomStay(roomStay RoomStay) error {
	if err := v.validateRoomType(roomStay.RoomType); err != nil {
		return err
	}

	if err := v.validateRatePlan(roomStay.RatePlan); err != nil {
		return err
	}

	if err := v.validateGuestCounts(roomStay.GuestCounts); err != nil {
		return err
	}

	if err := v.validateTimeSpan(roomStay.TimeSpan); err != nil {
		return err
	}

	if err := v.validateTotal(roomStay.Total); err != nil {
		return err
	}

	return nil
}

func (v ResRetrieveValidator) validateRoomType(roomType ResRoomType) error {
	if strings.TrimSpace(roomType.RoomTypeCode) == "" {
		return ErrMissingRoomTypeCode
	}

	if v.roomTypeCodes != nil { // TODO: oder len()>0
		if _, ok := v.roomTypeCodes[roomType.RoomTypeCode]; !ok {
			return ErrInvCodeNotFound(roomType.RoomTypeCode)
		}
	}

	return nil
}

func (v ResRetrieveValidator) validateRatePlan(ratePlan ResRatePlan) error {
	if strings.TrimSpace(ratePlan.RatePlanCode) == "" {
		return ErrMissingRatePlanCode
	}

	if c := ratePlan.Commission; c != nil {
		if err := v.validateCommission(*c); err != nil {
			return err
		}
	}

	if err := v.validateMealsIncluded(ratePlan.MealsIncluded); err != nil {
		return err
	}

	return nil
}

func (v ResRetrieveValidator) validateCommission(commission Commission) error {
	if commission.Percent != nil {
		if *commission.Percent > 100 {
			return ErrInvalidPercent
		}
	}

	return nil
}

func (v ResRetrieveValidator) validateMealsIncluded(mealsIncluded *MealsIncluded) error {
	if v.isReservation() && mealsIncluded == nil {
		return ErrMissingMealsIncluded
	}
	return nil
}

func (v ResRetrieveValidator) validateGuestCounts(guestCounts []GuestCount) error {
	if len(guestCounts) == 0 {
		return ErrMissingGuestCount
	}

	adultSeen := false
	for _, guestCount := range guestCounts {
		if guestCount.Age == nil && adultSeen {
			return ErrDuplicateAdultGuestCount
		}
		adultSeen = adultSeen || guestCount.Age == nil
	}

	return nil
}

func (v ResRetrieveValidator) validateTimeSpan(timeSpan TimeSpan) error {
	if v.isReservation() {
		if timeSpan.Start == nil {
			return ErrMissingStart
		}
		if timeSpan.End == nil {
			return ErrMissingEnd
		}
		if timeSpan.Start.After(*timeSpan.End) {
			return ErrStartAfterEnd
		}
	} else {
		if timeSpan.Duration == nil {
			return ErrMissingDuration
		}
		if err := v.validateStartDateWindow(timeSpan.StartDateWindow, *timeSpan.Duration); err != nil {
			return err
		}
	}
	return nil
}

func (v ResRetrieveValidator) validateStartDateWindow(w *StartDateWindow, nights duration.Nights) error {
	if w == nil {
		return ErrMissingStartDateWindow
	}

	if w.EarliestDate.After(w.LatestDate) {
		return ErrEarliestDateAfterLatestDate
	}

	if int(nights) <= w.LatestDate.DaysSince(w.EarliestDate) {
		return ErrDurationOutOfRange
	}

	return nil
}

func (v ResRetrieveValidator) validateTotal(total *Total) error {
	if v.isReservation() && total == nil {
		return ErrMissingTotal
	}
	return nil
}

func (v ResRetrieveValidator) validateCustomer(customer Customer) error {
	if err := v.validatePersonName(customer.PersonName); err != nil {
		return err
	}

	if customer.Email != nil {
		if err := v.validateEmail(*customer.Email); err != nil {
			return err
		}
	}

	if customer.Address != nil {
		if err := v.validateAddress(*customer.Address); err != nil {
			return err
		}
	}

	return nil
}

func (v ResRetrieveValidator) validatePersonName(personName PersonName) error {
	if personName.NamePrefix != nil && strings.TrimSpace(*personName.NamePrefix) == "" {
		return ErrInvalidNamePrefix
	}

	if strings.TrimSpace(personName.GivenName) == "" {
		return ErrMissingGivenName
	}

	if strings.TrimSpace(personName.Surname) == "" {
		return ErrMissingSurname
	}

	if personName.NameTitle != nil && strings.TrimSpace(*personName.NameTitle) == "" {
		return ErrInvalidNameTitle
	}

	return nil
}

func (v ResRetrieveValidator) validateEmail(email Email) error {
	_, err := mail.ParseAddress(email.Value)
	return err
}

func (v ResRetrieveValidator) validateAddress(address Address) error {
	if err := validateNonNilString(address.AddressLine); err != nil {
		return ErrInvalidAddressLine
	}

	if err := validateNonNilString(address.CityName); err != nil {
		return ErrInvalidCityName
	}

	if err := validateNonNilString(address.PostalCode); err != nil {
		return ErrInvalidPostalCode
	}

	if err := v.validateCountryName(address.CountryName); err != nil {
		return ErrInvalidCountryNameCode
	}

	return nil
}

func (v ResRetrieveValidator) validateCountryName(countryName *CountryName) error {
	if countryName == nil {
		return nil
	}
	return validateString(countryName.Code)
}

func (v ResRetrieveValidator) validateResGlobalInfo(globalInfo ResGlobalInfo) error {
	if err := v.validateComments(globalInfo.Comments); err != nil {
		return err
	}

	if v.isReservation() {
		if err := validateNonNilString(globalInfo.CancelPenalty); err != nil {
			return ErrInvalidPenaltyDescriptionText
		}
	}

	if err := v.validateHotelReservationID(globalInfo.HotelReservationID); err != nil {
		return err
	}

	if globalInfo.Profile != nil {
		if err := v.validateCompanyInfo(globalInfo.Profile.CompanyInfo); err != nil {
			return err
		}
	}

	if err := validateHotelCode(globalInfo.BasicPropertyInfo.HotelCode); err != nil {
		return err
	}

	return nil
}

func (v ResRetrieveValidator) validateComments(comments *[]Comment) error {
	if comments == nil {
		return nil
	}

	for _, comment := range *comments {
		for _, listItem := range comment.ListItems {
			if err := validateString(listItem.Value); err != nil {
				return ErrInvalidListItem
			}
		}
		if comment.Text != nil {
			if err := validateString(comment.Text.Value); err != nil {
				return ErrInvalidCommentText
			}
		}
	}

	return nil
}

func (v ResRetrieveValidator) validateHotelReservationID(id *HotelReservationID) error {
	if id == nil {
		return nil
	}

	if err := validateNonNilString(id.ResIDValue); err != nil {
		return ErrInvalidResIDValue
	}

	if err := validateNonNilString(id.ResIDSource); err != nil {
		return ErrInvalidResIDSource
	}

	if err := validateNonNilString(id.ResIDSourceContext); err != nil {
		return ErrInvalidResIDSourceContext
	}

	return nil
}

func (v ResRetrieveValidator) validateCompanyInfo(companyInfo CompanyInfo) error {
	if err := v.validateCompanyName(companyInfo.CompanyName); err != nil {
		return err
	}

	if companyInfo.AddressInfo != nil {
		if err := v.validateAddress(*companyInfo.AddressInfo); err != nil {
			return err
		}
	}

	if companyInfo.Email != nil {
		if err := v.validateEmail(*companyInfo.Email); err != nil {
			return ErrInvalidEmail
		}
	}

	return nil
}

func (v ResRetrieveValidator) validateCompanyName(companyName CompanyName) error {
	if err := validateString(companyName.Code); err != nil {
		return ErrInvalidCompanyNameCode
	}

	if err := validateString(companyName.Value); err != nil {
		return ErrInvalidCompanyNameValue
	}

	return nil
}

// Returns true if the current guest request being validated is a reservation.
func (v ResRetrieveValidator) isReservation() bool {
	status := v.guestRequests[len(v.guestRequests)-1]
	return status.IsReservation()
}
