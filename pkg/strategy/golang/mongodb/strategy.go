package mongodb

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/danilo-medeiros/fancybuild/engine/internal/templates"
	"github.com/danilo-medeiros/fancybuild/engine/pkg/entities"
)

type strategy struct {
	*entities.Definitions
	FileMap map[string]*entities.File
}

func (s *strategy) BuildFileMap() (map[string]*entities.File, error) {
	fileMap := map[string]*entities.File{
		"main": {
			FinalPath:    "main.go",
			TemplatePath: "go/main.tmpl",
		},
		"app": {
			FinalPath:    "internal/app/app.go",
			TemplatePath: "go/app.tmpl",
		},
		"validator": {
			FinalPath:    "internal/validator/validator.go",
			TemplatePath: "go/validator.tmpl",
		},
		"router": {
			FinalPath:    "internal/router/router.go",
			TemplatePath: "go/mongodb/router.tmpl",
		},
		"entities": {
			FinalPath:    "internal/entities/entities.go",
			TemplatePath: "go/mongodb/entities.tmpl",
		},
		"health": {
			FinalPath:    "internal/health/controller.go",
			TemplatePath: "go/mongodb/health.tmpl",
		},
		"database": {
			FinalPath:    "internal/database/database.go",
			TemplatePath: "go/mongodb/database.tmpl",
		},
		"error_handler": {
			FinalPath:    "internal/errors/handler.go",
			TemplatePath: "go/error_handler.tmpl",
		},
		"gitignore": {
			FinalPath:    ".gitignore",
			TemplatePath: "go/gitignore.tmpl",
		},
		"env": {
			FinalPath:    ".env",
			TemplatePath: "go/mongodb/env.tmpl",
		},
		"env_test": {
			FinalPath:    ".env.test",
			TemplatePath: "go/mongodb/env_test.tmpl",
		},
		"readme": {
			FinalPath:    "README.md",
			TemplatePath: "go/readme.tmpl",
		},
		"main_test": {
			FinalPath:    "test/main_test.go",
			TemplatePath: "go/main_test.tmpl",
		},
		"test_utils": {
			FinalPath:    "test/utils/utils.go",
			TemplatePath: "go/test_utils.tmpl",
		},
		"dockerfile": {
			FinalPath:    "Dockerfile",
			TemplatePath: "go/mongodb/dockerfile.tmpl",
		},
		"docker-compose": {
			FinalPath:    "docker-compose.yml",
			TemplatePath: "go/mongodb/docker-compose.tmpl",
		},
	}

	if s.Definitions.HasAuthentication() {
		fileMap["auth_controller"] = &entities.File{
			FinalPath:    "internal/auth/controller.go",
			TemplatePath: "go/auth_controller.tmpl",
		}
		fileMap["auth_handler"] = &entities.File{
			FinalPath:    "internal/auth/handler.go",
			TemplatePath: "go/auth_handler.tmpl",
		}
		fileMap["auth_service"] = &entities.File{
			FinalPath:    "internal/auth/service.go",
			TemplatePath: "go/auth_service.tmpl",
		}
		fileMap["auth_test"] = &entities.File{
			FinalPath:    "test/auth/auth_test.go",
			TemplatePath: "go/auth_test.tmpl",
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
				TemplatePath: "go/controller.tmpl",
				Data:         data,
			}

			// TODO: remove this later, currently we only have test coverage for create action
			if entity.HasAction("create") {
				fileMap[fmt.Sprintf("%s_controller_test", entity.Name)] = &entities.File{
					FinalPath:    fmt.Sprintf("test/%s/controller_test.go", entity.Name),
					TemplatePath: "go/controller_test.tmpl",
					Data:         data,
				}
			}
		}

		if entity.HasService() {
			fileMap[fmt.Sprintf("%s_service", entity.Name)] = &entities.File{
				FinalPath:    fmt.Sprintf("internal/%s/service.go", entity.Name),
				TemplatePath: "go/service.tmpl",
				Data:         data,
			}
		}

		if entity.HasRepository() {
			fileMap[fmt.Sprintf("%s_repository", entity.Name)] = &entities.File{
				FinalPath:    fmt.Sprintf("internal/%s/repository.go", entity.Name),
				TemplatePath: "go/mongodb/repository.tmpl",
				Data:         data,
			}
		}

		fileMap[fmt.Sprintf("%s_entity", entity.Name)] = &entities.File{
			FinalPath:    fmt.Sprintf("internal/entities/%s.go", entity.Name),
			TemplatePath: "go/mongodb/entity.tmpl",
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

	testCommand := exec.Command("go", "test", "./...")
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
	funcMap := templates.DefaultFuncMap()
	funcMap["buildValidations"] = buildValidations
	funcMap["mapSort"] = mapSort
	funcMap["jsonMarshal"] = jsonMarshal

	for key, file := range fileMap {
		result, err := templates.Render(&templates.Template{
			Path:    file.TemplatePath,
			Name:    key,
			Data:    file.Data,
			FuncMap: funcMap,
		})

		if err != nil {
			return err
		}

		file.Result = templates.SimpleFormat(result)
	}

	return nil
}

func NewStrategy(definitions *entities.Definitions) entities.Strategy {
	return &strategy{Definitions: definitions}
}
