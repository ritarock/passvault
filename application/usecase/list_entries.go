package usecase

import (
	"fmt"

	"github.com/ritarock/passvault/domain/entity"
	"github.com/ritarock/passvault/domain/repository"
)

type ListEntriesUsecase struct {
	vaultRepo repository.VaultRepository
}

func NewListEntriesUsecase(vaultRepo repository.VaultRepository) *ListEntriesUsecase {
	return &ListEntriesUsecase{
		vaultRepo: vaultRepo,
	}
}

func (uc *ListEntriesUsecase) Execute() ([]*entity.Entry, error) {
	vault, err := uc.vaultRepo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load vault: %w", err)
	}

	return vault.ListEntries(), nil
}
