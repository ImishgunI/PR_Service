package repository

import (
	"PullRequestService/internal/db"
	"PullRequestService/internal/models"
	"context"
)

type StatsRepository interface {
	GetReviewerStats(ctx context.Context) ([]models.ReviewersStats, error)
	GetPRStats(ctx context.Context) ([]models.PRStats, error)
}

type Statistic struct {
	db *db.DataBase
}

func NewStatisticRepository(db *db.DataBase) *Statistic {
	return &Statistic{db: db}
}

func (s *Statistic) GetReviewerStats(ctx context.Context) ([]models.ReviewersStats, error) {
	rows, err := s.db.Db.Query(ctx, `
		SELECT u.user_id, u.username, COUNT(pr.reviewer_id) AS assigned_count
        FROM users u
        LEFT JOIN pr_reviewers pr ON u.user_id = pr.reviewer_id
        GROUP BY u.user_id, u.username
        ORDER BY assigned_count DESC;
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.ReviewersStats{}
	for rows.Next() {
		var s models.ReviewersStats
		if err := rows.Scan(&s.UserID, &s.Username, &s.AssignedCount); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (s *Statistic) GetPRStats(ctx context.Context) ([]models.PRStats, error) {
	rows, err := s.db.Db.Query(ctx, `
		SELECT pr.pull_request_id, pr.pull_request_name, COUNT(pp.reviewer_id) as reviewers_count
		FROM pull_requests pr
		LEFT JOIN pr_reviewers pp ON pr.pull_request_id = pp.pull_request_id
		GROUP BY pr.pull_request_id, pr.pull_request_name
		ORDER BY reviewers_count DESC
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.PRStats{}
	for rows.Next() {
		var pr models.PRStats
		if err := rows.Scan(&pr.PrID, &pr.PrName, &pr.ReviewersCount); err != nil {
			return nil, err
		}
		result = append(result, pr)
	}
	return result, nil
}
