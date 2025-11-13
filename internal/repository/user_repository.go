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
	err := u.db.Db.QueryRow(ctx, "UPDATE users SET is_active = $2 WHERE user_id = $1 RETURNING user_id, username, team_name, is_active").Scan(&user.UserID, &user.Username, &user.TeamName, &user.Is_Active)
	if err != nil {
		return nil, errors.New("user_not_found")
	}
	return &user, nil
}
