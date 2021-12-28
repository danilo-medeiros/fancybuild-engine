package mysql

import (
	"testing"

	"github.com/danilo-medeiros/fancybuild/engine/pkg/entities"
)

func TestBuildIndexes(t *testing.T) {
	t.Run("single index should be created", func(t *testing.T) {
		entity := &entities.Entity{
			Name: "example",
			Fields: []*entities.Field{
				{
					Name: "name",
					Type: "string",
				},
			},
			Indexes: []*entities.Index{
				{
					Fields: []*entities.IndexField{
						{
							Name: "name",
						},
					},
				},
			},
		}

		indexMap := buildIndexes(entity)
		tag := indexMap["name"]

		expected := "index:idx_example_name"
		if tag != expected {
			t.Errorf("name expected %s, got: %s", expected, tag)
		}
	})

	t.Run("single index unique should be created", func(t *testing.T) {
		entity := &entities.Entity{
			Name: "user",
			Fields: []*entities.Field{
				{
					Name: "id",
					Type: "string",
				},
			},
			Indexes: []*entities.Index{
				{
					Fields: []*entities.IndexField{
						{
							Name: "id",
						},
					},
					Unique: true,
				},
			},
		}

		indexMap := buildIndexes(entity)
		tag := indexMap["id"]

		expected := "index:idx_user_id,unique"
		if tag != expected {
			t.Errorf("id expected %s, got: %s", expected, tag)
		}
	})

	t.Run("multiple indexes should be created", func(t *testing.T) {
		entity := &entities.Entity{
			Name: "user",
			Fields: []*entities.Field{
				{
					Name: "id",
					Type: "string",
				},
				{
					Name: "name",
					Type: "string",
				},
				{
					Name: "email",
					Type: "string",
				},
			},
			Indexes: []*entities.Index{
				{
					Fields: []*entities.IndexField{
						{
							Name: "id",
						},
					},
					Unique: true,
				},
				{
					Fields: []*entities.IndexField{
						{
							Name: "email",
						},
					},
					Unique: true,
				},
			},
		}

		indexMap := buildIndexes(entity)

		{
			tag := indexMap["id"]

			expected := "index:idx_user_id,unique"
			if tag != expected {
				t.Errorf("id expected %s, got: %s", expected, tag)
			}
		}

		{
			tag := indexMap["email"]

			expected := "index:idx_user_email,unique"
			if tag != expected {
				t.Errorf("email expected %s, got: %s", expected, tag)
			}
		}

		{
			tag, ok := indexMap["name"]

			if ok {
				t.Errorf("name expected empty value, got: %s", tag)
			}
		}
	})

	t.Run("multi-key indexes should be created", func(t *testing.T) {
		entity := &entities.Entity{
			Name: "user",
			Fields: []*entities.Field{
				{
					Name: "id",
					Type: "string",
				},
				{
					Name: "name",
					Type: "string",
				},
				{
					Name: "email",
					Type: "string",
				},
			},
			Indexes: []*entities.Index{
				{
					Fields: []*entities.IndexField{
						{
							Name: "id",
						},
					},
					Unique: true,
				},
				{
					Fields: []*entities.IndexField{
						{
							Name: "email",
						},
						{
							Name: "name",
						},
					},
					Unique: true,
				},
			},
		}

		indexMap := buildIndexes(entity)

		{
			tag := indexMap["id"]

			expected := "index:idx_user_id,unique"
			if tag != expected {
				t.Errorf("id expected %s, got: %s", expected, tag)
			}
		}

		{
			tag := indexMap["email"]

			expected := "index:idx_user_email_name,unique"
			if tag != expected {
				t.Errorf("email expected %s, got: %s", expected, tag)
			}
		}

		{
			tag := indexMap["name"]

			expected := "index:idx_user_email_name,unique"
			if tag != expected {
				t.Errorf("name expected %s, got: %s", expected, tag)
			}
		}
	})

	t.Run("multiple indexes in a field should be created", func(t *testing.T) {
		entity := &entities.Entity{
			Name: "user",
			Fields: []*entities.Field{
				{
					Name: "id",
					Type: "string",
				},
				{
					Name: "name",
					Type: "string",
				},
				{
					Name: "email",
					Type: "string",
				},
			},
			Indexes: []*entities.Index{
				{
					Fields: []*entities.IndexField{
						{
							Name: "id",
						},
					},
					Unique: true,
				},
				{
					Fields: []*entities.IndexField{
						{
							Name: "email",
						},
					},
				},
				{
					Fields: []*entities.IndexField{
						{
							Name: "email",
						},
						{
							Name: "name",
						},
					},
					Unique: true,
				},
			},
		}

		indexMap := buildIndexes(entity)

		{
			tag := indexMap["id"]

			expected := "index:idx_user_id,unique"
			if tag != expected {
				t.Errorf("id expected %s, got: %s", expected, tag)
			}
		}

		{
			tag := indexMap["email"]

			expected := "index:idx_user_email;index:idx_user_email_name,unique"
			if tag != expected {
				t.Errorf("email expected %s, got: %s", expected, tag)
			}
		}

		{
			tag := indexMap["name"]

			expected := "index:idx_user_email_name,unique"
			if tag != expected {
				t.Errorf("name expected %s, got: %s", expected, tag)
			}
		}
	})
}
