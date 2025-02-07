package guestrequests

import (
	"encoding/xml"

	"github.com/HGV/alpinebits/v_2018_10/common"
)

type NotifReportRQ struct {
	XMLName           xml.Name          `xml:"http://www.opentravel.org/OTA/2003/05 OTA_NotifReportRQ"`
	Version           string            `xml:"Version,attr"`
	Success           common.Success    `xml:"Success"`
	Warnings          *[]common.Warning `xml:"Warnings>Warning"`
	HotelReservations []Acknowledgement `xml:"NotifDetails>HotelNotifReport>HotelReservations>HotelReservation"`
}

type Acknowledgement struct {
	UniqueID UniqueID `xml:"UniqueID"`
}

type NotifReportRS struct {
	common.Response

	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_NotifReportRS"`
	Version string   `xml:"Version,attr"`
}
