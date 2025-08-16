package app

type action int

const (
	Generate action = iota + 1
	List
	Get
	Help
)

type generateData struct {
	title string
	url   string
}

type subCommand struct {
	action       action
	help         bool
	generateData generateData
	getData      int
}

func (a action) string() string {
	return []string{
		"GENERATE",
		"LIST",
		"GET",
		"HELP",
	}[a-1]
}
