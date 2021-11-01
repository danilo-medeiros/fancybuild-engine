package builder

import (
	"fmt"
	"os"
	"strings"

	"github.com/danilo-medeiros/fancybuild/engine/internal/entities"
)

type Builder interface {
	Build(*entities.Definitions, entities.Strategy) error
}

type builder struct{}

func (*builder) Build(definitions *entities.Definitions, strategy entities.Strategy) error {
	fileMap, err := strategy.BuildFileMap()

	if err != nil {
		return err
	}

	for _, file := range fileMap {
		pathSlices := strings.Split(file.FinalPath, "/")
		pathSlices = pathSlices[0 : len(pathSlices)-1]
		filePath := strings.Join(pathSlices, "/")

		err = os.MkdirAll(fmt.Sprintf("tmp/%s/%s", definitions.Id, filePath), 0744)

		if err != nil {
			return fmt.Errorf("on creating dir %s: %v", definitions.Id, err)
		}

		f, err := os.Create(fmt.Sprintf("tmp/%s/%s", definitions.Id, file.FinalPath))

		if err != nil {
			return fmt.Errorf("on creating file %s: %v", file.FinalPath, err)
		}

		defer f.Close()

		_, err = f.WriteString(file.Result)

		if err != nil {
			return fmt.Errorf("on writing file %s: %v", file.FinalPath, err)
		}
	}

	err = strategy.BuildPostActions()

	if err != nil {
		return fmt.Errorf("on post build actions: %s", err)
	}

	return nil
}

func NewBuilder() Builder {
	return &builder{}
}
