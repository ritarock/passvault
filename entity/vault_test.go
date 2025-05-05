package entity

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
