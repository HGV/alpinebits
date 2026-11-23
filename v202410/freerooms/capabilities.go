package freerooms

import freerooms202210 "github.com/HGV/alpinebits/v202210/freerooms"

// A server should accept deltas, complete sets or both (AlpineBits 2024-10
// §4.1).
const (
	CapAcceptRooms          = freerooms202210.CapAcceptRooms
	CapAcceptCategories     = freerooms202210.CapAcceptCategories
	CapAcceptDeltas         = freerooms202210.CapAcceptDeltas
	CapAcceptOutOfOrder     = freerooms202210.CapAcceptOutOfOrder
	CapAcceptOutOfMarket    = freerooms202210.CapAcceptOutOfMarket
	CapAcceptClosingSeasons = freerooms202210.CapAcceptClosingSeasons
	CapAcceptCompleteSet    = freerooms202210.CapAcceptCompleteSet
)
