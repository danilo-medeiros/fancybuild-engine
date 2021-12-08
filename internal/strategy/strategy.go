package strategy

import (
	"github.com/danilo-medeiros/fancybuild/engine/internal/strategy/golang"
	"github.com/danilo-medeiros/fancybuild/engine/pkg/entities"
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
