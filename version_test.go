package alpinebits

import "testing"

func TestNewVersionNilValidator(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewVersion with nil Validator did not panic")
		}
	}()
	NewVersion("2018-10", nil)
}
