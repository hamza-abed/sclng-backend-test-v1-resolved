package worker

import (
	"database/sql"
	"time"

	"github.com/Scalingo/sclng-backend-test-v1/service"
	"github.com/sirupsen/logrus"
)

type MainWorker interface {
	JobRoutine()
	job()
}

type GithubWorker struct {
	db            *sql.DB
	logger        logrus.FieldLogger
	githubToken   string
	githubService service.IGithubService
}

func NewMainWorker(db *sql.DB, logger logrus.FieldLogger, githubToken string) MainWorker {
	return &GithubWorker{
		db: db, logger: logger, githubService: service.NewGithubService(logger, githubToken),
	}
}
func (worker *GithubWorker) JobRoutine() {
	for {

		if worker.githubService.CanMakeACall() { // check ratelimit before creating a new goroutine
			go worker.job()
		}
		time.Sleep(time.Duration(160) * 1000 * time.Millisecond)

	}
}

func (worker *GithubWorker) job() {

	repos, err := worker.githubService.RetrieveRepos()
	if err != nil {
		worker.logger.Error("error while getting repos from github : ", err)
		return
	}
	for _, repo := range repos {
		repoWorker := newRepoWorker(worker.db, worker.logger, repo, worker.githubService)
		go repoWorker.manageRepo()
	}
}
