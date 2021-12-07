package entities

const (
	RelationshipTypeHasMany = "hasMany"
	RelationshipTypeHasOne  = "hasOne"
)

type Relationship struct {
	Item1      string `json:"item1"`
	Item2      string `json:"item2"`
	Type       string `json:"type"`
	Nested     bool   `json:"nested"`
	Visibility string `json:"visibility"`
}

func (r Relationship) IsTypeHasMany() bool {
	return r.Type == RelationshipTypeHasMany
}

func (r Relationship) IsTypeHasOne() bool {
	return r.Type == RelationshipTypeHasOne
}
