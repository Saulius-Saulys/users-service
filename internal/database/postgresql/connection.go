package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/Saulius-Saulys/users-service/internal/config"
	"github.com/Saulius-Saulys/users-service/internal/environment"
	"github.com/cenkalti/backoff/v3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"
	"time"
)

type GORMInstance struct {
	GORM *gorm.DB
}

type NewDBParams struct {
	Port            string
	Address         string
	User            string
	Password        string
	Name            string
	ConnMaxOpen     int
	ConnMaxIdle     int
	ConnMaxLifetime int
	ConnMaxIdleTime int
	SchemaName      string
}

func NewUsersDB(config *config.Config, env environment.Env, logger *zap.Logger) (*gorm.DB, error) {
	gormInstance, err := newGORMInstance(
		NewDBParams{
			Port:            config.Postgresql.DBPort,
			Address:         config.Postgresql.DBAddress,
			User:            env.PostgreSQLDBUser,
			Password:        env.PostgreSQLDBPassword,
			Name:            config.Postgresql.DBName,
			ConnMaxOpen:     config.Postgresql.ConnMaxOpen,
			ConnMaxIdle:     config.Postgresql.ConnMaxIdle,
			ConnMaxLifetime: config.Postgresql.ConnMaxLifetime,
			ConnMaxIdleTime: config.Postgresql.ConnMaxIdleTime,
			SchemaName:      config.Postgresql.SchemaName,
		},
		logger,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create GORM instance")
	}

	return gormInstance.GORM, nil
}

func newGORMInstance(params NewDBParams, logger *zap.Logger) (*GORMInstance, error) {
	connStr := getConnStrFromParams(params)

	gormConn, err := createGormInstance(params, connStr, logger)
	if err != nil {
		return &GORMInstance{}, err
	}

	logger.Info("Connected to Postgresql and GORM instance created")
	return &GORMInstance{
		GORM: gormConn,
	}, nil
}

func getConnStrFromParams(params NewDBParams) string {
	connStr := fmt.Sprintf(
		"port=%s host=%s user=%s password=%s dbname=%s sslmode=disable",
		params.Port,
		params.Address,
		params.User,
		params.Password,
		params.Name,
	)
	return connStr
}

func createGormInstance(params NewDBParams, connStr string, innerLogConfig *zap.Logger) (*gorm.DB, error) {
	logger := zapgorm2.New(innerLogConfig)
	logger.SlowThreshold = 300 * time.Millisecond

	gormConn, err := gorm.Open(
		postgres.Open(connStr),
		&gorm.Config{
			Logger: logger,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   params.SchemaName + ".",
				SingularTable: false,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("opening connection to DB: %w", err)
	}
	dbConn, err := gormConn.DB()
	if err != nil {
		return nil, fmt.Errorf("unable to get DB from GORM: %w", err)
	}
	configConn(dbConn, params)

	err = retry(time.Minute, func() error {
		if pingErr := dbConn.Ping(); pingErr != nil {
			return fmt.Errorf("failed to retry connection: %w", pingErr)
		}
		return nil
	})

	return gormConn, err
}

func configConn(conn *sql.DB, params NewDBParams) {
	conn.SetMaxOpenConns(params.ConnMaxOpen)
	conn.SetMaxIdleConns(params.ConnMaxIdle)
	conn.SetConnMaxLifetime(time.Duration(params.ConnMaxLifetime) * time.Second)
	conn.SetConnMaxIdleTime(time.Duration(params.ConnMaxIdleTime) * time.Second)
}

func retry(maxWait time.Duration, op func() error) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 5
	bo.MaxElapsedTime = maxWait
	if retryErr := backoff.Retry(op, bo); retryErr != nil {
		return fmt.Errorf("failed to make retry with backoff library: %w", retryErr)
	}

	return nil
}
