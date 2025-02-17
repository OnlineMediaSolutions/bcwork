# Before using following commands, change the values and add these enviroment variables:
# GOOSE_DRIVER="{goose_driver}" (for PostgreSQL use "postgres")
# GOOSE_DBSTRING="user={user} password={password} dbname={dbname} host={host} port={port} sslmode=disable"
# PGPASSWORD="{password}"

# apply migrations
migrate_up:
	goose -dir ./migrations up

# undo last migration
migrate_down:
	goose -dir ./migrations down

# reset migrations
migrate_reset:
	goose -dir ./migrations reset

# get current status which migrations were applied and which not
migrate_status:
	goose -dir ./migrations status

# create migration
migrate_new:
	goose -dir ./migrations create $(name) sql

# generate swagger docs
update_swagger:
	swag init -g cmd/api.go -o api/rest/docs --parseDependency github.com/volatiletech/null/v8

# update models according all changes in db (postgres)
update_models:
	PGPASSWORD="postgres" sqlboiler psql

# run api in local enviroment
run_api_local:
	go run main.go api --dbenv local --stenv local

# clean golang cache
clean_cache:
	go clean -cache

# run all tests except tests in package models
test: clean_cache
	go test $(shell go list ./... | grep -v /models)

# run linter
lint:
	golangci-lint run ./...