package api

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/sclng-backend-test-v1/repository"
	"github.com/sirupsen/logrus"
)

type Server struct {
	log            logrus.FieldLogger
	router         *handlers.Router
	db             *sql.DB
	repoRepository repository.IRepoRepository
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(logger logrus.FieldLogger, db *sql.DB) (*Server, error) {

	server := &Server{
		log:            logger,
		db:             db,
		repoRepository: repository.NewRepoRepository(db, logger),
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	server.log.Info("Initializing routes")
	router := handlers.NewRouter(server.log)
	router.HandleFunc("/ping", server.pongHandler)
	// Initialize web server and configure the following routes:
	// GET /repos
	router.HandleFunc("/repos", server.getRepositories).Methods("GET")
	// GET /stats
	router.HandleFunc("/stats", server.getRepositories).Methods("GET")

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	os.Setenv("PPROF_ENABLED", "true")
	log := server.log.WithField("port", address)
	log.Info("Listening...")
	err := http.ListenAndServe(address, server.router)
	if err != nil {
		log.WithError(err).Error("Fail to listen to the given port")
		os.Exit(2)
	}
	return err
}
