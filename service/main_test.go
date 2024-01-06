package service

import (
	"os"
	"testing"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/util"
	"github.com/stretchr/testify/require"
)

var gservice IGithubService
var cfg *util.Config

func TestMain(m *testing.M) {
	log := logger.Default()
	cfg, err := util.NewConfigTest()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(1)
	}
	gservice = NewGithubService(log, cfg.GithubToken)

	os.Exit(m.Run())
}

func TestRetrieveRepos(t *testing.T) {

	repos, err := gservice.RetrieveRepos()
	require.NoError(t, err, "should not return an error", err)
	require.NotEmpty(t, repos, "should not be empty")

}

/*
func TestRetrieveLanguagesFromURL(t *testing.T) {
	url1 := "https://api.github.com/repos/github/learn.github.com/languages"
	url2 := "https://api.github.com/repos/github/media/languages"

	languages, err := gservice.RetrieveLanguagesFromURL(url1)
	require.NoError(t, err, err)
	require.NotEmpty(t, languages, "! empty")

	emptyLanguages, err := gservice.RetrieveLanguagesFromURL(url2)
	require.NoError(t, err, err)
	require.Empty(t, emptyLanguages)
}*/
