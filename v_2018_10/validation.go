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

func validateLanguageUniqueness(descs []Description) error {
	seen := make(map[string]struct{})
	for _, desc := range descs {
		lang := strings.TrimSpace(desc.Language)
		key := lang + "|" + string(desc.TextFormat)
		if _, exists := seen[key]; exists {
			return ErrDuplicateLanguage
		}
		seen[key] = struct{}{}
	}
	return nil
}

func validateString(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("string is empty or contains only whitespace")
	}
	return nil
}

func validateNonNilString(s *string) error {
	if s == nil {
		return nil
	}
	return validateString(*s)
}
