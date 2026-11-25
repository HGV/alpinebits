package freerooms

import (
	"github.com/HGV/alpinebits"
	freerooms202010 "github.com/HGV/alpinebits/v202010/freerooms"
)

// A server should accept deltas, complete sets or both (AlpineBits 2022-10
// §4.1).
const (
	CapAcceptRooms          = freerooms202010.CapAcceptRooms
	CapAcceptCategories     = freerooms202010.CapAcceptCategories
	CapAcceptDeltas         = freerooms202010.CapAcceptDeltas
	CapAcceptOutOfOrder     = freerooms202010.CapAcceptOutOfOrder
	CapAcceptOutOfMarket    = freerooms202010.CapAcceptOutOfMarket
	CapAcceptClosingSeasons = freerooms202010.CapAcceptClosingSeasons

	CapAcceptCompleteSet alpinebits.Capability = "OTA_HotelInvCountNotif_accept_complete_set"
)
