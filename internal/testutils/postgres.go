package testutils

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/Saulius-Saulys/users-service/internal/database/postgresql/migrations"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
)

func StartDatabase() (string, *dockertest.Pool, *dockertest.Resource, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return "", nil, nil, errors.Errorf("unable to connect to docker: %v", err)
	}

	// Start database container.
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
			"POSTGRES_DB=users",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself even if the test fails
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		return "", nil, nil, errors.Errorf("unable to start postgres container: %v", err)
	}

	// Tell docker to hard kill the container in 120 seconds
	err = resource.Expire(120)
	if err != nil {
		return "", nil, nil, errors.Errorf("unable to expire resource: %v", err)
	}

	connectionString := fmt.Sprintf(
		"port=%s host=%s user=%s password=%s dbname=%s sslmode=disable",
		resource.GetPort("5432/tcp"),
		"localhost",
		"test",
		"test",
		"users",
	)

	pool.MaxWait = 80 * time.Second
	err = pool.Retry(func() (err error) {
		db, err := Connect(connectionString)
		if err != nil {
			return err
		}
		defer func() {
			cerr := db.Close()
			if err == nil {
				err = cerr
			}
		}()
		err = db.Ping()
		if err != nil {
			return errors.Wrap(err, "failed to ping database")
		}
		return nil
	})
	if err != nil {
		return "", nil, nil, errors.Errorf("unable to connect to postgres container: %v", err)
	}

	return connectionString, pool, resource, nil
}

func Connect(connectionString string) (*sql.DB, error) {
	c, err := pgx.ParseConfig(connectionString)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse connection string")
	}

	conn := stdlib.OpenDB(*c)
	return conn, nil
}

func PrepareUsersDB(connectionString string, runMigrations bool) (*sql.DB, error) {
	db, err := Connect(connectionString)

	if err != nil {
		return nil, errors.Errorf("unable to connect to db: %v", err)
	}

	err = db.Ping()

	if err != nil {
		return nil, errors.Errorf("failed to ping database: %v", err)
	}

	if runMigrations {
		err = DBSchema(db)
		if err != nil {
			return nil, errors.Errorf("failed to run migrations: %v", err)
		}
	}

	return db, nil
}

// DBSchema migrates the schema of the db to the specified version.
func DBSchema(db *sql.DB) error {
	assets := migrations.AssetNames()
	sort.Strings(assets)
	for _, assetName := range assets {
		asset, err := migrations.Asset(assetName)
		if err != nil {
			return errors.Wrap(err, "failed to read migration asset")
		}
		_, err = db.Exec(string(asset))
		if err != nil {
			return errors.Wrap(err, "failed to run migrations")
		}
	}
	return nil
}
