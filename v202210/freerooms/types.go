// Package freerooms holds the AlpineBits 2022-10 FreeRooms messages
// (OTA_HotelInvCountNotif, §4.1). Their structure is unchanged from 2020-10, so
// the types are aliases; the stricter value limits are enforced by the XSD.
package freerooms

import freerooms202010 "github.com/HGV/alpinebits/v202010/freerooms"

type (
	HotelInvCountNotifRQ     = freerooms202010.HotelInvCountNotifRQ
	HotelInvCountNotifRS     = freerooms202010.HotelInvCountNotifRS
	UniqueID                 = freerooms202010.UniqueID
	UniqueIDType             = freerooms202010.UniqueIDType
	Instance                 = freerooms202010.Instance
	Inventories              = freerooms202010.Inventories
	Inventory                = freerooms202010.Inventory
	StatusApplicationControl = freerooms202010.StatusApplicationControl
	InvCount                 = freerooms202010.InvCount
	HotelInventory           = freerooms202010.HotelInventory
	CountType                = freerooms202010.CountType
)

const (
	UniqueIDTypeReference             = freerooms202010.UniqueIDTypeReference
	UniqueIDTypePurgedMasterReference = freerooms202010.UniqueIDTypePurgedMasterReference
	InstanceCompleteSet               = freerooms202010.InstanceCompleteSet
	CountTypeBookable                 = freerooms202010.CountTypeBookable
	CountTypeOutOfOrder               = freerooms202010.CountTypeOutOfOrder
	CountTypeOutOfMarket              = freerooms202010.CountTypeOutOfMarket
)
