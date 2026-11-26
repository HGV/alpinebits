package freerooms

import freerooms202410 "github.com/HGV/alpinebits/v202410/freerooms"

// A server should accept deltas, complete sets or both (AlpineBits 2026-10
// §4.1).
const (
	CapabilityAcceptRooms          = freerooms202410.CapabilityAcceptRooms
	CapabilityAcceptCategories     = freerooms202410.CapabilityAcceptCategories
	CapabilityAcceptDeltas         = freerooms202410.CapabilityAcceptDeltas
	CapabilityAcceptOutOfOrder     = freerooms202410.CapabilityAcceptOutOfOrder
	CapabilityAcceptOutOfMarket    = freerooms202410.CapabilityAcceptOutOfMarket
	CapabilityAcceptClosingSeasons = freerooms202410.CapabilityAcceptClosingSeasons
	CapabilityAcceptCompleteSet    = freerooms202410.CapabilityAcceptCompleteSet
)
