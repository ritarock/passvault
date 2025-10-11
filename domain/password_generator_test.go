package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordGenerator_Generate(t *testing.T) {
	t.Parallel()

	const (
		defaultLength = 16
		lowercase     = "abcdefghijklmnopqrstuvwxyz"
		uppercase     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits        = "0123456789"
		symbols       = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	)

	tests := []struct {
		name       string
		length     int
		wantLength int
	}{
		{
			name:       "default length when zero",
			length:     0,
			wantLength: 16,
		},
		{
			name:       "default length when negative",
			length:     -1,
			wantLength: 16,
		},
		{
			name:       "custom length 8",
			length:     8,
			wantLength: 8,
		},
		{
			name:       "custom length 32",
			length:     32,
			wantLength: 32,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			pg := NewPasswordGenerator()
			charseet := lowercase + uppercase + digits + symbols
			got, err := pg.Generate(test.length)
			assert.NoError(t, err)
			assert.Len(t, got, test.wantLength)

			for _, c := range got {
				found := false
				for _, validChar := range charseet {
					if c == validChar {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Generate() contains invalid character: %c", c)
				}
			}
		})
	}
}
