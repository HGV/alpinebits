package freerooms

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHotelAvailNotifValidator(t *testing.T) {
	t.Run("DefaultValidator", func(t *testing.T) {
		validator := NewHotelAvailNotifValidator()
		assert.False(t, validator.supportsRooms)
		assert.False(t, validator.supportsCategories)
		assert.False(t, validator.supportsDeltas)
		assert.False(t, validator.supportsBookingThreshold)
		assert.Nil(t, validator.roomMapping)
		assert.Nil(t, validator.categoriesMapping)
	})

	t.Run("CustomValidator", func(t *testing.T) {
		roomMapping := map[string]map[string]struct{}{
			"DZ": {"101": {}, "102": {}},
		}
		categoryMapping := map[string]struct{}{
			"DZ": {},
		}

		validator := NewHotelAvailNotifValidator(
			WithRooms(true, &roomMapping),
			WithCategories(true, &categoryMapping),
			WithDeltas(true),
			WithBookingThreshold(true),
		)

		assert.True(t, validator.supportsRooms)
		assert.True(t, validator.supportsCategories)
		assert.True(t, validator.supportsDeltas)
		assert.True(t, validator.supportsBookingThreshold)
		assert.NotNil(t, validator.roomMapping)
		assert.NotNil(t, validator.categoriesMapping)
	})
}

func TestResRetrieveValidator_Validate(t *testing.T) {
	tests := []struct {
		file      string
		validator HotelAvailNotifValidator
	}{
		{
			file:      "test/data/FreeRooms-OTA_HotelAvailNotifRQ-empty.xml",
			validator: NewHotelAvailNotifValidator(),
		},
		{
			file:      "test/data/FreeRooms-OTA_HotelAvailNotifRQ.xml",
			validator: NewHotelAvailNotifValidator(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			data, err := os.ReadFile(tt.file)
			if err != nil {
				assert.NoError(t, err, "Failed to read file %s", tt.file)
			}

			var rq HotelAvailNotifRQ
			if err := xml.Unmarshal(data, &rq); err != nil {
				assert.NoError(t, err, "Failed to unmarshal data from file %s", tt.file)
			}

			assert.IsType(t, HotelAvailNotifRQ{}, rq)
			assert.NoError(t, tt.validator.Validate(rq))
		})
	}
}
