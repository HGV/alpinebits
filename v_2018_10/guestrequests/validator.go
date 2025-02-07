package guestrequests

import (
	"net/mail"
	"reflect"
	"strings"

	"github.com/HGV/alpinebits/duration"
	"github.com/HGV/alpinebits/v_2018_10/common"
	"github.com/HGV/alpinebits/v_2018_10/rateplans"
)

type ReadValidator struct{}

var _ common.Validatable[ReadRQ] = (*ReadValidator)(nil)

func (v ReadValidator) Validate(r ReadRQ) error {
	if err := common.ValidateHotelCode(r.HotelReadRequest.HotelCode); err != nil {
		return err
	}
	return nil
}

type ResRetrieveValidator struct {
	roomTypeCodes map[string]struct{}
	guestRequests []ResStatus // TODO: Naming
}

var _ common.Validatable[ResRetrieveRS] = (*ResRetrieveValidator)(nil)

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

func (v ResRetrieveValidator) validateUniqueID(uid UniqueID, resStatus ResStatus) error {
	switch resStatus {
	case ResStatusRequested, ResStatusReserved, ResStatusModify:
		if uid.Type != UniqueIDTypeReservation {
			return common.ErrInvalidUniqueID(string(resStatus), int(uid.Type))
		}
	case ResStatusCancelled:
		if uid.Type != UniqueIDTypeCancellation {
			return common.ErrInvalidUniqueID(string(resStatus), int(uid.Type))
		}
	}

	if strings.TrimSpace(uid.ID) == "" {
		return common.ErrMissingID
	}

	return nil
}

func (v ResRetrieveValidator) validateRoomStays(roomStays []RoomStay) error {
	if len(roomStays) == 0 && !v.isCancellation() {
		return common.ErrMissingRoomStay
	}

	for _, roomStay := range roomStays {
		if err := v.validateRoomStay(roomStay); err != nil {
			return err
		}

		// TODO: validateAlternativeRoomStay...
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
		return common.ErrMissingRoomTypeCode
	}

	if v.roomTypeCodes != nil {
		if _, ok := v.roomTypeCodes[roomType.RoomTypeCode]; !ok {
			return common.ErrInvCodeNotFound(roomType.RoomTypeCode)
		}
	}

	return nil
}

func (v ResRetrieveValidator) validateRatePlan(ratePlan ResRatePlan) error {
	if strings.TrimSpace(ratePlan.RatePlanCode) == "" {
		return common.ErrMissingRatePlanCode
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
			return common.ErrInvalidPercent
		}
	}

	return nil
}

func (v ResRetrieveValidator) validateMealsIncluded(mealsIncluded *rateplans.MealsIncluded) error {
	if v.isReservation() && mealsIncluded == nil {
		return common.ErrMissingMealsIncluded
	}
	return nil
}

func (v ResRetrieveValidator) validateGuestCounts(guestCounts []GuestCount) error {
	if len(guestCounts) == 0 {
		return common.ErrMissingGuestCount
	}

	adultSeen := false
	for _, guestCount := range guestCounts {
		if guestCount.Age == nil && adultSeen {
			return common.ErrDuplicateAdultGuestCount
		}
		adultSeen = adultSeen || guestCount.Age == nil
	}

	return nil
}

func (v ResRetrieveValidator) validateTimeSpan(timeSpan TimeSpan) error {
	if v.isReservation() {
		if timeSpan.Start == nil {
			return common.ErrMissingStart
		}
		if timeSpan.End == nil {
			return common.ErrMissingEnd
		}
		if timeSpan.Start.After(*timeSpan.End) {
			return common.ErrStartAfterEnd
		}
	} else {
		if timeSpan.Duration == nil {
			return common.ErrMissingDuration
		}
		if err := v.validateStartDateWindow(timeSpan.StartDateWindow, *timeSpan.Duration); err != nil {
			return err
		}
	}
	return nil
}

