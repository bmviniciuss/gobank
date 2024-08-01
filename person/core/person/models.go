package person

import (
	"time"

	"github.com/google/uuid"
)

type CreatePerson struct {
	Name     string
	Document string
}

type Person struct {
	ID        uuid.UUID
	Name      string
	Document  string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func newPerson(name, document string) (*Person, error) {
	pID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &Person{
		ID:        pID,
		Name:      name,
		Document:  document,
		Active:    true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}
