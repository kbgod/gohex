package types

import "github.com/google/uuid"

type ID = uuid.UUID

func NewID() ID {
	return uuid.Must(uuid.NewV7())
}
