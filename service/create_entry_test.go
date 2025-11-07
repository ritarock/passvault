package service

import (
	"errors"
	"testing"

	"github.com/ritarock/passvault/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateEntryUsecase_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		setup    func() *mockVaultRepository
		title    string
		username string
		password string
		url      string
		notes    string
		hasErr   bool
	}{
		{
			name: "succeed: create new entry",
			setup: func() *mockVaultRepository {
				return &mockVaultRepository{
					loadFunc: func() (*domain.Vault, error) {
						return domain.NewVault(), nil
					},
					saveFunc: func(vault *domain.Vault) error {
						return nil
					},
				}
			},
			title:    "test title",
			username: "test username",
			password: "test password",
			url:      "test url",
			notes:    "test notes",
			hasErr:   false,
		},
		{
			name: "failed: vault load error",
			setup: func() *mockVaultRepository {
				return &mockVaultRepository{
					loadFunc: func() (*domain.Vault, error) {
						return nil, errors.New("load error")
					},
				}
			},
			title:    "test title",
			username: "test username",
			password: "test password",
			url:      "test url",
			notes:    "test notes",
			hasErr:   true,
		},
		{
			name: "failed: vault save error",
			setup: func() *mockVaultRepository {
				return &mockVaultRepository{
					loadFunc: func() (*domain.Vault, error) {
						return domain.NewVault(), nil
					},
					saveFunc: func(vault *domain.Vault) error {
						return errors.New("save error")
					},
				}
			},
			title:    "test title",
			username: "test username",
			password: "test password",
			url:      "test url",
			notes:    "test notes",
			hasErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repo := test.setup()
			usecase := NewCreateEntryUsecase(repo)
			err := usecase.Execute(test.title, test.username, test.password, test.url, test.notes)
			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
