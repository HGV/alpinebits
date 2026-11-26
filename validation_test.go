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
