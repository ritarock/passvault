package service

import (
	"errors"
	"testing"

	"github.com/ritarock/passvault/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetEntryUsecase_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func() (*mockVaultRepository, string)
		hasErr bool
	}{
		{
			name: "succeed: get existing entry and mark as viewed",
			setup: func() (*mockVaultRepository, string) {
				vault := domain.NewVault()
				entry := domain.NewEntry("test title", "test username", "test password", "test url", "test notes")
				vault.Entries[entry.ID] = entry
				return &mockVaultRepository{
					loadFunc: func() (*domain.Vault, error) {
						return vault, nil
					},
					saveFunc: func(vault *domain.Vault) error {
						return nil
					},
				}, entry.ID
			},
			hasErr: false,
		},
		{
			name: "failed: vault load error",
			setup: func() (*mockVaultRepository, string) {
				return &mockVaultRepository{
					loadFunc: func() (*domain.Vault, error) {
						return nil, errors.New("load error")
					},
				}, "test-id"
			},
			hasErr: true,
		},
		{
			name: "failed: entry not found",
			setup: func() (*mockVaultRepository, string) {
				return &mockVaultRepository{
					loadFunc: func() (*domain.Vault, error) {
						return domain.NewVault(), nil
					},
				}, "non-existent-id"
			},
			hasErr: true,
		},
		{
			name: "failed: vault save error",
			setup: func() (*mockVaultRepository, string) {
				vault := domain.NewVault()
				entry := domain.NewEntry("test title", "test username", "test password", "test url", "test notes")
				vault.Entries[entry.ID] = entry
				return &mockVaultRepository{
					loadFunc: func() (*domain.Vault, error) {
						return vault, nil
					},
					saveFunc: func(vault *domain.Vault) error {
						return errors.New("save error")
					},
				}, entry.ID
			},
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repo, id := test.setup()
			usecase := NewGetEntryUsecase(repo)
			entry, err := usecase.Execute(id)
			if test.hasErr {
				assert.Error(t, err)
				assert.Nil(t, entry)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, entry)
				assert.Equal(t, id, entry.ID)
			}
		})
	}
}
