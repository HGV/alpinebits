package v_2018_10

import (
	"strings"

	"github.com/HGV/x/slicesx"
)

type HotelAvailNotifValidator struct {
	supportsRooms            bool
	roomMapping              *map[string]map[string]struct{}
	supportsCategories       bool
	categoriesMapping        *map[string]struct{}
	supportsDeltas           bool
	supportsBookingThreshold bool
}

var _ Validatable[HotelAvailNotifRQ] = (*HotelAvailNotifValidator)(nil)

type HotelAvailNotifValidatorFunc func(*HotelAvailNotifValidator)

func NewHotelAvailNotifValidator(opts ...HotelAvailNotifValidatorFunc) HotelAvailNotifValidator {
	var v HotelAvailNotifValidator
	for _, opt := range opts {
		opt(&v)
	}
	return v
}

func WithRooms(supports bool, mapping *map[string]map[string]struct{}) HotelAvailNotifValidatorFunc {
	return func(v *HotelAvailNotifValidator) {
		v.supportsRooms = supports
		v.roomMapping = mapping
	}
}

func WithCategories(supports bool, mapping *map[string]struct{}) HotelAvailNotifValidatorFunc {
	return func(v *HotelAvailNotifValidator) {
		v.supportsCategories = supports
		v.categoriesMapping = mapping
	}
}

func WithDeltas(supports bool) HotelAvailNotifValidatorFunc {
	return func(v *HotelAvailNotifValidator) {
		v.supportsDeltas = supports
	}
}

func WithBookingThreshold(supports bool) HotelAvailNotifValidatorFunc {
	return func(v *HotelAvailNotifValidator) {
		v.supportsBookingThreshold = supports
	}
}

func (v HotelAvailNotifValidator) Validate(r HotelAvailNotifRQ) error {
	if err := validateHotelCode(r.AvailStatusMessages.HotelCode); err != nil {
		return err
	}

	if err := v.validateUniqueID(r.UniqueID); err != nil {
		return err
	}

	if r.AvailStatusMessages.IsReset() {
		return nil
	}

	if err := v.validateAvailStatusMessages(r.AvailStatusMessages.AvailStatusMessages); err != nil {
		return err
	}

	if err := v.validateOverlaps(r.AvailStatusMessages.AvailStatusMessages); err != nil {
		return err
	}

	return nil
}

func (v HotelAvailNotifValidator) validateUniqueID(uid *UniqueID) error {
	if uid == nil && !v.supportsDeltas {
		return ErrDeltasNotSupported
	}
	return nil
}

func (v HotelAvailNotifValidator) validateAvailStatusMessages(msgs []AvailStatusMessage) error {
	for _, msg := range msgs {
		if err := v.validateAvailStatusMessage(msg); err != nil {
			return err
		}
	}
	return nil
}

func (v HotelAvailNotifValidator) validateAvailStatusMessage(msg AvailStatusMessage) error {
	availableRooms := msg.BookingLimit
	if availableRooms > 1 && v.supportsRooms {
		return ErrInvalidBookingLimit(availableRooms)
	}

	if !v.supportsBookingThreshold && msg.BookingThreshold > 0 {
		return ErrBookingThresholdNotSupported
	}

	if msg.BookingThreshold > availableRooms {
		return ErrBookingThresholdGreaterThanBookingLimit
	}

	if err := v.validateStatusApplicationControl(msg.StatusApplicationControl); err != nil {
		return err
	}

	return nil
}

func (v HotelAvailNotifValidator) validateStatusApplicationControl(s StatusApplicationControl) error {
	if strings.TrimSpace(s.InvTypeCode) == "" {
		return ErrMissingInvTypeCode
	}

	if v.supportsRooms {
		if strings.TrimSpace(s.InvCode) == "" {
			return ErrMissingInvCode
		}
		if v.roomMapping != nil {
			if _, ok := (*v.roomMapping)[s.InvTypeCode][s.InvCode]; !ok {
				return ErrInvCodeNotFound(s.InvCode)
			}
		}
	} else if v.supportsCategories {
		if v.categoriesMapping != nil {
			if _, ok := (*v.categoriesMapping)[s.InvTypeCode]; !ok {
				return ErrInvTypeCodeNotFound(s.InvTypeCode)
			}
		}
	}

	return nil
}

func (v HotelAvailNotifValidator) validateOverlaps(msgs []AvailStatusMessage) error {
	availsBy := slicesx.GroupByFunc(msgs, func(msg AvailStatusMessage) string {
		switch {
		case v.supportsRooms:
			return msg.StatusApplicationControl.InvCode
		case v.supportsCategories:
			return msg.StatusApplicationControl.InvTypeCode
		default:
			return ""
		}
	})

	for _, avails := range availsBy {
		if err := validateOverlaps(avails); err != nil {
			return err
		}
	}

	return nil
}
