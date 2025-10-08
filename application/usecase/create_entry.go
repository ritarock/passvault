package usecase

import (
	"fmt"

	"github.com/ritarock/passvault/domain/entity"
	"github.com/ritarock/passvault/domain/repository"
)

type CreateEntryUsecase struct {
	vaultRepo repository.VaultRepository
}

func NewCreateEntryUsecase(vaultRepo repository.VaultRepository) *CreateEntryUsecase {
	return &CreateEntryUsecase{
		vaultRepo: vaultRepo,
	}
}

func (uc *CreateEntryUsecase) Execute(title, password, url, notes string) error {
	vault, err := uc.vaultRepo.Load()
	if err != nil {
		return fmt.Errorf("failed to lead vault: %w", err)
	}

	en := entity.NewEntry(title, password, url, notes)

	if err := vault.CreateEntry(*en); err != nil {
		return fmt.Errorf("failed to create entry: %w", err)
	}

	if err := uc.vaultRepo.Save(vault); err != nil {
		return fmt.Errorf("failed to save vault: %w", err)
	}

	return nil
}
