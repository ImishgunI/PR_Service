package handler

import (
	"PullRequestService/internal/models"
	"PullRequestService/internal/repository"
	"PullRequestService/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamService interface {
	AddTeam(c *gin.Context)
	GetTeam(c *gin.Context)
}

type TeamHandler struct {
	repo *repository.TeamRepository
	log  logger.Logger
}

func NewTeamHandler(repo *repository.TeamRepository) *TeamHandler {
	log := logger.New()
	return &TeamHandler{repo: repo, log: log}
}

func (h *TeamHandler) AddTeam(c *gin.Context) {
	var req models.Team
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.CreateTeam(c, req)
	if err != nil {
		if err.Error() == "TEAM_EXISTS" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{"code": "TEAM_EXISTS", "message": "team already exists"},
			})
			return
		}
		h.log.Errorf("Failed to create team: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"team": req})
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	name := c.Query("team_name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team_name is required"})
		return
	}

	team, err := h.repo.GetTeam(c.Request.Context(), name)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "team not found"},
			})
			return
		}
		h.log.Errorf("Failed to get team: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}
