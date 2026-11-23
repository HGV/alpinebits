// Package rulestest holds the test helpers shared by the rule functions of all
// versions.
package rulestest

import (
	"slices"
	"testing"

	"github.com/HGV/alpinebits"
)

// Expect reports an error unless errs has exactly the messages want.
func Expect(t testing.TB, errs alpinebits.Errors, want []string) {
	t.Helper()
	var got []string
	for _, e := range errs {
		got = append(got, e.Message)
	}
	if !slices.Equal(got, want) {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}
