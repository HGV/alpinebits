package duration

import (
	"fmt"
	"regexp"
	"strconv"
)

type Days int

var regexpDays = regexp.MustCompile(`^P(P<days>[0-9]+)D$`)

func ParseDays(s string) (Days, error) {
	var d Days

	if !regexpDays.MatchString(s) {
		return d, fmt.Errorf("invalid duration format: %s, expected format is 'PxD'", s)
	}

	matches := regexpDays.FindStringSubmatch(s)
	for i, name := range regexpDays.SubexpNames() {
		switch match := matches[i]; name {
		case "days":
			days, err := strconv.Atoi(match)
			if err != nil {
				return d, err
			}
			d = Days(days)
		}
	}

	return d, nil
}

func (d Days) String() string {
	return fmt.Sprintf("P%dD", d)
}

func (d Days) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Days) UnmarshalText(data []byte) error {
	var err error
	*d, err = ParseDays(string(data))
	return err
}
