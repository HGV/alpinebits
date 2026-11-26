package freerooms

import freerooms202210 "github.com/HGV/alpinebits/v202210/freerooms"

// A server should accept deltas, complete sets or both (AlpineBits 2024-10
// §4.1).
const (
	CapabilityAcceptRooms          = freerooms202210.CapabilityAcceptRooms
	CapabilityAcceptCategories     = freerooms202210.CapabilityAcceptCategories
	CapabilityAcceptDeltas         = freerooms202210.CapabilityAcceptDeltas
	CapabilityAcceptOutOfOrder     = freerooms202210.CapabilityAcceptOutOfOrder
	CapabilityAcceptOutOfMarket    = freerooms202210.CapabilityAcceptOutOfMarket
	CapabilityAcceptClosingSeasons = freerooms202210.CapabilityAcceptClosingSeasons
	CapabilityAcceptCompleteSet    = freerooms202210.CapabilityAcceptCompleteSet
)
