package service

import (
	"fmt"

	"github.com/ritarock/passvault/domain"
)

type GetEntryUsecase struct {
	vaultRepo domain.VaultRepository
}

func NewGetEntryUsecase(vaultRepo domain.VaultRepository) *GetEntryUsecase {
	return &GetEntryUsecase{
		vaultRepo: vaultRepo,
	}
}

func (uc *GetEntryUsecase) Execute(id string) (*domain.Entry, error) {
	vault, err := uc.vaultRepo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load vault: %w", err)
	}

	en, err := vault.GetEntry(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get entry: %w", err)
	}

	en.MarkAsViewed()

	if err := uc.vaultRepo.Save(vault); err != nil {
		return nil, fmt.Errorf("failed to save vault: %w", err)
	}

	return en, nil
}
