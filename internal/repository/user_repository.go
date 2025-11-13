package repository

import (
	"PullRequestService/internal/db"
	"PullRequestService/internal/models"
	"context"
	"errors"
)

type UserRepository struct {
	db *db.DataBase
}

func NewUserHandler(db *db.DataBase) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	var user models.User
	err := u.db.Db.QueryRow(ctx, `UPDATE users SET is_active = $2 WHERE user_id = $1
		RETURNING user_id, username, team_name, is_active`, userID, isActive).Scan(&user.UserID, &user.Username, &user.TeamName, &user.Is_Active)
	if err != nil {
		return nil, errors.New("user_not_found")
	}
	return &user, nil
}

func (u *UserRepository) GetReview(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	rows, err := u.db.Db.Query(ctx, `
			SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
			FROM pull_requests pr
			JOIN pr_reviewers rr ON pr.pull_request_id = rr.pull_request_id
			WHERE rr.reviewer_id = $1
		`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []models.PullRequestShort
	for rows.Next() {
		var pr models.PullRequestShort
		if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, nil
}
