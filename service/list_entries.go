package service

import (
	"fmt"

	"github.com/ritarock/passvault/domain"
)

type ListEntriesUsecase struct {
	vaultRepo domain.VaultRepository
}

func NewListEntriesUsecase(vaultRepo domain.VaultRepository) *ListEntriesUsecase {
	return &ListEntriesUsecase{
		vaultRepo: vaultRepo,
	}
}

func (uc *ListEntriesUsecase) Execute() ([]*domain.Entry, error) {
	vault, err := uc.vaultRepo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load vault: %w", err)
	}

	return vault.ListEntries(), nil
}
