package v_2018_10

import (
	"encoding/xml"

	"github.com/HGV/alpinebits/internal"
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

func (mds MultimediaDescriptions) LongNames() []Description {
	for _, md := range mds {
		if md.InfoCode == InformationTypeLongName {
			return *md.TextItems
		}
	}
	return nil
}

func (mds MultimediaDescriptions) Descriptions() []Description {
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
	InfoCode   InformationType `xml:"InfoCode,attr"`
	TextItems  *[]Description  `xml:"TextItems>TextItem>Description"`
	ImageItems *[]ImageItem    `xml:"ImageItems>ImageItem"`
}

type InformationType int

const (
	InformationTypeDescription InformationType = 1
	InformationTypePictures    InformationType = 23
	InformationTypeLongName    InformationType = 25
)

type Description struct {
	TextFormat TextFormat `xml:"TextFormat,attr"`
	Language   string     `xml:"Language,attr"`
	Value      string     `xml:",innerxml"`
}

type TextFormat string

const (
	TextFormatPlainText = "PlainText"
	TextFormatHTML      = "HTML"
)

type ImageItem struct {
	Category     int           `xml:"Category,attr"`
	ImageFormat  ImageFormat   `xml:"ImageFormat"`
	Descriptions []Description `xml:"Description,omitempty"`
}

type ImageFormat struct {
	CopyrightNotice string `xml:"CopyrightNotice,attr,omitempty"`
	URL             URL    `xml:"URL"`
}

type URL struct {
	Value string `xml:",innerxml"`
}

type HotelDescriptiveContentNotifRS struct {
	response

	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelDescriptiveContentNotifRS"`
	Version string   `xml:"Version,attr"`
}
