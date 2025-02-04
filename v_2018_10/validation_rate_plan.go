package v_2018_10

import (
	"math"
	"regexp"

	"github.com/HGV/alpinebits/internal"
	"github.com/HGV/x/slicesx"
)

type RatePlanOccupancySettings struct {
	Min    *int
	Max    *int
	MinAge *int
}

type RoomTypeOccupancySettings struct {
	Min int
	Std int
	Max int
}

type HotelRatePlanNotifValidator struct {
	supportsArrivalDOW             bool
	supportsDepartureDOW           bool
	ratePlanMapping                map[string]any
	supportsRatePlanJoin           bool
	adultOccupancy                 RatePlanOccupancySettings
	childOccupancy                 RatePlanOccupancySettings
	supportsOverlay                bool
	supportsGenericBookingRules    bool
	supportsRoomTypeBokingRules    bool
	roomTypeMapping                map[string]RoomTypeOccupancySettings
	supplementMapping              map[string]struct{}
	supportsSupplements            bool
	supportsFreeNightOffer         bool
	supportsFamilyOffer            bool
	supportsOfferRuleBookingOffset bool
	supportsOfferRuleDOWLOS        bool
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

func WithArrivalDOW(supports bool) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsArrivalDOW = supports
	}
}

func WithDepartureDOW(supports bool) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsDepartureDOW = supports
	}
}

func WithRatePlanJoin(supports bool) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsRatePlanJoin = supports
	}
}

func WithGenericBookingRules(supports bool) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsGenericBookingRules = supports
	}
}

func WithRoomTypeBookingRules(supports bool) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsRoomTypeBokingRules = supports
	}
}

// func WithRoomTypeCodes(mapping map[string]struct{}) HotelRatePlanNotifValidatorFunc {
// 	return func(v *HotelRatePlanNotifValidator) {
// 		v.supportsRoomTypeBokingRules = true
// 		v.roomTypeMapping = mapping
// 	}
// }

func WithSupplements(supports bool) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsSupplements = supports
	}
}

func WithFreeNightOffer(supports bool) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsFreeNightOffer = supports
	}
}

func WithFamilyOffer(supports bool) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsFamilyOffer = supports
	}
}

func WithOfferRuleBookingOffset(supports bool) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsOfferRuleBookingOffset = supports
	}
}

func WithOfferRuleDOWLOS(supports bool) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsOfferRuleDOWLOS = supports
	}
}

