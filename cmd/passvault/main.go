package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ritarock/passvault/application/usecase"
	"github.com/ritarock/passvault/domain/entity"
	"github.com/ritarock/passvault/infra/crypto"
	"github.com/ritarock/passvault/infra/persistence"
	"github.com/ritarock/passvault/presentation/tui"
)

const (
	AppDir = ".passvault"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func run() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	baseDir := filepath.Join(homeDir, AppDir)

	keyManager := crypto.NewKeyManager(baseDir)
	cryptoSvc := crypto.NewAESEncryptor(keyManager)
	vaultRepo := persistence.NewFileVaultRepository(baseDir, cryptoSvc)

	if !cryptoSvc.KeyExists() {
		if err := initialize(cryptoSvc, vaultRepo); err != nil {
			return fmt.Errorf("failed to initialize: %w", err)
		}
	}

	if !vaultRepo.Exists() {
		vault := entity.NewVault()
		if err := vaultRepo.Save(vault); err != nil {
			return fmt.Errorf("failed to create vault: %w", err)
		}
	}

	listEntriesUc := usecase.NewListEntriesUsecase(vaultRepo)
	getEntryUc := usecase.NewGetEntryUsecase(vaultRepo)
	createEntryUc := usecase.NewCreateEntryUsecase(vaultRepo)
	updateEntryUc := usecase.NewUpdateEntryUsecase(vaultRepo)
	deleteEntryUc := usecase.NewDeleteEntryUsecase(vaultRepo)

	app := tui.NewApp(
		listEntriesUc,
		getEntryUc,
		createEntryUc,
		updateEntryUc,
		deleteEntryUc,
	)

	app.ShowList()

	return app.Run()
}

func initialize(cryptoSvc *crypto.AESEncryptor, vaultRepo *persistence.FileVaultRepository) error {
	fmt.Println("First time setup...")
	fmt.Println("Generating encryption key...")

	if err := cryptoSvc.InitializeKey(); err != nil {
		return fmt.Errorf("failed to initialize key: %w", err)
	}

	fmt.Println("Creating vault...")
	vault := entity.NewVault()
	if err := vaultRepo.Save(vault); err != nil {
		return fmt.Errorf("failed to create vault: %w", err)
	}

	fmt.Println("Setup complete!")

	return nil
}
