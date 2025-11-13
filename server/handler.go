package server

import (
	"avito-tech-internship/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		writeError(c, http.StatusBadRequest, "MISSING_PARAM", "user id is required")
		return
	}

	user, err := h.service.GetUser(userID)
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

func writeError(c *gin.Context, status int, code, message string) {
	errResponse := struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}{}
	errResponse.Error.Message = message
	errResponse.Error.Code = code
	c.JSON(status, errResponse)
}
