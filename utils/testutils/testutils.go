package testutils

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/m6yf/bcwork/utils/pointer"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
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

func SetupDockerTestPool(t *testing.T) *dockertest.Pool {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		t.Fatalf("could not connect to Docker: %s", err)
	}

	return pool
}

func SetupDB(t *testing.T, pool *dockertest.Pool) *dockertest.Resource {
	const (
		user     = "root"
		password = "root"
		dbName   = "example"
	)

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
		t.Fatalf("could not start postgresql resource: %s", err)
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
		t.Fatalf("could not connect to postgres: %s", err)
	}

	return pg
}

func SetupSuperTokens(t *testing.T, pool *dockertest.Pool) (*dockertest.Resource, *supertokens_module.SuperTokensClient) {
	st, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "registry.supertokens.io/supertokens/supertokens-postgresql",
		Tag:        "9.2.3",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		t.Fatalf("could not start supertokens resource: %s", err)
	}

	port := st.GetPort("3567/tcp")
	url := "http://localhost:" + port
	basePath := "/auth"
	antiCsrf := "NONE"

	if err := pool.Retry(func() error {
		err := supertokens.Init(supertokens.TypeInput{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: url,
				APIKey:        "",
			},
			AppInfo: supertokens.AppInfo{
				AppName:         "OMS-Test",
				APIDomain:       url,
				APIBasePath:     pointer.String(basePath),
				WebsiteDomain:   url,
				WebsiteBasePath: pointer.String(basePath),
			},
			RecipeList: []supertokens.Recipe{
				dashboard.Init(&dashboardmodels.TypeInput{
					ApiKey: "",
				}),
				session.Init(&sessmodels.TypeInput{
					AntiCsrf: &antiCsrf,
				}),
			},
		})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Fatalf("could not init to supertokens: %s", err)
	}

	client := supertokens_module.NewTestSuperTokensClient(url)

	return st, client
}
