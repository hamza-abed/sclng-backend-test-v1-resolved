package main

import (
	"fmt"
	"os"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/api"
	"github.com/Scalingo/sclng-backend-test-v1/migration"
	"github.com/Scalingo/sclng-backend-test-v1/util"
	"github.com/Scalingo/sclng-backend-test-v1/worker"
)

func main() {
	log := logger.Default()
	log.Info("Initializing app")
	cfg, err := util.NewConfig()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(1)
	}

	// connect DB
	db := util.ConnectPGDB(cfg, log)

	// migrate database if environment dev
	if cfg.Env == "dev" {
		migration.Migrate(db, log, !cfg.EraseDbWhenMigrate)
	}

	// start worker fetching
	mainWorker := worker.NewMainWorker(db, log, cfg.GithubToken)
	go mainWorker.JobRoutine()

	// start API server
	server, err := api.NewServer(log, db)
	if err != nil {
		log.Error("error while creating server", err)
	}
	server.Start(fmt.Sprintf(":%d", cfg.Port))
}
