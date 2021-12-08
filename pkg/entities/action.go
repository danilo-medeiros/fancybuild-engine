package entities

type Action struct {
	Type          string  `json:"type"`
	Authenticated bool    `json:"authenticated"`
	Input         Input   `json:"input"`
	Output        Output  `json:"output"`
	Entity        *Entity `json:"-" validate:"-"`
}

func (a Action) IsCreate() bool {
	return a.Type == "create"
}

func (a Action) IsGetAll() bool {
	return a.Type == "getAll"
}

func (a Action) IsGetOne() bool {
	return a.Type == "getOne"
}

func (a Action) IsUpdate() bool {
	return a.Type == "update"
}

func (a Action) IsDelete() bool {
	return a.Type == "delete"
}

func (a Action) HTTPMethod() string {
	switch a.Type {
	case "create":
		return "POST"
	case "getAll", "getOne":
		return "GET"
	case "update":
		return "PUT"
	case "delete":
		return "DELETE"
	}
	return "GET"
}
