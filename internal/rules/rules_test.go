package rules

import (
	"reflect"
	"testing"

	"github.com/HGV/alpinebits"
	"github.com/HGV/x/timex"
)

type span struct{ start, end, key string }

func (s span) DateRange() timex.DateRange {
	start, _ := timex.ParseDate(s.start)
	end, _ := timex.ParseDate(s.end)
	return timex.DateRange{Start: start, End: end}
}

func TestOverlaps(t *testing.T) {
	tests := []struct {
		name  string
		items []span
		want  []string
	}{
		{"shared last day", []span{{"2026-01-05", "2026-01-10", "a"}, {"2026-01-01", "2026-01-05", "a"}}, []string{"X[0]: overlaps X[1]"}},
		{"adjacent", []span{{"2026-01-01", "2026-01-05", "a"}, {"2026-01-06", "2026-01-10", "a"}}, nil},
		{"one long range", []span{{"2026-01-01", "2026-01-31", "a"}, {"2026-01-02", "2026-01-03", "a"}, {"2026-01-10", "2026-01-12", "a"}},
			[]string{"X[1]: overlaps X[0]", "X[2]: overlaps X[0]"}},
		{"different keys", []span{{"2026-01-01", "2026-01-05", "a"}, {"2026-01-01", "2026-01-05", "b"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(nil)
			idx := make([]int, len(tt.items))
			for i := range idx {
				idx[i] = i
			}
			Overlaps(c, "", "X", tt.items, idx, func(s span) string { return s.key })
			var got []string
			for _, e := range c.Errors() {
				got = append(got, e.Message)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNeedReportsOnce(t *testing.T) {
	c := New([]alpinebits.Capability{"b"})
	c.Need(alpinebits.At("X", 0), "a", "a")
	c.Need(alpinebits.At("X", 1), "a", "a")
	c.Need(alpinebits.At("X", 2), "b", "b")
	if got := c.Errors(); len(got) != 1 || got[0].Message != "X[0]: a not supported" {
		t.Errorf("got %v", got)
	}
}

func TestCurrency(t *testing.T) {
	tests := []struct {
		code   string
		amount alpinebits.Decimal
		want   []string
	}{
		{"EUR", 76.8, nil},
		{"EUR", 76.805, []string{"X: A has more than 2 decimal digits for EUR"}},
		{"JPY", 100, nil},
		{"JPY", 100.5, []string{"X: A has more than 0 decimal digits for JPY"}},
		{"BHD", 1.125, nil},
		{"eur", 1, []string{"X: CurrencyCode eur is not an ISO 4217 code"}},
		{"EUX", 1.12345, []string{"X: CurrencyCode EUX is not an ISO 4217 code"}},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			c := New(nil)
			c.Currency("X", tt.code)
			c.Amount("X", "A", tt.amount, tt.code)
			var got []string
			for _, e := range c.Errors() {
				got = append(got, e.Message)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
