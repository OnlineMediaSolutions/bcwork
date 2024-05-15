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

var db *sqlx.DB
var app *fiber.App

func TestBlockGetAllHandler(t *testing.T) {

	setup()
	connectToDB()
	createTables()

	t.Run("Valid Request", func(t *testing.T) {
		requestBody := `{"types: ["badv"], publisher": "badv:20356", "domain": "playpilot.com"}`

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
	// Define test cases
	testCases := []struct {
		name     string
		request  BlockGetRequest
		expected string
	}{
		{
			name:     "Empty request",
			request:  BlockGetRequest{},
			expected: " and 1=1 ",
		},
		{
			name: "Publisher and types provided",
			request: BlockGetRequest{
				Types:     []string{"badv", "bcat"},
				Publisher: "publisher",
				Domain:    "domain",
			},
			expected: "AND ( (metadata_queue.key = 'badv:publisher:domain') OR (metadata_queue.key = 'bcat:publisher:domain'))",
		},
		{
			name: "Publisher provided, no domain",
			request: BlockGetRequest{
				Types:     []string{"badv", "cbat"},
				Publisher: "publisher",
			},
			expected: "AND ( (metadata_queue.key = 'badv:publisher') OR (metadata_queue.key = 'cbat:publisher'))",
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		// Run the test
		t.Run(tc.name, func(t *testing.T) {
			result := createKeyForQuery(&tc.request)

			// Check if the result matches the expected value
			if result != tc.expected {
				t.Errorf("Test %s failed: expected '%s', got '%s'", tc.name, tc.expected, result)
			}
		})
	}
}
