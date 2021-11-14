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
