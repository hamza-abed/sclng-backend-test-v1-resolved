package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
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
		migration.Migrate(db, log, cfg.EraseDbWhenMigrate)
	}

	// start worker fetching
	mainWorker := worker.NewMainWorker(db, log, cfg.GithubToken)
	go mainWorker.JobRoutine()

	log.Info("Initializing routes")
	router := handlers.NewRouter(log)
	router.HandleFunc("/ping", pongHandler)
	// Initialize web server and configure the following routes:
	// GET /repos
	// GET /stats

	log = log.WithField("port", cfg.Port)
	log.Info("Listening...")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	if err != nil {
		log.WithError(err).Error("Fail to listen to the given port")
		os.Exit(2)
	}
}

func pongHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(map[string]string{"status": "pong"})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}
