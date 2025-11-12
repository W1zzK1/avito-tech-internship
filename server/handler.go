package server

import (
	"avito-tech-internship/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

type findUserRequest struct {
	Id int `uri:"id" binding:"required"`
}

func (h *Handler) GetUser(c *gin.Context) {
	var req findUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}

	user, err := h.service.GetUser(req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"pivo": err.Error()})
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
