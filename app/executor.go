package app

import (
	"errors"
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/ritarock/pa55vault/entity"
	"github.com/ritarock/pa55vault/infra"
)

type executor struct {
	sub   sub
	store infra.Store
}

func newExecutor(sub sub, filePath string) *executor {
	store := infra.NewStore(filePath)
	return &executor{
		sub:   sub,
		store: *store,
	}
}

func (e *executor) execute() error {
	switch e.sub.action {
	case Generate:
		if e.sub.help {
			help(e.sub.action)
			return nil
		}
		return generate(e.sub.generateData, &e.store)
	case List:
		if e.sub.help {
			help(e.sub.action)
			return nil
		}
		return list(&e.store)
	case Get:
		if e.sub.help {
			help(e.sub.action)
			return nil
		}
		return get(e.sub.getData, &e.store)
	case Help:
		help(Help)
	}
	return nil
}

func generate(generateData generateData, store *infra.Store) error {
	vault := entity.Vault{
		Title: generateData.title,
		Url:   generateData.url,
	}
	if err := vault.GenerateCode(); err != nil {
		return err
	}
	if err := vault.Encrypt(); err != nil {
		return err
	}
	vaults, err := store.Read()
	if err != nil {
		return err
	}
	vaults = append(vaults, vault)

	return store.Write(vaults)
}

func list(store *infra.Store) error {
	vaults, err := store.Read()
	if err != nil {
		return err
	}

	for i, v := range vaults {
		fmt.Printf("%d - Title: %s, URL: %s\n", i+1, v.Title, v.Url)
	}

	return nil
}

func get(getData int, store *infra.Store) error {
	vaults, err := store.Read()
	if err != nil {
		return err
	}

	index := getData - 1
	var vault entity.Vault
	for i, v := range vaults {
		if i == index {
			vault = v
		}
	}
	if vault.Code == "" {
		return errors.New("not found ID")
	}

	code, err := vault.Decrypt()
	if err != nil {
		return err
	}
	if err := clipboard.WriteAll(code); err != nil {
		return err
	}

	return nil
}

func help(action action) {
	switch action {
	case Generate:
		fmt.Print(GenerateHelp)
	case List:
		fmt.Print(HelpText)
	case Get:
		fmt.Print(GetHelp)
	case Help:
		fmt.Print(HelpText)
	}
}