func (v HotelRatePlanNotifValidator) Validate(r HotelRatePlanNotifRQ) error {
	if err := validateHotelCode(r.RatePlans.HotelCode); err != nil {
		return err
	}

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

	if err := v.validateDescriptions(ratePlan.Descriptions); err != nil {
		return err
	}

	if err := v.validateBookingRules(ratePlan.BookingRules); err != nil {
		return err
	}

	if err := v.validateRates(ratePlan.Rates); err != nil {
		return err
	}

	if err := v.validateSupplements(ratePlan.Supplements); err != nil {
		return err
	}

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

	if !v.supportsOfferRuleBookingOffset &&
		(offerRule.MinAdvancedBookingOffset != nil || offerRule.MaxAdvancedBookingOffset != nil) {
		return ErrOfferRuleBookingOffsetNotSupported
	}

	if len(offerRule.LengthsOfStay) > 0 ||
		offerRule.ArrivalDaysOfWeek != nil ||
		offerRule.DepartureDaysOfWeek != nil {
		return ErrOfferRuleDOWLOSNotSupported
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
		return ErrMinStayArrivalGratherThanMaxStayArrival(minArrival, maxArrival)
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
	if !v.supportsFreeNightOffer {
		return ErrFreeNightOfferNotSupported
	}

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

	if offer.Guest != nil {
		return ErrUnexpectedGuest
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateFamilyOffer(offer Offer) error {
	if !v.supportsFamilyOffer {
		return ErrFamilyOfferNotSupported
	}

	if offer.Guest.AgeQualifyingCode != AgeQualifyingCodeChild {
		return ErrInvalidGuestAgeQualifyngCode
	}

	if offer.Discount.NightsRequired > 0 {
		return ErrUnexpectedNightsRequired
	}

	if offer.Discount.NightsDiscounted > 0 {
		return ErrUnexpectedNightsDiscounted
	}

	if offer.Discount.DiscountPattern != "" {
		return ErrUnexpectedDiscountPattern
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateDescriptions(d RatePlanDescription) error {
	if err := validateLanguageUniqueness(d.Titles); err != nil {
		return err
	}

	if err := validateLanguageUniqueness(d.Intros); err != nil {
		return err
	}

	if err := validateLanguageUniqueness(d.Descriptions); err != nil {
		return err
	}

	for _, item := range d.Gallery {
		if err := validateLanguageUniqueness(item.Descriptions); err != nil {
			return err
		}
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateBookingRules(bookingRules []BookingRule) error {
	for _, bookgingRule := range bookingRules {
		if err := v.validateBookingRule(bookgingRule); err != nil {
			return err
		}
	}

	if err := v.validateBookingRuleOverlaps(bookingRules); err != nil {
		return err
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateBookingRule(bookingRule BookingRule) error {
	if v.supportsRoomTypeBokingRules {
		if err := validateString(bookingRule.Code); err != nil {
			return ErrMissingCode
		}
		if _, ok := v.roomTypeMapping[bookingRule.Code]; !ok {
			return ErrInvTypeCodeNotFound(bookingRule.Code)
		}
	} else if v.supportsGenericBookingRules {
		if bookingRule.Code != "" || bookingRule.CodeContext != "" {
			return ErrRoomTypeBookingRulesNotSupported
		}
	}

	if bookingRule.Start.After(bookingRule.End) {
		return ErrStartAfterEnd
	}

	if err := v.validateLengthsOfStay(bookingRule.LengthsOfStay); err != nil {
		return err
	}

	if !v.supportsArrivalDOW && bookingRule.ArrivalDaysOfWeek != nil {
		return ErrArrivalDOWNotSupported
	}

	if !v.supportsDepartureDOW && bookingRule.DepartureDaysOfWeek != nil {
		return ErrDepartureDOWNotSupported
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateLengthsOfStay(lengthsOfStay []LengthOfStay) error {
	minArrival, minThrough := 1, 1
	maxArrival, maxThrough := math.MaxInt32, math.MaxInt32

	for _, lengthOfStay := range lengthsOfStay {
		switch lengthOfStay.MinMaxMessageType {
		case StayTypeMinArrival:
			minArrival = lengthOfStay.Time
		case StayTypeMaxArrival:
			maxArrival = lengthOfStay.Time
		case StayTypeMinThrough:
			minThrough = lengthOfStay.Time
		case StayTypeMaxThrough:
			maxThrough = lengthOfStay.Time
		}
	}

	min := int(math.Max(float64(minArrival), float64(minThrough)))
	max := int(math.Min(float64(maxArrival), float64(maxThrough)))
	if min > max {
		return ErrMinStayGratherThanMaxStay(min, max)
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateBookingRuleOverlaps(bookingRules []BookingRule) error {
	if v.supportsRoomTypeBokingRules {
		bookingRulesByRoomType := slicesx.GroupByFunc(bookingRules, func(b BookingRule) string {
			return b.Code
		})
		for _, brs := range bookingRulesByRoomType {
			if err := validateOverlaps(brs); err != nil {
				return err
			}
		}
	} else if v.supportsGenericBookingRules {
		if err := validateOverlaps(bookingRules); err != nil {
			return err
		}
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateRates(rates []Rate) error {
	if len(rates) == 0 {
		return ErrMissingStaticRate
	}

	if err := v.validateStaticRate(rates[0]); err != nil {
		return err
	}

	if err := v.validateDateDependingRates(rates[1:]); err != nil {
		return err
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateStaticRate(rate Rate) error {
	if rate.RateTimeUnit != TimeUnitDay {
		return ErrInvalidRateTimeUnit
	}

	if rate.UnitMultiplier == 0 {
		return ErrMissingUnitMultiplier
	}

	switch len(rate.BaseByGuestAmts) {
	case 0:
		return ErrMissingBaseByGuestAmt
	case 1:
		b := rate.BaseByGuestAmts[0]
		if b.NumberOfGuests != nil {
			return ErrUnexpectedNumberOfGuests
		}
		if b.AgeQualifyingCode != nil {
			return ErrUnexpectedAgeQualifyingCode
		}
		if b.AmountAfterTax != nil {
			return ErrUnexpectedAmountAfterTax
		}
	default:
		return ErrUnexpectedBaseByGuestAmt
	}

	if rate.MealsIncluded == nil {
		return ErrMissingMealsIncluded
	}

	if rate.InvTypeCode != "" {
		return ErrUnexpectedInvTypeCode
	}

	if rate.Start != nil {
		return ErrUnexpectedStart
	}

	if rate.End != nil {
		return ErrUnexpectedEnd
	}

	if len(rate.AdditionalGuestAmounts) > 0 {
		return ErrUnexpectedAdditionalGuestAmounts
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateDateDependingRates(rates []Rate) error {
	for _, rate := range rates {
		if err := v.validateDateDependingRate(rate); err != nil {
			return err
		}
	}

	if err := v.validateDateDependingRateOverlaps(rates); err != nil {
		return err
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateDateDependingRate(rate Rate) error {
	if err := validateString(rate.InvTypeCode); err != nil {
		return ErrMissingInvTypeCode
	}

	roomTypeOccupancySettings, ok := v.roomTypeMapping[rate.InvTypeCode]
	if !ok {
		return ErrInvTypeCodeNotFound(rate.InvTypeCode)
	}

	if rate.Start == nil {
		return ErrMissingStart
	}

	if rate.End == nil {
		return ErrMissingEnd
	}

	if rate.Start.After(*rate.End) {
		return ErrStartAfterEnd
	}

	if err := v.validateBaseByGuestAmts(rate.BaseByGuestAmts, roomTypeOccupancySettings); err != nil {
		return err
	}

	if err := v.validateAdditionalGuestAmounts(rate.AdditionalGuestAmounts); err != nil {
		return err
	}

	if rate.RateTimeUnit != "" {
		return ErrUnexpectedRateTimeUnit
	}

	if rate.UnitMultiplier > 0 {
		return ErrUnexpectedUnitMultiplier
	}

	if rate.MealsIncluded != nil {
		return ErrUnexpectedMealsIncluded
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateBaseByGuestAmts(baseByGuestAmts []BaseByGuestAmt, roomTypeOccupancySettings RoomTypeOccupancySettings) error {
	numberOfGuestSeen := make(map[int]struct{})
	stdOccupancySeen := false
	for _, baseByGuestAmt := range baseByGuestAmts {
		if err := v.validateBaseByGuestAmt(baseByGuestAmt); err != nil {
			return err
		}

		numberOfGuests := *baseByGuestAmt.NumberOfGuests
		if _, exists := numberOfGuestSeen[numberOfGuests]; exists {
			return ErrDuplicateBaseByGuestAmt(numberOfGuests)
		}
		numberOfGuestSeen[numberOfGuests] = struct{}{}

		if numberOfGuests == roomTypeOccupancySettings.Std {
			stdOccupancySeen = true
		}
	}

	if !stdOccupancySeen {
		return ErrMissingBaseByGuestAmtWithStdOccupancy(roomTypeOccupancySettings.Std)
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateBaseByGuestAmt(baseByGuestAmt BaseByGuestAmt) error {
	if baseByGuestAmt.NumberOfGuests == nil {
		return ErrMissingNumberOfGuests
	}

	if baseByGuestAmt.AgeQualifyingCode == nil {
		return ErrMissingAgeQualifyingCode
	}

	if baseByGuestAmt.AmountAfterTax == nil {
		return ErrMissingAmountAfterTax
	}

	if baseByGuestAmt.Type != nil {
		return ErrUnexpectedType
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateAdditionalGuestAmounts(additionalGuestAmounts []AdditionalGuestAmount) error {
	adults := slicesx.Filter(additionalGuestAmounts, func(a AdditionalGuestAmount) bool {
		return a.AgeQualifyingCode != nil && *a.AgeQualifyingCode == AgeQualifyingCodeAdult
	})
	switch len(adults) {
	case 0:
		break
	case 1:
		if adults[0].Amount == nil {
			return ErrMissingAmount
		}
	default:
		return ErrDuplicateAdditionalGuestAmountAdult
	}

	children := slicesx.Filter(additionalGuestAmounts, func(a AdditionalGuestAmount) bool {
		return a.AgeQualifyingCode != nil && *a.AgeQualifyingCode == AgeQualifyingCodeChild
	})
	if v.childOccupancy.Max == nil && len(children) > 0 {
		return ErrChildrenNotAllowed
	}
	for _, child := range children {
		if child.MinAge == nil && child.MaxAge == nil {
			return ErrMissingMinAge
		}

		if child.MinAge != nil && child.MaxAge != nil && *child.MinAge >= *child.MaxAge {
			return ErrMinAgeGreaterThanOrEqualsThanMaxAge
		}

		if v.childOccupancy.MinAge != nil && child.MinAge != nil && *child.MinAge < *v.childOccupancy.MinAge {
			return ErrMinAgeOutOfRange(*child.MinAge, *v.childOccupancy.Min)
		}

		if v.adultOccupancy.MinAge != nil && child.MaxAge != nil && *child.MaxAge >= *v.adultOccupancy.MinAge {
			return ErrMaxAgeOutOfRange(*child.MaxAge, *v.adultOccupancy.MinAge)
		}

		if child.Amount == nil {
			return ErrMissingAmount
		}
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateDateDependingRateOverlaps(rates []Rate) error {
	ratesByInvTypeCode := slicesx.GroupByFunc(rates, func(r Rate) string {
		return r.InvTypeCode
	})
	for _, rates := range ratesByInvTypeCode {
		if err := validateOverlaps(rates); err != nil {
			return err
		}
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateSupplements(supplements []Supplement) error {
	if !v.supportsSupplements && len(supplements) > 0 {
		return ErrSupplementsNotSupported
	}

	staticSupplements := slicesx.Filter(supplements, Supplement.isStaticSupplement)
	if err := v.validateStaticSupplements(staticSupplements); err != nil {
		return err
	}

	dateDependingSupplements := slicesx.Filter(supplements, Supplement.isDateDependingSupplement)
	if err := v.validateDateDependingSupplements(dateDependingSupplements); err != nil {
		return err
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateStaticSupplements(supplements []Supplement) error {
	for _, supplement := range supplements {
		if err := v.validateStaticSupplement(supplement); err != nil {
			return err
		}
		v.supplementMapping[supplement.InvCode] = struct{}{}
	}
	return nil
}

func (v HotelRatePlanNotifValidator) validateStaticSupplement(supplement Supplement) error {
	if supplement.AddToBasicRateIndicator == nil {
		return ErrMissingAddToBasicRateIndicator
	}

	if supplement.MandatoryIndicator == nil {
		return ErrMissingMandatoryIndicator
	}

	if supplement.ChargeTypeCode == nil {
		return ErrMissingChargeTypeCode
	}

	if p := supplement.PrerequisiteInventory; p != nil {
		switch p.InvType {
		case PrerequisiteInventoryInvTypeAlpineBitsDOW:
			match, _ := regexp.MatchString("[0-1]{7}", p.InvCode)
			if !match {
				return ErrInvalidDOWString
			}
		default:
			return ErrInvalidInvType(string(p.InvType))
		}
	}

	if d := supplement.Descriptions; d != nil {
		if err := v.validateDescriptions(*d); err != nil {
			return err
		}
	}

	if supplement.Amount != nil {
		return ErrUnexpectedAmount
	}

	if supplement.Start != nil {
		return ErrUnexpectedStart
	}

	if supplement.End != nil {
		return ErrUnexpectedEnd
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateDateDependingSupplements(supplements []Supplement) error {
	for _, supplement := range supplements {
		if err := v.validateDateDependingSupplement(supplement); err != nil {
			return err
		}
	}

	if err := v.validateDateDependingSupplementsOverlaps(supplements); err != nil {
		return err
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateDateDependingSupplement(supplement Supplement) error {
	if err := validateString(supplement.InvCode); err != nil {
		return ErrMissingInvCode
	}

	if _, ok := v.supplementMapping[supplement.InvCode]; !ok {
		return ErrInvCodeNotFound(supplement.InvCode)
	}

	if supplement.Start == nil {
		return ErrMissingStart
	}

	if supplement.End == nil {
		return ErrMissingEnd
	}

	if supplement.Start.After(*supplement.End) {
		return ErrStartAfterEnd
	}

	if p := supplement.PrerequisiteInventory; p != nil {
		switch p.InvType {
		case PrerequisiteInventoryInvTypeRoomType:
			if _, ok := v.roomTypeMapping[p.InvCode]; !ok {
				return ErrInvCodeNotFound(p.InvCode)
			}
		default:
			return ErrInvalidInvType(string(p.InvType))
		}
	}

	if supplement.AddToBasicRateIndicator != nil {
		return ErrUnexpectedAddToBasicRateIndicator
	}

	if supplement.MandatoryIndicator != nil {
		return ErrUnexpectedAddToBasicRateIndicator
	}

	if supplement.ChargeTypeCode != nil {
		return ErrUnexpectedAddToBasicRateIndicator
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateDateDependingSupplementsOverlaps(supplements []Supplement) error {
	supplementsByInvCode := slicesx.GroupByFunc(supplements, func(s Supplement) string {
		if s.PrerequisiteInventory == nil {
			return ""
		}
		return s.PrerequisiteInventory.InvCode
	})
	for _, supplements := range supplementsByInvCode {
		if err := validateOverlaps(supplements); err != nil {
			return err
		}
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateRatePlanOverlay(ratePlan RatePlan) error {
	if !v.supportsOverlay {
		return ErrDeltasNotSupported
	}

	if err := v.validateBookingRules(ratePlan.BookingRules); err != nil {
		return err
	}

	if err := v.validateDateDependingRates(ratePlan.Rates); err != nil {
		return err
	}

	if err := v.validateDateDependingSupplements(ratePlan.Supplements); err != nil {
		return err
	}

	if len(ratePlan.Offers) == 0 {
		return ErrUnexpectedOffers
	}

	if !ratePlan.Descriptions.isZero() {
		return ErrUnexpectedDescription
	}

	return nil
}

func (v HotelRatePlanNotifValidator) validateRatePlanRemove(ratePlan RatePlan) error {
	if len(ratePlan.Offers) == 0 {
		return ErrUnexpectedOffers
	}

	if !ratePlan.Descriptions.isZero() {
		return ErrUnexpectedDescription
	}

	if len(ratePlan.BookingRules) > 0 {
		return ErrUnexpectedBookingRules
	}

	if len(ratePlan.Rates) > 0 {
		return ErrUnexpectedRates
	}

	if len(ratePlan.Supplements) > 0 {
		return ErrUnexpectedSupplements
	}

	return nil
}
