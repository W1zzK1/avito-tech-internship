package service

import (
	"avito-tech-internship/internal/domain"
	"avito-tech-internship/internal/storage"
)

type Service struct {
	repo storage.Repository
}

func NewService(repo storage.Repository) *Service {
	return &Service{repo: repo}
}

// Users
func (s *Service) AddNewUser(user *domain.User) (*domain.User, error) {
	user, err := s.repo.AddNewUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetUserByID(id string) (*domain.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) SetUserActive(userID string, isActive bool) (*domain.User, error) {
	err := s.repo.SetUserActive(userID, isActive)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(userID)
}

// Teams
func (s *Service) CreateNewTeam(team *domain.Team) error {
	err := s.repo.AddTeam(team)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetTeamByName(teamName string) (*domain.Team, error) {
	team, err := s.repo.GetTeamByName(teamName)
	if err != nil {
		return nil, err
	}
	return team, nil
}
