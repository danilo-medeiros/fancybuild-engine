package entities

type Entity struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Fields      []Field  `json:"fields"`
	Timestamps  bool     `json:"timestamps"`
	Actions     []Action `json:"actions"`
	Persisted   bool     `json:"persisted"`
}

func (e Entity) IsAuthenticated() bool {
	for _, action := range e.Actions {
		if !action.Authenticated {
			return false
		}
	}
	return true
}
