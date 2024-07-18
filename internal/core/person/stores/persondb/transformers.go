package persondb

import (
	"github.com/bmviniciuss/gobank/internal/core/person"
	"github.com/bmviniciuss/gobank/internal/core/person/stores/persondb/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

func toPerson(r generated.FindPersonByDocumentRow) *person.Person {
	return &person.Person{
		ID:        r.Uuid,
		Name:      r.Name,
		Document:  r.Document,
		CreatedAt: r.CreatedAt.Time,
		UpdatedAt: r.UpdatedAt.Time,
	}
}

func toInsertPersonParams(p *person.Person) generated.InsertPersonParams {
	return generated.InsertPersonParams{
		Uuid:     p.ID,
		Name:     p.Name,
		Document: p.Document,
		Active:   p.Active,
		CreatedAt: pgtype.Timestamptz{
			Time:  p.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  p.UpdatedAt,
			Valid: true,
		},
	}
}
