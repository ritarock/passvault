package app

import (
	"errors"
	"os/user"
	"strconv"
	"strings"
)

type handler struct {
	args []string
	sub  sub
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
	filePath := current.HomeDir + "/vault.json"

	executor := newExecutor(h.sub, filePath)
	return executor.execute()
}

func (h *handler) validate() error {
	if len(h.args) == 0 {
		return errors.New("not enough args")
	}

	action := strings.ToUpper(h.args[0])
	switch action {
	case Generate.string(), List.string(), Get.string(), Help.string():
	default:
		return errors.New("invalid args")
	}

	switch action {
	case Generate.string():
		if len(h.args) > 3 {
			return errors.New("require: generate Title URL")
		}
		if strings.ToUpper(h.args[1]) == Help.string() {
			return nil
		}
		if len(h.args) != 3 {
			return errors.New("require: generate Title URL")
		}
		return nil
	case List.string():
		if len(h.args) == 2 {
			if strings.ToUpper(h.args[1]) == Help.string() {
				return nil
			}
		}
		if len(h.args) != 1 {
			return errors.New("require: list")
		}
		return nil
	case Get.string():
		if len(h.args) > 2 {
			return errors.New("require: get ID")
		}
		if len(h.args) != 2 {
			return errors.New("require: get ID")
		}
		if strings.ToUpper(h.args[1]) == Help.string() {
			return nil
		}
		i, err := strconv.Atoi(h.args[1])
		if err != nil {
			return errors.New("require: get ID, ID should be a number")
		}
		if i == 0 {
			return errors.New("require: get ID, ID more than 1")
		}
		return nil
	case Help.string():
		return nil
	}

	return nil
}

func (h *handler) mapper() {
	action := strings.ToUpper(h.args[0])
	switch action {
	case Generate.string():
		h.sub.action = Generate
		if strings.ToUpper(h.args[1]) == Help.string() {
			h.sub.help = true
			return
		}
		h.sub.generateData.title = h.args[1]
		h.sub.generateData.url = h.args[2]
		return
	case List.string():
		h.sub.action = List
		if len(h.args) == 2 {
			h.sub.help = true
		}
		return
	case Get.string():
		h.sub.action = Get
		if strings.ToUpper(h.args[1]) == Help.string() {
			h.sub.help = true
			return
		}
		i, _ := strconv.Atoi(h.args[1])
		h.sub.getData = i
		return
	case Help.string():
		h.sub.action = Help
	}
}
