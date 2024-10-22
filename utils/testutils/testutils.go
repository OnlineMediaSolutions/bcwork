package testutils

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/m6yf/bcwork/utils/pointer"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
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

func SetupSuperTokens(t *testing.T, pool *dockertest.Pool) (*dockertest.Resource, supertokens_module.TokenManagementSystem) {
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

	baseURL := "http://localhost"
	basePort := "8000"
	basePath := "/auth"
	supertokenURL := baseURL + ":" + st.GetPort("3567/tcp")
	antiCsrf := "NONE"

	err = supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: supertokenURL,
			APIKey:        "",
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "OMS-Test",
			APIDomain:       baseURL + ":" + basePort,
			APIBasePath:     pointer.String(basePath),
			WebsiteDomain:   baseURL + ":" + basePort,
			WebsiteBasePath: pointer.String(basePath),
		},
		RecipeList: []supertokens.Recipe{
			thirdpartyemailpassword.Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &antiCsrf,
			}),
		},
	})
	if err != nil {
		t.Fatalf("could not init to supertokens: %s", err)
	}

	if err := pool.Retry(func() error {
		req, err := http.NewRequest(http.MethodGet, supertokenURL+"/hello", nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if string(data) != "Hello\n" {
			return fmt.Errorf("not expected response: %v", string(data))
		}

		return nil
	}); err != nil {
		t.Fatalf("could not healthcheck supertokens: %s", err)
	}

	client := supertokens_module.NewTestSuperTokensClient(baseURL + ":" + basePort + basePath)

	return st, client
}
