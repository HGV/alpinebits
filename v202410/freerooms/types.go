// Package freerooms holds the AlpineBits 2024-10 FreeRooms messages. Their
// structure is unchanged from 2022-10, so the types are aliases; the stricter
// value limits are enforced by the XSD. The rules are not implemented yet.
package freerooms

import freerooms202210 "github.com/HGV/alpinebits/v202210/freerooms"

type (
	HotelInvCountNotifRQ     = freerooms202210.HotelInvCountNotifRQ
	HotelInvCountNotifRS     = freerooms202210.HotelInvCountNotifRS
	UniqueID                 = freerooms202210.UniqueID
	UniqueIDType             = freerooms202210.UniqueIDType
	Instance                 = freerooms202210.Instance
	Inventories              = freerooms202210.Inventories
	Inventory                = freerooms202210.Inventory
	StatusApplicationControl = freerooms202210.StatusApplicationControl
	InvCount                 = freerooms202210.InvCount
	CountType                = freerooms202210.CountType
)

const (
	UniqueIDTypeReference             = freerooms202210.UniqueIDTypeReference
	UniqueIDTypePurgedMasterReference = freerooms202210.UniqueIDTypePurgedMasterReference
	InstanceCompleteSet               = freerooms202210.InstanceCompleteSet
	CountTypeBookable                 = freerooms202210.CountTypeBookable
	CountTypeOutOfOrder               = freerooms202210.CountTypeOutOfOrder
	CountTypeOutOfMarket              = freerooms202210.CountTypeOutOfMarket
)
