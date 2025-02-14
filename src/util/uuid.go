package util

import "github.com/google/uuid"

func NewUUID() uuid.UUID {
	return uuid.New()
}

func Parse(id string) uuid.UUID {
	i, _ := uuid.Parse(id)
	return i
}
