package models

import "time"

type PullRequest struct {
	PullRequestID     string     `db:"pull_request_id" json:"pull_request_id"`
	PullRequestName   string     `db:"pull_request_name" json:"pull_request_name"`
	AuthorID          string     `db:"author_id" json:"author_id"`
	Status            string     `db:"status" json:"status"`
	NeedMore          bool       `db:"need_more_reviewers" json:"needMoreReviewers"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `db:"created_at" json:"createdAt"`
	MergedAt          *time.Time `db:"merged_at" json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PullRequestID   string `db:"pull_request_id" json:"pull_request_id"`
	PullRequestName string `db:"pull_request_name" json:"pull_request_name"`
	AuthorID        string `db:"author_id" json:"author_id"`
	Status          string `db:"status" json:"status"`
}
