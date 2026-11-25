// Package rules holds what the message packages of all versions share.
package rules

import (
	"slices"
	"strconv"
	"strings"

	"github.com/HGV/alpinebits"
	"github.com/HGV/x/timex"
	"golang.org/x/text/currency"
)

// Checker collects the errors of a rule function.
type Checker struct {
	caps     []alpinebits.Capability
	reported []alpinebits.Capability
	errs     alpinebits.Errors
}

// New returns a Checker for a message sent with the negotiated caps.
func New(caps []alpinebits.Capability) *Checker {
	return &Checker{caps: caps}
}

func (c *Checker) Errors() alpinebits.Errors {
	return c.errs
}

func (c *Checker) Required(at alpinebits.Path, format string, args ...any) {
	c.errs.Addf(alpinebits.ErrCodeRequiredField, at, format, args...)
}

func (c *Checker) Invalid(at alpinebits.Path, format string, args ...any) {
	c.errs.Addf(alpinebits.ErrCodeInvalidValue, at, format, args...)
}

func (c *Checker) InvalidDateRange(at alpinebits.Path, format string, args ...any) {
	c.errs.Addf(alpinebits.ErrCodeInvalidDateRange, at, format, args...)
}

// Has reports whether capability was negotiated.
func (c *Checker) Has(capability alpinebits.Capability) bool {
	return slices.Contains(c.caps, capability)
}

// Need reports what as not supported unless capability was negotiated. Each
// capability is reported once, at its first use: messages can hold thousands
// of elements.
func (c *Checker) Need(at alpinebits.Path, capability alpinebits.Capability, what string) {
	if c.Has(capability) || slices.Contains(c.reported, capability) {
		return
	}
	c.reported = append(c.reported, capability)
	c.Invalid(at, "%s not supported", what)
}

// Hotel requires HotelCode in every version. AlpineBits made it mandatory in
// 2024-10 (§4.1.1, which the other sections refer to); earlier versions also
// allow HotelName alone, which no implementation is known to use.
func (c *Checker) Hotel(at alpinebits.Path, code string) {
	if Blank(code) {
		c.Required(at, "HotelCode required")
	}
}

// Texts checks the texts of one description: at most one per Language and
// TextFormat, and HTML only next to PlainText in the same Language. elem is
// the name of the text elements.
func (c *Checker) Texts(at alpinebits.Path, elem string, texts []alpinebits.Text) {
	type key struct {
		language string
		format   alpinebits.TextFormat
	}
	seen := make(map[key]bool)
	for k, t := range texts {
		if seen[key{t.Language, t.TextFormat}] {
			c.Invalid(at.At(elem, k), "duplicate %s for Language %s and TextFormat %s", elem, t.Language, t.TextFormat)
		}
		seen[key{t.Language, t.TextFormat}] = true
	}
	for k, t := range texts {
		if t.TextFormat == alpinebits.TextFormatHTML && !seen[key{t.Language, alpinebits.TextFormatPlainText}] {
			c.Required(at.At(elem, k), "HTML %s requires a PlainText %s in the same Language", elem, elem)
		}
	}
}

// DateRanged is implemented by elements that cover a period.
type DateRanged interface {
	DateRange() timex.DateRange
}

// Overlaps reports the items at idx whose period overlaps an earlier one with
// the same key, at parent.At(name, i). Start and End are inclusive.
func Overlaps[T DateRanged, K comparable](c *Checker, parent alpinebits.Path, name string, items []T, idx []int, key func(T) K) {
	var keys []K
	groups := make(map[K][]int)
	for _, i := range idx {
		k := key(items[i])
		if _, ok := groups[k]; !ok {
			keys = append(keys, k)
		}
		groups[k] = append(groups[k], i)
	}
	for _, k := range keys {
		g := groups[k]
		for _, p := range overlapping(g, func(i int) timex.DateRange { return items[i].DateRange() }) {
			c.InvalidDateRange(parent.At(name, p[1]), "overlaps %s[%d]", name, p[0])
		}
	}
}

// overlapping returns the pairs {earlier, later} of overlapping periods
// among idx, sorted by start.
func overlapping(idx []int, period func(int) timex.DateRange) [][2]int {
	if len(idx) < 2 {
		return nil
	}
	order := slices.Clone(idx)
	slices.SortStableFunc(order, func(a, b int) int {
		return period(a).Start.Compare(period(b).Start)
	})
	var pairs [][2]int
	reach := order[0] // the one with the latest End so far
	for _, i := range order[1:] {
		r := period(i)
		if !r.Start.After(period(reach).End) {
			pairs = append(pairs, [2]int{reach, i})
		}
		if r.End.After(period(reach).End) {
			reach = i
		}
	}
	return pairs
}

// Currency requires an ISO 4217 code, which AlpineBits allows since 2020-10
// (§4).
func (c *Checker) Currency(at alpinebits.Path, code string) {
	if _, ok := decimals(code); !ok {
		c.Invalid(at, "CurrencyCode %s is not an ISO 4217 code", code)
	}
}

// Amount requires amount to have no more decimal digits than its currency,
// given by code, allows. Invalid codes are left to Currency.
func (c *Checker) Amount(at alpinebits.Path, name string, amount alpinebits.Decimal, code string) {
	n, ok := decimals(code)
	if !ok {
		return
	}
	s := strconv.FormatFloat(float64(amount), 'f', -1, 64)
	if i := strings.IndexByte(s, '.'); i >= 0 && len(s)-i-1 > n {
		c.Invalid(at, "%s has more than %d decimal digits for %s", name, n, code)
	}
}

// decimals returns the decimal digits of the ISO 4217 currency code.
func decimals(code string) (int, bool) {
	u, err := currency.ParseISO(code)
	if err != nil || code != strings.ToUpper(code) {
		return 0, false
	}
	n, _ := currency.Standard.Rounding(u)
	return n, true
}

func Blank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Items returns the elements of an optional container. Containers are *[]T
// because encoding/xml writes an empty slice as an empty element.
func Items[T any](s *[]T) []T {
	if s == nil {
		return nil
	}
	return *s
}
