package golang

import (
	"github.com/danilo-medeiros/fancybuild/engine/pkg/entities"
	"github.com/danilo-medeiros/fancybuild/engine/pkg/strategy/golang/mongodb"
	"github.com/danilo-medeiros/fancybuild/engine/pkg/strategy/golang/mysql"
)

const (
	MongoDB = "mongodb"
	MySQL   = "mysql"
)

func NewStrategy(definitions *entities.Definitions) entities.Strategy {
	switch definitions.App.Stack.Database {
	case MongoDB:
		return mongodb.NewStrategy(definitions)
	case MySQL:
		return mysql.NewStrategy(definitions)
	}

	return nil
}
