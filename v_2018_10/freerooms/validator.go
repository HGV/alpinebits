package freerooms

import (
	"strings"

	"github.com/HGV/alpinebits/v_2018_10/common"
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

var _ common.Validatable[HotelAvailNotifRQ] = (*HotelAvailNotifValidator)(nil)

type HotelAvailNotifValidatorFunc func(*HotelAvailNotifValidator)

func NewHotelAvailNotifValidator(opts ...HotelAvailNotifValidatorFunc) HotelAvailNotifValidator {
	var v HotelAvailNotifValidator
	for _, opt := range opts {
		opt(&v)
	}
	return v
}

func WithRooms() HotelAvailNotifValidatorFunc {
	return func(v *HotelAvailNotifValidator) {
		v.supportsRooms = true
	}
}

func WithRoomMapping(mapping *map[string]map[string]struct{}) HotelAvailNotifValidatorFunc {
	return func(v *HotelAvailNotifValidator) {
		v.roomMapping = mapping
	}
}

func WithCategories() HotelAvailNotifValidatorFunc {
	return func(v *HotelAvailNotifValidator) {
		v.supportsCategories = true
	}
}

func WithCategoriesMapping(mapping *map[string]struct{}) HotelAvailNotifValidatorFunc {
	return func(v *HotelAvailNotifValidator) {
		v.categoriesMapping = mapping
	}
}

func WithDeltas() HotelAvailNotifValidatorFunc {
	return func(v *HotelAvailNotifValidator) {
		v.supportsDeltas = true
	}
}

func WithBookingThreshold() HotelAvailNotifValidatorFunc {
	return func(v *HotelAvailNotifValidator) {
		v.supportsBookingThreshold = true
	}
}

func (v HotelAvailNotifValidator) Validate(r HotelAvailNotifRQ) error {
	if err := common.ValidateHotelCode(r.AvailStatusMessages.HotelCode); err != nil {
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
		return common.ErrDeltasNotSupported
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
		return common.ErrInvalidBookingLimit(availableRooms)
	}

	if !v.supportsBookingThreshold && msg.BookingThreshold > 0 {
		return common.ErrBookingThresholdNotSupported
	}

	if msg.BookingThreshold > availableRooms {
		return common.ErrBookingThresholdGreaterThanBookingLimit
	}

	if err := v.validateStatusApplicationControl(msg.StatusApplicationControl); err != nil {
		return err
	}

	return nil
}

func (v HotelAvailNotifValidator) validateStatusApplicationControl(s StatusApplicationControl) error {
	if strings.TrimSpace(s.InvTypeCode) == "" {
		return common.ErrMissingInvTypeCode
	}

	if v.supportsRooms {
		if strings.TrimSpace(s.InvCode) == "" {
			return common.ErrMissingInvCode
		}
		if v.roomMapping != nil {
			if _, ok := (*v.roomMapping)[s.InvTypeCode][s.InvCode]; !ok {
				return common.ErrInvCodeNotFound(s.InvCode)
			}
		}
	} else if v.supportsCategories {
		if v.categoriesMapping != nil {
			if _, ok := (*v.categoriesMapping)[s.InvTypeCode]; !ok {
				return common.ErrInvTypeCodeNotFound(s.InvTypeCode)
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
		if err := common.ValidateOverlaps(avails); err != nil {
			return err
		}
	}

	return nil
}
