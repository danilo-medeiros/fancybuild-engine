package strategy

import (
	"github.com/danilo-medeiros/fancybuild/engine/internal/entities"
	"github.com/danilo-medeiros/fancybuild/engine/internal/strategy/golang"
)

const (
	GoLang = "go"
)

func NewStrategy(definitions *entities.Definitions) entities.Strategy {
	switch definitions.App.Stack.Language {
	case GoLang:
		return golang.NewStrategy(definitions)
	}

	return nil
}
