package models

type ReviewersStats struct {
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	AssignedCount int    `json:"assigned_count"`
}
