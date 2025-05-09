SHELL := /bin/bash


# Load .env file into Makefile context
include .env
export


# Wallets
# amara: 0x2B787621DD7A65270D0D995DD304787812211680
# magnus: 0x2fF11fB8f1aB5705a1E60456B5AF5B9182F735E0
# obinna: 0x742071730D6B1AB9223BA9DD9ef375C0C6f81F43
# Miner1: 0xDB38589dDe07C17E6499b40064C368C749D9ea8C
# Miner2: 0x6AEBdFFc872f572737a18B093B7002A12260cAaa
# Miner3: 0x0738cB2CAbcb49fD7d422d630309641020A86673

# Bookeeping transactions
# curl -il -X GET http://localhost:8080/v1/genesis/list
# curl -il -X GET http://localhost:8080/v1/accounts/list
# curl -il -X GET http://localhost:9080/v1/node/status
# curl -il -X GET http://localhost:8080/v1/tx/uncommitted/list
# curl -il -X GET http://localhost:8080/v1/blocks/list
# curl -il -X GET http://localhost:9080/v1/node/block/list/1/latest

NODE_BINARY=nodeApp
NODE_IMAGE := biostech/node-service:1.0.0

private_key:
	go run -mod=vendor ./cmd/utilities/scripts/private_key.go

build_push_node:
	cd ./cmd/app/handlers && docker build --no-cache -f Dockerfile -t $(NODE_IMAGE) . && docker push $(NODE_IMAGE)

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

down_kill:
	kill -INT $(shell ps | grep "main -race" | grep -v grep | sed -n 1,1p | cut -c1-5)

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_node 
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	@ docker compose up --build
	@echo "Docker images built and started!"

# build_node: builds the node binary as a linux executable
build_node:
	@echo "Building node binary..."
	@ env GOOS=linux CGO_ENABLED=0 go build  -o ${NODE_BINARY} ./cmd/app/handlers
	@echo "Done!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	@ docker compose down
	@echo "Done!"

# migrate: create a new migration file e.g make migrate schema=<migration_name>
MIGRATE_CMD = migrate create -ext sql -dir ./db/migrations -seq
MIGRATION_NAME ?= _
migrate: ## Create a new migration file e.g make migrate schema=<migration_name>
	@$(MIGRATE_CMD) $(MIGRATION_NAME)


#-------------------------------------------------------------TEST DB OPERATIONS - Using Docker----------------------------------------------------------------#
# migrate_up: ## Apply all migrations locally
migrate_up: ## Apply all migrations locally
	migrate -path ./db/migrations \
		-database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		-verbose up

# migrate_down: rollback all migrations locally
migrate_down: ## Rollback all migrations locally
	migrate -path ./db/migrations \
		-database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		-verbose down

# migrate_down_last: rollback the last migration locally
migrate_down_last: ## Rollback the last migration locally
	migrate -path ./db/migrations \
		-database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		-verbose down 1

# dropdb: drop the database
dropdb: ## Drop the database
	docker exec -it postgres dropdb -U admin nodedb

# createdb: create the database
createdb: ## Create the database
	docker exec -it postgres createdb --username=admin --owner=admin nodedb