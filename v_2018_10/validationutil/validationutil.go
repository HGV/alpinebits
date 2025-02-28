package validationutil

import (
	"slices"

	"github.com/HGV/alpinebits/v_2018_10"
	"github.com/HGV/alpinebits/v_2018_10/freerooms"
	"github.com/HGV/alpinebits/v_2018_10/inventory"
	"github.com/HGV/alpinebits/v_2018_10/rateplans"
)

func NewFreeRoomOptions(capabilities []string) []freerooms.HotelAvailNotifValidatorFunc {
	var options []freerooms.HotelAvailNotifValidatorFunc

	capabilityMap := map[v_2018_10.Capability]func() freerooms.HotelAvailNotifValidatorFunc{
		v_2018_10.CapabilityHotelAvailNotifAcceptRooms:            freerooms.WithRooms,
		v_2018_10.CapabilityHotelAvailNotifAcceptCategories:       freerooms.WithCategories,
		v_2018_10.CapabilityHotelAvailNotifAcceptDeltas:           freerooms.WithDeltas,
		v_2018_10.CapabilityHotelAvailNotifAcceptBookingThreshold: freerooms.WithBookingThreshold,
	}

	for cap, fn := range capabilityMap {
		if slices.Contains(capabilities, string(cap)) {
			options = append(options, fn())
		}
	}

	return options
}

func NewInventoryOptions(capabilities []string) []inventory.HotelDescriptiveContentNotifValidatorFunc {
	var options []inventory.HotelDescriptiveContentNotifValidatorFunc

	capabilityMap := map[v_2018_10.Capability]func() inventory.HotelDescriptiveContentNotifValidatorFunc{
		v_2018_10.CapabilityHotelDescriptiveContentNotifInventoryUseRooms:          inventory.WithRooms,
		v_2018_10.CapabilityHotelDescriptiveContentNotifInventoryOccupancyChildren: inventory.WithOccupancyChildren,
	}

	for cap, fn := range capabilityMap {
		if slices.Contains(capabilities, string(cap)) {
			options = append(options, fn())
		}
	}

	return options
}

func NewRatePlanOptions(capabilities []string) []rateplans.HotelRatePlanNotifValidatorFunc {
	var options []rateplans.HotelRatePlanNotifValidatorFunc

	capabilityMap := map[v_2018_10.Capability]func() rateplans.HotelRatePlanNotifValidatorFunc{
		v_2018_10.CapabilityHotelRatePlanNotifAcceptArrivalDOW:                  rateplans.WithArrivalDOW,
		v_2018_10.CapabilityHotelRatePlanNotifAcceptDepartureDOW:                rateplans.WithDepartureDOW,
		v_2018_10.CapabilityHotelRatePlanNotifAcceptRatePlanBookingRule:         rateplans.WithGenericBookingRules,
		v_2018_10.CapabilityHotelRatePlanNotifAcceptRatePlanRoomTypeBookingRule: rateplans.WithRoomTypeBookingRules,
		v_2018_10.CapabilityHotelRatePlanNotifAcceptSupplements:                 rateplans.WithSupplements,
		v_2018_10.CapabilityHotelRatePlanNotifAcceptFreeNightsOffers:            rateplans.WithFreeNightOffer,
		v_2018_10.CapabilityHotelRatePlanNotifAcceptFamilyOffers:                rateplans.WithFamilyOffer,
		v_2018_10.CapabilityHotelRatePlanNotifAcceptOverlay:                     rateplans.WithOverlay,
		v_2018_10.CapabilityHotelRatePlanNotifAcceptRatePlanJoin:                rateplans.WithRatePlanJoin,
		v_2018_10.CapabilityHotelRatePlanNotifAcceptOfferRuleBookingOffset:      rateplans.WithOfferRuleBookingOffset,
		v_2018_10.CapabilityHotelRatePlanNotifAcceptOfferRuleDOWLOS:             rateplans.WithOfferRuleDOWLOS,
	}

	for cap, fn := range capabilityMap {
		if slices.Contains(capabilities, string(cap)) {
			options = append(options, fn())
		}
	}

	return options
}
