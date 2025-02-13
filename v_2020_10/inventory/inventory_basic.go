package inventory

import (
	"encoding/xml"

	"github.com/HGV/alpinebits/internal"
	"github.com/HGV/alpinebits/v_2020_10/common"
)

type HotelDescriptiveContentNotifRQ struct {
	XMLName                 xml.Name                `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelDescriptiveContentNotifRQ"`
	Version                 string                  `xml:"Version,attr"`
	HotelDescriptiveContent HotelDescriptiveContent `xml:"HotelDescriptiveContents>HotelDescriptiveContent"`
}

type HotelDescriptiveContent struct {
	HotelCode  string      `xml:"HotelCode,attr"`
	HotelName  string      `xml:"HotelName,attr"`
	AreaID     int         `xml:"AreaID,attr,omitempty"`
	GuestRooms []GuestRoom `xml:"FacilityInfo>GuestRooms>GuestRoom"`
}

type GuestRoom struct {
	Code                   string                 `xml:"Code,attr"`
	MinOccupancy           int                    `xml:"MinOccupancy,attr"`
	MaxOccupancy           int                    `xml:"MaxOccupancy,attr"`
	MaxChildOccupancy      int                    `xml:"MaxChildOccupancy,attr,omitempty"`
	OldCode                string                 `xml:"ID,attr,omitempty"`
	TypeRoom               TypeRoom               `xml:"TypeRoom"`
	Amenities              *[]Amenity             `xml:"Amenities>Amenity"`
	MultimediaDescriptions MultimediaDescriptions `xml:"MultimediaDescriptions>MultimediaDescription"`
}

func (g GuestRoom) MinFull() int {
	return internal.CalculateMinFull(g.MaxChildOccupancy, g.TypeRoom.StandardOccupancy, g.MaxOccupancy)
}

type TypeRoom struct {
	StandardOccupancy      int    `xml:"StandardOccupancy,attr"`
	RoomClassificationCode int    `xml:"RoomClassificationCode,attr"`
	RoomType               int    `xml:"RoomType,attr,omitempty"`
	Size                   int    `xml:"Size,attr,omitempty"`
	RoomID                 string `xml:"RoomID,attr,omitempty"`
}

type Amenity struct {
	RoomAmenityCode int `xml:"RoomAmenityCode,attr"`
}

type MultimediaDescriptions []MultimediaDescription

func (mds MultimediaDescriptions) LongNames() []common.Description {
	for _, md := range mds {
		if md.InfoCode == InformationTypeLongName {
			return *md.TextItems
		}
	}
	return nil
}

func (mds MultimediaDescriptions) Descriptions() []common.Description {
	for _, md := range mds {
		if md.InfoCode == InformationTypeDescription {
			return *md.TextItems
		}
	}
	return nil
}

func (mds MultimediaDescriptions) Pictures() []ImageItem {
	for _, md := range mds {
		if md.InfoCode == InformationTypePictures {
			return *md.ImageItems
		}
	}
	return nil
}

type MultimediaDescription struct {
	InfoCode   InformationType       `xml:"InfoCode,attr"`
	TextItems  *[]common.Description `xml:"TextItems>TextItem>Description"`
	ImageItems *[]ImageItem          `xml:"ImageItems>ImageItem"`
}

type InformationType int

const (
	InformationTypeDescription InformationType = 1
	InformationTypePictures    InformationType = 23
	InformationTypeLongName    InformationType = 25
)

type ImageItem struct {
	Category     int                  `xml:"Category,attr"`
	ImageFormat  ImageFormat          `xml:"ImageFormat"`
	Descriptions []common.Description `xml:"Description,omitempty"`
}

type ImageFormat struct {
	CopyrightNotice string     `xml:"CopyrightNotice,attr,omitempty"`
	URL             common.URL `xml:"URL"`
}

type HotelDescriptiveContentNotifRS struct {
	common.Response

	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelDescriptiveContentNotifRS"`
	Version string   `xml:"Version,attr"`
}
