include .env
export

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker-compose
	docker-compose up --build -d postgres redis worker server master surrealdb
.PHONY: compose-up

compose-up-integration-test: ### Run docker-compose with integration test
	docker-compose up --build --abort-on-container-exit --exit-code-from integration
.PHONY: compose-up-integration-test

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

swag-v1: ### swag init
	swag init -g internal/app/config-http-server.go
.PHONY: swag-v1

# run: swag-v1 ### swag run
# 	go mod tidy && go mod download && \
# 	DISABLE_SWAGGER_HTTP_HANDLER='' CGO_ENABLED=0 go run -tags migrate ./cmd/app
# .PHONY: run

run:
	air -c .air.toml
.PHONY: run

docker-rm-volume: ### remove docker volume
	docker volume rm isling-be_pg-data isling-be_worker_data isling-be_server_data isling-be_master_data isling-be_gorse_log isling-be_surreal_data
.PHONY: docker-rm-volume

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

lint-fix: ### check and fix lint by golangci linter
	golangci-lint run --fix
.PHONY: lint-fix

linter-hadolint: ### check by hadolint linter
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

test: ### run test
	go test -v -cover -race ./internal/...
.PHONY: test

integration-test: ### run integration-test
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

mock: ### run mockery
	mockery --all -r --case snake
.PHONY: mock

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations '$(name)'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

migrate-down: ### migration down
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' down
.PHONY: migrate-down

psql: ### access database by psql
	psql -h localhost -p 5432 -d postgres -U user
.PHONY: psql
