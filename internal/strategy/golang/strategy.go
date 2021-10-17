package golang

import (
	"github.com/danilo-medeiros/fancybuild/engine/internal/entities"
	"github.com/danilo-medeiros/fancybuild/engine/internal/strategy/golang/mongodb"
)

const (
	MongoDB = "mongodb"
)

func NewStrategy(definitions *entities.Definitions) entities.Strategy {
	switch definitions.App.Stack.Database {
	case MongoDB:
		return mongodb.NewStrategy(definitions)
	}

	return nil
}
