package models

type User struct {
	UserID    string `json:"user_id" db:"user_id"`
	Username  string `json:"username" db:"username"`
	TeamName  string `json:"team_name" db:"team_name"`
	Is_Active bool   `json:"is_active" db:"is_active"`
}
