package repository

import (
	"PullRequestService/internal/db"
	"PullRequestService/internal/models"
	"PullRequestService/pkg/logger"
	"context"
	"errors"
)

type TeamRepository struct {
	pool *db.DataBase
	log  logger.Logger
}

func NewTeamRepository(db *db.DataBase) *TeamRepository {
	log := logger.New()
	return &TeamRepository{pool: db, log: log}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team models.Team) error {
	tx, err := r.pool.Db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var existTeam bool
	err = tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM teams WHERE team_name=$1)", team.TeamName).Scan(&existTeam)
	if err != nil {
		return err
	}
	if existTeam {
		return errors.New("TEAM_EXISTS")
	}
	r.log.Info("Team does not exist")
	_, err = tx.Exec(ctx, "INSERT INTO teams (team_name) VALUES ($1)", team.TeamName)
	if err != nil {
		return err
	}
	for _, m := range team.Members {
		_, err = tx.Exec(ctx, `
			INSERT INTO users (user_id, username, team_name, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id)
			DO UPDATE SET username = EXCLUDED.username, is_active = EXCLUDED.is_active, team_name = EXCLUDED.team_name`,
			m.UserID, m.Username, team.TeamName, m.IsActive)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	r.log.Info("Successfully created team")
	return nil
}

func (r *TeamRepository) GetTeam(ctx context.Context, name string) (*models.Team, error) {
	rows, err := r.pool.Db.Query(ctx, `
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []models.TeamMember{}
	for rows.Next() {
		var m models.TeamMember
		if err := rows.Scan(&m.UserID, &m.Username, &m.IsActive); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	if len(members) == 0 {
		r.log.Error("Length of members slice equal 0")
		return nil, errors.New("NOT_FOUND")
	}

	return &models.Team{
		TeamName: name,
		Members:  members,
	}, nil
}
