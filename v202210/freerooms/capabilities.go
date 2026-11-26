package freerooms

import (
	"github.com/HGV/alpinebits"
	freerooms202010 "github.com/HGV/alpinebits/v202010/freerooms"
)

// A server should accept deltas, complete sets or both (AlpineBits 2022-10
// §4.1).
const (
	CapabilityAcceptRooms          = freerooms202010.CapabilityAcceptRooms
	CapabilityAcceptCategories     = freerooms202010.CapabilityAcceptCategories
	CapabilityAcceptDeltas         = freerooms202010.CapabilityAcceptDeltas
	CapabilityAcceptOutOfOrder     = freerooms202010.CapabilityAcceptOutOfOrder
	CapabilityAcceptOutOfMarket    = freerooms202010.CapabilityAcceptOutOfMarket
	CapabilityAcceptClosingSeasons = freerooms202010.CapabilityAcceptClosingSeasons

	CapabilityAcceptCompleteSet alpinebits.Capability = "OTA_HotelInvCountNotif_accept_complete_set"
)
