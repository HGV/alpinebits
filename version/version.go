package version

import (
	"fmt"
	"regexp"

	"github.com/HGV/x/timex"
)

type (
	Version[A Action] interface {
		fmt.Stringer

		ValidateXML(xml string) error
	}
	Action interface {
		fmt.Stringer

		Unmarshal(b []byte) (any, error)
	}
	HotelCodeProvider interface {
		HotelCode() string
	}
	DateRangeProvider interface {
		DateRange() timex.DateRange
	}
)

func ValidateVersionString(s string) error {
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}\w?$`, s)
	if !matched {
		return fmt.Errorf("invalid protocol version string: %s", s)
	}
	return nil
}
