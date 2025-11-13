package storage

import (
	"avito-tech-internship/domain"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	// Users
	GetUserByID(userId string) (*domain.User, error)
	SetUserActive(userId string, isActive bool) error

	//AddTeam()
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Users
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
