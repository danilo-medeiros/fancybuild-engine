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

type Entity struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Fields      []Field `json:"fields"`
	Timestamps  bool    `json:"timestamps"`
}

type Relationship struct {
	Item1    string `json:"item1"`
	Item2    string `json:"item2"`
	Type     string `json:"type"`
	Embedded bool   `json:"embedded"`
}

type Endpoint struct {
	Method string `json:"method"`
	URI    string `json:"uri"`
	Action string `json:"action"`
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
	Name     string   `json:"name"`
	Stack    Stack    `json:"stack"`
	Version  string   `json:"version"`
	Type     string   `json:"type"`
	Entities []Entity `json:"entities"`
}

type Definitions struct {
	Version string `json:"version"`
	App     App    `json:"app"`
}

type Strategy interface {
	BuildFileMap() map[string]string
}
