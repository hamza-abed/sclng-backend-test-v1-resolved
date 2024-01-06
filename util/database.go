package util

import (
	"database/sql"

	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func ConnectPGDB(cfg *Config, logger logrus.FieldLogger) *sql.DB {
	/******** infrastructure.DB INFO ******/
	var psqlInfo string
	sslMode := "require"

	/******** infrastructure.DB INFO ******/
	if cfg.Env == "dev" {
		sslMode = "disable"
	}
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, sslMode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
