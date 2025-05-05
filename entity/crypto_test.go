package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVault_Encrypt(t *testing.T) {
	t.Parallel()
	initCode := "1234567890ab"
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
			err := test.vault.Encrypt()
			assert.NoError(t, err)
			assert.NotEqual(t, initCode, test.vault.Code)
		})
	}
}

func TestVault_Decrypt(t *testing.T) {
	t.Parallel()
	initCode := "1234567890ab"
	tests := []struct {
		name  string
		vault Vault
		want  string
	}{
		{
			name:  "pass",
			vault: Vault{Code: initCode},
			want:  initCode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, initCode, test.vault.Code)
			err := test.vault.Encrypt()
			assert.NoError(t, err)
			assert.NoError(t, err)
			assert.NotEqual(t, initCode, test.vault.Code)

			got, err := test.vault.Decrypt()
			assert.NoError(t, err)
			assert.NotEqual(t, initCode, test.vault.Code)
			assert.Equal(t, test.want, got)
		})
	}
}
