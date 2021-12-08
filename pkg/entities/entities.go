package entities

// Single validation specification for a field. E.g. required, max, min
type Validation struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Input entity that should be used to send data to an action. E.g. some cases
// the entity of the action should not be used to post/put information, instead another entity
// that have some of the fields can be used.
type Input struct {
	Entity string `json:"entity"`
}

// Ouput entity that should be used to return from a action. E.g. some cases
// the entity of the action should not be returned, instead another entity
// that have some of the fields can be used.
type Output struct {
	Entity string `json:"entity"`
}

// The authentication and authorization specifications of the project
type Authentication struct {
	Entity string `json:"entity"`
}

// Defines some specifications of the implementation of the project
type Stack struct {
	Language string `json:"language"` // The language to be used (e.g. go, node, etc...)
	Database string `json:"database"` // The database to be used (e.g. mongodb, mysql, etc...)
}

// Has information about a file that will be mapped in the final project. It is used by the strategy
// to retrieve the information of the template, render with some data and hold the output.
type File struct {
	FinalPath    string
	TemplatePath string
	Result       string
	Data         interface{}
}

// Build a project file map and execute commands in it in order to format, test and do some other actions
type Strategy interface {
	BuildFileMap() (map[string]*File, error)
	BuildPostActions(string) error
}
