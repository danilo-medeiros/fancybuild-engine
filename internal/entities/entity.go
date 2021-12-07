package entities

type Entity struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Fields      []Field      `json:"fields"`
	Timestamps  bool         `json:"timestamps"`
	Actions     []*Action    `json:"actions"`
	Persisted   bool         `json:"persisted"`
	Indexes     []*Index     `json:"indexes"`
	Definitions *Definitions `json:"-"`
}

func (e Entity) IsNested() bool {
	for _, r := range e.Definitions.App.Relationships {
		if r.Item2 == e.Name && r.Type == "hasMany" && r.Nested {
			return true
		}
	}
	return false
}

func (e Entity) HasController() bool {
	return e.Persisted && !e.IsNested()
}

func (e Entity) HasService() bool {
	return e.Persisted && !e.IsNested()
}

func (e Entity) HasRepository() bool {
	return e.Persisted && !e.IsNested()
}

func (e Entity) HasAction(action string) bool {
	for _, act := range e.Actions {
		if act.Type == action {
			return true
		}
	}
	return false
}

func (e Entity) BelongsToAuthenticatedEntity() bool {
	for _, owner := range e.BelongsTo() {
		if owner.IsUsedForAuthentication() {
			return true
		}
	}
	return false
}

func (e Entity) IsAuthenticated() bool {
	for _, action := range e.Actions {
		if !action.Authenticated {
			return false
		}
	}
	return true
}

func (e Entity) HasMany() []*Entity {
	result := make([]*Entity, 0)

	for _, rel := range e.Definitions.App.Relationships {
		if rel.Item1 == e.Name && rel.IsTypeHasMany() {
			result = append(result, e.Definitions.FindEntity(rel.Item2))
		}
	}

	return result
}

func (e Entity) IsNestedIn(entity *Entity) bool {
	for _, rel := range e.Definitions.App.Relationships {
		if rel.Item1 == entity.Name && rel.Item2 == e.Name && rel.Nested && rel.IsTypeHasMany() {
			return true
		}
	}
	return false
}

func (e Entity) BelongsTo() []*Entity {
	result := make([]*Entity, 0)
	for _, rel := range e.Definitions.App.Relationships {
		if rel.Item2 == e.Name && !rel.Nested {
			result = append(result, e.Definitions.FindEntity(rel.Item1))
		}
	}
	return result
}

func (e Entity) IsUsedForAuthentication() bool {
	return e.Name == e.Definitions.AuthEntity().Name
}

func (e Entity) HasIndexes() bool {
	return len(e.Indexes) > 0
}
