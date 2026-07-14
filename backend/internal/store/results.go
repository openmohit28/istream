package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

// TestResult persists a completed quiz: the raw answers, the computed
// scores, and a snapshot of the matches (so history stays stable even if
// the job catalog changes later).
type TestResult struct {
	ID        string          `json:"id"`
	UserID    string          `json:"-"`
	Answers   json.RawMessage `json:"answers"`
	Scores    json.RawMessage `json:"scores"`
	Matches   json.RawMessage `json:"matches"`
	CreatedAt time.Time       `json:"createdAt"`
}

type Results struct {
	DB *sql.DB
}

func (s *Results) Create(ctx context.Context, userID string, answers, scores, matches []byte) (TestResult, error) {
	var r TestResult
	err := s.DB.QueryRowContext(ctx,
		`INSERT INTO test_results (user_id, answers, scores, matches)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, answers, scores, matches, created_at`,
		userID, answers, scores, matches,
	).Scan(&r.ID, &r.UserID, &r.Answers, &r.Scores, &r.Matches, &r.CreatedAt)
	return r, err
}

// ByIDForUser fetches one result, scoped to its owner so users can never
// read each other's results.
func (s *Results) ByIDForUser(ctx context.Context, id, userID string) (TestResult, error) {
	var r TestResult
	err := s.DB.QueryRowContext(ctx,
		`SELECT id, user_id, answers, scores, matches, created_at
		 FROM test_results WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&r.ID, &r.UserID, &r.Answers, &r.Scores, &r.Matches, &r.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return TestResult{}, ErrNotFound
	}
	return r, err
}

func (s *Results) ListByUser(ctx context.Context, userID string) ([]TestResult, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, user_id, answers, scores, matches, created_at
		 FROM test_results WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []TestResult{}
	for rows.Next() {
		var r TestResult
		if err := rows.Scan(&r.ID, &r.UserID, &r.Answers, &r.Scores, &r.Matches, &r.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}
