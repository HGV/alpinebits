// Package inventory holds the AlpineBits 2026-10 Inventory/Basic messages
// (OTA_HotelDescriptiveContentNotif:Inventory, §4.3.1). GuestRoom gained
// MaxAdultOccupancy and HotelDescriptiveContent AdditionalCodes, so
// HotelDescriptiveContentNotifRQ and the types leading to them are its own;
// the others are aliases of 2024-10.
package inventory

import (
	"encoding/xml"

	"github.com/HGV/alpinebits/internal/rules"
	inventory202410 "github.com/HGV/alpinebits/v202410/inventory"
)

type (
	HotelDescriptiveContentNotifRS = inventory202410.HotelDescriptiveContentNotifRS
	TypeRoom                       = inventory202410.TypeRoom
	RoomType                       = inventory202410.RoomType
	Amenity                        = inventory202410.Amenity
	InfoCode                       = inventory202410.InfoCode
	MultimediaDescription          = inventory202410.MultimediaDescription
	TextItem                       = inventory202410.TextItem
	ImageItem                      = inventory202410.ImageItem
	ImageFormat                    = inventory202410.ImageFormat
)

const (
	RoomTypeRoom                 = inventory202410.RoomTypeRoom
	RoomTypeApartment            = inventory202410.RoomTypeApartment
	RoomTypeMobileHome           = inventory202410.RoomTypeMobileHome
	RoomTypeBungalow             = inventory202410.RoomTypeBungalow
	RoomTypeHolidayHome          = inventory202410.RoomTypeHolidayHome
	RoomTypeCampingGround        = inventory202410.RoomTypeCampingGround
	RoomTypePitch                = inventory202410.RoomTypePitch
	RoomTypeCampingGroundOrPitch = inventory202410.RoomTypeCampingGroundOrPitch
	RoomTypeRestingPlace         = inventory202410.RoomTypeRestingPlace

	InfoCodeDescription = inventory202410.InfoCodeDescription
	InfoCodePictures    = inventory202410.InfoCodePictures
	InfoCodeLongName    = inventory202410.InfoCodeLongName
)

type HotelDescriptiveContentNotifRQ struct {
	XMLName                 xml.Name                `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelDescriptiveContentNotifRQ"`
	Version                 string                  `xml:"Version,attr"`
	HotelDescriptiveContent HotelDescriptiveContent `xml:"HotelDescriptiveContents>HotelDescriptiveContent"`
}

// HotelCode returns the code of the hotel the request is for.
func (r HotelDescriptiveContentNotifRQ) HotelCode() string {
	return r.HotelDescriptiveContent.HotelCode
}

// GuestRooms returns the GuestRoom elements in message order.
func (r HotelDescriptiveContentNotifRQ) GuestRooms() []GuestRoom {
	if f := r.HotelDescriptiveContent.FacilityInfo; f != nil {
		return f.GuestRooms
	}
	return nil
}

type HotelDescriptiveContent struct {
	HotelCode       string            `xml:"HotelCode,attr"`
	HotelName       string            `xml:"HotelName,attr,omitempty"`
	FacilityInfo    *FacilityInfo     `xml:"FacilityInfo,omitempty"`
	AdditionalCodes *[]AdditionalCode `xml:"TPA_Extensions>AdditionalCodes>AdditionalCode"`
}

// AdditionalCode identifies the accommodation as national or EU rules
// require, such as Type CIN in Italy (AlpineBits 2026-10 §4.3.1).
type AdditionalCode struct {
	Type  string `xml:"Type,attr"`
	Value string `xml:",chardata"`
}

// FacilityInfo holds all room categories of the hotel; without GuestRooms it
// removes them all (AlpineBits 2026-10 §4.3.1).
type FacilityInfo struct {
	GuestRooms []GuestRoom `xml:"GuestRooms>GuestRoom"`
}

// GuestRoom defines a room category if it is the first with its Code, and
// lists a room of that category otherwise; see Partition.
type GuestRoom struct {
	Code              string `xml:"Code,attr"`
	MinOccupancy      int    `xml:"MinOccupancy,attr,omitempty"`
	MaxOccupancy      int    `xml:"MaxOccupancy,attr,omitempty"`
	MaxChildOccupancy int    `xml:"MaxChildOccupancy,attr,omitempty"`
	MaxAdultOccupancy int    `xml:"MaxAdultOccupancy,attr,omitempty"`
	// OldCode is the previous Code of a renamed category; see Renames.
	OldCode                string                   `xml:"ID,attr,omitempty"`
	TypeRoom               TypeRoom                 `xml:"TypeRoom"`
	Amenities              *[]Amenity               `xml:"Amenities>Amenity"`
	MultimediaDescriptions *[]MultimediaDescription `xml:"MultimediaDescriptions>MultimediaDescription"`
}

// MinFull is the minimum number of guests that pay the full rate in the
// category (§4.4.2, computing the cost of a stay).
func (g GuestRoom) MinFull() int {
	if g.MaxChildOccupancy == 0 {
		return g.TypeRoom.StandardOccupancy
	}
	return min(g.MaxOccupancy-g.MaxChildOccupancy, g.TypeRoom.StandardOccupancy)
}

// MultimediaDescription returns the MultimediaDescription with the given
// InfoCode, or nil.
func (g GuestRoom) MultimediaDescription(code InfoCode) *MultimediaDescription {
	mds := rules.Items(g.MultimediaDescriptions)
	for i := range mds {
		if mds[i].InfoCode == code {
			return &mds[i]
		}
	}
	return nil
}
