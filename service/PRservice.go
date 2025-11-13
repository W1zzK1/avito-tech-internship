package service

import (
	"avito-tech-internship/domain"
	"errors"
	"fmt"
	"math/rand"
)

func (s *Service) CreatePullRequest(req *domain.CreatePRRequest) (*domain.PullRequest, error) {
	// Проверяем существование PR
	exists, err := s.repo.PRExists(req.PRID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("PR_EXISTS")
	}

	// Проверяем существование автора и получаем его команду
	authorTeamID, err := s.repo.GetAuthorTeam(req.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("author not found")
	}

	// Создаем PR
	pr := &domain.PullRequest{
		ID:       req.PRID,
		Name:     req.Name,
		AuthorId: req.AuthorID,
		Status:   "OPEN",
	}

	err = s.repo.CreatePullRequest(pr)
	if err != nil {
		return nil, err
	}

	// Назначаем ревьюеров
	reviewers, err := s.assignReviewers(authorTeamID, req.AuthorID)
	if err != nil {
		return nil, err
	}

	// Сохраняем ревьюеров
	if len(reviewers) > 0 {
		err = s.repo.AssignReviewers(pr.ID, reviewers)
		if err != nil {
			return nil, err
		}
	}

	// Возвращаем созданный PR с ревьюерами
	return s.repo.GetPullRequestByID(pr.ID)
}

func (s *Service) MergePullRequest(prID string) (*domain.PullRequest, error) {
	// Получаем текущее состояние PR
	pr, err := s.repo.GetPullRequestByID(prID)
	if err != nil {
		return nil, err
	}

	// Если уже мерджен - возвращаем как есть (идемпотентность)
	if pr.Status == "MERGED" {
		return pr, nil
	}

	// Мерджим PR
	err = s.repo.MergePullRequest(prID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetPullRequestByID(prID)
}

func (s *Service) ReassignReviewer(prID, oldReviewerID string) (*domain.PullRequest, string, error) {
	// Получаем PR
	pr, err := s.repo.GetPullRequestByID(prID)
	if err != nil {
		return nil, "", err
	}

	// Проверяем что PR не мерджен
	if pr.Status == "MERGED" {
		return nil, "", errors.New("PR_MERGED")
	}

	// Проверяем что старый ревьюер назначен на этот PR
	isAssigned := false
	for _, reviewer := range pr.AssignedReviewers {
		if reviewer == oldReviewerID {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		return nil, "", errors.New("NOT_ASSIGNED")
	}

	// Получаем команду старого ревьюера
	oldReviewerTeamID, err := s.repo.GetAuthorTeam(oldReviewerID)
	if err != nil {
		return nil, "", fmt.Errorf("reviewer not found")
	}

	// Ищем нового ревьюера из той же команды
	excludeIDs := append(pr.AssignedReviewers, pr.AuthorId)
	newReviewer, err := s.repo.GetRandomActiveTeamMember(oldReviewerTeamID, excludeIDs)
	if err != nil {
		return nil, "", err
	}
	if newReviewer == nil {
		return nil, "", errors.New("NO_CANDIDATE")
	}

	// Заменяем ревьюера
	err = s.repo.ReplaceReviewer(prID, oldReviewerID, newReviewer.UserId)
	if err != nil {
		return nil, "", err
	}

	// Возвращаем обновленный PR
	updatedPR, err := s.repo.GetPullRequestByID(prID)
	return updatedPR, newReviewer.UserId, err
}

func (s *Service) GetUserAssignedPRs(userID string) ([]domain.PullRequestShort, error) {
	// Проверяем что пользователь существует
	_, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserAssignedPRs(userID)
}

// Вспомогательный метод для назначения ревьюеров
func (s *Service) assignReviewers(teamID, excludeUserID string) ([]string, error) {
	members, err := s.repo.GetActiveTeamMembers(teamID, excludeUserID)
	if err != nil {
		return nil, err
	}

	// Если нет доступных ревьюеров
	if len(members) == 0 {
		return []string{}, nil
	}

	// Перемешиваем для случайного выбора
	shuffled := make([]domain.User, len(members))
	copy(shuffled, members)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Выбираем до 2 ревьюеров
	var reviewerIDs []string
	maxReviewers := min(2, len(shuffled))
	for i := 0; i < maxReviewers; i++ {
		reviewerIDs = append(reviewerIDs, shuffled[i].UserId)
	}

	return reviewerIDs, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
