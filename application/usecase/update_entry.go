package usecase

import (
	"fmt"

	"github.com/ritarock/passvault/domain/repository"
)

type UpdateEntryUsecase struct {
	vaultRepo repository.VaultRepository
}

func NewUpdateEntryUsecase(vaultRepo repository.VaultRepository) *UpdateEntryUsecase {
	return &UpdateEntryUsecase{
		vaultRepo: vaultRepo,
	}
}

func (uc *UpdateEntryUsecase) Execute(id, name, password, url, notes string) error {
	vault, err := uc.vaultRepo.Load()
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	en, err := vault.GetEntry(id)
	if err != nil {
		return fmt.Errorf("failed to get entry: %w", err)
	}

	en.Update(name, password, url, notes)

	if err := vault.UpdateEntry(*en); err != nil {
		return fmt.Errorf("failed to update entry: %w", err)
	}

	if err := uc.vaultRepo.Save(vault); err != nil {
		return fmt.Errorf("failed to save vault: %w", err)
	}

	return nil
}
