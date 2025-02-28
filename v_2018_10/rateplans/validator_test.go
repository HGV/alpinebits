package rateplans

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHotelRatePlanNotifValidator_Validate(t *testing.T) {
	tests := []struct {
		file      string
		validator HotelRatePlanNotifValidator
	}{
		{
			file: "test/data/RatePlans-OTA_HotelRatePlanNotifRQ.xml",
			validator: NewHotelRatePlanNotifValidator(
				WithArrivalDOW(),
				WithDepartureDOW(),
				WithRoomTypeCodes(map[string]RoomTypeOccupancySettings{
					"double": {Std: 2},
				}),
				WithSupplements(),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			data, err := os.ReadFile(tt.file)
			if err != nil {
				assert.NoError(t, err, "Failed to read file %s", tt.file)
			}

			var rq HotelRatePlanNotifRQ
			if err := xml.Unmarshal(data, &rq); err != nil {
				assert.NoError(t, err, "Failed to unmarshal data from file %s", tt.file)
			}

			assert.IsType(t, HotelRatePlanNotifRQ{}, rq)
			assert.NoError(t, tt.validator.Validate(rq))
		})
	}
}
