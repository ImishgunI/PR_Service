package routes

import (
	"PullRequestService/internal/db"
	"PullRequestService/internal/handler"
	"PullRequestService/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine, db *db.DataBase) {
	rt := repository.NewTeamRepository(db)
	ru := repository.NewUserHandler(db)
	h := handler.NewTeamHandler(rt)
	uh := handler.NewUserHandler(ru)
	pr := repository.NewPRRepository(db)
	ph := handler.NewPRHandler(pr)
	sr := repository.NewStatisticRepository(db)
	sh := handler.NewStatHandler(sr)
	r.POST("/team/add", h.AddTeam)
	r.GET("/team/get", h.GetTeam)
	r.POST("/users/setIsActive", uh.SetActive)
	r.POST("/pullRequest/create", ph.CreatePR)
	r.POST("/pullRequest/merge", ph.MergePR)
	r.POST("/pullRequest/reassign", ph.ReassignPR)
	r.GET("/users/getReview", uh.GetReview)
	r.GET("/stats/reviewers", sh.GetReviewerStats)
}
