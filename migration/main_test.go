package migration

import (
	"os"
	"testing"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/util"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

func TestMigrate(t *testing.T) {
	log := logger.Default()
	cfg, err := util.NewConfigTest()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(1)
	}

	// connect DB
	db := util.ConnectPGDB(cfg, log)

	// migrate database
	err = Migrate(db, log, cfg.EraseDbWhenMigrate)
	require.NoError(t, err)
}
