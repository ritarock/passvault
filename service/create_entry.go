package service

import (
	"fmt"

	"github.com/ritarock/passvault/domain"
)

type CreateEntryUsecase struct {
	vaultRepo domain.VaultRepository
}

func NewCreateEntryUsecase(vaultRepo domain.VaultRepository) *CreateEntryUsecase {
	return &CreateEntryUsecase{
		vaultRepo: vaultRepo,
	}
}

func (uc *CreateEntryUsecase) Execute(title, username, password, url, notes string) error {
	vault, err := uc.vaultRepo.Load()
	if err != nil {
		return fmt.Errorf("failed to lead vault: %w", err)
	}

	en := domain.NewEntry(title, username, password, url, notes)

	if err := vault.CreateEntry(*en); err != nil {
		return fmt.Errorf("failed to create entry: %w", err)
	}

	if err := uc.vaultRepo.Save(vault); err != nil {
		return fmt.Errorf("failed to save vault: %w", err)
	}

	return nil
}
