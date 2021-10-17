package mongodb

import "github.com/danilo-medeiros/fancybuild/engine/internal/entities"

type strategy struct {
}

func (s *strategy) BuildFileMap() map[string]string {
	return make(map[string]string)
}

func NewStrategy(definitions *entities.Definitions) entities.Strategy {
	return &strategy{}
}
