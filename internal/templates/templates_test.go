package templates

import "testing"

func TestPluralize(t *testing.T) {
	nounMap := map[string]string{
		"project":       "projects",
		"authorization": "authorizations",
		"user":          "users",
		"knife":         "knives",
		"life":          "lives",
		"wife":          "wives",
		"calf":          "calves",
		"leaf":          "leaves",
		"foot":          "feet",
		"tooth":         "teeth",
		"goose":         "geese",
		"man":           "men",
		"woman":         "women",
		"cat":           "cats",
		"house":         "houses",
		"truss":         "trusses",
		"bus":           "buses",
		"marsh":         "marshes",
		"lunch":         "lunches",
		"tax":           "taxes",
		"blitz":         "blitzes",
		"city":          "cities",
		"puppy":         "puppies",
		"ray":           "rays",
		"boy":           "boys",
		"photo":         "photos",
		"piano":         "pianos",
		"halo":          "halos",
		"cactus":        "cacti",
		"focus":         "foci",
		"potato":        "potatoes",
		"tomato":        "tomatoes",
		"analysis":      "analyses",
		"ellipsis":      "ellipses",
		"phenomenon":    "phenomena",
		"criterion":     "criteria",
		"child":         "children",
	}

	for key, value := range nounMap {
		res := Pluralize(key)
		if res != value {
			t.Errorf("Pluralize(%s) wanted %s, got %s", key, value, res)
		}
	}
}

const testSimpleFormatInput = `package test
import (
	"fmt"
	"github.com/stretchr/testify/assert"
)
func test(a string, b string) {
	fmt.Printf(
		"a: %s, b: %s",
		a,
		b,
	)
}
func main() {
	fmt.Println("Hello World!")
	for i := 0; i < 10; i++ {
		fmt.Printf("Hello #%s", count)
		if i == 5 {
			fmt.Println("five!")
			test(
				"Hello",
				"Five",
			)
		}
		fmt.Printf("Hello #%s", count)
	}
	test(
		"Hello",
		"World",
	)
}`

const testSimpleFormatExpected = `package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
)

func test(a string, b string) {
	fmt.Printf(
		"a: %s, b: %s",
		a,
		b,
	)
}

func main() {
	fmt.Println("Hello World!")

	for i := 0; i < 10; i++ {
		fmt.Printf("Hello #%s", count)

		if i == 5 {
			fmt.Println("five!")

			test(
				"Hello",
				"Five",
			)
		}

		fmt.Printf("Hello #%s", count)
	}

	test(
		"Hello",
		"World",
	)
}`

func TestSimpleFormat(t *testing.T) {
	input := testSimpleFormatInput
	expected := testSimpleFormatExpected
	actual := SimpleFormat(input)

	if actual != expected {
		t.Errorf("SimpleFormat wanted:\n%s\nBut got:\n%s\n", expected, actual)
	}
}
