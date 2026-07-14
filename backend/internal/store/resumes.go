package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type Resume struct {
	ID        string          `json:"id"`
	UserID    string          `json:"-"`
	Title     string          `json:"title"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type Resumes struct {
	DB *sql.DB
}

func (s *Resumes) Create(ctx context.Context, userID, title string, data []byte) (Resume, error) {
	var r Resume
	err := s.DB.QueryRowContext(ctx,
		`INSERT INTO resumes (user_id, title, data)
		 VALUES ($1, $2, $3)
		 RETURNING id, user_id, title, data, created_at, updated_at`,
		userID, title, data,
	).Scan(&r.ID, &r.UserID, &r.Title, &r.Data, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func (s *Resumes) Update(ctx context.Context, id, userID, title string, data []byte) (Resume, error) {
	var r Resume
	err := s.DB.QueryRowContext(ctx,
		`UPDATE resumes SET title = $3, data = $4, updated_at = now()
		 WHERE id = $1 AND user_id = $2
		 RETURNING id, user_id, title, data, created_at, updated_at`,
		id, userID, title, data,
	).Scan(&r.ID, &r.UserID, &r.Title, &r.Data, &r.CreatedAt, &r.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Resume{}, ErrNotFound
	}
	return r, err
}

func (s *Resumes) ByIDForUser(ctx context.Context, id, userID string) (Resume, error) {
	var r Resume
	err := s.DB.QueryRowContext(ctx,
		`SELECT id, user_id, title, data, created_at, updated_at
		 FROM resumes WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&r.ID, &r.UserID, &r.Title, &r.Data, &r.CreatedAt, &r.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Resume{}, ErrNotFound
	}
	return r, err
}

func (s *Resumes) ListByUser(ctx context.Context, userID string) ([]Resume, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, user_id, title, data, created_at, updated_at
		 FROM resumes WHERE user_id = $1 ORDER BY updated_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resumes := []Resume{}
	for rows.Next() {
		var r Resume
		if err := rows.Scan(&r.ID, &r.UserID, &r.Title, &r.Data, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		resumes = append(resumes, r)
	}
	return resumes, rows.Err()
}

func (s *Resumes) Delete(ctx context.Context, id, userID string) error {
	res, err := s.DB.ExecContext(ctx,
		`DELETE FROM resumes WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
