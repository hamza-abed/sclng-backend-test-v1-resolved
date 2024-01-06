package repository

import (
	"os"
	"testing"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/model"
	"github.com/Scalingo/sclng-backend-test-v1/util"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestCreateRepository(t *testing.T) {
	log := logger.Default()
	cfg, err := util.NewConfigTest()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(1)
	}

	// connect DB
	db := util.ConnectPGDB(cfg, log)

	repo := model.Repository{ID: 55, FullName: "repotest", Name: "test", CreatedAt: time.Now(),
		Owner:     &model.Owner{ID: 55, Name: "Jhon Doe"},
		Licence:   &model.Licence{Key: "MIT", Name: "MIT Licence"},
		Languages: []*model.Language{{Name: "Ruby", Bytes: 6579}, {Name: "HTML", Bytes: 11106}}}

	repoRository := NewRepoRepository(db, log)
	err = repoRository.SaveRepository(&repo)
	require.NoError(t, err)

	// check repository saved
	repoSaved, err := repoRository.IsRepositorySaved(repo.ID)
	require.NoError(t, err)
	require.Equal(t, true, repoSaved, "repo saved : ", repoSaved)

	// check owner saved
	ownerSaved, err := repoRository.isOwnerSaved(repo.Owner.ID)
	require.NoError(t, err)
	require.Equal(t, true, ownerSaved)

	for _, lang := range repo.Languages {
		// check all languages saved
		langSaved, err := repoRository.isLanguageSaved(lang.Name)
		require.NoError(t, err)
		require.Equal(t, true, langSaved, "lang saved ", langSaved)
	}

}
