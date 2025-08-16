package app

import (
	"errors"
	"os/user"
	"strconv"
	"strings"
)

type handler struct {
	args       []string
	subCommand subCommand
}

func NewHandler(args []string) *handler {
	return &handler{
		args: args,
	}
}

func (h *handler) Run() error {
	if err := h.validate(); err != nil {
		return err
	}
	h.mapper()

	current, err := user.Current()
	if err != nil {
		return err
	}
	filepath := current.HomeDir + "/vault.json"

	executor := newExecutor(h.subCommand, filepath)

	return executor.execute()
}

func (h *handler) validate() error {
	if len(h.args) == 0 {
		return errors.New("not enough args")
	}
	action := strings.ToUpper(h.args[0])

	switch action {
	case Generate.string():
		if len(h.args) > 3 {
			return errors.New("require: generate <title> <url>")
		}
		if len(h.args) == 1 {
			return errors.New("require: generate <title> <url>")
		}
		if isHelp(h.args[1]) {
			return nil
		}
		if len(h.args) != 3 {
			return errors.New("require: generate <title> <url>")
		}
		return nil
	case List.string():
		if len(h.args) == 2 {
			if isHelp(h.args[1]) {
				return nil
			}
		}
		if len(h.args) != 1 {
			return errors.New("require: list")
		}
		return nil
	case Get.string():
		if len(h.args) > 2 {
			return errors.New("require: get <id>")
		}
		if len(h.args) != 2 {
			return errors.New("require: get <id>")
		}
		if isHelp(h.args[1]) {
			return nil
		}
		i, err := strconv.Atoi(h.args[1])
		if err != nil {
			return errors.New("require: get <id>, id should be a number")
		}
		if i == 0 {
			return errors.New("require: get ID, ID more than 1")
		}
		return nil
	case Help.string():
		return nil
	default:
		return errors.New("invalid args")
	}
}

func (h *handler) mapper() {
	action := strings.ToUpper(h.args[0])

	switch action {
	case Generate.string():
		h.subCommand.action = Generate
		if isHelp(h.args[1]) {
			h.subCommand.help = true
			return
		}
		h.subCommand.generateData.title = h.args[1]
		h.subCommand.generateData.url = h.args[2]
		return
	case List.string():
		h.subCommand.action = List
		if len(h.args) == 2 {
			h.subCommand.help = true
		}
		return
	case Get.string():
		h.subCommand.action = Get
		if isHelp(h.args[1]) {
			h.subCommand.help = true
			return
		}
		i, _ := strconv.Atoi(h.args[1])
		h.subCommand.getData = i
		return
	case Help.string():
		h.subCommand.action = Help
	}
}

func isHelp(arg string) bool {
	return strings.ToUpper(arg) == Help.string()
}