func (v ResRetrieveValidator) validateStartDateWindow(w *StartDateWindow, nights duration.Nights) error {
	if w == nil {
		return common.ErrMissingStartDateWindow
	}

	if w.EarliestDate.After(w.LatestDate) {
		return common.ErrEarliestDateAfterLatestDate
	}

	if int(nights) <= w.LatestDate.DaysSince(w.EarliestDate) {
		return common.ErrDurationOutOfRange
	}

	return nil
}

func (v ResRetrieveValidator) validateTotal(total *Total) error {
	if v.isReservation() && total == nil {
		return common.ErrMissingTotal
	}
	return nil
}

func (v ResRetrieveValidator) validateCustomer(customer Customer) error {
	if reflect.DeepEqual(customer, Customer{}) && v.isCancellation() {
		return nil
	}

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
		return common.ErrInvalidNamePrefix
	}

	if strings.TrimSpace(personName.GivenName) == "" {
		return common.ErrMissingGivenName
	}

	if strings.TrimSpace(personName.Surname) == "" {
		return common.ErrMissingSurname
	}

	if personName.NameTitle != nil && strings.TrimSpace(*personName.NameTitle) == "" {
		return common.ErrInvalidNameTitle
	}

	return nil
}

func (v ResRetrieveValidator) validateEmail(email Email) error {
	_, err := mail.ParseAddress(email.Value)
	return err
}

func (v ResRetrieveValidator) validateAddress(address Address) error {
	if err := common.ValidateNonNilString(address.AddressLine); err != nil {
		return common.ErrInvalidAddressLine
	}

	if err := common.ValidateNonNilString(address.CityName); err != nil {
		return common.ErrInvalidCityName
	}

	if err := common.ValidateNonNilString(address.PostalCode); err != nil {
		return common.ErrInvalidPostalCode
	}

	if err := v.validateCountryName(address.CountryName); err != nil {
		return common.ErrInvalidCountryNameCode
	}

	return nil
}

func (v ResRetrieveValidator) validateCountryName(countryName *CountryName) error {
	if countryName == nil {
		return nil
	}
	return common.ValidateString(countryName.Code)
}

func (v ResRetrieveValidator) validateResGlobalInfo(globalInfo ResGlobalInfo) error {
	if err := v.validateComments(globalInfo.Comments); err != nil {
		return err
	}

	if v.isReservation() {
		if err := common.ValidateNonNilString(globalInfo.CancelPenalty); err != nil {
			return common.ErrInvalidPenaltyDescriptionText
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

	if err := common.ValidateHotelCode(globalInfo.BasicPropertyInfo.HotelCode); err != nil && !v.isCancellation() {
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
			if err := common.ValidateString(listItem.Value); err != nil {
				return common.ErrInvalidListItem
			}
		}
		if comment.Text != nil {
			if err := common.ValidateString(comment.Text.Value); err != nil {
				return common.ErrInvalidCommentText
			}
		}
	}

	return nil
}

func (v ResRetrieveValidator) validateHotelReservationID(id *HotelReservationID) error {
	if id == nil {
		return nil
	}

	if err := common.ValidateNonNilString(id.ResIDValue); err != nil {
		return common.ErrInvalidResIDValue
	}

	if err := common.ValidateNonNilString(id.ResIDSource); err != nil {
		return common.ErrInvalidResIDSource
	}

	if err := common.ValidateNonNilString(id.ResIDSourceContext); err != nil {
		return common.ErrInvalidResIDSourceContext
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
			return common.ErrInvalidEmail
		}
	}

	return nil
}

func (v ResRetrieveValidator) validateCompanyName(companyName CompanyName) error {
	if err := common.ValidateString(companyName.Code); err != nil {
		return common.ErrInvalidCompanyNameCode
	}

	if err := common.ValidateString(companyName.Value); err != nil {
		return common.ErrInvalidCompanyNameValue
	}

	return nil
}

// Returns true if the current guest request being validated is a reservation.
func (v ResRetrieveValidator) isReservation() bool {
	status := v.guestRequests[len(v.guestRequests)-1]
	return status.IsReservation()
}

// Returns true if the current guest request being validated is a cancellation.
func (v ResRetrieveValidator) isCancellation() bool {
	status := v.guestRequests[len(v.guestRequests)-1]
	return status == ResStatusCancelled
}
