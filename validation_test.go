package alpinebits

import (
	"reflect"
	"testing"
)

func TestErrorsAddf(t *testing.T) {
	var errs Errors
	errs.Addf(ErrCodeRequiredField, "", "missing HotelCode")
	errs.Addf(ErrCodeInvalidValue, At("RatePlan", 0).At("Supplement", 2).Field("PrerequisiteInventory"), "InvCode %s", "required")
	want := Errors{
		{Type: ErrTypeApplicationError, Code: ErrCodeRequiredField, Message: "missing HotelCode"},
		{Type: ErrTypeApplicationError, Code: ErrCodeInvalidValue, Message: "RatePlan[0].Supplement[2].PrerequisiteInventory: InvCode required"},
	}
	if !reflect.DeepEqual(errs, want) {
		t.Errorf("got %+v\nwant %+v", errs, want)
	}
}

func TestReject(t *testing.T) {
	var none Errors
	if err := Reject(none...); err != nil {
		t.Errorf("Reject() = %v, want nil", err)
	}
	err := Reject(Error{Message: "x"})
	if rej, ok := err.(*Rejection); !ok || len(rej.Errors) != 1 {
		t.Errorf("Reject(x) = %#v, want *Rejection with one error", err)
	}
}
