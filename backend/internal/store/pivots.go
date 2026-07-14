package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

// PivotThread is one exploration through the pivot decision tree. Forked
// threads remember their origin so the UI can show the relationship.
type PivotThread struct {
	ID         string          `json:"id"`
	UserID     string          `json:"-"`
	Steps      json.RawMessage `json:"steps"`
	ForkedFrom *string         `json:"forkedFrom,omitempty"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

type Pivots struct {
	DB *sql.DB
}

const pivotColumns = `id, user_id, steps, forked_from, created_at, updated_at`

func scanPivot(row interface{ Scan(...any) error }) (PivotThread, error) {
	var t PivotThread
	err := row.Scan(&t.ID, &t.UserID, &t.Steps, &t.ForkedFrom, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (s *Pivots) Create(ctx context.Context, userID string, steps []byte, forkedFrom *string) (PivotThread, error) {
	return scanPivot(s.DB.QueryRowContext(ctx,
		`INSERT INTO pivot_threads (user_id, steps, forked_from)
		 VALUES ($1, $2, $3)
		 RETURNING `+pivotColumns,
		userID, steps, forkedFrom,
	))
}

func (s *Pivots) UpdateSteps(ctx context.Context, id, userID string, steps []byte) (PivotThread, error) {
	t, err := scanPivot(s.DB.QueryRowContext(ctx,
		`UPDATE pivot_threads SET steps = $3, updated_at = now()
		 WHERE id = $1 AND user_id = $2
		 RETURNING `+pivotColumns,
		id, userID, steps,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return PivotThread{}, ErrNotFound
	}
	return t, err
}

func (s *Pivots) ByIDForUser(ctx context.Context, id, userID string) (PivotThread, error) {
	t, err := scanPivot(s.DB.QueryRowContext(ctx,
		`SELECT `+pivotColumns+` FROM pivot_threads WHERE id = $1 AND user_id = $2`,
		id, userID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return PivotThread{}, ErrNotFound
	}
	return t, err
}

func (s *Pivots) ListByUser(ctx context.Context, userID string) ([]PivotThread, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT `+pivotColumns+` FROM pivot_threads
		 WHERE user_id = $1 ORDER BY updated_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := []PivotThread{}
	for rows.Next() {
		t, err := scanPivot(rows)
		if err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}
	return threads, rows.Err()
}

func (s *Pivots) Delete(ctx context.Context, id, userID string) error {
	res, err := s.DB.ExecContext(ctx,
		`DELETE FROM pivot_threads WHERE id = $1 AND user_id = $2`, id, userID)
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
