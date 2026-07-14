package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"istream/backend/internal/models"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrDuplicateEmail = errors.New("email already registered")
)

type Users struct {
	DB *sql.DB
}

func (s *Users) Create(ctx context.Context, email, name, passwordHash string) (models.User, error) {
	var u models.User
	err := s.DB.QueryRowContext(ctx,
		`INSERT INTO users (email, name, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id, email, name, password_hash, created_at`,
		email, name, passwordHash,
	).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.User{}, ErrDuplicateEmail
		}
		return models.User{}, err
	}
	return u, nil
}

func (s *Users) ByEmail(ctx context.Context, email string) (models.User, error) {
	return s.scanOne(s.DB.QueryRowContext(ctx,
		`SELECT id, email, name, password_hash, created_at FROM users WHERE email = $1`, email))
}

func (s *Users) ByID(ctx context.Context, id string) (models.User, error) {
	return s.scanOne(s.DB.QueryRowContext(ctx,
		`SELECT id, email, name, password_hash, created_at FROM users WHERE id = $1`, id))
}

func (s *Users) scanOne(row *sql.Row) (models.User, error) {
	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}