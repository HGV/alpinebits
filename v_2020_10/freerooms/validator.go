package freerooms

import (
	"errors"
	"slices"
	"strings"

	"github.com/HGV/alpinebits/v_2020_10/common"
	"github.com/HGV/x/slicesx"
)

type HotelInvCountNotifValidator struct {
	supportsRooms          bool
	roomMapping            *map[string]map[string]struct{}
	supportsCategories     bool
	categoriesMapping      *map[string]struct{}
	supportsDeltas         bool
	supportsOutOfOrder     bool
	supportsOutOfMarket    bool
	supportsClosingSeasons bool
}

var _ common.Validatable[HotelInvCountNotifRQ] = (*HotelInvCountNotifValidator)(nil)

type HotelInvCountNotifValidatorFunc func(*HotelInvCountNotifValidator)

func NewHotelInvCountNotifValidator(opts ...HotelInvCountNotifValidatorFunc) HotelInvCountNotifValidator {
	var v HotelInvCountNotifValidator
	for _, opt := range opts {
		opt(&v)
	}
	return v
}

func WithRooms() HotelInvCountNotifValidatorFunc {
	return func(v *HotelInvCountNotifValidator) {
		v.supportsRooms = true
	}
}

func WithRoomMapping(mapping map[string]map[string]struct{}) HotelInvCountNotifValidatorFunc {
	return func(v *HotelInvCountNotifValidator) {
		v.roomMapping = &mapping
	}
}

func WithCategories() HotelInvCountNotifValidatorFunc {
	return func(v *HotelInvCountNotifValidator) {
		v.supportsCategories = true
	}
}

func WithCategoriesMapping(mapping map[string]struct{}) HotelInvCountNotifValidatorFunc {
	return func(v *HotelInvCountNotifValidator) {
		v.categoriesMapping = &mapping
	}
}

func WithDeltas() HotelInvCountNotifValidatorFunc {
	return func(v *HotelInvCountNotifValidator) {
		v.supportsDeltas = true
	}
}

func WithOutOfOrder() HotelInvCountNotifValidatorFunc {
	return func(v *HotelInvCountNotifValidator) {
		v.supportsOutOfOrder = true
	}
}

func WithOutOfMarket() HotelInvCountNotifValidatorFunc {
	return func(v *HotelInvCountNotifValidator) {
		v.supportsOutOfMarket = true
	}
}

func WithClosingSeasons() HotelInvCountNotifValidatorFunc {
	return func(v *HotelInvCountNotifValidator) {
		v.supportsClosingSeasons = true
	}
}

func (v HotelInvCountNotifValidator) Validate(r HotelInvCountNotifRQ) error {
	if err := common.ValidateHotelCode(r.Inventories.HotelCode); err != nil {
		return err
	}

	if err := v.validateUniqueID(r.UniqueID); err != nil {
		return err
	}

	if r.Inventories.IsReset() {
		return nil
	}

	if err := v.validateInventories(r.Inventories.Inventories); err != nil {
		return err
	}

	return nil
}

func (v HotelInvCountNotifValidator) validateUniqueID(uid *UniqueID) error {
	if uid == nil && !v.supportsDeltas {
		return common.ErrDeltasNotSupported
	}
	return nil
}

func (v HotelInvCountNotifValidator) validateInventories(invs []Inventory) error {
	avails := slicesx.Filter(invs, Inventory.isAvailability)
	if err := v.validateAvailabilities(avails); err != nil {
		return err
	}

	closingSeasons := slicesx.Filter(invs, Inventory.isClosingSeason)
	if err := v.validateClosingSeasons(closingSeasons); err != nil {
		return err
	}

	if v.supportsClosingSeasons {
		if err := v.validateClosingSeasonsOverlapBookableAvailabilities(avails, closingSeasons); err != nil {
			return err
		}
	}

	return nil
}

func (v HotelInvCountNotifValidator) validateAvailabilities(avails []Inventory) error {
	for _, avail := range avails {
		if err := v.validateAvailability(avail); err != nil {
			return err
		}
	}

	if err := v.validateAvailabilitiesOverlap(avails); err != nil {
		return err
	}

	return nil
}

func (v HotelInvCountNotifValidator) validateAvailability(avail Inventory) error {
	if err := v.validateStatusApplicationControl(avail.StatusApplicationControl); err != nil {
		return err
	}

	if err := v.validateInvCounts(avail.InvCounts); err != nil {
		return err
	}

	return nil
}

func (v HotelInvCountNotifValidator) validateStatusApplicationControl(s *StatusApplicationControl) error {
	if s == nil {
		return errors.New("")
	}

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

func (v HotelInvCountNotifValidator) validateInvCounts(invCounts *[]InvCount) error {
	if invCounts == nil {
		return nil
	}

	if n := len(*invCounts); v.supportsRooms && n > 1 {
		return common.ErrInvalidInvCounts(n)
	}

	for _, invCount := range *invCounts {
		if v.supportsRooms && invCount.Count > 1 {
			return common.ErrInvalidCount(invCount.Count)
		}

		switch ct := invCount.CountType; ct {
		case CountTypeBookable:
		case CountTypeOutOfOrder:
			if !v.supportsOutOfOrder {
				return common.ErrOutOfOrderNotSupported
			}
		case CountTypeFree:
			if !v.supportsOutOfMarket {
				return common.ErrOutOfMarketNotSupported
			}
		}
	}

	return nil
}

func (v HotelInvCountNotifValidator) validateAvailabilitiesOverlap(invs []Inventory) error {
	availsBy := v.groupAvailabilitiesByCapability(invs)
	for _, avails := range availsBy {
		if err := common.ValidateOverlaps(avails); err != nil {
			return err
		}
	}
	return nil
}

func (v HotelInvCountNotifValidator) validateClosingSeasons(closingSeasons []Inventory) error {
	if !v.supportsClosingSeasons && len(closingSeasons) > 0 {
		return common.ErrClosingSeasonsNotSupported
	}

	for _, closingSeason := range closingSeasons {
		if err := v.validateClosingSeason(closingSeason); err != nil {
			return err
		}
	}

	if err := common.ValidateOverlaps(closingSeasons); err != nil {
		return err
	}

	return nil
}

func (v HotelInvCountNotifValidator) validateClosingSeason(closingSeason Inventory) error {
	if closingSeason.InvCounts != nil {
		return common.ErrUnexpectedInvCounts
	}
	return nil
}

func (v HotelInvCountNotifValidator) validateClosingSeasonsOverlapBookableAvailabilities(avails []Inventory, closingSeasons []Inventory) error {
	bookableAvails := slicesx.Filter(avails, func(a Inventory) bool {
		return a.InvCounts != nil && slices.ContainsFunc(*a.InvCounts, func(i InvCount) bool {
			return i.CountType == CountTypeBookable
		})
	})

	bookableAvailsBy := v.groupAvailabilitiesByCapability(bookableAvails)
	for _, avails := range bookableAvailsBy {
		if err := common.ValidateOverlaps(append(avails, closingSeasons...)); err != nil {
			return common.ErrAvailabilitiesOverlapClosingSeasons
		}
	}

	return nil
}

func (v HotelInvCountNotifValidator) groupAvailabilitiesByCapability(avails []Inventory) map[string][]Inventory {
	return slicesx.GroupByFunc(avails, func(inv Inventory) string {
		switch {
		case v.supportsRooms:
			return inv.StatusApplicationControl.InvCode
		case v.supportsCategories:
			return inv.StatusApplicationControl.InvTypeCode
		default:
			return ""
		}
	})
}
