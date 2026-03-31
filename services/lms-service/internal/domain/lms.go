package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// LMSConnection represents a user's connection to a Canvas LMS instance.
type LMSConnection struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	SchoolID       uuid.UUID
	LMSType        string
	InstitutionURL string
	AccessToken    string
	RefreshToken   string
	TokenExpiresAt time.Time
	LastSyncedAt   *time.Time
	SyncStatus     string
	CreatedAt      time.Time
}

// LMSConnectionResponse is the public-facing DTO — never exposes raw tokens.
type LMSConnectionResponse struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	SchoolID       uuid.UUID  `json:"school_id"`
	LMSType        string     `json:"lms_type"`
	InstitutionURL string     `json:"institution_url"`
	TokenExpiresAt time.Time  `json:"token_expires_at"`
	LastSyncedAt   *time.Time `json:"last_synced_at"`
	SyncStatus     string     `json:"sync_status"`
	CreatedAt      time.Time  `json:"created_at"`
}

// ToResponse strips sensitive token fields from LMSConnection.
func (c *LMSConnection) ToResponse() *LMSConnectionResponse {
	return &LMSConnectionResponse{
		ID:             c.ID,
		UserID:         c.UserID,
		SchoolID:       c.SchoolID,
		LMSType:        c.LMSType,
		InstitutionURL: c.InstitutionURL,
		TokenExpiresAt: c.TokenExpiresAt,
		LastSyncedAt:   c.LastSyncedAt,
		SyncStatus:     c.SyncStatus,
		CreatedAt:      c.CreatedAt,
	}
}

// LMSRepository abstracts persistence for LMS connection data.
type LMSRepository interface {
	UpsertConnection(ctx context.Context, conn *LMSConnection) error
	FindConnectionByUser(ctx context.Context, userID, schoolID uuid.UUID) (*LMSConnection, error)
	DeleteConnection(ctx context.Context, userID, schoolID uuid.UUID) error
}
