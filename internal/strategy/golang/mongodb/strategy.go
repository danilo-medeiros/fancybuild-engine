package mongodb

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/danilo-medeiros/fancybuild/engine/internal/entities"
)

type strategy struct {
}

func (s *strategy) BuildFileMap(definitions *entities.Definitions) (map[string]*entities.File, error) {
	fileMap := map[string]*entities.File{
		"main": {
			FinalPath:    "main.go",
			TemplatePath: "go_main.tmpl",
		},
		"validator": {
			FinalPath:    "internal/validator/validator.go",
			TemplatePath: "go_validator.tmpl",
		},
		"router": {
			FinalPath:    "./internal/router/router.go",
			TemplatePath: "go_mongodb_router.tmpl",
		},
		"entities": {
			FinalPath:    "./internal/entities.go",
			TemplatePath: "go_mongodb_entities.tmpl",
		},
		"healthcheck": {
			FinalPath:    "./internal/healthcheck/controller.go",
			TemplatePath: "go_mongodb_healthcheck.tmpl",
		},
		"database": {
			FinalPath:    "./internal/database/database.go",
			TemplatePath: "go_mongodb_database.tmpl",
		},
		"auth_controller": {
			FinalPath:    "./internal/auth/controller.go",
			TemplatePath: "go_auth_controller.tmpl",
		},
		"auth_handler": {
			FinalPath:    "./internal/auth/handler.go",
			TemplatePath: "go_auth_handler.tmpl",
		},
	}

	/* entityMap := map[string]*entities.File{
		"controller": {FinalPath: "./internal/{{lower .Name}}/controller.go", TemplatePath: "go_controller.tmpl"},
		"repository": {FinalPath: "./internal/{{lower .Name}}/repository.go", TemplatePath: "go_mongodb_repository.tmpl"},
		"service":    {FinalPath: "./internal/{{lower .Name}}/service.go", TemplatePath: "go_mongodb_service.tmpl"},
	}

	endpointMap := map[string]*entities.File{
		"controller": {FinalPath: "./internal/{{lower .Name}}/controller.go", TemplatePath: "go_controller.tmpl"},
	} */

	err := renderFileMap(definitions, fileMap)

	if err != nil {
		return nil, fmt.Errorf("error rendering file map: %v", err)
	}

	return fileMap, nil
}

func renderFileMap(definitions *entities.Definitions, fileMap map[string]*entities.File) error {
	funcMap := template.FuncMap{
		"capitalize": func(text string) string {
			splitted := strings.Split(text, "")
			splitted[0] = strings.ToUpper(splitted[0])
			return strings.Join(splitted, "")
		},
		"camelize": func(text ...string) string {
			for index, part := range text {
				if index == 0 {
					continue
				}

				text[index] = fmt.Sprintf("%s%s", strings.ToUpper(string(part[0])), part[1:])
			}

			return strings.Join(text, "")
		},
		"slice": func(text string, start int, end int) string {
			return text[start:end]
		},
		"split": strings.Split,
		"pluralize": func(text string) string {
			// TODO: Implement pluralization rules: https://www.grammarly.com/blog/plural-nouns/
			return fmt.Sprintf("%ss", text)
		},
		"replaceAll": strings.ReplaceAll,
		"hasController": func(entity string) bool {
			for _, r := range definitions.App.Relationships {
				if r.Nested && r.Item2 == entity && r.Type == "hasMany" {
					return false
				}
			}
			return true
		},
	}

	for key, file := range fileMap {
		sb := strings.Builder{}
		path := fmt.Sprintf("internal/templates/%s", file.TemplatePath)
		templateContent, err := os.ReadFile(path)

		if err != nil {
			return fmt.Errorf("reading template file %s: %v", key, err)
		}

		parsedTemplate, err := template.New(key).Funcs(funcMap).Parse(string(templateContent))

		if err != nil {
			return fmt.Errorf("error parsing template file %s: %v", key, err)
		}

		err = parsedTemplate.Execute(&sb, definitions)

		if err != nil {
			return fmt.Errorf("error rendering template %s: %v", key, err)
		}

		result := sb.String()
		file.Result = result
	}

	return nil
}

func NewStrategy(definitions *entities.Definitions) entities.Strategy {
	return &strategy{}
}
