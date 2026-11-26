package v202210

import (
	"testing"

	"github.com/HGV/alpinebits/v202210/freerooms"
)

// The 2022-10 kit changed only the dates of the 2020-10 samples.
const freeRoomsSamples = "../v202010/freerooms/testdata/FreeRooms-"

func TestFreeRoomsTypesRoundTrip(t *testing.T) {
	for _, name := range []string{"", "-closing_seasons", "-delta", "-empty"} {
		t.Run("RQ"+name, func(t *testing.T) {
			roundTrip(t, readSample[freerooms.HotelInvCountNotifRQ](t, freeRoomsSamples+"OTA_HotelInvCountNotifRQ"+name+".xml"))
		})
	}
	for _, name := range []string{"success", "advisory", "warning", "error"} {
		t.Run("RS "+name, func(t *testing.T) {
			roundTrip(t, readSample[freerooms.HotelInvCountNotifRS](t, freeRoomsSamples+"OTA_HotelInvCountNotifRS-"+name+".xml"))
		})
	}
}
