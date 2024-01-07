package api

import (
	"os"
	"testing"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/migration"
	"github.com/Scalingo/sclng-backend-test-v1/model"
	"github.com/Scalingo/sclng-backend-test-v1/repository"
	"github.com/Scalingo/sclng-backend-test-v1/util"
	"github.com/stretchr/testify/require"
)

var server *Server

func newTestServer(t *testing.T) *Server {
	if server != nil {
		t.Log("server exists")
		return server
	}
	log := logger.Default()
	log.Info("Initializing app")
	cfg, err := util.NewConfigTest()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(1)
	}

	// connect DB
	db := util.ConnectPGDB(cfg, log)

	// migrate database
	migration.Migrate(db, log, cfg.EraseDbWhenMigrate)

	// store some datas in DB
	repo := model.Repository{ID: 55, FullName: "repotest", Name: "test", CreatedAt: time.Now(),
		Owner:     &model.Owner{ID: 55, Name: "Jhon Doe"},
		Licence:   &model.Licence{Key: "MIT", Name: "MIT Licence"},
		Languages: []*model.Language{{Name: "Ruby", Bytes: 6579}, {Name: "HTML", Bytes: 11106}}}

	repoRository := repository.NewRepoRepository(db, log)
	err = repoRository.SaveRepository(&repo)

	require.NoError(t, err)
	// start API server
	server, err := NewServer(log, db)
	if err != nil {
		log.Error("error while creating server", err)
	}

	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
