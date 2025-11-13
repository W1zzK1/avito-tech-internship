package service

import (
	"avito-tech-internship/domain"
	"avito-tech-internship/storage"
)

type Service struct {
	repo storage.Repository
}

func NewService(repo storage.Repository) *Service {
	return &Service{repo: repo}
}

// Users
func (service *Service) GetUser(id string) (*domain.User, error) {
	user, err := service.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *Service) SetUserActive(userID string, isActive bool) (*domain.User, error) {
	err := service.repo.SetUserActive(userID, isActive)
	if err != nil {
		return nil, err
	}
	return service.repo.GetUserByID(userID)
}

// Teams
func (service *Service) CreateNewTeam(team *domain.Team) error {
	err := service.repo.AddTeam(team)
	if err != nil {
		return err
	}
	return nil
}
