package golang

import (
	"github.com/danilo-medeiros/fancybuild/engine/internal/strategy/golang/mongodb"
	"github.com/danilo-medeiros/fancybuild/engine/pkg/entities"
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
