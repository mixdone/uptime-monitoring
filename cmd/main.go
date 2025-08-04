package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mixdone/uptime-monitoring/internal/config"
	"github.com/mixdone/uptime-monitoring/internal/database"
	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/repository"
	"github.com/mixdone/uptime-monitoring/internal/services"
	"github.com/mixdone/uptime-monitoring/internal/transport"
	"github.com/mixdone/uptime-monitoring/pkg/logger"
)

// @title Uptime Monitoring API
// @version 1.0
// @description API for uptime monitoring service
// @host localhost:8080
// @BasePath /

// @securitydefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <your-token>
func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Error loading config: " + err.Error())
	}

	log, err := logger.InitStructuredLogger(cfg)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	db, err := database.NewDB(cfg)
	if err != nil {
		log.WithError(err).Error("Failed db connection")
		return
	}
	defer db.Close()

	log.Info("Connected to PostgreSQL")

	repository := repository.NewRepository(db)
	services := services.NewServices(repository, *cfg, log)
	handlers := transport.NewHandler(services, log)

	srv := new(models.ServerApi)

	go func() {
		if err := srv.Run(cfg.Server.Host, cfg.Server.Port, handlers.InitRoutes()); err != nil {
			log.WithError(err).Error("Server stopped with error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server shutdown failed")
	} else {
		log.Info("Server shutdown successfully")
	}
}
