SHELL := /bin/bash

# Wallets
# amara: 0x58faAD3334A942BeD2cAd3D5aA8A688478933084
# magnus: 0x436720591bb8FD9C1f2dC3595795aA1135f61e5c
# obinna: 0x4641270AA92a075a1Bc79aa327b03fdb4c2D28d8
# Miner1: 0xc6E21c47FE24b3aDF477f7149F1caAA3a915624B
# Miner2: 0xBdC65974AA77521303FB11f0B162024257e01004
# Miner3: 0x63918C76b7D264878dFa4b43eAFFfA9763E8Bc5E

# Bookeeping transactions
# curl -il -X GET http://localhost:8080/v1/genesis/list
# curl -il -X GET http://localhost:9080/v1/node/status
# curl -il -X GET http://localhost:8080/v1/accounts/list
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
MIGRATION_NAME ?= migration
migrate: ## Create a new migration file e.g make migrate schema=<migration_name>
	@$(MIGRATE_CMD) $(MIGRATION_NAME)

# 1. create migration files - up and down files
# migrate create -ext sql -dir migrations/db -seq init_schema

#-------------------------------------------------------------TEST DB OPERATIONS - Using Docker----------------------------------------------------------------#
# migrate_up: apply all migrations locally
migrate_up: ## Apply all migrations locally
	migrate -path ./db/migrations -database "postgresql://admin:password@localhost:5433/nodedb?sslmode=disable" -verbose up

# migrate_down: rollback all migrations locally
migrate_down: ## Rollback all migrations locally
	migrate -path ./db/migrations -database "postgresql://admin:password@localhost:5433/nodedb?sslmode=disable" -verbose down

# migrate_down_last: rollback the last migration locally
migrate_down_last: ## Rollback the last migration locally
	migrate -path ./db/migrations -database "postgresql://admin:password@localhost:5433/nodedb?sslmode=disable" -verbose down 1

# dropdb: drop the database
dropdb: ## Drop the database
	docker exec -it postgres dropdb -U admin nodedb

# createdb: create the database
createdb: ## Create the database
	docker exec -it postgres createdb --username=admin --owner=admin nodedb