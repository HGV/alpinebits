package freerooms

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHotelInvCountNotifValidator_Validate(t *testing.T) {
	tests := []struct {
		file      string
		validator HotelInvCountNotifValidator
	}{
		{
			file: "test/data/FreeRooms-OTA_HotelInvCountNotifRQ-closing_seasons.xml",
			validator: NewHotelInvCountNotifValidator(
				WithClosingSeasons(true),
			),
		},
		{
			file: "test/data/FreeRooms-OTA_HotelInvCountNotifRQ-delta.xml",
			validator: NewHotelInvCountNotifValidator(
				WithDeltas(true),
			),
		},
		{
			file:      "test/data/FreeRooms-OTA_HotelInvCountNotifRQ-empty.xml",
			validator: NewHotelInvCountNotifValidator(),
		},
		{
			file:      "test/data/FreeRooms-OTA_HotelInvCountNotifRQ.xml",
			validator: NewHotelInvCountNotifValidator(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			data, err := os.ReadFile(tt.file)
			if err != nil {
				assert.NoError(t, err, "Failed to read file %s", tt.file)
			}

			var rq HotelInvCountNotifRQ
			if err := xml.Unmarshal(data, &rq); err != nil {
				assert.NoError(t, err, "Failed to unmarshal data from file %s", tt.file)
			}

			assert.IsType(t, HotelInvCountNotifRQ{}, rq)
			assert.NoError(t, tt.validator.Validate(rq))
		})
	}
}
