package reader

import (
	"encoding/json"
	"fmt"

	"github.com/danilo-medeiros/fancybuild/engine/internal/entities"
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

	return nil
}

func NewReader() Reader {
	return &reader{}
}
