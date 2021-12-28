package mysql

import (
	"fmt"
	"strings"

	"github.com/danilo-medeiros/fancybuild/engine/pkg/entities"
)

func buildIndexes(entity *entities.Entity) map[string]string {
	result := make(map[string]string)
	indexByField := make(map[string][]string)

	for _, index := range entity.Indexes {
		indexName := fmt.Sprintf("idx_%s", entity.Name)

		for _, field := range index.Fields {
			indexName = fmt.Sprintf("%s_%s", indexName, field.Name)
		}

		var tag string

		if index.Unique {
			tag = fmt.Sprintf("index:%s,unique", indexName)
		} else {
			tag = fmt.Sprintf("index:%s", indexName)
		}

		for _, field := range index.Fields {
			if len(indexByField[field.Name]) == 0 {
				indexByField[field.Name] = make([]string, 0)
			}
			indexByField[field.Name] = append(indexByField[field.Name], tag)
		}
	}

	for key, indexes := range indexByField {
		result[key] = strings.Join(indexes, ";")
	}

	return result
}
