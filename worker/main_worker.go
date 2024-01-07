package worker

import (
	"database/sql"
	"sync"
	"time"

	"github.com/Scalingo/sclng-backend-test-v1/repository"
	"github.com/Scalingo/sclng-backend-test-v1/service"
	"github.com/sirupsen/logrus"
)

type MainWorker interface {
	JobRoutine()
	job()
}

type GithubWorker struct {
	db             *sql.DB
	logger         logrus.FieldLogger
	githubToken    string
	githubService  service.IGithubService
	repoRepository repository.IRepoRepository // this should be created here to manage concurrent access to DB
}

func NewMainWorker(db *sql.DB, logger logrus.FieldLogger, githubToken string) MainWorker {
	return &GithubWorker{
		db: db, logger: logger, githubService: service.NewGithubService(logger, githubToken), repoRepository: repository.NewRepoRepository(db, logger),
	}
}
func (worker *GithubWorker) JobRoutine() {
	for {
		if worker.githubService.CanMakeACall() { // check ratelimit before creating a new goroutine
			worker.job()
		}
		time.Sleep(time.Duration(1000) * 1000 * time.Millisecond) // once job finished wait a second !
	}
}

func (worker *GithubWorker) job() {
	var waitGroup sync.WaitGroup
	sinceRepoId, err := worker.repoRepository.GetMaxRepositoryId()
	if err != nil {
		worker.logger.Error("error while getting max repo id DB : ", err)
	}
	if sinceRepoId == -1 {
		sinceRepoId = 712574319 // some recent repo
	}
	repos, err := worker.githubService.FetchRepos(sinceRepoId)
	if err != nil {
		worker.logger.Error("error while getting repos from github : ", err)
		return
	}
	waitGroup.Add(len(repos))
	for _, repo := range repos {
		repoWorker := newRepoWorker(worker.db, worker.logger, repo, worker.githubService, worker.repoRepository)
		go repoWorker.manageRepo(&waitGroup)
	}
	waitGroup.Wait()
}
