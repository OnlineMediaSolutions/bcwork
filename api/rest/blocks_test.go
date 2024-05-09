package rest

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
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
	"time"
)

// var db *gorm.DB
var db *sqlx.DB
var app *fiber.App

type MetadataQueue struct {
	Transaction_id     uint `gorm:"primary_key"`
	Key                string
	Version            string
	Commited_instances string
}

func TestBlockGetAllHandler(t *testing.T) {

	setup()
	connectToDB()
	createTables()

	t.Run("Valid Request", func(t *testing.T) {
		requestBody := `{"publisher": "badv:20356", "domain": "playpilot.com"}`

		req := httptest.NewRequest(http.MethodPost, "/block/get", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

}

func connectToDB() {
	dbinfo := fmt.Sprintf("host=localhost port=%s user=root password=root dbname=example sslmode=disable", os.Getenv("POSTGRES_PORT"))
	var err error
	db, err = sqlx.Connect("postgres", dbinfo)

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
			"POSTGRES_PASSWORD=root",
			"POSTGRES_USER=root",
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
	pg.Expire(1000)

	// Set this so our app can use it
	postgresPort := pg.GetPort("5432/tcp")
	os.Setenv("POSTGRES_PORT", postgresPort)

	// Wait for the HTTP endpoint to be ready
	if err := pool.Retry(func() error {
		dsn := fmt.Sprintf("host=localhost user=root password=root dbname=example port=%s sslmode=disable", postgresPort)
		_, connErr := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if connErr != nil {
			return connErr
		}

		return nil
	}); err != nil {
		panic("Could not connect to postgres: " + err.Error())
	}

	//setup Fiber
	app = fiber.New()
	app.Post("/block/get", BlockGetAllHandler)
}

func createTables() {

	log.Println("Starting creating new Table")
	tx := db.MustBegin()
	tx.MustExec("CREATE TABLE metadata_queue (transaction_id varchar(36), key varchar(256), version varchar(16),value varchar(512),commited_instances integer, created_at timestamp, updated_at timestamp)")
	tx.MustExec("INSERT INTO metadata_queue (transaction_id, key, version, value, commited_instances, created_at, updated_at) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7)", "f2b8833e-e0e4-57e0-a68b-6792e337ab4d", "badv:20223:realgm.com", nil, "[\"safesysdefender.xyz\"]", nil, time.Now(), time.Now())
	tx.MustExec("INSERT INTO metadata_queue (transaction_id, key, version, value, commited_instances, created_at, updated_at) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7)", "c53c4dd2-6f68-5b62-b613-999a5239ad36", "badv:20356:playpilot.com", nil, "[\"fraction-content.com\"]", nil, time.Now(), time.Now())
	tx.Commit()
	log.Println("Finished Creating DB")
}

func TestCreateKeyForQuery(t *testing.T) {

	// Test cases
	testCases := []struct {
		name          string
		request       *BlockGetRequest
		expectedQuery string
	}{
		{
			name: "Both Publisher and Domain Provided",
			request: &BlockGetRequest{
				Publisher: "example_publisher",
				Domain:    "example_domain",
			},
			expectedQuery: " and metadata_queue.key = 'example_publisher:example_domain'",
		},
		{
			name: "Only Publisher Provided",
			request: &BlockGetRequest{
				Publisher: "example_publisher",
				Domain:    "",
			},
			expectedQuery: " and last.key = 'example_publisher'",
		},
		{
			name: "Neither Publisher nor Domain Provided",
			request: &BlockGetRequest{
				Publisher: "",
				Domain:    "",
			},
			expectedQuery: " and 1=1 ",
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualQuery := createKeyForQuery(tc.request)
			if actualQuery != tc.expectedQuery {
				t.Errorf("Expected query '%s', but got '%s'", tc.expectedQuery, actualQuery)
			}
		})
	}
}
