package inventory

import (
	"encoding/xml"

	"github.com/HGV/alpinebits"
)

type HotelDescriptiveContentNotifRQ struct {
	XMLName                 xml.Name                `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelDescriptiveContentNotifRQ"`
	Version                 string                  `xml:"Version,attr"`
	HotelDescriptiveContent HotelDescriptiveContent `xml:"HotelDescriptiveContents>HotelDescriptiveContent"`
}

// HotelCode implements alpinebits.HotelCoded.
func (r HotelDescriptiveContentNotifRQ) HotelCode() string {
	return r.HotelDescriptiveContent.HotelCode
}

type HotelDescriptiveContent struct {
	HotelCode  string      `xml:"HotelCode,attr"`
	HotelName  string      `xml:"HotelName,attr"`
	AreaID     int         `xml:"AreaID,attr,omitempty"`
	GuestRooms []GuestRoom `xml:"FacilityInfo>GuestRooms>GuestRoom"`
}

type GuestRoom struct {
	Code                   string                         `xml:"Code,attr"`
	MinOccupancy           int                            `xml:"MinOccupancy,attr,omitempty"`
	MaxOccupancy           int                            `xml:"MaxOccupancy,attr,omitempty"`
	MaxChildOccupancy      int                            `xml:"MaxChildOccupancy,attr,omitempty"`
	OldCode                string                         `xml:"ID,attr,omitempty"`
	TypeRoom               TypeRoom                       `xml:"TypeRoom"`
	Amenities              []Amenity                      `xml:"Amenities>Amenity"`
	MultimediaDescriptions alpinebits.MultimediaDescriptions `xml:"MultimediaDescriptions>MultimediaDescription"`
}

type TypeRoom struct {
	StandardOccupancy      int    `xml:"StandardOccupancy,attr,omitempty"`
	RoomClassificationCode int    `xml:"RoomClassificationCode,attr,omitempty"`
	RoomType               int    `xml:"RoomType,attr,omitempty"`
	Size                   int    `xml:"Size,attr,omitempty"`
	RoomID                 string `xml:"RoomID,attr,omitempty"`
}

func (g GuestRoom) MinFull() int {
	if g.MaxChildOccupancy == 0 {
		return g.TypeRoom.StandardOccupancy
	}
	return min(g.MaxOccupancy-g.MaxChildOccupancy, g.TypeRoom.StandardOccupancy)
}

type Amenity struct {
	RoomAmenityCode int `xml:"RoomAmenityCode,attr"`
}

type HotelDescriptiveContentNotifRS struct {
	XMLName  xml.Name             `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelDescriptiveContentNotifRS"`
	Version  string               `xml:"Version,attr"`
	Success  *alpinebits.Success  `xml:"Success"`
	Warnings *alpinebits.Warnings `xml:"Warnings>Warning"`
	Errors   *alpinebits.Errors   `xml:"Errors>Error"`
}
