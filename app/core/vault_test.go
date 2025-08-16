package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVault_GenerateCode(t *testing.T) {
	t.Parallel()
	initCode := ""
	tests := []struct {
		name  string
		vault Vault
	}{
		{
			name:  "pass",
			vault: Vault{Code: initCode},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, initCode, test.vault.Code)
			err := test.vault.GenerateCode()
			assert.NoError(t, err)
			assert.NotEqual(t, initCode, test.vault.Code)
			assert.Equal(t, len(test.vault.Code), TokenLength)
		})
	}
}

func Test_getRandomChar(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		charset  string
		hasError bool
	}{
		{
			name:     "pass",
			charset:  "abc",
			hasError: false,
		},
		{
			name:     "failed",
			charset:  "",
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := getRandomChar(test.charset)
			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_shuffle(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		token []byte
	}{
		{
			name:  "pass",
			token: []byte("ABCDEFG"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			original := make([]byte, len(test.token))
			copy(original, test.token)

			err := shuffle(test.token)
			assert.NoError(t, err)
			assert.Equal(t, len(original), len(test.token))
		})
	}
}
