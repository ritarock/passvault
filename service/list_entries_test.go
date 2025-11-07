package service

import (
	"errors"
	"testing"

	"github.com/ritarock/passvault/domain"
	"github.com/stretchr/testify/assert"
)

func TestListEntriesUsecase_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func() *mockVaultRepository
		wantLength int
		hasErr     bool
	}{
		{
			name: "succeed: empty vault",
			setup: func() *mockVaultRepository {
				return &mockVaultRepository{
					loadFunc: func() (*domain.Vault, error) {
						return domain.NewVault(), nil
					},
				}
			},
			wantLength: 0,
			hasErr:     false,
		},
		{
			name: "succeed: vault with single entry",
			setup: func() *mockVaultRepository {
				vault := domain.NewVault()
				entry := domain.NewEntry("test title", "test username", "test password", "test url", "test notes")
				vault.Entries[entry.ID] = entry
				return &mockVaultRepository{
					loadFunc: func() (*domain.Vault, error) {
						return vault, nil
					},
				}
			},
			wantLength: 1,
			hasErr:     false,
		},
		{
			name: "succeed: vault with multiple entries",
			setup: func() *mockVaultRepository {
				vault := domain.NewVault()
				entry1 := domain.NewEntry("title1", "username1", "password1", "test url1", "notes1")
				entry2 := domain.NewEntry("title2", "username2", "password2", "test url2", "notes2")
				entry3 := domain.NewEntry("title3", "username3", "password3", "test url3", "notes3")
				vault.Entries[entry1.ID] = entry1
				vault.Entries[entry2.ID] = entry2
				vault.Entries[entry3.ID] = entry3
				return &mockVaultRepository{
					loadFunc: func() (*domain.Vault, error) {
						return vault, nil
					},
				}
			},
			wantLength: 3,
			hasErr:     false,
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
			wantLength: 0,
			hasErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repo := test.setup()
			usecase := NewListEntriesUsecase(repo)
			entries, err := usecase.Execute()
			if test.hasErr {
				assert.Error(t, err)
				assert.Nil(t, entries)
			} else {
				assert.NoError(t, err)
				assert.Len(t, entries, test.wantLength)
			}
		})
	}
}
