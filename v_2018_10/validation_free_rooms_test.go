package v_2018_10

import (
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
