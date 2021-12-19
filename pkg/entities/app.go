package entities

type App struct {
	Name           string          `json:"name" validate:"min=3,max=50"`
	Description    string          `json:"description" validate:"max=200"`
	Version        string          `json:"version"`
	Repository     string          `json:"repository"`
	Type           string          `json:"type"`
	Stack          Stack           `json:"stack" validate:"dive"`
	Entities       []*Entity       `json:"entities" validate:"dive"`
	Relationships  []*Relationship `json:"relationships" validate:"dive"`
	Authentication Authentication  `json:"authentication" validate:"dive"`
}
