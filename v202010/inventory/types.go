// Package inventory holds the AlpineBits 2020-10 Inventory/Basic messages
// (OTA_HotelDescriptiveContentNotif:Inventory, §4.4.1). Their structure and
// rules are unchanged from 2018-10, so the types are aliases.
package inventory

import inventory201810 "github.com/HGV/alpinebits/v201810/inventory"

type (
	HotelDescriptiveContentNotifRQ = inventory201810.HotelDescriptiveContentNotifRQ
	HotelDescriptiveContentNotifRS = inventory201810.HotelDescriptiveContentNotifRS
	HotelDescriptiveContent        = inventory201810.HotelDescriptiveContent
	FacilityInfo                   = inventory201810.FacilityInfo
	GuestRoom                      = inventory201810.GuestRoom
	TypeRoom                       = inventory201810.TypeRoom
	RoomType                       = inventory201810.RoomType
	Amenity                        = inventory201810.Amenity
	InfoCode                       = inventory201810.InfoCode
	MultimediaDescription          = inventory201810.MultimediaDescription
	TextItem                       = inventory201810.TextItem
	ImageItem                      = inventory201810.ImageItem
	ImageFormat                    = inventory201810.ImageFormat
)

const (
	RoomTypeRoom                 = inventory201810.RoomTypeRoom
	RoomTypeApartment            = inventory201810.RoomTypeApartment
	RoomTypeMobileHome           = inventory201810.RoomTypeMobileHome
	RoomTypeBungalow             = inventory201810.RoomTypeBungalow
	RoomTypeHolidayHome          = inventory201810.RoomTypeHolidayHome
	RoomTypeCampingGround        = inventory201810.RoomTypeCampingGround
	RoomTypePitch                = inventory201810.RoomTypePitch
	RoomTypeCampingGroundOrPitch = inventory201810.RoomTypeCampingGroundOrPitch
	RoomTypeRestingPlace         = inventory201810.RoomTypeRestingPlace

	InfoCodeDescription = inventory201810.InfoCodeDescription
	InfoCodePictures    = inventory201810.InfoCodePictures
	InfoCodeLongName    = inventory201810.InfoCodeLongName
)
