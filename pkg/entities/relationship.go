package entities

const (
	RelationshipTypeHasMany = "hasMany"
	RelationshipTypeHasOne  = "hasOne"
)

type Relationship struct {
	Item1      string `json:"item1" validate:"required"`
	Item2      string `json:"item2" validate:"required"`
	Type       string `json:"type" validate:"oneof=hasMany hasOne,required"`
	Nested     bool   `json:"nested"`
	Visibility string `json:"visibility"`
}

func (r Relationship) IsTypeHasMany() bool {
	return r.Type == RelationshipTypeHasMany
}

func (r Relationship) IsTypeHasOne() bool {
	return r.Type == RelationshipTypeHasOne
}
