package alpinebits

import (
	"fmt"
	"slices"

	"github.com/HGV/x/timex"
)

// RuleResult holds validation warnings and errors.
type RuleResult struct {
	Warnings Warnings
	Errors   Errors
}

// Ok returns true if there are no errors.
func (r RuleResult) Ok() bool {
	return len(r.Errors) == 0
}

// WarningsPtr returns a pointer to the warnings slice, or nil if empty.
func (r RuleResult) WarningsPtr() *Warnings {
	if len(r.Warnings) == 0 {
		return nil
	}
	return &r.Warnings
}

// Merge combines another RuleResult into this one.
func (r *RuleResult) Merge(other RuleResult) {
	r.Warnings = append(r.Warnings, other.Warnings...)
	r.Errors = append(r.Errors, other.Errors...)
}

// Validate runs rules in order, accumulating warnings and stopping on first error.
func Validate[RQ any](rq RQ, rules ...func(RQ) RuleResult) RuleResult {
	var result RuleResult
	for _, rule := range rules {
		r := rule(rq)
		result.Merge(r)
		if len(r.Errors) > 0 {
			return result
		}
	}
	return result
}

// When wraps a rule to only run when condition is true.
func When[RQ any](cond func(RQ) bool, rule func(RQ) RuleResult) func(RQ) RuleResult {
	return func(rq RQ) RuleResult {
		if !cond(rq) {
			return RuleResult{}
		}
		return rule(rq)
	}
}

// RequiredHotelCode returns a rule that checks for a non-empty hotel code.
func RequiredHotelCode[RQ HotelCoded](rq RQ) RuleResult {
	if rq.HotelCode() == "" {
		return RuleResult{
			Errors: []Error{
				ApplicationError(ErrCodeRequiredField, "missing HotelCode"),
			},
		}
	}
	return RuleResult{}
}

// DateRanged is implemented by types that have a date range.
type DateRanged interface {
	DateRange() timex.DateRange
}

// CheckOverlaps checks for overlapping date ranges in a slice.
// Returns an error describing the first overlap found, or nil if no overlaps.
func CheckOverlaps[T DateRanged](items []T) error {
	if len(items) <= 1 {
		return nil
	}

	// Sort by start date
	slices.SortFunc(items, func(a, b T) int {
		return a.DateRange().Start.Compare(b.DateRange().Start)
	})

	// Check consecutive ranges for overlap
	for i := 0; i < len(items)-1; i++ {
		r1 := items[i].DateRange()
		r2 := items[i+1].DateRange()
		if r1.End.After(r2.Start) {
			return fmt.Errorf("date range overlap: %s-%s and %s-%s",
				r1.Start, r1.End, r2.Start, r2.End)
		}
	}
	return nil
}
