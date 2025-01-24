package v_2018_10

import (
	"github.com/HGV/alpinebits/internal"
	"github.com/HGV/x/slicesx"
)

type OccupancyRules struct {
	Min    *int
	Max    *int
	MinAge *int
}

type HotelRatePlanNotifValidator struct {
	ratePlanMapping      map[string]any
	supportsRatePlanJoin bool
	adultOccupancy       OccupancyRules
	childOccupancy       OccupancyRules
}

var _ Validatable[HotelRatePlanNotifRQ] = (*HotelRatePlanNotifValidator)(nil)

type HotelRatePlanNotifValidatorFunc func(*HotelRatePlanNotifValidator)

func NewHotelRatePlanNotifValidator(opts ...HotelRatePlanNotifValidatorFunc) HotelRatePlanNotifValidator {
	var v HotelRatePlanNotifValidator
	for _, opt := range opts {
		opt(&v)
	}
	return v
}

func WithRatePlanJoin(supports bool) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsRatePlanJoin = supports
	}
}

func (v HotelRatePlanNotifValidator) Validate(r HotelRatePlanNotifRQ) error {
	if err := validateHotelCode(r.RatePlans.HotelCode); err != nil {
		return err
	}

	// TODO: uniqueid?

	if err := v.validateRatePlans(r.RatePlans.RatePlans); err != nil {
		return err
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateRatePlans(ratePlans []RatePlan) error {
	for _, ratePlan := range ratePlans {
		if err := v.validateRatePlan(ratePlan); err != nil {
			return err
		}
	}
	return nil
}

func (v HotelRatePlanNotifValidator) validateRatePlan(ratePlan RatePlan) error {
	if err := v.validateRatePlanCode(ratePlan.RatePlanCode); err != nil {
		return err
	}

	if err := v.validateCurrencyCode(ratePlan.CurrencyCode); err != nil {
		return err
	}

	if !v.supportsRatePlanJoin && (ratePlan.RatePlanID != "" || ratePlan.RatePlanQualifier) {
		return ErrRatePlanJoinNotSupported
	}

	switch ratePlan.RatePlanNotifType {
	case RatePlanNotifTypeNew:
		return v.validateRatePlanNew(ratePlan)
	case RatePlanNotifTypeOverlay:
		return v.validateRatePlanOverlay(ratePlan)
	case RatePlanNotifTypeRemove:
		return v.validateRatePlanRemove(ratePlan)
	default:
		return nil
	}
}

func (v HotelRatePlanNotifValidator) validateRatePlanCode(code string) error {
	if err := validateString(code); err != nil {
		return ErrMissingRatePlanCode
	}
	return nil
}

func (v HotelRatePlanNotifValidator) validateCurrencyCode(code string) error {
	if err := validateString(code); err != nil {
		return ErrMissingCurrencyCode
	}
	return nil
}

func (v HotelRatePlanNotifValidator) validateRatePlanNew(ratePlan RatePlan) error {
	if err := v.validateOffers(ratePlan.Offers); err != nil {
		return err
	}

	// TODO: br
	// TODO: rt
	// TODO: sp

	return nil
}

func (v HotelRatePlanNotifValidator) validateOffers(offers []Offer) error {
	if len(offers) == 0 {
		return ErrMissingOfferRule
	}

	if err := v.validateOfferRule(offers[0].OfferRule); err != nil {
		return err
	}

	if err := v.validateAdditionalOffers(offers[1:]); err != nil {
		return err
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateOfferRule(offerRule *OfferRule) error {
	if offerRule == nil {
		return ErrMissingOfferRule
	}

	if err := v.validateOfferRuleLengthsOfStay(offerRule); err != nil {
		return err
	}

	if err := v.validateOccupancies(offerRule.Occupancies); err != nil {
		return err
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateOfferRuleLengthsOfStay(offerRule *OfferRule) error {
	var minArrival int
	var maxArrival int
	for _, los := range offerRule.LengthsOfStay {
		switch los.MinMaxMessageType {
		case StayTypeMinArrival:
			minArrival = los.Time
		case StayTypeMaxArrival:
			maxArrival = los.Time
		case StayTypeMinThrough, StayTypeMaxThrough:
			return ErrStayThroughNotAllowedInOfferRule
		}
	}

	if maxArrival > 0 && minArrival > maxArrival {
		return ErrMinArrivalGratherThanMaxStayArrival(minArrival, maxArrival)
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateOccupancies(occupancies []Occupancy) error {
	adults := slicesx.Filter(occupancies, Occupancy.isAdult)
	switch len(adults) {
	case 0:
		return ErrMissingAdultOccupancy
	case 1:
		adultOccupancy := adults[0]
		if err := v.validateOccupancy(adultOccupancy); err != nil {
			return err
		}
		v.populateAdultOccupancy(adultOccupancy)
	}

	children := slicesx.Filter(occupancies, Occupancy.isChild)
	switch len(children) {
	case 0:
		break
	case 1:
		childOccupancy := children[0]
		if err := v.validateOccupancy(childOccupancy); err != nil {
			return err
		}
		v.populateChildOccupancy(childOccupancy)
	default:
		return ErrDuplicateChildOccupancy
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateOccupancy(o Occupancy) error {
	if min := o.MinOccupancy; min != nil && *min > 99 {
		return ErrInvalidMinOccupancy
	}
	if max := o.MaxOccupancy; max != nil && *max > 99 {
		return ErrInvalidMaxOccupancy
	}
	return nil
}

func (v *HotelRatePlanNotifValidator) populateAdultOccupancy(occupancy Occupancy) {
	v.adultOccupancy.Min = occupancy.MinOccupancy
	v.adultOccupancy.Max = occupancy.MaxOccupancy
	v.adultOccupancy.MinAge = occupancy.MinAge
}

func (v *HotelRatePlanNotifValidator) populateChildOccupancy(occupancy Occupancy) {
	v.childOccupancy.Min = occupancy.MinOccupancy
	v.childOccupancy.Max = occupancy.MaxOccupancy
	v.childOccupancy.MinAge = occupancy.MinAge
}

func (v HotelRatePlanNotifValidator) validateAdditionalOffers(offers []Offer) error {
	freeNightOffers := slicesx.Filter(offers, Offer.IsFreeNightOffer)
	switch len(freeNightOffers) {
	case 0:
		break
	case 1:
		if err := v.validateFreeNightOffer(freeNightOffers[0]); err != nil {
			return err
		}
	default:
		return ErrDuplicateFreeNightOffer
	}

	familyOffers := slicesx.Filter(offers, Offer.IsFamilyOffer)
	switch len(familyOffers) {
	case 0:
		break
	case 1:
		if err := v.validateFamilyOffer(familyOffers[0]); err != nil {
			return err
		}
	default:
		return ErrDuplicateFamilyOffer
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateFreeNightOffer(offer Offer) error {
	if offer.Discount.NightsRequired == 0 {
		return ErrMissingNightsRequired
	}

	if offer.Discount.NightsDiscounted == 0 {
		return ErrMissingNightsDiscounted
	}

	if pattern := offer.Discount.DiscountPattern; pattern != "" {
		expectedPattern := internal.CalculateDiscountPattern(
			offer.Discount.NightsRequired,
			offer.Discount.NightsDiscounted,
		)
		if pattern != expectedPattern {
			return ErrInvalidDiscountPattern
		}
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateFamilyOffer(offer Offer) error {
	if offer.Guest.AgeQualifyingCode != AgeQualifyingCodeChild {
		return ErrInvalidGuestAgeQualifyngCode
	}
	return nil
}

// func (v HotelRatePlanNotifValidator) validateBookingRules(bookingRules []BookingRule) error {
// 	return nil
// }

// func (v HotelRatePlanNotifValidator) validateRates(rates []Rate) error {
// 	return nil
// }

// func (v HotelRatePlanNotifValidator) validateSupplements(supplements []Supplement) error {
// 	return nil
// }

func (v HotelRatePlanNotifValidator) validateRatePlanOverlay(ratePlan RatePlan) error {
	return nil
}

func (v HotelRatePlanNotifValidator) validateRatePlanRemove(ratePlan RatePlan) error {
	return nil
}
