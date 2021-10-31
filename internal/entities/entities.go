package entities

type Validation struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Field struct {
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Validations []Validation `json:"validations"`
}

type Action struct {
	Type string `json:"type"`
}

type Entity struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Fields      []Field  `json:"fields"`
	Timestamps  bool     `json:"timestamps"`
	Actions     []Action `json:"actions"`
}

type Relationship struct {
	Item1  string `json:"item1"`
	Item2  string `json:"item2"`
	Type   string `json:"type"`
	Nested bool   `json:"nested"`
}

type Authentication struct {
	Type   string `json:"type"`
	Entity Entity `json:"entity"`
}

type Stack struct {
	Language string `json:"language"`
	Database string `json:"database"`
}

type App struct {
	Name          string         `json:"name" validate:"min=3,max=50"`
	Description   string         `json:"description" validate:"max=200"`
	Version       string         `json:"version"`
	Repository    string         `json:"repository"`
	Type          string         `json:"type"`
	Stack         Stack          `json:"stack" validate:"dive"`
	Entities      []Entity       `json:"entities" validate:"dive"`
	Relationships []Relationship `json:"relationships" validate:"dive"`
}

type Definitions struct {
	Id      string `json:"id"`
	Version string `json:"version"`
	App     App    `json:"app"`
}

type File struct {
	FinalPath    string
	TemplatePath string
	Result       string
}

type Strategy interface {
	BuildFileMap(*Definitions) (map[string]*File, error)
}
