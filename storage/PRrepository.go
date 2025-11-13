package storage

import (
	"avito-tech-internship/domain"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// PR методы
func (r *PostgresRepository) CreatePullRequest(pr *domain.PullRequest) error {
	query := `
        INSERT INTO pull_requests (id, name, author_id, status) 
        VALUES ($1, $2, $3, 'OPEN')
    `
	_, err := r.db.Exec(query, pr.ID, pr.Name, pr.AuthorId)
	return err
}

func (r *PostgresRepository) GetPullRequestByID(prID string) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	query := `
        SELECT 
            id as pull_request_id,
            name as pull_request_name, 
            author_id,
            status, 
            created_at, 
            merged_at 
        FROM pull_requests WHERE id = $1
    `
	err := r.db.Get(&pr, query, prID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("pull request not found")
	}
	if err != nil {
		return nil, err
	}

	// Получаем ревьюеров
	reviewers, err := r.GetPRReviewers(prID)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *PostgresRepository) MergePullRequest(prID string) error {
	query := `
        UPDATE pull_requests 
        SET status = 'MERGED', merged_at = NOW() 
        WHERE id = $1 AND status != 'MERGED'
    `
	result, err := r.db.Exec(query, prID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		// PR уже мерджен или не существует
		return nil
	}

	return nil
}

func (r *PostgresRepository) PRExists(prID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE id = $1)`
	err := r.db.Get(&exists, query, prID)
	return exists, err
}

// Reviewer методы
func (r *PostgresRepository) AssignReviewers(prID string, reviewerIDs []string) error {
	if len(reviewerIDs) == 0 {
		return nil
	}

	query := `INSERT INTO pull_request_reviewers (pull_request_id, user_id) VALUES ($1, $2)`
	for _, reviewerID := range reviewerIDs {
		_, err := r.db.Exec(query, prID, reviewerID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresRepository) GetPRReviewers(prID string) ([]string, error) {
	var reviewers []string
	query := `SELECT user_id FROM pull_request_reviewers WHERE pull_request_id = $1`
	err := r.db.Select(&reviewers, query, prID)
	return reviewers, err
}

func (r *PostgresRepository) ReplaceReviewer(prID, oldReviewerID, newReviewerID string) error {
	query := `
        UPDATE pull_request_reviewers 
        SET user_id = $1 
        WHERE pull_request_id = $2 AND user_id = $3
    `
	result, err := r.db.Exec(query, newReviewerID, prID, oldReviewerID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("reviewer not assigned to this PR")
	}
	return nil
}

func (r *PostgresRepository) GetUserAssignedPRs(userID string) ([]domain.PullRequestShort, error) {
	var prs []domain.PullRequestShort
	query := `
        SELECT 
            pr.id,
            pr.name,
            pr.author_id,
            pr.status
        FROM pull_requests pr
        JOIN pull_request_reviewers prr ON pr.id = prr.pull_request_id
        WHERE prr.user_id = $1
        ORDER BY pr.created_at DESC
    `
	err := r.db.Select(&prs, query, userID)
	return prs, err
}

// Business logic helpers
func (r *PostgresRepository) GetActiveTeamMembers(teamID string, excludeUserID string) ([]domain.User, error) {
	var users []domain.User
	query := `
        SELECT 
            id,
            username, 
            is_active,
            team_id
        FROM users 
        WHERE team_id = $1 AND is_active = true AND id != $2
        ORDER BY id
    `
	err := r.db.Select(&users, query, teamID, excludeUserID)
	return users, err
}

func (r *PostgresRepository) GetRandomActiveTeamMember(teamID string, excludeUserIDs []string) (*domain.User, error) {
	if len(excludeUserIDs) == 0 {
		excludeUserIDs = []string{""}
	}

	placeholders := make([]string, len(excludeUserIDs))
	args := []interface{}{teamID}
	for i, id := range excludeUserIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args = append(args, id)
	}

	query := fmt.Sprintf(`
        SELECT 
            id,
            username,
            is_active,
            team_id
        FROM users 
        WHERE team_id = $1 AND is_active = true AND id NOT IN (%s)
        ORDER BY RANDOM()
        LIMIT 1
    `, strings.Join(placeholders, ", "))

	var user domain.User
	err := r.db.Get(&user, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // Нет доступных кандидатов
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Получить команду автора PR
func (r *PostgresRepository) GetAuthorTeam(authorID string) (string, error) {
	var teamID string
	query := `SELECT team_id FROM users WHERE id = $1`
	err := r.db.Get(&teamID, query, authorID)
	if err != nil {
		return "", fmt.Errorf("author not found")
	}
	return teamID, nil
}
