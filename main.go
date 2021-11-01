package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/danilo-medeiros/fancybuild/engine/internal/builder"
	"github.com/danilo-medeiros/fancybuild/engine/internal/entities"
	"github.com/danilo-medeiros/fancybuild/engine/internal/reader"
	"github.com/danilo-medeiros/fancybuild/engine/internal/strategy"
)

const (
	OutputFolder = "/home/danilo/development/personal/generated-projects"
)

func main() {
	data, err := os.ReadFile("./test/definition.json")

	log.Println("reading definition file")

	if err != nil {
		log.Fatalf("error on reading definition: %s", err)
	}

	log.Println("parsing definition file")

	var definition entities.Definitions
	r := reader.NewReader()

	err = r.Read(data, &definition)

	if err != nil {
		log.Fatalf("error on parsing definition: %s", err)
	}

	definition.Id = fmt.Sprintf("%v", time.Now().Unix())

	log.Println("definition parsed successfully")

	validationErrs := r.Validate(&definition)

	if validationErrs != nil {
		if validationErr, ok := err.(*reader.ValidationError); ok {
			for _, fieldErr := range validationErr.Errors {
				log.Printf("validation error on: %s %s %s", fieldErr.Field, fieldErr.Tag, fieldErr.Value)
			}
			return
		}
	}

	log.Println("definition validated successfully")

	stgy := strategy.NewStrategy(&definition)

	if stgy == nil {
		log.Fatalf("error: strategy not found for %v", definition)
	}

	b := builder.NewBuilder(OutputFolder)

	err = b.Build(&definition, stgy)

	if err != nil {
		log.Fatalf("error on building project: %s", err)
	}
}
