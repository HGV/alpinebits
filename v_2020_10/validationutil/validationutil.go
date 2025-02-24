package validationutil

import (
	"slices"

	"github.com/HGV/alpinebits/v_2020_10"
	"github.com/HGV/alpinebits/v_2020_10/freerooms"
	"github.com/HGV/alpinebits/v_2020_10/inventory"
	"github.com/HGV/alpinebits/v_2020_10/rateplans"
)

func NewFreeRoomOptions(capabilities []string) []freerooms.HotelInvCountNotifValidatorFunc {
	var options []freerooms.HotelInvCountNotifValidatorFunc

	capabilityMap := map[v_2020_10.Capability]func() freerooms.HotelInvCountNotifValidatorFunc{
		v_2020_10.CapabilityHotelInvCountNotifAcceptRooms:          freerooms.WithRooms,
		v_2020_10.CapabilityHotelInvCountNotifAcceptRoomCategories: freerooms.WithCategories,
		v_2020_10.CapabilityHotelInvCountNotifAcceptDeltas:         freerooms.WithDeltas,
		v_2020_10.CapabilityHotelInvCountNotifAcceptOutOfOrder:     freerooms.WithOutOfOrder,
		v_2020_10.CapabilityHotelInvCountNotifAcceptOutOfMarket:    freerooms.WithOutOfMarket,
		v_2020_10.CapabilityHotelInvCountNotifAcceptClosingSeasons: freerooms.WithClosingSeasons,
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

	capabilityMap := map[v_2020_10.Capability]func() inventory.HotelDescriptiveContentNotifValidatorFunc{
		v_2020_10.CapabilityHotelDescriptiveContentNotifInventoryUseRooms:          inventory.WithRooms,
		v_2020_10.CapabilityHotelDescriptiveContentNotifInventoryOccupancyChildren: inventory.WithOccupancyChildren,
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

	capabilityMap := map[v_2020_10.Capability]func() rateplans.HotelRatePlanNotifValidatorFunc{
		v_2020_10.CapabilityHotelRatePlanNotifAcceptArrivalDOW:                  rateplans.WithArrivalDOW,
		v_2020_10.CapabilityHotelRatePlanNotifAcceptDepartureDOW:                rateplans.WithDepartureDOW,
		v_2020_10.CapabilityHotelRatePlanNotifAcceptRatePlanBookingRule:         rateplans.WithGenericBookingRules,
		v_2020_10.CapabilityHotelRatePlanNotifAcceptRatePlanRoomTypeBookingRule: rateplans.WithRoomTypeBookingRules,
		v_2020_10.CapabilityHotelRatePlanNotifAcceptSupplements:                 rateplans.WithSupplements,
		v_2020_10.CapabilityHotelRatePlanNotifAcceptFreeNightsOffers:            rateplans.WithFreeNightOffer,
		v_2020_10.CapabilityHotelRatePlanNotifAcceptFamilyOffers:                rateplans.WithFamilyOffer,
		v_2020_10.CapabilityHotelRatePlanNotifAcceptOverlay:                     rateplans.WithOverlay,
		v_2020_10.CapabilityHotelRatePlanNotifAcceptRatePlanJoin:                rateplans.WithRatePlanJoin,
		v_2020_10.CapabilityHotelRatePlanNotifAcceptOfferRuleBookingOffset:      rateplans.WithOfferRuleBookingOffset,
		v_2020_10.CapabilityHotelRatePlanNotifAcceptOfferRuleDOWLOS:             rateplans.WithOfferRuleDOWLOS,
	}

	for cap, fn := range capabilityMap {
		if slices.Contains(capabilities, string(cap)) {
			options = append(options, fn())
		}
	}

	return options
}
