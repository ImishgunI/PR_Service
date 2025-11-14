package handler

import (
	"PullRequestService/internal/repository"
	"PullRequestService/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PullRequestService interface {
	CreatePR(c *gin.Context)
	MergePR(c *gin.Context)
	ReassignPR(c *gin.Context)
}

type PRHandler struct {
	rp  *repository.PRRepository
	log logger.Logger
}

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

func NewPRHandler(repo *repository.PRRepository) *PRHandler {
	log := logger.New()
	return &PRHandler{rp: repo, log: log}
}

func (h *PRHandler) CreatePR(c *gin.Context) {
	var req CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_BODY",
				"message": err.Error(),
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	pr, reviewers, err := h.rp.CreatePR(ctx, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		switch err.Error() {
		case "PR_EXISTS":
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "PR_EXISTS",
					"message": "PR id already exists",
				},
			})
			return
		case "AUTHOR_NOT_FOUND":
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "author or team not found",
				},
			})
			return
		default:
			h.log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"pr": gin.H{
			"pull_request_id":    pr.PullRequestID,
			"pull_request_name":  pr.PullRequestName,
			"author_id":          pr.AuthorID,
			"status":             pr.Status,
			"assigned_reviewers": reviewers,
			"needMoreReviewers":  pr.NeedMore,
			"createdAt":          pr.CreatedAt.Format(time.RFC3339),
		},
	})
}

func (h *PRHandler) MergePR(c *gin.Context) {
	var req struct {
		PullRequestID string `json:"pull_request_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_BODY", "message": err.Error()},
		})
		return
	}

	pr, err := h.rp.UpdateStatus(c.Request.Context(), req.PullRequestID)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "PR not found"},
			})
			return
		}
		h.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pr": pr})
}

func (h *PRHandler) ReassignPR(c *gin.Context) {
	var req struct {
		PullRequestID string `json:"pull_request_id" binding:"required"`
		OldReviewerID string `json:"old_user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_BODY", "message": err.Error()}})
		return
	}

	pr, newUserID, err := h.rp.ReassignPR(c.Request.Context(), req.PullRequestID, req.OldReviewerID)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "PR not found"}})
		case "PR_MERGED":
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "PR_MERGED", "message": "cannot reassign on merged PR"}})
		case "NOT_ASSIGNED":
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "NOT_ASSIGNED", "message": "reviewer is not assigned to this PR"}})
		case "NO_CANDIDATE":
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "NO_CANDIDATE", "message": "no active replacement candidate in team"}})
		default:
			h.log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr":          pr,
		"replaced_by": newUserID,
	})
}
