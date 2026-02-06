package domain

import (
	"strings"
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

func TestDefaultPasswordOptions(t *testing.T) {
	t.Parallel()

	opts := DefaultPasswordOptions()

	assert.Equal(t, 16, opts.Length)
	assert.True(t, opts.IncludeLowercase)
	assert.True(t, opts.IncludeUppercase)
	assert.True(t, opts.IncludeDigits)
	assert.False(t, opts.IncludeSymbols)
}

func TestPasswordOptions_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    PasswordOptions
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid default options",
			opts:    DefaultPasswordOptions(),
			wantErr: false,
		},
		{
			name: "length too short",
			opts: PasswordOptions{
				Length:           5,
				IncludeLowercase: true,
			},
			wantErr: true,
			errMsg:  "password length must be at least 8",
		},
		{
			name: "length too long",
			opts: PasswordOptions{
				Length:           100,
				IncludeLowercase: true,
			},
			wantErr: true,
			errMsg:  "password length must be at most 64",
		},
		{
			name: "no character types selected",
			opts: PasswordOptions{
				Length:           16,
				IncludeLowercase: false,
				IncludeUppercase: false,
				IncludeDigits:    false,
				IncludeSymbols:   false,
			},
			wantErr: true,
			errMsg:  "at least one character type must be selected",
		},
		{
			name: "minimum valid length",
			opts: PasswordOptions{
				Length:           8,
				IncludeLowercase: true,
			},
			wantErr: false,
		},
		{
			name: "maximum valid length",
			opts: PasswordOptions{
				Length:           64,
				IncludeLowercase: true,
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.opts.Validate()
			if test.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPasswordGenerator_GenerateWithOptions(t *testing.T) {
	t.Parallel()

	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits    = "0123456789"
		symbols   = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	)

	tests := []struct {
		name          string
		opts          PasswordOptions
		validChars    string
		invalidChars  string
		wantErr       bool
	}{
		{
			name: "lowercase only",
			opts: PasswordOptions{
				Length:           16,
				IncludeLowercase: true,
				IncludeUppercase: false,
				IncludeDigits:    false,
				IncludeSymbols:   false,
			},
			validChars:   lowercase,
			invalidChars: uppercase + digits + symbols,
			wantErr:      false,
		},
		{
			name: "uppercase only",
			opts: PasswordOptions{
				Length:           16,
				IncludeLowercase: false,
				IncludeUppercase: true,
				IncludeDigits:    false,
				IncludeSymbols:   false,
			},
			validChars:   uppercase,
			invalidChars: lowercase + digits + symbols,
			wantErr:      false,
		},
		{
			name: "digits only",
			opts: PasswordOptions{
				Length:           16,
				IncludeLowercase: false,
				IncludeUppercase: false,
				IncludeDigits:    true,
				IncludeSymbols:   false,
			},
			validChars:   digits,
			invalidChars: lowercase + uppercase + symbols,
			wantErr:      false,
		},
		{
			name: "symbols only",
			opts: PasswordOptions{
				Length:           16,
				IncludeLowercase: false,
				IncludeUppercase: false,
				IncludeDigits:    false,
				IncludeSymbols:   true,
			},
			validChars:   symbols,
			invalidChars: lowercase + uppercase + digits,
			wantErr:      false,
		},
		{
			name: "all character types",
			opts: PasswordOptions{
				Length:           16,
				IncludeLowercase: true,
				IncludeUppercase: true,
				IncludeDigits:    true,
				IncludeSymbols:   true,
			},
			validChars:   lowercase + uppercase + digits + symbols,
			invalidChars: "",
			wantErr:      false,
		},
		{
			name: "invalid - no character types",
			opts: PasswordOptions{
				Length:           16,
				IncludeLowercase: false,
				IncludeUppercase: false,
				IncludeDigits:    false,
				IncludeSymbols:   false,
			},
			wantErr: true,
		},
		{
			name: "invalid - length too short",
			opts: PasswordOptions{
				Length:           5,
				IncludeLowercase: true,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			pg := NewPasswordGenerator()
			got, err := pg.GenerateWithOptions(test.opts)

			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, got, test.opts.Length)

			for _, c := range got {
				assert.True(t, strings.ContainsRune(test.validChars, c),
					"password contains invalid character: %c", c)
			}

			if test.invalidChars != "" {
				for _, c := range got {
					assert.False(t, strings.ContainsRune(test.invalidChars, c),
						"password contains excluded character: %c", c)
				}
			}
		})
	}
}
