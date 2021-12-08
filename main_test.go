package engine

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/danilo-medeiros/fancybuild/engine/pkg/builder"
	"github.com/danilo-medeiros/fancybuild/engine/pkg/entities"
	"github.com/danilo-medeiros/fancybuild/engine/pkg/reader"
	"github.com/danilo-medeiros/fancybuild/engine/pkg/strategy"
)

const OutputFolder = "/home/danilo/development/personal/generated-projects"

func subTest(file string) func(t *testing.T) {
	return func(t *testing.T) {
		data, err := os.ReadFile(fmt.Sprintf("./_examples/%s", file))

		if err != nil {
			panic(fmt.Sprintf("error on reading definition: %s", err))
		}

		var definition entities.Definitions
		r := reader.NewReader()

		err = r.Read(data, &definition)

		if err != nil {
			panic(fmt.Sprintf("error on parsing definition: %s", err))
		}

		definition.Id = fmt.Sprintf("%v", time.Now().Unix())
		validationErrs := r.Validate(&definition)

		if validationErrs != nil {
			for _, fieldErr := range validationErrs.Errors {
				t.Fatalf("validation error on field: %s, error: %s, value: %s", fieldErr.Field, fieldErr.Tag, fieldErr.Value)
			}
		}

		stgy := strategy.NewStrategy(&definition)

		if stgy == nil {
			t.Fatalf("error: strategy not found for %v", definition)
		}

		b := builder.NewBuilder(OutputFolder)
		err = b.Build(&definition, stgy)

		if err != nil {
			t.Fatalf("error on building project: %s", err)
		}
	}
}

func TestExamples(t *testing.T) {
	files := []string{
		"blog.json",
		"todoapp.json",
		"ecommerce.json",
	}

	for _, file := range files {
		t.Run(file, subTest(file))
	}
}
