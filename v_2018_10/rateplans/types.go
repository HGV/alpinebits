package rateplans

import (
	"encoding/xml"

	"github.com/HGV/alpinebits"
)

type HotelRatePlanNotifRQ struct {
	XMLName   xml.Name  `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelRatePlanNotifRQ"`
	RatePlans RatePlans `xml:"RatePlans"`
}

// HotelCode implements alpinebits.HotelCoded.
func (r HotelRatePlanNotifRQ) HotelCode() string {
	return r.RatePlans.HotelCode
}

type RatePlans struct {
	HotelCode string `xml:"HotelCode,attr"`
}

type HotelRatePlanNotifRS struct {
	XMLName  xml.Name             `xml:"http://www.opentravel.org/OTA/2003/05 OTA_HotelRatePlanNotifRS"`
	Version  string               `xml:"Version,attr"`
	Success  *alpinebits.Success  `xml:"Success"`
	Warnings *alpinebits.Warnings `xml:"Warnings>Warning"`
	Errors   *alpinebits.Errors   `xml:"Errors>Error"`
}
