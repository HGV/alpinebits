// Package inventory holds the AlpineBits 2024-10 Inventory/Basic messages
// (OTA_HotelDescriptiveContentNotif:Inventory, §4.3.1). Their structure and
// rules are unchanged from 2022-10, so the types are aliases; RoomAmenityCode
// may now be an ABR code, and the stricter value limits are enforced by the
// XSD.
package inventory

import inventory202210 "github.com/HGV/alpinebits/v202210/inventory"

type (
	HotelDescriptiveContentNotifRQ = inventory202210.HotelDescriptiveContentNotifRQ
	HotelDescriptiveContentNotifRS = inventory202210.HotelDescriptiveContentNotifRS
	HotelDescriptiveContent        = inventory202210.HotelDescriptiveContent
	FacilityInfo                   = inventory202210.FacilityInfo
	GuestRoom                      = inventory202210.GuestRoom
	TypeRoom                       = inventory202210.TypeRoom
	RoomType                       = inventory202210.RoomType
	Amenity                        = inventory202210.Amenity
	InfoCode                       = inventory202210.InfoCode
	MultimediaDescription          = inventory202210.MultimediaDescription
	TextItem                       = inventory202210.TextItem
	ImageItem                      = inventory202210.ImageItem
	ImageFormat                    = inventory202210.ImageFormat
)

const (
	RoomTypeRoom                 = inventory202210.RoomTypeRoom
	RoomTypeApartment            = inventory202210.RoomTypeApartment
	RoomTypeMobileHome           = inventory202210.RoomTypeMobileHome
	RoomTypeBungalow             = inventory202210.RoomTypeBungalow
	RoomTypeHolidayHome          = inventory202210.RoomTypeHolidayHome
	RoomTypeCampingGround        = inventory202210.RoomTypeCampingGround
	RoomTypePitch                = inventory202210.RoomTypePitch
	RoomTypeCampingGroundOrPitch = inventory202210.RoomTypeCampingGroundOrPitch
	RoomTypeRestingPlace         = inventory202210.RoomTypeRestingPlace

	InfoCodeDescription = inventory202210.InfoCodeDescription
	InfoCodePictures    = inventory202210.InfoCodePictures
	InfoCodeLongName    = inventory202210.InfoCodeLongName
)
