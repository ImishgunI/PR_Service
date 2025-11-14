package main

import (
	"PullRequestService/internal/db"
	"PullRequestService/internal/handler"
	"PullRequestService/internal/repository"
	"PullRequestService/pkg/config"
	"PullRequestService/pkg/logger"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode("release")
	log := logger.New()
	r := gin.Default()
	config.InitConfig()
	db := db.New()
	if db == nil {
		log.Warn("Object db equal nil")
		panic("DataBase connection failed")
	}
	defer db.Close()
	port := config.GetString("PORT")

	rt := repository.NewTeamRepository(db)
	ru := repository.NewUserHandler(db)
	h := handler.NewTeamHandler(rt)
	uh := handler.NewUserHandler(ru)
	pr := repository.NewPRRepository(db)
	ph := handler.NewPRHandler(pr)
	r.POST("/team/add", h.AddTeam)
	r.GET("/team/get", h.GetTeam)
	r.POST("/users/setIsActive", uh.SetActive)
	r.POST("/pullRequest/create", ph.CreatePR)
	r.POST("/pullRequest/merge", ph.MergePR)
	r.POST("/pullRequest/reassign", ph.ReassignPR)
	r.GET("/users/getReview", uh.GetReview)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Infof("Server started on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Server shutdown error: %v", err)
	} else {
		log.Info("Server stopped gracefully")
	}
}
