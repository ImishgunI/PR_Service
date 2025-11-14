package handler

import (
	"PullRequestService/internal/repository"
	"PullRequestService/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	SetActive(c *gin.Context)
	GetReview(c *gin.Context)
}

type UserHandler struct {
	repo *repository.UserRepository
	log  logger.Logger
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	log := logger.New()
	return &UserHandler{repo: repo, log: log}
}

func (u *UserHandler) SetActive(c *gin.Context) {
	var req struct {
		UserId   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()},
		})
		return
	}
	ctx := c.Request.Context()
	user, err := u.repo.SetIsActive(ctx, req.UserId, req.IsActive)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "user not found"},
			})
			return
		}
		u.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (u *UserHandler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_PARAM", "message": "user_id is required"},
		})
		return
	}

	prs, err := u.repo.GetReview(c.Request.Context(), userID)
	if err != nil {
		u.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
