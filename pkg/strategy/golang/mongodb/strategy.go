package mongodb

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/danilo-medeiros/fancybuild/engine/internal/templates"
	"github.com/danilo-medeiros/fancybuild/engine/pkg/entities"
)

type strategy struct {
	*entities.Definitions
	FileMap map[string]*entities.File
}

func (s *strategy) BuildFileMap() (map[string]*entities.File, error) {
	for _, entity := range s.Definitions.App.Entities {
		entity.Definitions = s.Definitions

		for _, action := range entity.Actions {
			action.Entity = entity
		}
	}

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
		"health": {
			FinalPath:    "internal/health/controller.go",
			TemplatePath: "go_mongodb_health.tmpl",
		},
		"database": {
			FinalPath:    "internal/database/database.go",
			TemplatePath: "go_mongodb_database.tmpl",
		},
		"error_handler": {
			FinalPath:    "internal/errors/handler.go",
			TemplatePath: "go_error_handler.tmpl",
		},
		"gitignore": {
			FinalPath:    ".gitignore",
			TemplatePath: "go_gitignore.tmpl",
		},
		"env": {
			FinalPath:    ".env",
			TemplatePath: "go_mongodb_env.tmpl",
		},
		"env_test": {
			FinalPath:    ".env.test",
			TemplatePath: "go_mongodb_env_test.tmpl",
		},
		"readme": {
			FinalPath:    "README.md",
			TemplatePath: "go_readme.tmpl",
		},
		"main_test": {
			FinalPath:    "main_test.go",
			TemplatePath: "go_main_test.tmpl",
		},
		"dockerfile": {
			FinalPath:    "Dockerfile",
			TemplatePath: "go_mongodb_dockerfile.tmpl",
		},
		"docker-compose": {
			FinalPath:    "docker-compose.yml",
			TemplatePath: "go_mongodb_docker-compose.tmpl",
		},
	}

	if len(s.Definitions.App.Authentication.Entity) != 0 {
		fileMap["auth_controller"] = &entities.File{
			FinalPath:    "internal/auth/controller.go",
			TemplatePath: "go_auth_controller.tmpl",
		}
		fileMap["auth_handler"] = &entities.File{
			FinalPath:    "internal/auth/handler.go",
			TemplatePath: "go_auth_handler.tmpl",
		}
		fileMap["auth_service"] = &entities.File{
			FinalPath:    "internal/auth/service.go",
			TemplatePath: "go_auth_service.tmpl",
		}
	}

	for _, file := range fileMap {
		file.Data = &s.Definitions
	}

	for _, entity := range s.Definitions.App.Entities {
		var data struct {
			*entities.Definitions
			*entities.Entity
		}

		data.Definitions = s.Definitions
		data.Entity = entity

		if entity.HasController() {
			fileMap[fmt.Sprintf("%s_controller", entity.Name)] = &entities.File{
				FinalPath:    fmt.Sprintf("internal/%s/controller.go", entity.Name),
				TemplatePath: "go_controller.tmpl",
				Data:         data,
			}
		}

		if entity.HasService() {
			fileMap[fmt.Sprintf("%s_service", entity.Name)] = &entities.File{
				FinalPath:    fmt.Sprintf("internal/%s/service.go", entity.Name),
				TemplatePath: "go_service.tmpl",
				Data:         data,
			}
		}

		if entity.HasRepository() {
			fileMap[fmt.Sprintf("%s_repository", entity.Name)] = &entities.File{
				FinalPath:    fmt.Sprintf("internal/%s/repository.go", entity.Name),
				TemplatePath: "go_mongodb_repository.tmpl",
				Data:         data,
			}
		}

		fileMap[fmt.Sprintf("%s_entity", entity.Name)] = &entities.File{
			FinalPath:    fmt.Sprintf("internal/entities/%s.go", entity.Name),
			TemplatePath: "go_mongodb_entity.tmpl",
			Data:         data,
		}
	}

	err := s.renderFileMap(fileMap)

	if err != nil {
		return nil, fmt.Errorf("error rendering file map: %v", err)
	}

	s.FileMap = fileMap

	return fileMap, nil
}

func (s *strategy) BuildPostActions(projectPath string) error {
	commands := make([]*exec.Cmd, 0)
	isGoFileRegexp := regexp.MustCompile(".go$")

	for _, file := range s.FileMap {
		isGoFile := isGoFileRegexp.MatchString(file.FinalPath)

		if isGoFile {
			cmd := exec.Command("go", "fmt", fmt.Sprintf("%s/%s", projectPath, file.FinalPath))
			commands = append(commands, cmd)
		}
	}

	modCommand := exec.Command("go", "mod", "init", s.App.Repository)
	modCommand.Dir = projectPath
	commands = append(commands, modCommand)

	tidyCommand := exec.Command("go", "mod", "tidy")
	tidyCommand.Dir = projectPath
	commands = append(commands, tidyCommand)

	buildCommand := exec.Command("go", "build")
	buildCommand.Dir = projectPath
	commands = append(commands, buildCommand)

	testCommand := exec.Command("go", "test", ".")
	testCommand.Dir = projectPath
	commands = append(commands, testCommand)

	for _, command := range commands {
		var errb bytes.Buffer
		command.Stderr = &errb
		err := command.Run()

		if err != nil {
			return fmt.Errorf("on running command %s: %v: %s", command.String(), err, errb.String())
		}
	}

	return nil
}

func (s *strategy) renderFileMap(fileMap map[string]*entities.File) error {
	funcMap := template.FuncMap{
		"capitalize":       templates.Capitalize,
		"camelize":         templates.Camelize,
		"slice":            templates.Slice,
		"split":            strings.Split,
		"pluralize":        templates.Pluralize,
		"replaceAll":       strings.ReplaceAll,
		"buildValidations": buildValidations,
		"empty":            templates.Empty,
		"join":             strings.Join,
		"mapSort":          mapSort,
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
		file.Result = templates.SimpleFormat(result)
	}

	return nil
}

func buildValidations(field entities.Field) string {
	validations := make([]string, 0)

	for _, validation := range field.Validations {
		switch validation.Name {
		case "required":
			validations = append(validations, "required")
		case "type":
			validations = append(validations, validation.Value)
		default:
			validations = append(validations, fmt.Sprintf("%s=%s", validation.Name, validation.Value))
		}
	}

	return strings.Join(validations, ",")
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

func NewStrategy(definitions *entities.Definitions) entities.Strategy {
	return &strategy{Definitions: definitions}
}
