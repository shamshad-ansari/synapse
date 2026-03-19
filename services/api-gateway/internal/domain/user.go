package domain

import (
	"time"

	"github.com/google/uuid"
)

// School mirrors the schools table from migration 000001.
type School struct {
	ID        uuid.UUID
	Name      string
	Domain    *string
	CreatedAt time.Time
}

// User mirrors the users table from migration 000001.
type User struct {
	ID           uuid.UUID
	SchoolID     uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
