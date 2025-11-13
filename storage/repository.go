package storage

import (
	"avito-tech-internship/domain"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Repository interface {
	// Users
	GetUserByID(userId string) (*domain.User, error)
	SetUserActive(userId string, isActive bool) error
	AddNewUser(user *domain.User) (*domain.User, error)

	//Teams
	AddTeam(team *domain.Team) error
	GetTeamByName(name string) (*domain.Team, error)
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Users
func (r *PostgresRepository) AddNewUser(user *domain.User) (*domain.User, error) {
	var teamID string
	err := r.db.Get(&teamID, "SELECT id FROM teams WHERE name = $1", user.TeamName)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	query := `
        INSERT INTO users (id, username, is_active, team_id) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, username, is_active, team_id
    `

	var newUser domain.User
	err = r.db.QueryRow(query, user.UserId, user.Username, user.IsActive, teamID).
		Scan(&newUser.UserId, &newUser.Username, &newUser.IsActive, &newUser.TeamId)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, fmt.Errorf("user already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	newUser.TeamName = user.TeamName

	return &newUser, nil
}

func (r *PostgresRepository) GetUserByID(userID string) (*domain.User, error) {
	var user domain.User
	query := "SELECT * FROM users WHERE id = $1"
	err := r.db.Get(&user, query, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, err
}

func (r *PostgresRepository) SetUserActive(userID string, iaActive bool) error {
	query := "UPDATE users SET is_active = $1 WHERE id = $2"
	result, err := r.db.Exec(query, iaActive, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user was not found")
	}
	return nil
}

// Teams
func (r *PostgresRepository) AddTeam(team *domain.Team) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Создаем команду (ID сгенерируется автоматически)
	var teamID string
	err = tx.QueryRow(
		"INSERT INTO teams (name) VALUES ($1) RETURNING id",
		team.TeamName,
	).Scan(&teamID)

	if err != nil {
		// Проверяем на уникальность имени команды
		if strings.Contains(err.Error(), "unique constraint") {
			return errors.New("TEAM_EXISTS")
		}
		return err
	}

	// Создаем/обновляем пользователей
	for _, member := range team.Members {
		_, err := tx.Exec(`
            INSERT INTO users (id, username, is_active, team_id) 
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) 
            DO UPDATE SET username = $2, is_active = $3, team_id = $4
        `, member.UserId, member.Username, member.IsActive, teamID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetTeamByName(teamName string) (*domain.Team, error) {
	var teamExists bool
	err := r.db.Get(&teamExists, "SELECT exists(SELECT 1 FROM teams WHERE name = $1)", teamName)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if !teamExists {
		return nil, errors.New("team not found")
	}

	var members []*domain.User
	query := `
        SELECT 
            u.id, 
            u.username, 
            u.is_active,
            t.name as team_name,
            u.team_id
        FROM users u 
        JOIN teams t ON u.team_id = t.id 
        WHERE t.name = $1
        ORDER BY u.username
    `
	err = r.db.Select(&members, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	team := &domain.Team{
		TeamName: teamName,
		Members:  members,
	}
	return team, nil
}
