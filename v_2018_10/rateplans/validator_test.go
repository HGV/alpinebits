package rateplans

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"

	"github.com/HGV/alpinebits/v_2018_10/common"
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

func TestHotelRatePlanNotifValidator_ValidatePurge(t *testing.T) {
	const tmpl = `<OTA_HotelRatePlanNotifRQ xmlns="http://www.opentravel.org/OTA/2003/05" Version="1.000">
%s
<RatePlans HotelCode="123">
%s
</RatePlans>
</OTA_HotelRatePlanNotifRQ>`
	const completeSet = `<UniqueID Type="16" ID="1" Instance="CompleteSet"/>`

	tests := []struct {
		name      string
		uniqueID  string
		ratePlans string
		wantPurge bool
		wantErr   error
	}{
		{
			name:      "complete set with single empty rate plan purges",
			uniqueID:  completeSet,
			ratePlans: `<RatePlan/>`,
			wantPurge: true,
		},
		{
			name:      "complete set with rate plan codes is a reset, not a purge",
			uniqueID:  completeSet,
			ratePlans: `<RatePlan RatePlanCode="A"/><RatePlan RatePlanCode="B"/>`,
		},
		{
			name:      "complete set mixing empty and coded rate plans is rejected",
			uniqueID:  completeSet,
			ratePlans: `<RatePlan/><RatePlan RatePlanCode="A"/>`,
			wantErr:   common.ErrMissingRatePlanCode,
		},
		{
			name:      "complete set with multiple empty rate plans is rejected",
			uniqueID:  completeSet,
			ratePlans: `<RatePlan/><RatePlan/>`,
			wantErr:   common.ErrMissingRatePlanCode,
		},
		{
			name:      "empty rate plan without complete set is rejected",
			ratePlans: `<RatePlan/>`,
			wantErr:   common.ErrMissingRatePlanCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rq HotelRatePlanNotifRQ
			payload := fmt.Sprintf(tmpl, tt.uniqueID, tt.ratePlans)
			assert.NoError(t, xml.Unmarshal([]byte(payload), &rq))
			assert.Equal(t, tt.wantPurge, rq.IsPurge())

			v := NewHotelRatePlanNotifValidator()
			assert.ErrorIs(t, v.Validate(rq), tt.wantErr)
		})
	}
}
