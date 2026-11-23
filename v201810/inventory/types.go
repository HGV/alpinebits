// Package inventory holds the AlpineBits 2018-10 Inventory messages
// (OTA_HotelDescriptiveContentNotif, §4.4).
package inventory

import (
	"encoding/xml"

	"github.com/HGV/alpinebits"
	"github.com/HGV/alpinebits/internal/rules"
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

// GuestRooms returns the GuestRoom elements in message order.
func (r HotelDescriptiveContentNotifRQ) GuestRooms() []GuestRoom {
	if f := r.HotelDescriptiveContent.FacilityInfo; f != nil {
		return f.GuestRooms
	}
	return nil
}

type HotelDescriptiveContent struct {
	HotelCode    string        `xml:"HotelCode,attr,omitempty"`
	HotelName    string        `xml:"HotelName,attr,omitempty"`
	FacilityInfo *FacilityInfo `xml:"FacilityInfo,omitempty"`
}

// FacilityInfo holds all room categories of the hotel; without GuestRooms it
// removes them all (AlpineBits 2018-10 §4.4.1).
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
	// OldCode is the previous Code of a renamed category; see Renames.
	OldCode                string                   `xml:"ID,attr,omitempty"`
	TypeRoom               TypeRoom                 `xml:"TypeRoom"`
	Amenities              *[]Amenity               `xml:"Amenities>Amenity"`
	MultimediaDescriptions *[]MultimediaDescription `xml:"MultimediaDescriptions>MultimediaDescription"`
}

// MinFull is the minimum number of guests that pay the full rate in the
// category (§4.5, computing the cost of a stay).
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

type TypeRoom struct {
	StandardOccupancy      int      `xml:"StandardOccupancy,attr,omitempty"`
	RoomClassificationCode int      `xml:"RoomClassificationCode,attr,omitempty"`
	RoomType               RoomType `xml:"RoomType,attr,omitempty"`
	Size                   int      `xml:"Size,attr,omitempty"`
	RoomID                 string   `xml:"RoomID,attr,omitempty"`
}

// RoomType refines RoomClassificationCode; each should only be used with the
// RoomClassificationCode in its comment (§4.4.1).
type RoomType int

const (
	RoomTypeRoom                 RoomType = 1 // 42
	RoomTypeApartment            RoomType = 2 // 13
	RoomTypeMobileHome           RoomType = 3 // 13
	RoomTypeBungalow             RoomType = 4 // 13
	RoomTypeHolidayHome          RoomType = 5 // 13
	RoomTypeCampingGround        RoomType = 6 // 5
	RoomTypePitch                RoomType = 7 // 5
	RoomTypeCampingGroundOrPitch RoomType = 8 // 5
	RoomTypeRestingPlace         RoomType = 9 // 42
)

type Amenity struct {
	RoomAmenityCode int `xml:"RoomAmenityCode,attr,omitempty"`
}

// InfoCode identifies the content of a MultimediaDescription (OTA list INF).
type InfoCode int

const (
	InfoCodeDescription InfoCode = 1
	InfoCodePictures    InfoCode = 23
	InfoCodeLongName    InfoCode = 25
)

// MultimediaDescription holds texts (InfoCode 1 and 25) or pictures (InfoCode
// 23).
type MultimediaDescription struct {
	InfoCode   InfoCode     `xml:"InfoCode,attr,omitempty"`
	TextItem   *TextItem    `xml:"TextItems>TextItem"`
	ImageItems *[]ImageItem `xml:"ImageItems>ImageItem"`
}

type TextItem struct {
	Descriptions []alpinebits.Text `xml:"Description"`
}

type ImageItem struct {
	Category     int               `xml:"Category,attr"`
	ImageFormat  ImageFormat       `xml:"ImageFormat"`
	Descriptions []alpinebits.Text `xml:"Description,omitempty"`
}

type ImageFormat struct {
	CopyrightNotice string `xml:"CopyrightNotice,attr,omitempty"`
	URL             string `xml:"URL"`
}

type HotelDescriptiveContentNotifRS struct {
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelDescriptiveContentNotifRS"`
	alpinebits.Response
}
