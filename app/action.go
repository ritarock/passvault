package app

type action int

const (
	Generate action = iota + 1
	List
	Get
	Help
)

func (a action) string() string {
	return []string{
		"GENERATE",
		"LIST",
		"GET",
		"HELP",
	}[a-1]
}

type generateData struct {
	title string
	url   string
}

type sub struct {
	action       action
	help         bool
	generateData generateData
	getData      int
}
