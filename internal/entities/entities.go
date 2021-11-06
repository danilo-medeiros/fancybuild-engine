package entities

type Validation struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Field struct {
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Validations []Validation `json:"validations"`
	Secret      bool         `json:"secret"`
	Hashed      bool         `json:"hashed"`
}

type Input struct {
	Entity string `json:"entity"`
}

type Output struct {
	Entity string `json:"entity"`
}

type Action struct {
	Type          string `json:"type"`
	Authenticated bool   `json:"authenticated"`
	Input         Input  `json:"input"`
	Output        Output `json:"output"`
}

type Relationship struct {
	Item1  string `json:"item1"`
	Item2  string `json:"item2"`
	Type   string `json:"type"`
	Nested bool   `json:"nested"`
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
