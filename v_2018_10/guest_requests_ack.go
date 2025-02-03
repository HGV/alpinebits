package v_2018_10

import "encoding/xml"

type NotifReportRQ struct {
	XMLName           xml.Name            `xml:"http://www.opentravel.org/OTA/2003/05 OTA_NotifReportRQ"`
	Version           string              `xml:"Version,attr"`
	Success           Success             `xml:"Success"`
	Warnings          *[]Warning          `xml:"Warnings>Warning"`
	HotelReservations []HotelReservation2 `xml:"NotifDetails>HotelNotifReport>HotelReservations>HotelReservation"`
}

// TODO: Rename
type HotelReservation2 struct {
	UniqueID UniqueID2 `xml:"UniqueID"`
}

type NotifReportRS struct {
	response

	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_NotifReportRS"`
	Version string   `xml:"Version,attr"`
}
