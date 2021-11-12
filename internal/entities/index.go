package entities

type IndexField struct {
	Name string `json:"name"`
	Sort string `json:"sort"`
}

type Index struct {
	Fields []*IndexField `json:"fields"`
	Unique bool          `json:"unique"`
}
