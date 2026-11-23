package freerooms

import freerooms202410 "github.com/HGV/alpinebits/v202410/freerooms"

// A server should accept deltas, complete sets or both (AlpineBits 2026-10
// §4.1).
const (
	CapAcceptRooms          = freerooms202410.CapAcceptRooms
	CapAcceptCategories     = freerooms202410.CapAcceptCategories
	CapAcceptDeltas         = freerooms202410.CapAcceptDeltas
	CapAcceptOutOfOrder     = freerooms202410.CapAcceptOutOfOrder
	CapAcceptOutOfMarket    = freerooms202410.CapAcceptOutOfMarket
	CapAcceptClosingSeasons = freerooms202410.CapAcceptClosingSeasons
	CapAcceptCompleteSet    = freerooms202410.CapAcceptCompleteSet
)
