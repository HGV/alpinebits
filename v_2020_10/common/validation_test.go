package common

import (
	"testing"
	"time"

	"github.com/HGV/x/timex"
	"github.com/stretchr/testify/assert"
)

func TestValidateHotelCode(t *testing.T) {
	assert.ErrorIs(t, ValidateHotelCode(" "), ErrMissingHotelCode)
	assert.Nil(t, ValidateHotelCode("9000"), ErrMissingHotelCode)
}

type mockRange struct {
	dateRange timex.DateRange
}

func (m mockRange) DateRange() timex.DateRange {
	return m.dateRange
}

func newDate(year, month, day int) timex.Date {
	return timex.Date{
		Year:  year,
		Month: time.Month(month),
		Day:   day,
	}
}

func TestValidateOverlaps(t *testing.T) {
	tests := []struct {
		name        string
		ranges      []mockRange
		expectError bool
	}{
		{
			name:        "No ranges",
			ranges:      []mockRange{},
			expectError: false,
		},
		{
			name: "Single range",
			ranges: []mockRange{
				{dateRange: timex.DateRange{Start: newDate(2023, 1, 1), End: newDate(2023, 1, 10)}},
			},
			expectError: false,
		},
		{
			name: "Non-overlapping ranges",
			ranges: []mockRange{
				{dateRange: timex.DateRange{Start: newDate(2023, 1, 1), End: newDate(2023, 1, 10)}},
				{dateRange: timex.DateRange{Start: newDate(2023, 1, 11), End: newDate(2023, 1, 20)}},
			},
			expectError: false,
		},
		{
			name: "Overlapping ranges",
			ranges: []mockRange{
				{dateRange: timex.DateRange{Start: newDate(2023, 1, 1), End: newDate(2023, 1, 10)}},
				{dateRange: timex.DateRange{Start: newDate(2023, 1, 9), End: newDate(2023, 1, 20)}},
			},
			expectError: true,
		},
		{
			name: "Adjacent ranges",
			ranges: []mockRange{
				{dateRange: timex.DateRange{Start: newDate(2023, 1, 1), End: newDate(2023, 1, 10)}},
				{dateRange: timex.DateRange{Start: newDate(2023, 1, 10), End: newDate(2023, 1, 20)}},
			},
			expectError: true,
		},
		{
			name: "Multiple overlaps",
			ranges: []mockRange{
				{dateRange: timex.DateRange{Start: newDate(2023, 1, 1), End: newDate(2023, 1, 10)}},
				{dateRange: timex.DateRange{Start: newDate(2023, 1, 9), End: newDate(2023, 1, 15)}},
				{dateRange: timex.DateRange{Start: newDate(2023, 1, 14), End: newDate(2023, 1, 20)}},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOverlaps(tt.ranges)
			if tt.expectError {
				assert.Error(t, err, "expected an error")
			} else {
				assert.NoError(t, err, "did not expect an error")
			}
		})
	}
}
