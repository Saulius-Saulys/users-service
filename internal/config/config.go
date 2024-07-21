package config

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/config"
	"go.uber.org/zap"
)

var (
	baseConfigFile  = "./internal/config/base.yaml"
	localConfigFile = "./internal/config/local.yaml"
)

type Postgresql struct {
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int    `yaml:"conn_max_idle_time"`
	ConnMaxOpen     int    `yaml:"conn_max_open"`
	ConnMaxIdle     int    `yaml:"conn_max_idle"`
	DBAddress       string `yaml:"db_address"`
	DBPort          string `yaml:"db_port"`
	DBName          string `yaml:"db_name"`
	SchemaName      string `yaml:"schema_name"`
}

type RabbitMQ struct {
	Address               string `yaml:"address"`
	UnusedQueueTTLSeconds int    `yaml:"unused_queue_ttl_seconds"`
}

type HTTPClient struct {
	TimeoutSeconds int `yaml:"timeout_seconds"`
}

type ContentSafety struct {
	APIVersion string `yaml:"api_version"`
	Address    string `yaml:"address"`
}

type RateLimits struct {
	CompletionsServerPerSecond float64 `yaml:"completions_server_per_second"`
}

type Config struct {
	DBQueryTimeoutMS int           `yaml:"db_query_timeout_ms"`
	HTTPPort         string        `yaml:"http_port"`
	BasePath         string        `yaml:"base_path"`
	SecretMode       string        `yaml:"secret_mode"`
	RabbitMQ         RabbitMQ      `yaml:"rabbit_mq"`
	Postgresql       Postgresql    `yaml:"postgresql"`
	HTTPClient       HTTPClient    `yaml:"http_client"`
	ContentSafety    ContentSafety `yaml:"content_safety"`
	RateLimits       RateLimits    `yaml:"rate_limits"`
}

func NewConfig(logger *zap.Logger) (*Config, error) {
	appConfig, err := loadConfig(logger)
	if err != nil {
		return &Config{}, errors.Wrap(err, "could not load app config")
	}

	return &appConfig, nil
}

func loadConfig(logger *zap.Logger) (Config, error) {
	configFiles := []config.YAMLOption{config.File(baseConfigFile)}

	// If there would be more environments, we could add them here by switching on the env string
	configFiles = append(configFiles, config.File(localConfigFile))

	configProvider, err := config.NewYAML(configFiles...)
	if err != nil {
		return Config{}, errors.Wrap(err, "could not locate config files")
	}

	var appConfig Config
	err = configProvider.Get("").Populate(&appConfig)
	if err != nil {
		return Config{}, errors.Wrap(err, "could not construct app config")
	}
	configInJSON, err := json.Marshal(appConfig)
	if err != nil {
		return Config{}, errors.Wrap(err, "unable to marshal configuration")
	}
	logger.Sugar().Infow("Application configuration", "config", string(configInJSON))
	return appConfig, nil
}