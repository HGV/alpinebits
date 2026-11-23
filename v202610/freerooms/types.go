// Package freerooms holds the AlpineBits 2026-10 FreeRooms messages. Their
// structure is unchanged from 2024-10, so the types are aliases; the stricter
// value limits are enforced by the XSD. The rules are not implemented yet.
package freerooms

import freerooms202410 "github.com/HGV/alpinebits/v202410/freerooms"

type (
	HotelInvCountNotifRQ     = freerooms202410.HotelInvCountNotifRQ
	HotelInvCountNotifRS     = freerooms202410.HotelInvCountNotifRS
	UniqueID                 = freerooms202410.UniqueID
	UniqueIDType             = freerooms202410.UniqueIDType
	Instance                 = freerooms202410.Instance
	Inventories              = freerooms202410.Inventories
	Inventory                = freerooms202410.Inventory
	StatusApplicationControl = freerooms202410.StatusApplicationControl
	InvCount                 = freerooms202410.InvCount
	CountType                = freerooms202410.CountType
)

const (
	UniqueIDTypeReference             = freerooms202410.UniqueIDTypeReference
	UniqueIDTypePurgedMasterReference = freerooms202410.UniqueIDTypePurgedMasterReference
	InstanceCompleteSet               = freerooms202410.InstanceCompleteSet
	CountTypeBookable                 = freerooms202410.CountTypeBookable
	CountTypeOutOfOrder               = freerooms202410.CountTypeOutOfOrder
	CountTypeOutOfMarket              = freerooms202410.CountTypeOutOfMarket
)
