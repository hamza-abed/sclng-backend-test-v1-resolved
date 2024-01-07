package worker

import (
	"database/sql"
	"sync"

	"github.com/Scalingo/sclng-backend-test-v1/factory"
	"github.com/Scalingo/sclng-backend-test-v1/model"
	"github.com/Scalingo/sclng-backend-test-v1/repository"
	"github.com/Scalingo/sclng-backend-test-v1/service"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

type IRepoWorker interface {
	manageRepo(waitgroup *sync.WaitGroup)
}

type GithubRepoWorker struct {
	logger         logrus.FieldLogger
	repo           *github.Repository
	repoRepository repository.IRepoRepository
	githubService  service.IGithubService
}

func newRepoWorker(db *sql.DB, logger logrus.FieldLogger, repo *github.Repository, githubService service.IGithubService, repoRepository repository.IRepoRepository) IRepoWorker {
	return &GithubRepoWorker{
		logger: logger, repo: repo, repoRepository: repoRepository, githubService: githubService,
	}
}

func (worker *GithubRepoWorker) manageRepo(waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	// if repository already saved do nothing
	isSaved, err := worker.repoRepository.IsRepositorySaved(worker.repo.GetID())
	if isSaved || err != nil {
		if err != nil {
			worker.logger.Error("error when trying to check if repo is already saved ", err)
		}
		return
	}
	worker.repo, err = worker.githubService.GetFullRepos(worker.repo)
	if err != nil {
		worker.logger.Error("error when getting full repo info ", err)
		return
	}
	var languages []*model.Language

	languages, err = worker.githubService.FetchLanguagesByRepo(worker.repo)
	if err != nil {
		worker.logger.Error("error while getting languages", err)
		return
	}
	localRepo := factory.CreateRepository(worker.repo, languages, worker.logger)
	// persist repo in db
	err = worker.repoRepository.SaveRepository(&localRepo)
}
