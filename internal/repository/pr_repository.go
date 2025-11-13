package repository

import (
	"PullRequestService/internal/db"
	"PullRequestService/internal/models"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
)

type PRRepository struct {
	db *db.DataBase
}

func NewPRRepository(db *db.DataBase) *PRRepository {
	return &PRRepository{db: db}
}

func (p *PRRepository) CreatePR(ctx context.Context, pull_request_id, pull_request_name, author_id string) (*models.PullRequest, []string, error) {
	var exists bool
	err := p.db.Db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`, pull_request_id).Scan(&exists)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, errors.New("PR_EXISTS")
	}

	var teamName string
	err = p.db.Db.QueryRow(ctx, `SELECT team_name FROM users WHERE user_id = $1`, author_id).Scan(&teamName)
	if err != nil {
		return nil, nil, errors.New("AUTHOR_NOT_FOUND")
	}

	rows, err := p.db.Db.Query(ctx, `
			SELECT user_id
			FROM users
			WHERE team_name = $1 AND is_active = TRUE AND user_id <> $2
		`, teamName, author_id)
	if err != nil {
		return nil, nil, fmt.Errorf("select candidates: %w", err)
	}
	defer rows.Close()

	candidates := make([]string, 0)
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, nil, fmt.Errorf("scan candidate: %w", err)
		}
		candidates = append(candidates, uid)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("rows err: %w", err)
	}

	reviewers := selectRandomReviewers(candidates, 2)
	needMore := len(reviewers) < 2

	tx, err := p.db.Db.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
			INSERT INTO pull_requests
				(pull_request_id, pull_request_name, author_id, status, need_more_reviewers, created_at)
			VALUES ($1, $2, $3, 'OPEN', $4, $5)
		`, pull_request_id, pull_request_name, author_id, needMore, time.Now().UTC())
	if err != nil {
		return nil, nil, fmt.Errorf("insert pr: %w", err)
	}

	for _, rid := range reviewers {
		_, err := tx.Exec(ctx, `
				INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)
			`, pull_request_id, rid)
		if err != nil {
			return nil, nil, fmt.Errorf("insert reviewer: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("commit tx: %w", err)
	}

	pr := &models.PullRequest{
		PullRequestID:   pull_request_id,
		PullRequestName: pull_request_name,
		AuthorID:        author_id,
		Status:          "OPEN",
		NeedMore:        needMore,
		CreatedAt:       time.Now().UTC(),
		MergedAt:        nil,
	}
	return pr, reviewers, nil
}

func selectRandomReviewers(users []string, n int) []string {
	if len(users) == 0 {
		return nil
	}
	if len(users) <= n {
		out := make([]string, len(users))
		copy(out, users)
		return out
	}
	rand := randSource()
	perm := rand.Perm(len(users))
	out := make([]string, 0, n)
	for i := range n {
		out = append(out, users[perm[i]])
	}
	return out
}

func randSource() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (p *PRRepository) UpdateStatus(ctx context.Context, prID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	err := p.db.Db.QueryRow(ctx, `
		UPDATE pull_requests
		SET status = 'MERGED', merged_at = COALESCE(merged_at, NOW())
		WHERE pull_request_id = $1
		RETURNING pull_request_id, pull_request_name, author_id, status, need_more_reviewers, created_at, merged_at`, prID).
		Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.NeedMore, &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("NOT_FOUND")
		}
		return nil, err
	}

	rows, err := p.db.Db.Query(ctx, `
			SELECT reviewer_id
			FROM pr_reviewers
			WHERE pull_request_id = $1
		`, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviewers := []string{}
	for rows.Next() {
		var r string
		if err := rows.Scan(&r); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, r)
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}
