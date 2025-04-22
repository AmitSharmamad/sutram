package util

import (
	"github.com/google/uuid"
)

// Parse converts a string to a UUID
func Parse(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}

// NewUUID generates a new UUID
func NewUUID() uuid.UUID {
	return uuid.New()
}

// ParsePtr converts a string to a UUID pointer
func ParsePtr(s string) *uuid.UUID {
	id := Parse(s)
	return &id
}

// ToString converts a UUID to a string
func ToString(id uuid.UUID) string {
	return id.String()
}
