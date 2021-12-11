package entities

import (
	"testing"
)

type entityExampleTestCase struct {
	Description string
	Entity      *Entity
}

func (e *entityExampleTestCase) IsValid() bool {
	for _, f := range e.Entity.Fields {
		testCase := newFieldExampleTestCase(f)
		exampleValue := f.Example()

		if !testCase.IsValid(exampleValue) {
			return false
		}
	}
	return true
}

func TestEntityExample(t *testing.T) {
	testCases := []*entityExampleTestCase{
		{
			Description: "customer entity example",
			Entity: &Entity{
				Name: "customer",
				Fields: []*Field{
					{
						Name: "firstName",
						Type: "string",
						Validations: []*Validation{
							{
								Name:  "min",
								Value: "2",
							},
							{
								Name:  "max",
								Value: "100",
							},
							{
								Name: "required",
							},
						},
					},
					{
						Name: "lastName",
						Type: "string",
						Validations: []*Validation{
							{
								Name:  "min",
								Value: "2",
							},
							{
								Name:  "max",
								Value: "100",
							},
							{
								Name: "required",
							},
						},
					},
					{
						Name: "email",
						Type: "string",
						Validations: []*Validation{
							{
								Name: "email",
							},
							{
								Name: "required",
							},
						},
					},
					{
						Name: "score",
						Type: "float64",
						Validations: []*Validation{
							{
								Name:  "lte",
								Value: "1",
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			if !testCase.IsValid() {
				t.Errorf("%s: not valid", testCase.Description)
			}
		})
	}
}
