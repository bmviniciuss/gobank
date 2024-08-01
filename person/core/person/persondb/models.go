package persondb

import (
	"time"

	"github.com/bmviniciuss/gobank/person/core/person"
	"github.com/google/uuid"
)

type findPersonByDocumentRow struct {
	UUID      uuid.UUID `db:"uuid"`
	Name      string    `db:"name"`
	Document  string    `db:"document"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r findPersonByDocumentRow) toPerson() person.Person {
	return person.Person{
		ID:        r.UUID,
		Name:      r.Name,
		Document:  r.Document,
		Active:    r.Active,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
