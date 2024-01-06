package util

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Env                            string `envconfig:"ENV" default:"dev"`
	Port                           int    `envconfig:"PORT" default:"5000"`
	DBHost                         string `envconfig:"PORT" default:"database"`
	DBPort                         int    `envconfig:"DB_PORT" default:"5432"`
	DBUser                         string `envconfig:"DB_USER" default:"postgres"`
	DBPassword                     string `envconfig:"DB_PASSWORD" default:"postgres"`
	DBName                         string `envconfig:"DB_NAME" default:"test"`
	EraseDbWhenMigrate             bool   `envconfig:"ERASE_DB_WHEN_MIGRATE" default:"true"`
	GithubToken                    string `envconfig:"GITHUB_TOKEN" default:"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"`
	WorkerWakeIntervalInSeconds    int    `envconfig:"WORKER_WAKE_INTERVAL_IN_SECONDS" default:"3600"`
	WorkerAuthorizedRequestsNumber int    `envconfig:"WORKER_AUTHORIZED_REQUESTS_NUMBER" default:"60"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to build config from env")
	}
	return &cfg, nil
}

func NewConfigTest() (*Config, error) {
	cfg, err := NewConfig()
	if cfg != nil {
		cfg.DBHost = "localhost"
	}

	return cfg, err
}
