package domain

import "time"

type User struct {
	UserId   string `db:"id" json:"userId"`
	Username string `db:"username" json:"username"`
	IsActive bool   `db:"is_active" json:"isActive"`
	TeamName string `db:"team_name" json:"teamName,omitempty"`
	TeamId   int    `db:"team_id" json:"-"`
}

type Team struct {
	TeamName string  `json:"team_name,omitempty"`
	Members  []*User `json:"members"`
}

type PullRequest struct {
	ID                string     `db:"pull_request_id" json:"pull_request_id"`
	Name              string     `db:"pull_request_name" json:"pull_request_name"`
	AuthorId          string     `db:"author_id" json:"author_id"`
	Status            string     `db:"status" json:"status"`
	AssignedReviewers []string   `db:"-" json:"assigned_reviewers"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	MergedAt          *time.Time `db:"merged_at" json:"merged_at,omitempty"`
}
