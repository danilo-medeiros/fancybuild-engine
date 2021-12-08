package builder

import (
	"fmt"
	"os"
	"strings"

	"github.com/danilo-medeiros/fancybuild/engine/pkg/entities"
)

type Builder interface {
	Build(*entities.Definitions, entities.Strategy) error
}

type builder struct {
	OutputFolder string
}

func (b *builder) Build(definitions *entities.Definitions, strategy entities.Strategy) error {
	fileMap, err := strategy.BuildFileMap()

	if err != nil {
		return err
	}

	projectPath := fmt.Sprintf("%s/%s/%s", b.OutputFolder, definitions.Id, definitions.App.Name)

	for _, file := range fileMap {
		pathSlices := strings.Split(file.FinalPath, "/")
		pathSlices = pathSlices[0 : len(pathSlices)-1]
		filePath := strings.Join(pathSlices, "/")

		err = os.MkdirAll(fmt.Sprintf("%s/%s", projectPath, filePath), 0744)

		if err != nil {
			return fmt.Errorf("on creating dir: %v", err)
		}

		f, err := os.Create(fmt.Sprintf("%s/%s", projectPath, file.FinalPath))

		if err != nil {
			return fmt.Errorf("on creating file %s: %v", file.FinalPath, err)
		}

		defer f.Close()

		_, err = f.WriteString(file.Result)

		if err != nil {
			return fmt.Errorf("on writing file %s: %v", file.FinalPath, err)
		}
	}

	err = strategy.BuildPostActions(projectPath)

	if err != nil {
		return fmt.Errorf("on build post actions: %s", err)
	}

	return nil
}

func NewBuilder(outputFolder string) Builder {
	return &builder{outputFolder}
}
