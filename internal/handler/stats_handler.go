package handler

import (
	"PullRequestService/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatisticFactory interface {
	GetReviewerStats(c *gin.Context)
}

type StatHandler struct {
	repo *repository.Statistic
}

func NewStatHandler(statRepo *repository.Statistic) *StatHandler {
	return &StatHandler{repo: statRepo}
}

func (h *StatHandler) GetReviewerStats(c *gin.Context) {
	stats, err := h.repo.GetReviewerStats(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reviewer_stats": stats})
}
