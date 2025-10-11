package service

import (
	"errors"
	"testing"

	"github.com/ritarock/passvault/domain"
	"github.com/stretchr/testify/assert"
)

func TestUpdateEntryUsecase_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		setup    func() (*mockVaultRepository, string)
		title    string
		password string
		url      string
		notes    string
		hasErr   bool
	}{
		{
			name: "succeed: update existing entry",
			setup: func() (*mockVaultRepository, string) {
				vault := domain.NewVault()
				entry := domain.NewEntry("old title", "old password", "old url", "old notes")
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
			title:    "new title",
			password: "new password",
			url:      "new url",
			notes:    "new notes",
			hasErr:   false,
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
			title:    "new title",
			password: "new password",
			url:      "new url",
			notes:    "new notes",
			hasErr:   true,
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
			title:    "new title",
			password: "new password",
			url:      "new url",
			notes:    "new notes",
			hasErr:   true,
		},
		{
			name: "failed: vault save error",
			setup: func() (*mockVaultRepository, string) {
				vault := domain.NewVault()
				entry := domain.NewEntry("old title", "old password", "old url", "old notes")
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
			title:    "new title",
			password: "new password",
			url:      "new url",
			notes:    "new notes",
			hasErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repo, id := test.setup()
			usecase := NewUpdateEntryUsecase(repo)
			err := usecase.Execute(id, test.title, test.password, test.url, test.notes)
			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
