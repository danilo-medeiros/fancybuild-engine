package entities

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

type fieldExampleTestCase interface {
	IsValid(value string) bool
	Field() *Field
	Description() string
}

type stringTestCase struct {
	description string
	field       *Field
}

func parseToIntOrPanic(value string) int {
	res, err := strconv.Atoi(value)

	if err != nil {
		panic(fmt.Sprintf("on parsing value %s: %s", value, err))
	}

	return res
}

func newFieldExampleTestCase(f *Field) fieldExampleTestCase {
	switch f.Type {
	case "string":
		return &stringTestCase{
			field: f,
		}
	default:
		return &intTestCase{
			field: f,
		}
	}
}

func (s *stringTestCase) IsValid(value string) bool {
	valid := true

	for _, validation := range s.field.Validations {
		switch validation.Name {
		case "min", "gte":
			min := parseToIntOrPanic(validation.Value)

			valid = valid && len(value) >= min
		case "max", "lte":
			max := parseToIntOrPanic(validation.Value)

			valid = valid && len(value) <= max
		case "len", "eq":
			eq := parseToIntOrPanic(validation.Value)

			valid = valid && len(value) == eq
		case "gt":
			gt := parseToIntOrPanic(validation.Value)

			valid = valid && len(value) > gt
		case "lt":
			lt := parseToIntOrPanic(validation.Value)
			valid = valid && len(value) < lt
		case "required":
			valid = valid && value != ""
		case "email":
			valid = valid && value == "example.email@example.com"
		case "oneof":
			valid = valid && strings.Split(validation.Value, " ")[0] == value
		default:
			valid = false
		}
	}

	return valid
}

func (s *stringTestCase) Field() *Field {
	return s.field
}

func (s *stringTestCase) Description() string {
	return s.description
}

type intTestCase struct {
	description string
	field       *Field
}

func (i *intTestCase) IsValid(value string) bool {
	valid := true
	exampleValue, err := strconv.Atoi(value)

	if err != nil {
		panic(fmt.Sprintf("on parsing example value %s: %s", value, err))
	}

	for _, validation := range i.field.Validations {
		switch validation.Name {
		case "min", "gte":
			min := parseToIntOrPanic(validation.Value)

			valid = valid && exampleValue >= min
		case "max", "lte":
			max := parseToIntOrPanic(validation.Value)

			valid = valid && exampleValue <= max
		case "len", "eq":
			len := parseToIntOrPanic(validation.Value)

			valid = valid && exampleValue == len
		case "gt":
			gt := parseToIntOrPanic(validation.Value)

			valid = valid && exampleValue > gt
		case "lt":
			lt := parseToIntOrPanic(validation.Value)

			valid = valid && exampleValue < lt
		case "required":
			valid = valid && exampleValue > 0
		case "oneof":
			valid = valid && strings.Split(validation.Value, " ")[0] == value
		default:
			valid = false
		}
	}
	return valid
}

func (i *intTestCase) Field() *Field {
	return i.field
}

func (i *intTestCase) Description() string {
	return i.description
}

func TestFieldExample(t *testing.T) {
	testCases := []fieldExampleTestCase{
		&intTestCase{
			description: "field type int with validation tag max",
			field: &Field{
				Name: "example",
				Type: "int",
				Validations: []*Validation{
					{
						Name: "required",
					},
					{
						Name:  "max",
						Value: "10",
					},
				},
			},
		},
		&intTestCase{
			description: "field type int with validation tag min",
			field: &Field{
				Name: "example",
				Type: "int",
				Validations: []*Validation{
					{
						Name: "required",
					},
					{
						Name:  "min",
						Value: "10",
					},
				},
			},
		},
		&intTestCase{
			description: "field type int with validation tag gte",
			field: &Field{
				Name: "example",
				Type: "int",
				Validations: []*Validation{
					{
						Name: "required",
					},
					{
						Name:  "gte",
						Value: "50",
					},
				},
			},
		},
		&intTestCase{
			description: "field type int with validation tag gt",
			field: &Field{
				Name: "example",
				Type: "int",
				Validations: []*Validation{
					{
						Name:  "gt",
						Value: "50",
					},
				},
			},
		},
		&intTestCase{
			description: "field type int with validation tag lte",
			field: &Field{
				Name: "example",
				Type: "int",
				Validations: []*Validation{
					{
						Name:  "lte",
						Value: "50",
					},
				},
			},
		},
		&intTestCase{
			description: "field type int with validation tag lt",
			field: &Field{
				Name: "example",
				Type: "int",
				Validations: []*Validation{
					{
						Name:  "lt",
						Value: "50",
					},
				},
			},
		},
		&intTestCase{
			description: "field type int with validation tag eq",
			field: &Field{
				Name: "example",
				Type: "int",
				Validations: []*Validation{
					{
						Name:  "lt",
						Value: "50",
					},
				},
			},
		},
		&intTestCase{
			description: "field type int with validation tag oneof",
			field: &Field{
				Name: "example",
				Type: "int",
				Validations: []*Validation{
					{
						Name:  "oneof",
						Value: "5 10",
					},
				},
			},
		},
		&stringTestCase{
			description: "field type string with validation tag max",
			field: &Field{
				Name: "example",
				Type: "string",
				Validations: []*Validation{
					{
						Name:  "max",
						Value: "50",
					},
				},
			},
		},
		&stringTestCase{
			description: "field type string with validation tag min",
			field: &Field{
				Name: "example",
				Type: "string",
				Validations: []*Validation{
					{
						Name: "required",
					},
					{
						Name:  "min",
						Value: "10",
					},
				},
			},
		},
		&stringTestCase{
			description: "field type string with validation tag gte",
			field: &Field{
				Name: "example",
				Type: "string",
				Validations: []*Validation{
					{
						Name:  "gte",
						Value: "10",
					},
				},
			},
		},
		&stringTestCase{
			description: "field type string with validation tag gt",
			field: &Field{
				Name: "example",
				Type: "string",
				Validations: []*Validation{
					{
						Name:  "gt",
						Value: "10",
					},
				},
			},
		},
		&stringTestCase{
			description: "field type string with validation tag lte",
			field: &Field{
				Name: "example",
				Type: "string",
				Validations: []*Validation{
					{
						Name:  "lte",
						Value: "60",
					},
				},
			},
		},
		&stringTestCase{
			description: "field type string with validation tag lt",
			field: &Field{
				Name: "example",
				Type: "string",
				Validations: []*Validation{
					{
						Name:  "lt",
						Value: "10",
					},
				},
			},
		},
		&stringTestCase{
			description: "field type string with validation tag email",
			field: &Field{
				Name: "example",
				Type: "string",
				Validations: []*Validation{
					{
						Name: "email",
					},
				},
			},
		},
		&stringTestCase{
			description: "field type string with validation tag len",
			field: &Field{
				Name: "example",
				Type: "string",
				Validations: []*Validation{
					{
						Name:  "len",
						Value: "89",
					},
				},
			},
		},
		&stringTestCase{
			description: "field type string with validation tag oneof",
			field: &Field{
				Name: "example",
				Type: "string",
				Validations: []*Validation{
					{
						Name:  "oneof",
						Value: "foo bar",
					},
				},
			},
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.Description(), func(t *testing.T) {
			exampleValue := tCase.Field().Example()
			if !tCase.IsValid(exampleValue) {
				t.Errorf("%s: example value %v is not valid", tCase.Description(), exampleValue)
			}
		})
	}
}
