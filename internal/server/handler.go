package server

import (
	"avito-tech-internship/internal/domain"
	"avito-tech-internship/internal/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) AddNewUser(c *gin.Context) {
	var user domain.User
	body, _ := c.GetRawData()
	if err := json.Unmarshal(body, &user); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_JSON", fmt.Sprintf("Invalid JSON: %v", err))
		return
	}
	createdUser, err := h.service.AddNewUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "team not found") {
			writeError(c, http.StatusNotFound, "TEAM_NOT_FOUND", "team not found")
			return
		}
		if strings.Contains(err.Error(), "user already exists") {
			writeError(c, http.StatusConflict, "USER_EXISTS", "user already exists")
			return
		}
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, createdUser)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		writeError(c, http.StatusBadRequest, "MISSING_PARAM", "user id is required")
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(c, http.StatusNotFound, "NOT_FOUND", "user not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *Handler) SetUserActive(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		IsActive bool   `json:"is_active" binding:"required"`
	}

	body, _ := c.GetRawData()
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_JSON", fmt.Sprintf("Invalid JSON: %v", err))
		return
	}

	user, err := h.service.SetUserActive(req.UserID, req.IsActive)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(c, http.StatusNotFound, "NOT_FOUND", "user not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *Handler) createNewTeam(c *gin.Context) {
	var team domain.Team
	body, _ := c.GetRawData()
	if err := json.Unmarshal(body, &team); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_JSON", fmt.Sprintf("Invalid JSON: %v", err))
		return
	}
	err := h.service.CreateNewTeam(&team)
	if err != nil {
		if err.Error() == "TEAM_EXISTS" {
			writeError(c, http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists")
			return
		}
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response := domain.Team{
		TeamName: team.TeamName,
		Members:  team.Members,
	}
	c.JSON(http.StatusOK, gin.H{
		"team": response,
	})
}

func (h *Handler) GetTeamByName(c *gin.Context) {
	teamName := c.Param("teamName")
	if teamName == "" {
		writeError(c, http.StatusBadRequest, "MISSING_PARAM", "user id is required")
		return
	}

	team, err := h.service.GetTeamByName(teamName)
	if err != nil {
		writeError(c, http.StatusNotFound, "NOT_FOUND", "resource not found")
		return
	}

	c.JSON(http.StatusOK, team)

}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

func writeError(c *gin.Context, status int, code, message string) {
	errResponse := struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}{}
	errResponse.Error.Message = message
	errResponse.Error.Code = code
	slog.ErrorContext(c, message, errResponse)
	c.JSON(status, errResponse)
}
