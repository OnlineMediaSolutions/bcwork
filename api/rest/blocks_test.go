package rest

import (
	"database/sql"
	"fmt"
	"github.com/danhawkins/go-dockertest-example/database"
	"github.com/gofiber/fiber/v2"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var db *sql.DB
var app *fiber.App

func TestBlockGetAllHandler(t *testing.T) {
	setup()
	connect()
	database.CreatePerson()

	count := database.CountPeople()

	fmt.Println("count is: %", count)
	createTables()

	t.Run("Valid Request", func(t *testing.T) {
		requestBody := `{"publisher": "publisher", "domain": "domain"}`

		req := httptest.NewRequest(http.MethodPost, "/block/get", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

}

func connect() {
	log.Println("Setting up the database")

	//pgUrl := fmt.Sprintf("postgresql://postgres@127.0.0.1:%s/example", os.Getenv("POSTGRES_PORT"))
	pgUrl := fmt.Sprintf("postgresql://postgres@127.0.0.1:5432/example")
	log.Printf("Connecting to %s\n", pgUrl)
	var err error

	db, err = gorm.Open(postgres.Open(pgUrl), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
}

func setup() {

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// Uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	pg, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_DB=example",
			"POSTGRES_HOST_AUTH_METHOD=trust",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	pg.Expire(10)

	// Set this so our app can use it
	postgresPort := pg.GetPort("5432/tcp")
	os.Setenv("POSTGRES_PORT", postgresPort)

	// Wait for the HTTP endpoint to be ready
	if err := pool.Retry(func() error {
		_, connErr := gorm.Open(postgres.Open(fmt.Sprintf("postgresql://postgres@localhost:%s/example", postgresPort)), &gorm.Config{})
		if connErr != nil {
			return connErr
		}

		return nil
	}); err != nil {
		panic("Could not connect to postgres: " + err.Error())
	}

	app = fiber.New()
	app.Post("/block/get", BlockGetAllHandler)
}

func createTables() {
	_, err := db.Exec(`CREATE TABLE metadata_queue (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT,
		value TEXT,
		created_at DATETIME,
		commited_instances INTEGER
	)`)
	if err != nil {
		log.Fatalf("Error creating metadata_queue table: %v", err)
	}
}

func tearDown() {
	// Close database
	if db != nil {
		db.Close()
	}
}

//func TestCreateKeyForQuery(t *testing.T) {
//
//	// Test cases
//	testCases := []struct {
//		name          string
//		request       *BlockGetRequest
//		expectedQuery string
//	}{
//		{
//			name: "Both Publisher and Domain Provided",
//			request: &BlockGetRequest{
//				Publisher: "example_publisher",
//				Domain:    "example_domain",
//			},
//			expectedQuery: " and metadata_queue.key = 'example_publisher:example_domain'",
//		},
//		{
//			name: "Only Publisher Provided",
//			request: &BlockGetRequest{
//				Publisher: "example_publisher",
//				Domain:    "",
//			},
//			expectedQuery: " and last.key = 'example_publisher'",
//		},
//		{
//			name: "Neither Publisher nor Domain Provided",
//			request: &BlockGetRequest{
//				Publisher: "",
//				Domain:    "",
//			},
//			expectedQuery: " and 1=1 ",
//		},
//	}
//
//	// Run tests
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			actualQuery := createKeyForQuery(tc.request)
//			if actualQuery != tc.expectedQuery {
//				t.Errorf("Expected query '%s', but got '%s'", tc.expectedQuery, actualQuery)
//			}
//		})
//	}
//}
