package common

import (
	"errors"
	"slices"
	"strings"

	"github.com/HGV/alpinebits/version"
)

type Validatable[T any] interface {
	Validate(v T) error
}

func ValidateHotelCode(hotelCode string) error {
	if strings.TrimSpace(hotelCode) == "" {
		return ErrMissingHotelCode
	}
	return nil
}

func ValidateOverlaps[T version.DateRangeProvider](ranges []T) error {
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

func ValidateLanguageUniqueness(descs []Description) error {
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

func ValidateString(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("string is empty or contains only whitespace")
	}
	return nil
}

func ValidateNonNilString(s *string) error {
	if s == nil {
		return nil
	}
	return ValidateString(*s)
}
