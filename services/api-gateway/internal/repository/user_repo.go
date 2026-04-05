package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/google/uuid"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

type PostgresUserRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepo(pool *pgxpool.Pool) *PostgresUserRepo {
	return &PostgresUserRepo{pool: pool}
}

func (r *PostgresUserRepo) CreateSchool(ctx context.Context, name, domainStr string) (*domain.School, error) {
	var s domain.School
	err := r.pool.QueryRow(ctx,
		`INSERT INTO schools (name, domain)
		 VALUES ($1, $2)
		 ON CONFLICT (domain) DO UPDATE SET name = EXCLUDED.name
		 RETURNING id, name, domain, created_at`,
		name, domainStr,
	).Scan(&s.ID, &s.Name, &s.Domain, &s.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("CreateSchool: %w", err)
	}
	return &s, nil
}

func (r *PostgresUserRepo) FindSchoolByDomain(ctx context.Context, domainStr string) (*domain.School, error) {
	var s domain.School
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, domain, created_at FROM schools WHERE lower(domain) = lower($1)`,
		domainStr,
	).Scan(&s.ID, &s.Name, &s.Domain, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("FindSchoolByDomain: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("FindSchoolByDomain: %w", err)
	}
	return &s, nil
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, schoolID uuid.UUID, name, email, passwordHash string) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (school_id, name, email, password_hash)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, school_id, name, email, password_hash, created_at`,
		schoolID, name, email, passwordHash,
	).Scan(&u.ID, &u.SchoolID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("CreateUser: %w", domain.ErrConflict)
		}
		return nil, fmt.Errorf("CreateUser: %w", err)
	}
	return &u, nil
}

func (r *PostgresUserRepo) FindUserByEmail(ctx context.Context, email string, schoolID uuid.UUID) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx,
		`SELECT u.id, u.school_id, s.name, u.name, u.email, u.password_hash, u.created_at
		 FROM users u
		 INNER JOIN schools s ON s.id = u.school_id
		 WHERE lower(u.email) = lower($1) AND u.school_id = $2`,
		email, schoolID,
	).Scan(&u.ID, &u.SchoolID, &u.SchoolName, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("FindUserByEmail: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("FindUserByEmail: %w", err)
	}
	return &u, nil
}

func (r *PostgresUserRepo) FindUserByID(ctx context.Context, userID uuid.UUID, schoolID uuid.UUID) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx,
		`SELECT u.id, u.school_id, s.name, u.name, u.email, u.password_hash, u.created_at
		 FROM users u
		 INNER JOIN schools s ON s.id = u.school_id
		 WHERE u.id = $1 AND u.school_id = $2`,
		userID, schoolID,
	).Scan(&u.ID, &u.SchoolID, &u.SchoolName, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("FindUserByID: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("FindUserByID: %w", err)
	}
	return &u, nil
}
