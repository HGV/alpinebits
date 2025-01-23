package duration

import (
	"fmt"
	"regexp"
	"strconv"
)

type Nights int

var regexpNights = regexp.MustCompile(`^P(P<nights>[0-9]+)N$`)

func ParseNights(s string) (Nights, error) {
	matches := regexpNights.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("invalid duration format: %s, expected format is 'PxN'", s)
	}
	nightsStr := matches[regexpNights.SubexpIndex("nights")]
	nights, _ := strconv.Atoi(nightsStr)
	return Nights(nights), nil
}

func (d Nights) String() string {
	return fmt.Sprintf("P%dN", d)
}

func (d Nights) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Nights) UnmarshalText(data []byte) error {
	var err error
	*d, err = ParseNights(string(data))
	return err
}
