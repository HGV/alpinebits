// Package inventory holds the AlpineBits 2022-10 Inventory/Basic messages
// (OTA_HotelDescriptiveContentNotif:Inventory, §4.4.1). Their structure and
// rules are unchanged from 2020-10, so the types are aliases; the stricter
// value limits are enforced by the XSD.
package inventory

import inventory202010 "github.com/HGV/alpinebits/v202010/inventory"

type (
	HotelDescriptiveContentNotifRQ = inventory202010.HotelDescriptiveContentNotifRQ
	HotelDescriptiveContentNotifRS = inventory202010.HotelDescriptiveContentNotifRS
	HotelDescriptiveContent        = inventory202010.HotelDescriptiveContent
	FacilityInfo                   = inventory202010.FacilityInfo
	GuestRoom                      = inventory202010.GuestRoom
	TypeRoom                       = inventory202010.TypeRoom
	RoomType                       = inventory202010.RoomType
	Amenity                        = inventory202010.Amenity
	InfoCode                       = inventory202010.InfoCode
	MultimediaDescription          = inventory202010.MultimediaDescription
	TextItem                       = inventory202010.TextItem
	ImageItem                      = inventory202010.ImageItem
	ImageFormat                    = inventory202010.ImageFormat
)

const (
	RoomTypeRoom                 = inventory202010.RoomTypeRoom
	RoomTypeApartment            = inventory202010.RoomTypeApartment
	RoomTypeMobileHome           = inventory202010.RoomTypeMobileHome
	RoomTypeBungalow             = inventory202010.RoomTypeBungalow
	RoomTypeHolidayHome          = inventory202010.RoomTypeHolidayHome
	RoomTypeCampingGround        = inventory202010.RoomTypeCampingGround
	RoomTypePitch                = inventory202010.RoomTypePitch
	RoomTypeCampingGroundOrPitch = inventory202010.RoomTypeCampingGroundOrPitch
	RoomTypeRestingPlace         = inventory202010.RoomTypeRestingPlace

	InfoCodeDescription = inventory202010.InfoCodeDescription
	InfoCodePictures    = inventory202010.InfoCodePictures
	InfoCodeLongName    = inventory202010.InfoCodeLongName
)
