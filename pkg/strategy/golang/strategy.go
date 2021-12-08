package golang

import (
	"github.com/danilo-medeiros/fancybuild/engine/pkg/entities"
	"github.com/danilo-medeiros/fancybuild/engine/pkg/strategy/golang/mongodb"
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
