package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateVersionString(t *testing.T) {
	tests := []struct {
		version string
		isValid bool
	}{
		{"2017-10b", true},
		{"2018-10", true},
		{"2020-10", true},
		{"2018_10", false},
		{"2017-AB", false},
		{"abcd-ef", false},
		{"2018-b10", false},
		{"202010", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			err := ValidateVersionString(tt.version)
			if tt.isValid {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
