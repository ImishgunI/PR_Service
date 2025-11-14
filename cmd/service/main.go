package main

import (
	"PullRequestService/internal/db"
	"PullRequestService/internal/routes"
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

	routes.SetRoutes(r, db)

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
