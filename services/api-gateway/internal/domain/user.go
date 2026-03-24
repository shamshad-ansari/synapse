package domain

import (
	"context"
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

// UserResponse is the public-facing DTO that never exposes password_hash.
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	SchoolID  uuid.UUID `json:"school_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse strips sensitive fields from User.
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		SchoolID:  u.SchoolID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

// UserRepository abstracts persistence for users and schools.
type UserRepository interface {
	CreateSchool(ctx context.Context, name, domain string) (*School, error)
	FindSchoolByDomain(ctx context.Context, domain string) (*School, error)
	CreateUser(ctx context.Context, schoolID uuid.UUID, name, email, passwordHash string) (*User, error)
	FindUserByEmail(ctx context.Context, email string, schoolID uuid.UUID) (*User, error)
	FindUserByID(ctx context.Context, userID uuid.UUID, schoolID uuid.UUID) (*User, error)
}
