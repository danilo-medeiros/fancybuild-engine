package reader

import (
	"encoding/json"
	"fmt"

	"github.com/danilo-medeiros/fancybuild/engine/pkg/entities"
	"github.com/go-playground/validator/v10"
)

const (
	ErrorMessage = "Validation error"
)

type FieldError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

type ValidationError struct {
	Message string        `json:"message"`
	Errors  []*FieldError `json:"errors"`
}

func (*ValidationError) Error() string {
	return ErrorMessage
}

type Reader interface {
	Read([]byte, *entities.Definitions) error
	Validate(*entities.Definitions) *ValidationError
}

type reader struct{}

func (r *reader) Read(data []byte, output *entities.Definitions) error {
	err := json.Unmarshal(data, output)

	if err != nil {
		return fmt.Errorf("error while unmarshaling data: %w", err)
	}

	for _, entity := range output.App.Entities {
		entity.Definitions = output

		for _, action := range entity.Actions {
			action.Entity = entity
		}
	}

	return nil
}

func (r *reader) Validate(definitions *entities.Definitions) *ValidationError {
	errors := make([]*FieldError, 0)

	validate := validator.New()
	err := validate.Struct(definitions)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element FieldError
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()

			errors = append(errors, &element)
		}

		if len(errors) > 0 {
			return &ValidationError{
				Message: ErrorMessage,
				Errors:  errors,
			}
		}
	}

	customValidationError := customValidation(definitions)

	if customValidationError != nil {
		return customValidationError
	}

	return nil
}

func customValidation(definitions *entities.Definitions) *ValidationError {
	errors := make([]*FieldError, 0)

	// Validate entities
	for index, entity := range definitions.App.Entities {
		if !entity.IsNested() && entity.Persisted && len(entity.Actions) == 0 {
			errors = append(errors, &FieldError{
				Field: fmt.Sprintf("app.entities[%v].actions", index),
				Tag:   "not empty",
			})
		}
	}

	if len(errors) > 0 {
		return &ValidationError{
			Message: "validation error",
			Errors:  errors,
		}
	}

	return nil
}

func NewReader() Reader {
	return &reader{}
}
