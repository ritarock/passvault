package service

import (
	"fmt"

	"github.com/ritarock/passvault/domain"
)

type UpdateEntryUsecase struct {
	vaultRepo domain.VaultRepository
}

func NewUpdateEntryUsecase(vaultRepo domain.VaultRepository) *UpdateEntryUsecase {
	return &UpdateEntryUsecase{
		vaultRepo: vaultRepo,
	}
}

func (uc *UpdateEntryUsecase) Execute(id, title, username, password, url, notes string) error {
	vault, err := uc.vaultRepo.Load()
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	en, err := vault.GetEntry(id)
	if err != nil {
		return fmt.Errorf("failed to get entry: %w", err)
	}

	en.Update(title, username, password, url, notes)

	if err := vault.UpdateEntry(*en); err != nil {
		return fmt.Errorf("failed to update entry: %w", err)
	}

	if err := uc.vaultRepo.Save(vault); err != nil {
		return fmt.Errorf("failed to save vault: %w", err)
	}

	return nil
}
