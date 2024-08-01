package personapp

import (
	"encoding/json"

	"github.com/bmviniciuss/gobank/person/app/sdk/errs"
	"github.com/bmviniciuss/gobank/person/core/person"
	"github.com/bmviniciuss/gobank/person/foundation/utc"
)

type NewPerson struct {
	Name     string `json:"name" validate:"required"`
	Document string `json:"document" validate:"required"`
}

func (np *NewPerson) Decode(data []byte) error {
	return json.Unmarshal(data, &np)
}

func (np *NewPerson) Validate() error {
	return errs.Validate(np)
}

type Person struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Active    bool     `json:"active"`
	CreatedAt utc.Time `json:"created_at"`
	UpdatedAt utc.Time `json:"updated_at"`
}

func (p *Person) FromPerson(pp *person.Person) {
	*p = Person{
		ID:        pp.ID.String(),
		Name:      pp.Name,
		Active:    pp.Active,
		CreatedAt: utc.NewFromTime(pp.CreatedAt),
		UpdatedAt: utc.NewFromTime(pp.UpdatedAt),
	}
}
