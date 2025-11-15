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

// Stats
func (s *Service) GetStats() (*domain.StatsResponse, error) {
	userStats, err := s.repo.GetPRReviewersStats()
	if err != nil {
		return nil, err
	}

	prStats, err := s.repo.GetDetailedPRStats()
	if err != nil {
		return nil, err
	}

	teamStats, err := s.repo.GetTeamStats()
	if err != nil {
		return nil, err
	}

	summary, err := s.calculateSummary(userStats, prStats, teamStats)
	if err != nil {
		return nil, err
	}

	return &domain.StatsResponse{
		UserStats: userStats,
		PRStats:   prStats,
		TeamStats: teamStats,
		Summary:   summary,
	}, nil
}

func (s *Service) calculateSummary(userStats []*domain.UserStats, prStats []*domain.PRStats, teamStats []*domain.TeamStats) (*domain.StatsSummary, error) {
	summary := &domain.StatsSummary{}

	// Подсчет пользователей
	userSet := make(map[string]bool)
	for _, user := range userStats {
		userSet[user.UserID] = true
	}
	summary.TotalUsers = len(userSet)

	// Подсчет команд
	summary.TotalTeams = len(teamStats)

	// Подсчет PR
	summary.TotalPRs = len(prStats)
	for _, pr := range prStats {
		if pr.Status == "OPEN" {
			summary.OpenPRs++
		} else if pr.Status == "MERGED" {
			summary.MergedPRs++
		}
	}

	// Подсчет ревью
	totalReviews := 0
	for _, user := range userStats {
		totalReviews += user.PRCount
	}
	summary.TotalReviews = totalReviews

	// Среднее количество ревьюеров на PR
	if summary.TotalPRs > 0 {
		summary.AvgReviewsPerPR = totalReviews / summary.TotalPRs
	}

	return summary, nil
}
