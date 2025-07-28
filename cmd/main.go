package main

import (
	"github.com/mixdone/uptime-monitoring/internal/config"
	"github.com/mixdone/uptime-monitoring/internal/database"
	"github.com/mixdone/uptime-monitoring/pkg/logger"
)

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
	}
	defer db.Close()

	log.Info("Connected to PostgreSQL")

}
