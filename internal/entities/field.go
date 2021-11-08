package entities

type Field struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Validations []*Validation `json:"validations"`
	Secret      bool          `json:"secret"`
	Hashed      bool          `json:"hashed"`
}
