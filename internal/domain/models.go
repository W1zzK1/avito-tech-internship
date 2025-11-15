package domain

import "time"

type User struct {
	UserId   string `db:"id" json:"user_id"`
	Username string `db:"username" json:"username"`
	IsActive bool   `db:"is_active" json:"isActive"`
	TeamName string `db:"team_name" json:"teamName,omitempty"`
	TeamId   string `db:"team_id" json:"-"`
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

type PullRequestShort struct {
	ID       string `db:"id" json:"pull_request_id"`
	Name     string `db:"name" json:"pull_request_name"`
	AuthorID string `db:"author_id" json:"author_id"`
	Status   string `db:"status" json:"status"`
}

// Request/Response структуры
type CreatePRRequest struct {
	PRID     string `json:"pull_request_id" binding:"required"`
	Name     string `json:"pull_request_name" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
}

type MergePRRequest struct {
	PRID string `json:"pull_request_id" binding:"required"`
}

type ReassignRequest struct {
	PRID          string `json:"pull_request_id" binding:"required"`
	OldReviewerID string `json:"old_reviewer_id" binding:"required"`
}

// stats
type UserStats struct {
	UserID        string `db:"user_id" json:"user_id"`
	Username      string `db:"username" json:"username"`
	TeamName      string `db:"team_name" json:"team_name"`
	PRCount       int    `db:"pr_count" json:"pr_count"`
	MergedPRCount int    `db:"merged_pr_count" json:"merged_pr_count"`
}

type PRStats struct {
	PRID          string     `db:"pull_request_id" json:"pull_request_id"`
	PRName        string     `db:"pull_request_name" json:"pull_request_name"`
	Status        string     `db:"status" json:"status"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	MergedAt      *time.Time `db:"merged_at" json:"merged_at,omitempty"`
	AuthorName    string     `db:"author_name" json:"author_name"`
	AuthorTeam    string     `db:"author_team" json:"author_team"`
	ReviewerCount int        `db:"reviewer_count" json:"reviewer_count"`
	ReviewerNames string     `db:"reviewer_names" json:"reviewer_names"`
}

type TeamStats struct {
	TeamName        string `db:"team_name" json:"team_name"`
	MemberCount     int    `db:"member_count" json:"member_count"`
	AuthoredPRCount int    `db:"authored_pr_count" json:"authored_pr_count"`
	ReviewedPRCount int    `db:"reviewed_pr_count" json:"reviewed_pr_count"`
	MergedPRCount   int    `db:"merged_pr_count" json:"merged_pr_count"`
}

type StatsResponse struct {
	UserStats []*UserStats  `json:"user_stats"`
	PRStats   []*PRStats    `json:"pr_stats"`
	TeamStats []*TeamStats  `json:"team_stats"`
	Summary   *StatsSummary `json:"summary"`
}

type StatsSummary struct {
	TotalUsers      int `json:"total_users"`
	TotalTeams      int `json:"total_teams"`
	TotalPRs        int `json:"total_prs"`
	OpenPRs         int `json:"open_prs"`
	MergedPRs       int `json:"merged_prs"`
	TotalReviews    int `json:"total_reviews"`
	AvgReviewsPerPR int `json:"avg_reviews_per_pr"`
}
