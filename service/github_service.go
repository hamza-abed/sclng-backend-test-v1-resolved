package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Scalingo/sclng-backend-test-v1/model"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

type IGithubService interface {
	CanMakeACall() bool
	RetrieveRepos() ([]*github.Repository, error)
	RetrieveLanguagesByRepo(repo *github.Repository) ([]*model.Language, error)
}

type GithubService struct {
	logger          logrus.FieldLogger
	token           string
	limiter         *rate.Limiter
	lastReservation *rate.Reservation
}

func NewGithubService(logger logrus.FieldLogger, token string) IGithubService {

	limiter := rate.NewLimiter(rate.Every(time.Hour), 60)
	if token != "" {
		limiter = rate.NewLimiter(rate.Every(time.Hour), 5000)
	}
	return &GithubService{
		logger: logger, token: token, limiter: limiter,
	}
}

func (githubService *GithubService) CanMakeACall() bool {
	if githubService.lastReservation == nil || githubService.lastReservation.Delay() <= 0 {
		return true
	}
	return false
}

func (githubService *GithubService) RetrieveRepos() ([]*github.Repository, error) {
	if !githubService.CanMakeACall() {
		return nil, fmt.Errorf("RateLimit attended !")
	}

	// option for public repositories
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 1},
		Type:        "public",
	}

	client := githubService.getClient()
	githubService.lastReservation = githubService.limiter.Reserve()
	// get a list of public repos
	repos, _, err := client.Repositories.ListByOrg(context.Background(), "github", opt)

	return repos, err
}

func (githubService *GithubService) RetrieveLanguagesByRepo(repo *github.Repository) ([]*model.Language, error) {
	githubService.logger.Info("getting repo languages owner: ", repo.GetOwner().Login, " repo name : ", repo.GetName())
	if !githubService.CanMakeACall() {
		return nil, fmt.Errorf("RateLimit attended !")
	}
	var result []*model.Language
	client := githubService.getClient()
	githubService.lastReservation = githubService.limiter.Reserve()
	//githubService.logger.Debug("getting repo languages owner: ", repo.Owner.GetLogin(), " repo name : ", repo.GetName())
	data, _, err := client.Repositories.ListLanguages(context.Background(), repo.Owner.GetLogin(), repo.GetName())

	if err != nil {
		githubService.logger.Error("error while getting repo languages", err)
		return nil, err
	}

	for key, value := range data {
		result = append(result, &model.Language{Name: key, Bytes: int64(value)})
	}

	return result, err
}

func (githubService *GithubService) getClient() *github.Client {
	client := github.NewClient(nil)

	ctx := context.Background()
	var ts oauth2.TokenSource
	if githubService.token != "" {
		ts = oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubService.token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}
	return client
}
