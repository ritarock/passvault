package usecase

import (
	"fmt"

	"github.com/ritarock/passvault/domain/repository"
)

type DeleteEntryUsecase struct {
	vaultRepo repository.VaultRepository
}

func NewDeleteEntryUsecase(vaultRepo repository.VaultRepository) *DeleteEntryUsecase {
	return &DeleteEntryUsecase{
		vaultRepo: vaultRepo,
	}
}

func (uc *DeleteEntryUsecase) Execute(id string) error {
	vault, err := uc.vaultRepo.Load()
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	if err := vault.DeleteEntry(id); err != nil {
		return fmt.Errorf("failed to delete entry: %w", err)
	}

	if err := uc.vaultRepo.Save(vault); err != nil {
		return fmt.Errorf("failed to save vault: %w", err)
	}

	return nil
}
