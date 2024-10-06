package testutils

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

type AppSetup struct {
	Endpoints []EndpointSetup
}

type EndpointSetup struct {
	Method   string
	Path     string
	Handlers []fiber.Handler
}

func SetupApp(setup *AppSetup) *fiber.App {
	app := fiber.New()
	for _, endpoint := range setup.Endpoints {
		app.Add(endpoint.Method, endpoint.Path, endpoint.Handlers...)
	}

	return app
}

func SetupDB(t *testing.T) (*sqlx.DB, *dockertest.Pool, *dockertest.Resource) {
	const (
		user     = "root"
		password = "root"
		dbName   = "example"
	)

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		t.Fatalf("Could not connect to Docker: %s", err)
	}

	pg, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", dbName),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
			fmt.Sprintf("POSTGRES_USER=%s", user),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	port := pg.GetPort("5432/tcp")
	dsn := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable",
		user, password, dbName, port,
	)

	if err := pool.Retry(func() error {
		err := bcdb.InitTestDB(dsn)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Fatalf("Could not connect to postgres: %s", err)
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Fatal("failed to connect database")
	}

	return db, pool, pg
}
