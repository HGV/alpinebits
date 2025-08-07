package rateplans

import (
	"cmp"
	"math"
	"regexp"

	"github.com/HGV/alpinebits/internal"
	"github.com/HGV/alpinebits/v_2018_10/common"
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
	ratePlanMapping                map[string]map[string]struct{}
	supportsRatePlanJoin           bool
	adultOccupancy                 RatePlanOccupancySettings
	childOccupancy                 *RatePlanOccupancySettings
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
	ratePlanNotifType              RatePlanNotifType
}

var _ common.Validatable[HotelRatePlanNotifRQ] = (*HotelRatePlanNotifValidator)(nil)

type HotelRatePlanNotifValidatorFunc func(*HotelRatePlanNotifValidator)

func NewHotelRatePlanNotifValidator(opts ...HotelRatePlanNotifValidatorFunc) HotelRatePlanNotifValidator {
	v := HotelRatePlanNotifValidator{
		supplementMapping: map[string]struct{}{},
	}

	for _, opt := range opts {
		opt(&v)
	}
	return v
}

func WithArrivalDOW() HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsArrivalDOW = true
	}
}

func WithDepartureDOW() HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsDepartureDOW = true
	}
}

func WithRatePlanMapping(mapping map[string]map[string]struct{}) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.ratePlanMapping = mapping
	}
}

func WithAdultOccupancy(occupancy RatePlanOccupancySettings) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.adultOccupancy = occupancy
	}
}

func WithChildOccupancy(occupancy RatePlanOccupancySettings) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.childOccupancy = &occupancy
	}
}

func WithRatePlanJoin() HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsRatePlanJoin = true
	}
}

func WithOverlay() HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsOverlay = true
	}
}

func WithGenericBookingRules() HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsGenericBookingRules = true
	}
}

func WithRoomTypeBookingRules() HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsRoomTypeBokingRules = true
	}
}

func WithRoomTypeCodes(mapping map[string]RoomTypeOccupancySettings) HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.roomTypeMapping = mapping
	}
}

func WithSupplements() HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsSupplements = true
	}
}

func WithFreeNightOffer() HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsFreeNightOffer = true
	}
}

func WithFamilyOffer() HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsFamilyOffer = true
	}
}

func WithOfferRuleBookingOffset() HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsOfferRuleBookingOffset = true
	}
}

func WithOfferRuleDOWLOS() HotelRatePlanNotifValidatorFunc {
	return func(v *HotelRatePlanNotifValidator) {
		v.supportsOfferRuleDOWLOS = true
	}
}

