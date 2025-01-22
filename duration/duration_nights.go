package duration

import (
	"fmt"
	"regexp"
	"strconv"
)

type Nights int

var regexpNights = regexp.MustCompile(`^P(P<nights>[0-9]+)N$`)

func ParseNights(s string) (Nights, error) {
	var d Nights

	if !regexpNights.MatchString(s) {
		return d, fmt.Errorf("invalid duration format: %s, expected format is 'PxN'", s)
	}

	matches := regexpNights.FindStringSubmatch(s)
	for i, name := range regexpNights.SubexpNames() {
		switch match := matches[i]; name {
		case "nights":
			nights, err := strconv.Atoi(match)
			if err != nil {
				return d, err
			}
			d = Nights(nights)
		}
	}

	return d, nil
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
