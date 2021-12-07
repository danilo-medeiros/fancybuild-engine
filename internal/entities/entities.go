package entities

type Validation struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Input struct {
	Entity string `json:"entity"`
}

type Output struct {
	Entity string `json:"entity"`
}

type Authentication struct {
	Entity string `json:"entity"`
}

type Stack struct {
	Language string `json:"language"`
	Database string `json:"database"`
}

type File struct {
	FinalPath    string
	TemplatePath string
	Result       string
	Data         interface{}
}

type Strategy interface {
	BuildFileMap() (map[string]*File, error)
	BuildPostActions(string) error
}
