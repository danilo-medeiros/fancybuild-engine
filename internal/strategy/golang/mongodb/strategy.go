package mongodb

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/danilo-medeiros/fancybuild/engine/internal/entities"
)

type strategy struct {
	*entities.Definitions
	FileMap map[string]*entities.File
}

func (s *strategy) BuildFileMap() (map[string]*entities.File, error) {
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
			FinalPath:    "internal/router/router.go",
			TemplatePath: "go_mongodb_router.tmpl",
		},
		"entities": {
			FinalPath:    "internal/entities/entities.go",
			TemplatePath: "go_mongodb_entities.tmpl",
		},
		"healthcheck": {
			FinalPath:    "internal/healthcheck/controller.go",
			TemplatePath: "go_mongodb_healthcheck.tmpl",
		},
		"database": {
			FinalPath:    "internal/database/database.go",
			TemplatePath: "go_mongodb_database.tmpl",
		},
		"auth_controller": {
			FinalPath:    "internal/auth/controller.go",
			TemplatePath: "go_auth_controller.tmpl",
		},
		"auth_handler": {
			FinalPath:    "internal/auth/handler.go",
			TemplatePath: "go_auth_handler.tmpl",
		},
	}

	for _, file := range fileMap {
		file.Data = s.Definitions
	}

	for _, entity := range s.Definitions.App.Entities {
		var data struct {
			*entities.Definitions
			entities.Entity
		}

		data.Definitions = s.Definitions
		data.Entity = entity

		if hasController(entity.Name, s.Definitions) {
			fileMap[fmt.Sprintf("%s_controller", entity.Name)] = &entities.File{
				FinalPath:    fmt.Sprintf("internal/%s/controller.go", entity.Name),
				TemplatePath: "go_controller.tmpl",
				Data:         data,
			}
		}

		if hasService(entity.Name, s.Definitions) {
			fileMap[fmt.Sprintf("%s_service", entity.Name)] = &entities.File{
				FinalPath:    fmt.Sprintf("internal/%s/service.go", entity.Name),
				TemplatePath: "go_service.tmpl",
				Data:         data,
			}
		}

		if hasRepository(entity.Name, s.Definitions) {
			fileMap[fmt.Sprintf("%s_repository", entity.Name)] = &entities.File{
				FinalPath:    fmt.Sprintf("internal/%s/repository.go", entity.Name),
				TemplatePath: "go_mongodb_repository.tmpl",
				Data:         data,
			}
		}
	}

	/* entityMap := map[string]*entities.File{
		"controller": {FinalPath: "./internal/{{lower .Name}}/controller.go", TemplatePath: "go_controller.tmpl"},
		"repository": {FinalPath: "./internal/{{lower .Name}}/repository.go", TemplatePath: "go_mongodb_repository.tmpl"},
		"service":    {FinalPath: "./internal/{{lower .Name}}/service.go", TemplatePath: "go_mongodb_service.tmpl"},
	}

	endpointMap := map[string]*entities.File{
		"controller": {FinalPath: "./internal/{{lower .Name}}/controller.go", TemplatePath: "go_controller.tmpl"},
	} */

	err := s.renderFileMap(fileMap)

	if err != nil {
		return nil, fmt.Errorf("error rendering file map: %v", err)
	}

	s.FileMap = fileMap

	return fileMap, nil
}

func (s *strategy) BuildPostActions() error {
	for _, file := range s.FileMap {
		cmd := exec.Command("go", "fmt", fmt.Sprintf("tmp/%s/%s", s.Definitions.Id, file.FinalPath))
		var errb bytes.Buffer
		cmd.Stderr = &errb
		err := cmd.Run()

		if err != nil {
			return fmt.Errorf("on running command: %v: %s", err, errb.String())
		}
	}
	return nil
}

func (s *strategy) renderFileMap(fileMap map[string]*entities.File) error {
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
			return hasController(entity, s.Definitions)
		},
		"belongsToAuthenticatedEntity": func(entity string) bool {
			return belongsToAuthenticatedEntity(entity, s.Definitions)
		},
		"entityHasAction": func(entity string, action string) bool {
			for _, ent := range s.Definitions.App.Entities {
				if ent.Name == entity {
					for _, act := range ent.Actions {
						if act.Type == action {
							return true
						}
					}
				}
			}
			return false
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
			return fmt.Errorf("parsing template file %s: %v", key, err)
		}

		err = parsedTemplate.Execute(&sb, file.Data)

		if err != nil {
			return fmt.Errorf("rendering template %s: %v", key, err)
		}

		result := sb.String()
		file.Result = result
	}

	return nil
}

func hasController(entity string, definitions *entities.Definitions) bool {
	for _, r := range definitions.App.Relationships {
		if r.Nested && r.Item2 == entity && r.Type == "hasMany" {
			return false
		}
	}
	return true
}

func hasService(entity string, definitions *entities.Definitions) bool {
	for _, r := range definitions.App.Relationships {
		if r.Nested && r.Item2 == entity && r.Type == "hasMany" {
			return false
		}
	}
	return true
}

func hasRepository(entity string, definitions *entities.Definitions) bool {
	for _, r := range definitions.App.Relationships {
		if r.Nested && r.Item2 == entity && r.Type == "hasMany" {
			return false
		}
	}
	return true
}

func belongsToAuthenticatedEntity(entity string, definitions *entities.Definitions) bool {
	for _, r := range definitions.App.Relationships {
		if r.Item2 == entity && r.Type == "privateHasMany" {
			return true
		}
	}
	return false
}

func NewStrategy(definitions *entities.Definitions) entities.Strategy {
	return &strategy{Definitions: definitions}
}
