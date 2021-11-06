package entities

type Definitions struct {
	Id      string `json:"id"`
	Version string `json:"version"`
	App     *App   `json:"app"`
}

func (d Definitions) HasAuthentication() bool {
	return len(d.App.Authentication.Entity) > 0
}

func (d Definitions) AuthEntity() *Entity {
	for _, e := range d.App.Entities {
		if e.Name == d.App.Authentication.Entity {
			return e
		}
	}
	return nil
}

func (d Definitions) FindEntity(entity string) *Entity {
	for _, e := range d.App.Entities {
		if e.Name == entity {
			return e
		}
	}
	return nil
}
