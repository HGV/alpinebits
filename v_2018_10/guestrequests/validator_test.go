package guestrequests

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResRetrieveValidator_Validate(t *testing.T) {
	tests := []struct {
		file      string
		validator ResRetrieveValidator
	}{
		{
			file:      "test/data/GuestRequests-OTA_ResRetrieveRS-cancellation.xml",
			validator: NewResRetrieveValidator(),
		},
		{
			file:      "test/data/GuestRequests-OTA_ResRetrieveRS-error.xml",
			validator: NewResRetrieveValidator(),
		},
		// {
		// 	file:      "test/data/GuestRequests-OTA_ResRetrieveRS-request-with-roomtype.xml",
		// 	validator: NewResRetrieveValidator(),
		// },
		{
			file:      "test/data/GuestRequests-OTA_ResRetrieveRS-reservation-empty.xml",
			validator: NewResRetrieveValidator(),
		},
		{
			file:      "test/data/GuestRequests-OTA_ResRetrieveRS-reservation.xml",
			validator: NewResRetrieveValidator(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			data, err := os.ReadFile(tt.file)
			if err != nil {
				assert.NoError(t, err, "Failed to read file %s", tt.file)
			}

			var rq ResRetrieveRS
			if err := xml.Unmarshal(data, &rq); err != nil {
				assert.NoError(t, err, "Failed to unmarshal data from file %s", tt.file)
			}

			assert.IsType(t, ResRetrieveRS{}, rq)
			assert.NoError(t, tt.validator.Validate(rq))
		})
	}
}