func (v *HotelRatePlanNotifValidator) Validate(r HotelRatePlanNotifRQ) error {
	if err := common.ValidateHotelCode(r.RatePlans.HotelCode); err != nil {
		return err
	}

	if r.IsReset() {
		if err := v.validateRatePlansReset(r.RatePlans.RatePlans); err != nil {
			return err
		}
	} else {
		if err := v.validateRatePlans(r.RatePlans.RatePlans); err != nil {
			return err
		}
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateRatePlansReset(ratePlans []RatePlan) error {
	for _, ratePlan := range ratePlans {
		if err := v.validateRatePlanCode(ratePlan.RatePlanCode); err != nil {
			return err
		}
	}
	return nil
}

func (v *HotelRatePlanNotifValidator) validateRatePlans(ratePlans []RatePlan) error {
	for _, ratePlan := range ratePlans {
		if err := v.validateRatePlan(ratePlan); err != nil {
			return err
		}
	}
	return nil
}

func (v *HotelRatePlanNotifValidator) validateRatePlan(ratePlan RatePlan) error {
	if err := v.validateRatePlanCode(ratePlan.RatePlanCode); err != nil {
		return err
	}

	if err := v.validateCurrencyCode(ratePlan.CurrencyCode); err != nil {
		return err
	}

	if !v.supportsRatePlanJoin && !ratePlan.IsMaster() {
		return common.ErrRatePlanJoinNotSupported
	}

	switch v.ratePlanNotifType = ratePlan.RatePlanNotifType; v.ratePlanNotifType {
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

func (v *HotelRatePlanNotifValidator) validateRatePlanCode(code string) error {
	if err := common.ValidateString(code); err != nil {
		return common.ErrMissingRatePlanCode
	}
	return nil
}

func (v *HotelRatePlanNotifValidator) validateCurrencyCode(code string) error {
	if err := common.ValidateString(code); err != nil {
		return common.ErrMissingCurrencyCode
	}
	return nil
}

func (v *HotelRatePlanNotifValidator) validateRatePlanNew(ratePlan RatePlan) error {
	if ratePlan.IsMaster() {
		return v.validateRatePlanNewMaster(ratePlan)
	}
	return v.validateRatePlanNewDerived(ratePlan)
}

func (v *HotelRatePlanNotifValidator) validateRatePlanNewMaster(ratePlan RatePlan) error {
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

func (v *HotelRatePlanNotifValidator) validateRatePlanNewDerived(ratePlan RatePlan) error {
	code := cmp.Or(ratePlan.RatePlanID, ratePlan.RatePlanCode)
	if _, ok := v.ratePlanMapping[code]; !ok {
		return common.ErrRatePlanNotFound(code)
	}

	if err := v.validateBookingRules(ratePlan.BookingRules); err != nil {
		return err
	}

	if err := v.validateRates(ratePlan.Rates); err != nil {
		return err
	}

	if err := v.validateDateDependingSupplements(ratePlan.Supplements); err != nil {
		return err
	}

	if len(ratePlan.Offers) > 0 {
		return common.ErrUnexpectedOffers
	}

	if !ratePlan.Descriptions.isZero() {
		return common.ErrUnexpectedDescription
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateOffers(offers []Offer) error {
	if len(offers) == 0 {
		return common.ErrMissingOfferRule
	}

	if err := v.validateOfferRule(offers[0].OfferRule); err != nil {
		return err
	}

	if err := v.validateAdditionalOffers(offers[1:]); err != nil {
		return err
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateOfferRule(offerRule *OfferRule) error {
	if offerRule == nil {
		return common.ErrMissingOfferRule
	}

	if !v.supportsOfferRuleBookingOffset &&
		(offerRule.MinAdvancedBookingOffset != nil || offerRule.MaxAdvancedBookingOffset != nil) {
		return common.ErrOfferRuleBookingOffsetNotSupported
	}

	if len(offerRule.LengthsOfStay) > 0 ||
		offerRule.ArrivalDaysOfWeek != nil ||
		offerRule.DepartureDaysOfWeek != nil {
		return common.ErrOfferRuleDOWLOSNotSupported
	}

	if err := v.validateOfferRuleLengthsOfStay(offerRule); err != nil {
		return err
	}

	if err := v.validateOccupancies(offerRule.Occupancies); err != nil {
		return err
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateOfferRuleLengthsOfStay(offerRule *OfferRule) error {
	var minArrival int
	var maxArrival int
	for _, los := range offerRule.LengthsOfStay {
		switch los.MinMaxMessageType {
		case StayTypeMinArrival:
			minArrival = los.Time
		case StayTypeMaxArrival:
			maxArrival = los.Time
		case StayTypeMinThrough, StayTypeMaxThrough:
			return common.ErrStayThroughNotAllowedInOfferRule
		}
	}

	if maxArrival > 0 && minArrival > maxArrival {
		return common.ErrMinStayArrivalGratherThanMaxStayArrival(minArrival, maxArrival)
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateOccupancies(occupancies []Occupancy) error {
	adults := slicesx.Filter(occupancies, Occupancy.isAdult)
	switch len(adults) {
	case 0:
		return common.ErrMissingAdultOccupancy
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
		return common.ErrDuplicateChildOccupancy
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateOccupancy(o Occupancy) error {
	if min := o.MinOccupancy; min != nil && *min > 99 {
		return common.ErrInvalidMinOccupancy
	}
	if max := o.MaxOccupancy; max != nil && *max > 99 {
		return common.ErrInvalidMaxOccupancy
	}
	return nil
}

func (v *HotelRatePlanNotifValidator) populateAdultOccupancy(occupancy Occupancy) {
	v.adultOccupancy.Min = occupancy.MinOccupancy
	v.adultOccupancy.Max = occupancy.MaxOccupancy
	v.adultOccupancy.MinAge = occupancy.MinAge
}

func (v *HotelRatePlanNotifValidator) populateChildOccupancy(occupancy Occupancy) {
	v.childOccupancy = &RatePlanOccupancySettings{
		Min:    occupancy.MinOccupancy,
		Max:    occupancy.MaxOccupancy,
		MinAge: occupancy.MinAge,
	}
}

func (v *HotelRatePlanNotifValidator) validateAdditionalOffers(offers []Offer) error {
	freeNightOffers := slicesx.Filter(offers, Offer.IsFreeNightOffer)
	switch len(freeNightOffers) {
	case 0:
		break
	case 1:
		if err := v.validateFreeNightOffer(freeNightOffers[0]); err != nil {
			return err
		}
	default:
		return common.ErrDuplicateFreeNightOffer
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
		return common.ErrDuplicateFamilyOffer
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateFreeNightOffer(offer Offer) error {
	if !v.supportsFreeNightOffer {
		return common.ErrFreeNightOfferNotSupported
	}

	if offer.Discount.NightsRequired == 0 {
		return common.ErrMissingNightsRequired
	}

	if offer.Discount.NightsDiscounted == 0 {
		return common.ErrMissingNightsDiscounted
	}

	if pattern := offer.Discount.DiscountPattern; pattern != "" {
		expectedPattern := internal.CalculateDiscountPattern(
			offer.Discount.NightsRequired,
			offer.Discount.NightsDiscounted,
		)
		if pattern != expectedPattern {
			return common.ErrInvalidDiscountPattern
		}
	}

	if offer.Guest != nil {
		return common.ErrUnexpectedGuest
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateFamilyOffer(offer Offer) error {
	if !v.supportsFamilyOffer {
		return common.ErrFamilyOfferNotSupported
	}

	if offer.Guest.AgeQualifyingCode != AgeQualifyingCodeChild {
		return common.ErrInvalidGuestAgeQualifyngCode
	}

	if offer.Discount.NightsRequired > 0 {
		return common.ErrUnexpectedNightsRequired
	}

	if offer.Discount.NightsDiscounted > 0 {
		return common.ErrUnexpectedNightsDiscounted
	}

	if offer.Discount.DiscountPattern != "" {
		return common.ErrUnexpectedDiscountPattern
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateDescriptions(d RatePlanDescription) error {
	if err := common.ValidateLanguageUniqueness(d.Titles); err != nil {
		return err
	}

	if err := common.ValidateLanguageUniqueness(d.Intros); err != nil {
		return err
	}

	if err := common.ValidateLanguageUniqueness(d.Descriptions); err != nil {
		return err
	}

	for _, item := range d.Gallery {
		if err := common.ValidateLanguageUniqueness(item.Descriptions); err != nil {
			return err
		}
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateBookingRules(bookingRules []BookingRule) error {
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

func (v *HotelRatePlanNotifValidator) validateBookingRule(bookingRule BookingRule) error {
	if v.supportsRoomTypeBokingRules {
		if err := common.ValidateString(bookingRule.Code); err != nil {
			return common.ErrMissingCode
		}
		if _, ok := v.roomTypeMapping[bookingRule.Code]; !ok {
			return common.ErrInvTypeCodeNotFound(bookingRule.Code)
		}
	} else if v.supportsGenericBookingRules {
		if bookingRule.Code != "" || bookingRule.CodeContext != "" {
			return common.ErrRoomTypeBookingRulesNotSupported
		}
	}

	if bookingRule.Start.After(bookingRule.End) {
		return common.ErrStartAfterEnd
	}

	if err := v.validateLengthsOfStay(bookingRule.LengthsOfStay); err != nil {
		return err
	}

	if !v.supportsArrivalDOW && bookingRule.ArrivalDaysOfWeek != nil {
		return common.ErrArrivalDOWNotSupported
	}

	if !v.supportsDepartureDOW && bookingRule.DepartureDaysOfWeek != nil {
		return common.ErrDepartureDOWNotSupported
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateLengthsOfStay(lengthsOfStay []LengthOfStay) error {
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
		return common.ErrMinStayGratherThanMaxStay(min, max)
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateBookingRuleOverlaps(bookingRules []BookingRule) error {
	if v.supportsRoomTypeBokingRules {
		bookingRulesByRoomType := slicesx.GroupByFunc(bookingRules, func(b BookingRule) string {
			return b.Code
		})
		for _, brs := range bookingRulesByRoomType {
			if err := common.ValidateOverlaps(brs); err != nil {
				return err
			}
		}
	} else if v.supportsGenericBookingRules {
		if err := common.ValidateOverlaps(bookingRules); err != nil {
			return err
		}
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateRates(rates []Rate) error {
	if len(rates) == 0 {
		return common.ErrMissingStaticRate
	}

	if err := v.validateStaticRate(rates[0]); err != nil {
		return err
	}

	if err := v.validateDateDependingRates(rates[1:]); err != nil {
		return err
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateStaticRate(rate Rate) error {
	if rate.RateTimeUnit != nil && *rate.RateTimeUnit != TimeUnitDay {
		return common.ErrInvalidRateTimeUnit
	}

	switch len(rate.BaseByGuestAmts) {
	case 0:
		return common.ErrMissingBaseByGuestAmt
	case 1:
		b := rate.BaseByGuestAmts[0]
		if b.NumberOfGuests != nil {
			return common.ErrUnexpectedNumberOfGuests
		}
		if b.AgeQualifyingCode != nil {
			return common.ErrUnexpectedAgeQualifyingCode
		}
		if b.AmountAfterTax != nil {
			return common.ErrUnexpectedAmountAfterTax
		}
	default:
		return common.ErrUnexpectedBaseByGuestAmt
	}

	if rate.MealsIncluded == nil {
		return common.ErrMissingMealsIncluded
	}

	if rate.InvTypeCode != "" {
		return common.ErrUnexpectedInvTypeCode
	}

	if rate.Start != nil {
		return common.ErrUnexpectedStart
	}

	if rate.End != nil {
		return common.ErrUnexpectedEnd
	}

	if len(rate.AdditionalGuestAmounts) > 0 {
		return common.ErrUnexpectedAdditionalGuestAmounts
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateDateDependingRates(rates []Rate) error {
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

func (v *HotelRatePlanNotifValidator) validateDateDependingRate(rate Rate) error {
	if err := common.ValidateString(rate.InvTypeCode); err != nil {
		return common.ErrMissingInvTypeCode
	}

	roomTypeOccupancySettings, ok := v.roomTypeMapping[rate.InvTypeCode]
	if !ok {
		return common.ErrInvTypeCodeNotFound(rate.InvTypeCode)
	}

	if rate.Start == nil {
		return common.ErrMissingStart
	}

	if rate.End == nil {
		return common.ErrMissingEnd
	}

	if rate.Start.After(*rate.End) {
		return common.ErrStartAfterEnd
	}

	if err := v.validateBaseByGuestAmts(rate.BaseByGuestAmts, roomTypeOccupancySettings); err != nil {
		return err
	}

	if err := v.validateAdditionalGuestAmounts(rate.AdditionalGuestAmounts); err != nil {
		return err
	}

	if rate.RateTimeUnit != nil {
		return common.ErrUnexpectedRateTimeUnit
	}

	if rate.UnitMultiplier > 0 {
		return common.ErrUnexpectedUnitMultiplier
	}

	if rate.MealsIncluded != nil {
		return common.ErrUnexpectedMealsIncluded
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateBaseByGuestAmts(baseByGuestAmts []BaseByGuestAmt, roomTypeOccupancySettings RoomTypeOccupancySettings) error {
	numberOfGuestSeen := make(map[int]struct{})
	stdOccupancySeen := false
	for _, baseByGuestAmt := range baseByGuestAmts {
		if err := v.validateBaseByGuestAmt(baseByGuestAmt); err != nil {
			return err
		}

		numberOfGuests := *baseByGuestAmt.NumberOfGuests
		if _, exists := numberOfGuestSeen[numberOfGuests]; exists {
			return common.ErrDuplicateBaseByGuestAmt(numberOfGuests)
		}
		numberOfGuestSeen[numberOfGuests] = struct{}{}

		if numberOfGuests == roomTypeOccupancySettings.Std {
			stdOccupancySeen = true
		}
	}

	isStdOccupancyRequired := v.ratePlanNotifType == RatePlanNotifTypeNew
	if isStdOccupancyRequired && !stdOccupancySeen {
		return common.ErrMissingBaseByGuestAmtWithStdOccupancy(roomTypeOccupancySettings.Std)
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateBaseByGuestAmt(baseByGuestAmt BaseByGuestAmt) error {
	if baseByGuestAmt.NumberOfGuests == nil {
		return common.ErrMissingNumberOfGuests
	}

	if baseByGuestAmt.AgeQualifyingCode == nil {
		return common.ErrMissingAgeQualifyingCode
	}

	if baseByGuestAmt.AmountAfterTax == nil {
		return common.ErrMissingAmountAfterTax
	}

	if baseByGuestAmt.Type != nil {
		return common.ErrUnexpectedType
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateAdditionalGuestAmounts(additionalGuestAmounts []AdditionalGuestAmount) error {
	adults := slicesx.Filter(additionalGuestAmounts, AdditionalGuestAmount.IsAdult)
	switch len(adults) {
	case 0:
		break
	case 1:
		if adults[0].Amount == nil {
			return common.ErrMissingAmount
		}
	default:
		return common.ErrDuplicateAdditionalGuestAmountAdult
	}

	children := slicesx.Filter(additionalGuestAmounts, AdditionalGuestAmount.IsChild)
	if v.childOccupancy == nil && len(children) > 0 {
		return common.ErrChildrenNotAllowed
	}
	for _, child := range children {
		if child.MinAge == nil && child.MaxAge == nil {
			return common.ErrMissingMinAge
		}

		if child.MinAge != nil && child.MaxAge != nil && *child.MinAge >= *child.MaxAge {
			return common.ErrMinAgeGreaterThanOrEqualsThanMaxAge
		}

		if v.childOccupancy.MinAge != nil && child.MinAge != nil && *child.MinAge < *v.childOccupancy.MinAge {
			return common.ErrMinAgeOutOfRange(*child.MinAge, *v.childOccupancy.Min)
		}

		if v.adultOccupancy.MinAge != nil && child.MaxAge != nil && *child.MaxAge > *v.adultOccupancy.MinAge {
			return common.ErrMaxAgeOutOfRange(*child.MaxAge, *v.adultOccupancy.MinAge)
		}

		if child.Amount == nil {
			return common.ErrMissingAmount
		}
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateDateDependingRateOverlaps(rates []Rate) error {
	ratesByInvTypeCode := slicesx.GroupByFunc(rates, func(r Rate) string {
		return r.InvTypeCode
	})
	for _, rates := range ratesByInvTypeCode {
		if err := common.ValidateOverlaps(rates); err != nil {
			return err
		}
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateSupplements(supplements []Supplement) error {
	if !v.supportsSupplements && len(supplements) > 0 {
		return common.ErrSupplementsNotSupported
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

func (v *HotelRatePlanNotifValidator) validateStaticSupplement(supplement Supplement) error {
	if supplement.AddToBasicRateIndicator == nil {
		return common.ErrMissingAddToBasicRateIndicator
	}

	if supplement.MandatoryIndicator == nil {
		return common.ErrMissingMandatoryIndicator
	}

	if supplement.ChargeTypeCode == nil {
		return common.ErrMissingChargeTypeCode
	}

	if p := supplement.PrerequisiteInventory; p != nil {
		switch p.InvType {
		case PrerequisiteInventoryInvTypeAlpineBitsDOW:
			match, _ := regexp.MatchString("[0-1]{7}", p.InvCode)
			if !match {
				return common.ErrInvalidDOWString
			}
		default:
			return common.ErrInvalidInvType(string(p.InvType))
		}
	}

	if d := supplement.Descriptions; d != nil {
		if err := v.validateDescriptions(*d); err != nil {
			return err
		}
	}

	if supplement.Amount != nil {
		return common.ErrUnexpectedAmount
	}

	if supplement.Start != nil {
		return common.ErrUnexpectedStart
	}

	if supplement.End != nil {
		return common.ErrUnexpectedEnd
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateDateDependingSupplements(supplements []Supplement) error {
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

func (v *HotelRatePlanNotifValidator) validateDateDependingSupplement(supplement Supplement) error {
	if err := common.ValidateString(supplement.InvCode); err != nil {
		return common.ErrMissingInvCode
	}

	if _, ok := v.supplementMapping[supplement.InvCode]; !ok {
		return common.ErrInvCodeNotFound(supplement.InvCode)
	}

	if supplement.Start == nil {
		return common.ErrMissingStart
	}

	if supplement.End == nil {
		return common.ErrMissingEnd
	}

	if supplement.Start.After(*supplement.End) {
		return common.ErrStartAfterEnd
	}

	if p := supplement.PrerequisiteInventory; p != nil {
		switch p.InvType {
		case PrerequisiteInventoryInvTypeRoomType:
			if _, ok := v.roomTypeMapping[p.InvCode]; !ok {
				return common.ErrInvCodeNotFound(p.InvCode)
			}
		default:
			return common.ErrInvalidInvType(string(p.InvType))
		}
	}

	if supplement.AddToBasicRateIndicator != nil {
		return common.ErrUnexpectedAddToBasicRateIndicator
	}

	if supplement.MandatoryIndicator != nil {
		return common.ErrUnexpectedAddToBasicRateIndicator
	}

	if supplement.ChargeTypeCode != nil {
		return common.ErrUnexpectedAddToBasicRateIndicator
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateDateDependingSupplementsOverlaps(supplements []Supplement) error {
	type key struct {
		SupplementCode string
		InvTypeCode    string
	}
	supplementsByInvCode := slicesx.GroupByFunc(supplements, func(s Supplement) key {
		key := key{
			SupplementCode: s.InvCode,
		}
		if s.PrerequisiteInventory != nil {
			switch s.PrerequisiteInventory.InvType {
			case PrerequisiteInventoryInvTypeRoomType:
				key.InvTypeCode = s.PrerequisiteInventory.InvCode
			}
		}
		return key
	})
	for _, supplements := range supplementsByInvCode {
		if err := common.ValidateOverlaps(supplements); err != nil {
			return err
		}
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateRatePlanOverlay(ratePlan RatePlan) error {
	if !v.supportsOverlay {
		return common.ErrDeltasNotSupported
	}

	if ratePlan.RatePlanID != "" {
		if _, ok := v.ratePlanMapping[ratePlan.RatePlanID]; !ok {
			return common.ErrRatePlanNotFound(ratePlan.RatePlanID)
		}
	}

	mealPlanSeen := false
	for _, v := range v.ratePlanMapping {
		if _, ok := v[ratePlan.RatePlanCode]; ok {
			mealPlanSeen = true
			break
		}
	}
	if !mealPlanSeen {
		return common.ErrRatePlanNotFound(ratePlan.RatePlanCode)
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

	if len(ratePlan.Offers) > 0 {
		return common.ErrUnexpectedOffers
	}

	if !ratePlan.Descriptions.isZero() {
		return common.ErrUnexpectedDescription
	}

	return nil
}

func (v *HotelRatePlanNotifValidator) validateRatePlanRemove(ratePlan RatePlan) error {
	if len(ratePlan.Offers) > 0 {
		return common.ErrUnexpectedOffers
	}

	if !ratePlan.Descriptions.isZero() {
		return common.ErrUnexpectedDescription
	}

	if len(ratePlan.BookingRules) > 0 {
		return common.ErrUnexpectedBookingRules
	}

	if len(ratePlan.Rates) > 0 {
		return common.ErrUnexpectedRates
	}

	if len(ratePlan.Supplements) > 0 {
		return common.ErrUnexpectedSupplements
	}

	return nil
}
