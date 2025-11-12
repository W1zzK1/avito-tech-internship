package storage

import (
	"avito-tech-internship/domain"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetUserByID(userId int) (*domain.User, error)
	//GetUserWithTeam(userId string) (*domain.User, error)
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetUserByID(userID int) (*domain.User, error) {
	var user domain.User
	println(userID)
	query := "SELECT * FROM users WHERE user_id = $1"
	err := r.db.Get(&user, query, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, err
}
