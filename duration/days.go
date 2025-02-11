package duration

import (
	"fmt"
	"regexp"
	"strconv"
)

type Days int

var regexpDays = regexp.MustCompile(`^P(?P<days>[0-9]+)D$`)

func ParseDays(s string) (Days, error) {
	matches := regexpDays.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("invalid duration format: %s, expected format is 'PxD'", s)
	}
	daysStr := matches[regexpDays.SubexpIndex("days")]
	days, _ := strconv.Atoi(daysStr)
	return Days(days), nil
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
