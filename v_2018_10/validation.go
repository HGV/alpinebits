package v_2018_10

import (
	"errors"
	"slices"
	"strings"

	"github.com/HGV/alpinebits/version"
)

type Validatable[T any] interface {
	Validate(v T) error
}

func validateHotelCode(hotelCode string) error {
	if strings.TrimSpace(hotelCode) == "" {
		return ErrMissingHotelCode
	}
	return nil
}

func validateOverlaps[T version.DateRangeProvider](ranges []T) error {
	if len(ranges) <= 1 {
		return nil
	}

	slices.SortFunc(ranges, func(a, b T) int {
		return a.DateRange().Start.Compare(b.DateRange().Start)
	})

	for i := 0; i < len(ranges)-1; i++ {
		if ranges[i].DateRange().End.After(ranges[i+1].DateRange().Start) {
			return errors.New("overlap")
		}
	}

	return nil
}
