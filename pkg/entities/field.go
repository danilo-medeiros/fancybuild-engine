package entities

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Field struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Validations []*Validation `json:"validations"`
	Secret      bool          `json:"secret"`
	Hashed      bool          `json:"hashed"`
}

func randomString(chars string, size int) string {
	result := ""

	for {
		if len(result) == size {
			return result
		}

		result += string(chars[rand.Intn(len(chars))])
	}
}

// Generates an example value for this field.
// The value is generated within the validation constraints
func (f Field) Example() string {
	isNumber := false
	var result string

	switch f.Type {
	case "int", "uint", "int32", "int64", "float32", "float64":
		isNumber = true
	default:
		isNumber = false
	}

	max := 20
	min := 0
	required := false

	for _, validation := range f.Validations {
		if validation.Name == "max" || validation.Name == "lte" {
			m, err := strconv.Atoi(validation.Value)
			if err != nil {
				panic(fmt.Sprintf("on parsing \"max\" validation: %s", err))
			}
			max = m
		}

		if validation.Name == "lt" {
			m, err := strconv.Atoi(validation.Value)
			if err != nil {
				panic(fmt.Sprintf("on parsing \"lt\" validation: %s", err))
			}
			max = m - 1
		}

		if validation.Name == "min" || validation.Name == "gte" {
			m, err := strconv.Atoi(validation.Value)
			if err != nil {
				panic(fmt.Sprintf("on parsing \"min\" validation: %s", err))
			}
			min = m
		}

		if validation.Name == "gt" {
			m, err := strconv.Atoi(validation.Value)
			if err != nil {
				panic(fmt.Sprintf("on parsing \"gt\" validation: %s", err))
			}
			min = m + 1
		}

		if validation.Name == "eq" {
			return validation.Value
		}

		if validation.Name == "oneof" {
			return strings.Split(validation.Value, " ")[0]
		}

		if validation.Name == "len" {
			m, err := strconv.Atoi(validation.Value)
			if err != nil {
				panic(fmt.Sprintf("on parsing \"len\" validation: %s", err))
			}
			min = m
			max = m
		}

		if validation.Name == "email" {
			return fmt.Sprintf("example.%s@example.com", randomString("abcdefghijklmnopqrstuvwxyz1234567890", 5))
		}

		if validation.Name == "required" {
			required = true
		}
	}

	if max <= min {
		max += min
	}

	if required && min == 0 {
		min = 1
	}

	if isNumber {
		value := rand.Intn(max-min) + min
		return fmt.Sprintf("%v", value)
	}

	size := rand.Intn(max-min) + min
	result = randomString("abcdefghijklmnopqrstuvwxyz", size)

	return result
}
