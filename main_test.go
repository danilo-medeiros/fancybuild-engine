package main

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/danilo-medeiros/fancybuild/engine/internal/builder"
	"github.com/danilo-medeiros/fancybuild/engine/internal/entities"
	"github.com/danilo-medeiros/fancybuild/engine/internal/reader"
	"github.com/danilo-medeiros/fancybuild/engine/internal/strategy"
)

func TestExamples(t *testing.T) {
	files := []string{
		"blog.json",
		"todoapp.json",
	}

	for _, file := range files {
		data, err := os.ReadFile(fmt.Sprintf("./examples/%s", file))

		if err != nil {
			panic(fmt.Sprintf("error on reading definition: %s", err))
		}

		var definition entities.Definitions
		r := reader.NewReader()

		err = r.Read(data, &definition)

		if err != nil {
			log.Fatalf("error on parsing definition: %s", err)
		}

		definition.Id = fmt.Sprintf("%v", time.Now().Unix())

		validationErrs := r.Validate(&definition)

		if validationErrs != nil {
			if validationErr, ok := err.(*reader.ValidationError); ok {
				for _, fieldErr := range validationErr.Errors {
					panic(fmt.Sprintf("validation error on: %s %s %s", fieldErr.Field, fieldErr.Tag, fieldErr.Value))
				}
				return
			}
		}

		stgy := strategy.NewStrategy(&definition)

		if stgy == nil {
			panic(fmt.Sprintf("error: strategy not found for %v", definition))
		}

		b := builder.NewBuilder(OutputFolder)

		err = b.Build(&definition, stgy)

		if err != nil {
			panic(fmt.Sprintf("error on building project: %s", err))
		}
	}
}
