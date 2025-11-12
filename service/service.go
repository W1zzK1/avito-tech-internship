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

func (service *Service) GetUser(id int) (*domain.User, error) {
	user, err := service.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
