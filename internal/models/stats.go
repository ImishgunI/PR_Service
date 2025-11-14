package models

type ReviewersStats struct {
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	AssignedCount int    `json:"assigned_count"`
}

type PRStats struct {
	PrID           string `json:"pull_request_id"`
	PrName         string `json:"pull_request_name"`
	ReviewersCount int    `json:"reviewers_count"`
}
