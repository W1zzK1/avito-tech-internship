package server

import (
	"avito-tech-internship/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// PR handlers
func (h *Handler) CreatePullRequest(c *gin.Context) {
	var req domain.CreatePRRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	pr, err := h.service.CreatePullRequest(&req)
	if err != nil {
		switch err.Error() {
		case "PR_EXISTS":
			writeError(c, http.StatusConflict, "PR_EXISTS", "PR id already exists")
		case "author not found":
			writeError(c, http.StatusNotFound, "NOT_FOUND", "author not found")
		default:
			writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"pr": pr,
	})
}

func (h *Handler) MergePullRequest(c *gin.Context) {
	var req domain.MergePRRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	pr, err := h.service.MergePullRequest(req.PRID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(c, http.StatusNotFound, "NOT_FOUND", "pull request not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr": pr,
	})
}

func (h *Handler) ReassignReviewer(c *gin.Context) {
	var req domain.ReassignRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	pr, newReviewerID, err := h.service.ReassignReviewer(req.PRID, req.OldReviewerID)
	if err != nil {
		switch err.Error() {
		case "PR_MERGED":
			writeError(c, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
		case "NOT_ASSIGNED":
			writeError(c, http.StatusConflict, "NOT_ASSIGNED", "reviewer is not assigned to this PR")
		case "NO_CANDIDATE":
			writeError(c, http.StatusConflict, "NO_CANDIDATE", "no active replacement candidate in team")
		case "pull request not found", "reviewer not found":
			writeError(c, http.StatusNotFound, "NOT_FOUND", "resource not found")
		default:
			writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr":          pr,
		"replaced_by": newReviewerID,
	})
}

func (h *Handler) GetUserReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		writeError(c, http.StatusBadRequest, "MISSING_PARAM", "user_id parameter is required")
		return
	}

	prs, err := h.service.GetUserAssignedPRs(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(c, http.StatusNotFound, "NOT_FOUND", "user not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
