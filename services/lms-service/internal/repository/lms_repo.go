package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shamshad-ansari/synapse/services/lms-service/internal/domain"
)

type PostgresLMSRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresLMSRepo(pool *pgxpool.Pool) *PostgresLMSRepo {
	return &PostgresLMSRepo{pool: pool}
}

func (r *PostgresLMSRepo) UpsertConnection(ctx context.Context, conn *domain.LMSConnection) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO lms_connections (user_id, school_id, lms_type, institution_url, access_token, refresh_token, token_expires_at, sync_status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (user_id, lms_type) DO UPDATE SET
		   institution_url  = EXCLUDED.institution_url,
		   access_token     = EXCLUDED.access_token,
		   refresh_token    = EXCLUDED.refresh_token,
		   token_expires_at = EXCLUDED.token_expires_at,
		   sync_status      = EXCLUDED.sync_status`,
		conn.UserID, conn.SchoolID, conn.LMSType, conn.InstitutionURL,
		conn.AccessToken, conn.RefreshToken, conn.TokenExpiresAt, conn.SyncStatus,
	)
	if err != nil {
		return fmt.Errorf("UpsertConnection: %w", err)
	}
	return nil
}

func (r *PostgresLMSRepo) FindConnectionByUser(ctx context.Context, userID, schoolID uuid.UUID) (*domain.LMSConnection, error) {
	var c domain.LMSConnection
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, school_id, lms_type, institution_url,
		        access_token, refresh_token, token_expires_at,
		        last_synced_at, sync_status, created_at
		 FROM lms_connections
		 WHERE user_id = $1 AND school_id = $2`,
		userID, schoolID,
	).Scan(
		&c.ID, &c.UserID, &c.SchoolID, &c.LMSType, &c.InstitutionURL,
		&c.AccessToken, &c.RefreshToken, &c.TokenExpiresAt,
		&c.LastSyncedAt, &c.SyncStatus, &c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("FindConnectionByUser: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("FindConnectionByUser: %w", err)
	}
	return &c, nil
}

func (r *PostgresLMSRepo) DeleteConnection(ctx context.Context, userID, schoolID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM lms_connections WHERE user_id = $1 AND school_id = $2`,
		userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("DeleteConnection: %w", err)
	}
	return nil
}
