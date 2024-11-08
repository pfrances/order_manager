package id

import (
	"github.com/google/uuid"
)

type ID struct {
	uuid.UUID
}

func New() ID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	return ID(ID{id})
}

func NilID() ID {
	return ID(ID{uuid.Nil})
}

func (id ID) IsNil() bool {
	return id == NilID()
}

func (id ID) String() string {
	return id.UUID.String()
}
