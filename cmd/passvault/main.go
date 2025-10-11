package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ritarock/passvault/domain"
	"github.com/ritarock/passvault/service"
	"github.com/ritarock/passvault/storage"
	"github.com/ritarock/passvault/tui"
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

	keyManager := storage.NewKeyManager(baseDir)
	cryptoSvc := storage.NewAESEncryptor(keyManager)
	vaultRepo := storage.NewFileVaultRepository(baseDir, cryptoSvc)

	if !cryptoSvc.KeyExists() {
		if err := initialize(cryptoSvc, vaultRepo); err != nil {
			return fmt.Errorf("failed to initialize: %w", err)
		}
	}

	if !vaultRepo.Exists() {
		vault := domain.NewVault()
		if err := vaultRepo.Save(vault); err != nil {
			return fmt.Errorf("failed to create vault: %w", err)
		}
	}

	listEntriesUc := service.NewListEntriesUsecase(vaultRepo)
	getEntryUc := service.NewGetEntryUsecase(vaultRepo)
	createEntryUc := service.NewCreateEntryUsecase(vaultRepo)
	updateEntryUc := service.NewUpdateEntryUsecase(vaultRepo)
	deleteEntryUc := service.NewDeleteEntryUsecase(vaultRepo)

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

func initialize(cryptoSvc *storage.AESEncryptor, vaultRepo *storage.FileVaultRepository) error {
	fmt.Println("First time setup...")
	fmt.Println("Generating encryption key...")

	if err := cryptoSvc.InitializeKey(); err != nil {
		return fmt.Errorf("failed to initialize key: %w", err)
	}

	fmt.Println("Creating vault...")
	vault := domain.NewVault()
	if err := vaultRepo.Save(vault); err != nil {
		return fmt.Errorf("failed to create vault: %w", err)
	}

	fmt.Println("Setup complete!")

	return nil
}
