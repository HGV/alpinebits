package inventory

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHotelDescriptiveContentNotifValidator_Validate(t *testing.T) {
	tests := []struct {
		file      string
		validator HotelDescriptiveContentNotifValidator
	}{
		{
			file:      "test/data/Inventory-OTA_HotelDescriptiveContentNotifRQ-delete-all.xml",
			validator: NewHotelDescriptiveContentNotifValidator(),
		},
		{
			file:      "test/data/Inventory-OTA_HotelDescriptiveContentNotifRQ-with-roomtype.xml",
			validator: NewHotelDescriptiveContentNotifValidator(),
		},
		{
			file: "test/data/Inventory-OTA_HotelDescriptiveContentNotifRQ.xml",
			validator: NewHotelDescriptiveContentNotifValidator(
				WithHotelDescriptiveContentNotifOccupancyChildren(true),
				WithHotelDescriptiveContentNotifRooms(true),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			data, err := os.ReadFile(tt.file)
			if err != nil {
				assert.NoError(t, err, "Failed to read file %s", tt.file)
			}

			var rq HotelDescriptiveContentNotifRQ
			if err := xml.Unmarshal(data, &rq); err != nil {
				assert.NoError(t, err, "Failed to unmarshal data from file %s", tt.file)
			}

			assert.IsType(t, HotelDescriptiveContentNotifRQ{}, rq)
			assert.NoError(t, tt.validator.Validate(rq))
		})
	}
}
