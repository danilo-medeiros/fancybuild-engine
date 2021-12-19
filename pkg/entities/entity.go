package entities

// Single entity specification of the project. E.g. User, Sale, Product, etc...
// Defines metadata about the entity, how the values should be stored and the actions that should be implemented.
type Entity struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Fields      []*Field     `json:"fields"`
	Timestamps  bool         `json:"timestamps"`
	Actions     []*Action    `json:"actions"`
	Persisted   bool         `json:"persisted"`
	Indexes     []*Index     `json:"indexes"`
	Definitions *Definitions `json:"-" validate:"-"`
}

// Check if the entity is nested to another
func (e Entity) IsNested() bool {
	for _, r := range e.Definitions.App.Relationships {
		if r.Item2 == e.Name && r.Nested {
			return true
		}
	}
	return false
}

// Checks if the entity should have a controller file
func (e Entity) HasController() bool {
	return e.Persisted && !e.IsNested()
}

// Checks if the entity should have a service file
func (e Entity) HasService() bool {
	return e.Persisted && !e.IsNested()
}

// Checks if the entity should have a repository file
func (e Entity) HasRepository() bool {
	return e.Persisted && !e.IsNested()
}

// Checks if the entity has a specific action
func (e Entity) HasAction(action string) bool {
	for _, act := range e.Actions {
		if act.Type == action {
			return true
		}
	}
	return false
}

// Checks if the entity belongs to another entity that is authenticated.
// Used to add the logged user metadata to the dto
func (e Entity) BelongsToAuthenticatedEntity() bool {
	for _, owner := range e.BelongsTo() {
		if owner.IsUsedForAuthentication() {
			return true
		}
	}
	return false
}

// Checks if all the actions of the entity require authentication.
func (e Entity) IsAuthenticated() bool {
	for _, action := range e.Actions {
		if !action.Authenticated {
			return false
		}
	}
	return true
}

// Returns all the entities to which this entity has a "hasMany" relationship.
func (e Entity) HasMany() []*Entity {
	result := make([]*Entity, 0)

	for _, rel := range e.Definitions.App.Relationships {
		if rel.Item1 == e.Name && rel.IsTypeHasMany() {
			result = append(result, e.Definitions.FindEntity(rel.Item2))
		}
	}

	return result
}

// Returns all the entities to which this entity has a "hasOne" relationship.
func (e Entity) HasOne() []*Entity {
	result := make([]*Entity, 0)

	for _, rel := range e.Definitions.App.Relationships {
		if rel.Item1 == e.Name && rel.IsTypeHasOne() {
			result = append(result, e.Definitions.FindEntity(rel.Item2))
		}
	}

	return result
}

// Checks if has a nested relationship to the nested entity.
// Used for mongodb cases
func (e Entity) IsNestedIn(entity *Entity) bool {
	for _, rel := range e.Definitions.App.Relationships {
		if rel.Item1 == entity.Name && rel.Item2 == e.Name && rel.Nested {
			return true
		}
	}
	return false
}

// Returns all the entities that belongs to this entity (have any kind relationship with it and is not nested)
func (e Entity) BelongsTo() []*Entity {
	result := make([]*Entity, 0)
	for _, rel := range e.Definitions.App.Relationships {
		if rel.Item2 == e.Name && !rel.Nested {
			result = append(result, e.Definitions.FindEntity(rel.Item1))
		}
	}
	return result
}

// Checks if is used for authentication
func (e Entity) IsUsedForAuthentication() bool {
	return e.Name == e.Definitions.AuthEntity().Name
}

// Checks if has defined indexes
func (e Entity) HasIndexes() bool {
	return len(e.Indexes) > 0
}

// Generates an map with example values for this entity.
// The key is the field (in lowercase) and the value is
// an example value generated within the validation constratins
func (e Entity) Example() map[string]string {
	result := make(map[string]string)

	for _, field := range e.Fields {
		result[field.Name] = field.Example()
	}

	return result
}
