package mongodb

import (
	"fmt"
	"strings"

	"github.com/danilo-medeiros/fancybuild/engine/internal/templates"
	"github.com/danilo-medeiros/fancybuild/engine/pkg/entities"
)

func jsonMarshal(entity *entities.Entity) (string, error) {
	funcMap := templates.DefaultFuncMap()
	funcMap["jsonMarshalField"] = jsonMarshalField

	return templates.Render(&templates.Template{
		Path:    "go/json_marshal.tmpl",
		Name:    "json_marshal",
		Data:    entity,
		FuncMap: funcMap,
	})
}

func jsonMarshalField(field *entities.Field) string {
	switch field.Type {
	case "int", "uint", "int32", "int64", "float32", "float64":
		return field.Example()
	case "string":
		return fmt.Sprintf("\"%s\"", field.Example())
	}
	return ""
}

func buildValidations(field *entities.Field, includeRequired bool) string {
	validations := make([]string, 0)

	if !includeRequired {
		validations = append(validations, "omitempty")
	}

	for _, validation := range field.Validations {
		if validation.Name == "required" && !includeRequired {
			continue
		}
		switch validation.Name {
		case "required", "email":
			validations = append(validations, validation.Name)
		default:
			validations = append(validations, fmt.Sprintf("%s=%s", validation.Name, validation.Value))
		}
	}

	result := strings.Join(validations, ",")

	if result == "" {
		return ""
	} else {
		return fmt.Sprintf("validate:\"%s\"", result)
	}
}

func mapSort(sort string) int {
	switch sort {
	case "asc":
		return 1
	case "desc":
		return -1
	}
	return 1
}
