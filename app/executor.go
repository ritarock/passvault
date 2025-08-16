package app

import (
	"errors"
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/ritarock/passvault/app/core"
	"github.com/ritarock/passvault/infra"
)

type executor struct {
	subCommand subCommand
	store      infra.Store[core.Vault]
}

func newExecutor(subCommand subCommand, filpath string) *executor {
	store := infra.NewStore[core.Vault](filpath)
	return &executor{
		subCommand: subCommand,
		store:      *store,
	}
}

func (e *executor) execute() error {
	switch e.subCommand.action {
	case Generate:
		if e.subCommand.help {
			help(e.subCommand.action)
			return nil
		}
		return generate(e.subCommand.generateData, &e.store)
	case List:
		if e.subCommand.help {
			help(e.subCommand.action)
			return nil
		}
		return list(&e.store)
	case Get:
		if e.subCommand.help {
			help(e.subCommand.action)
			return nil
		}
		return get(e.subCommand.getData, &e.store)
	case Help:
		help(Help)
	}
	return nil
}

func generate(generateData generateData, store *infra.Store[core.Vault]) error {
	vault := core.Vault{
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

func list(store *infra.Store[core.Vault]) error {
	vaults, err := store.Read()
	if err != nil {
		return err
	}

	for i, v := range vaults {
		fmt.Printf("%d - Title: %s, URL: %s\n", i+1, v.Title, v.Url)
	}

	return nil
}

func get(getData int, store *infra.Store[core.Vault]) error {
	vaults, err := store.Read()
	if err != nil {
		return err
	}

	index := getData - 1
	var vault core.Vault
	for i, v := range vaults {
		if i == index {
			vault = v
		}
	}
	if vault.Code == "" {
		return errors.New("not found id")
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
		fmt.Println(GenerateHelp)
	case List:
		fmt.Println(HelpText)
	case Get:
		fmt.Println(GetHelp)
	case Help:
		fmt.Println(HelpText)
	}
}
