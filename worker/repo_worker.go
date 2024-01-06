package worker

import (
	"database/sql"

	"github.com/Scalingo/sclng-backend-test-v1/factory"
	"github.com/Scalingo/sclng-backend-test-v1/repository"
	"github.com/Scalingo/sclng-backend-test-v1/service"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

type IRepoWorker interface {
	manageRepo()
}

type GithubRepoWorker struct {
	logger         logrus.FieldLogger
	repo           *github.Repository
	repoRepository repository.IRepoRepository
	githubService  service.IGithubService
}

func newRepoWorker(db *sql.DB, logger logrus.FieldLogger, repo *github.Repository, githubService service.IGithubService) IRepoWorker {
	return &GithubRepoWorker{
		logger: logger, repo: repo, repoRepository: repository.NewRepoRepository(db, logger), githubService: githubService,
	}
}

func (worker *GithubRepoWorker) manageRepo() {
	// if repository already saved do nothing
	isSaved, err := worker.repoRepository.IsRepositorySaved(worker.repo.GetID())
	if isSaved || err != nil {
		worker.logger.Error(err)
		return
	}
	// else get all languages
	languages, err := worker.githubService.RetrieveLanguagesByRepo(worker.repo)
	if err != nil {
		worker.logger.Error("error while getting languages", err)
		return
	}
	localRepo := factory.CreateRepository(worker.repo, languages, worker.logger)
	// persist repo in db
	err = worker.repoRepository.SaveRepository(&localRepo)
}
