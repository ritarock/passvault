package service

import "github.com/ritarock/passvault/domain"

type mockVaultRepository struct {
	loadFunc   func() (*domain.Vault, error)
	saveFunc   func(vault *domain.Vault) error
	existsFunc func() bool
}

func (m *mockVaultRepository) Load() (*domain.Vault, error) {
	if m.loadFunc != nil {
		return m.loadFunc()
	}
	return domain.NewVault(), nil
}

func (m *mockVaultRepository) Save(vault *domain.Vault) error {
	if m.saveFunc != nil {
		return m.saveFunc(vault)
	}
	return nil
}

func (m *mockVaultRepository) Exists() bool {
	if m.existsFunc != nil {
		return m.existsFunc()
	}
	return true
}
