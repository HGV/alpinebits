package alpinebits

import (
	"bytes"
	"fmt"
	"strconv"
)

// HotelCoded is implemented by request types that contain a hotel code.
type HotelCoded interface {
	HotelCode() string
}

// Decimal is an xs:decimal, such as an amount. Unlike a float64 it is never
// written in exponent form, which xs:decimal does not allow.
type Decimal float64

func (d Decimal) MarshalText() ([]byte, error) {
	return strconv.AppendFloat(nil, float64(d), 'f', -1, 64), nil
}

func (d *Decimal) UnmarshalText(data []byte) error {
	f, err := strconv.ParseFloat(string(bytes.TrimSpace(data)), 64)
	*d = Decimal(f)
	return err
}

// Days represents a duration in days (ISO 8601 format PxD).
type Days int

func ParseDays(s string) (Days, error) {
	n, err := parseDuration(s, 'D')
	return Days(n), err
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

// Nights represents a duration in nights (ISO 8601 format PxN).
type Nights int

func ParseNights(s string) (Nights, error) {
	n, err := parseDuration(s, 'N')
	return Nights(n), err
}

func (n Nights) String() string {
	return fmt.Sprintf("P%dN", n)
}

func (n Nights) MarshalText() ([]byte, error) {
	return []byte(n.String()), nil
}

func (n *Nights) UnmarshalText(data []byte) error {
	var err error
	*n, err = ParseNights(string(data))
	return err
}

// parseDuration parses P<num><suffix> format.
func parseDuration(s string, suffix byte) (int, error) {
	if len(s) < 3 || s[0] != 'P' || s[len(s)-1] != suffix {
		return 0, fmt.Errorf("invalid duration format: %s, expected P<num>%c", s, suffix)
	}
	n, err := strconv.Atoi(s[1 : len(s)-1])
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %s, expected P<num>%c", s, suffix)
	}
	return n, nil
}

// TextFormat specifies the format of description text.
type TextFormat string

const (
	TextFormatPlainText TextFormat = "PlainText"
	TextFormatHTML      TextFormat = "HTML"
)

// Text is a text in one language and format.
type Text struct {
	TextFormat TextFormat `xml:"TextFormat,attr"`
	Language   string     `xml:"Language,attr,omitempty"`
	Value      string     `xml:",chardata"`
}
