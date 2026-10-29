.PHONY: help proto proto-lint proto-breaking lint fmt vet test build up down logs ci

SVC ?=

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*##' $(MAKEFILE_LIST) | sort | sed 's/:.*##/ /'

proto: ## generate protobuf Go code
	@echo "Generating protobufs"
	@buf generate

proto-lint: ## lint protobuf definitions
	@buf lint

proto-breaking: ## check protobuf breaking changes
	@buf breaking --against '.git#branch=main' || true

lint: ## run golangci-lint
	@golangci-lint run ./...

fmt: ## format Go files
	@go fmt ./...

vet: ## run go vet
	@go vet ./...

build: ## build all workspace modules
	@go build ./...

test: ## run Go tests
	@go test -race -shuffle=on ./...

ci: proto proto-lint proto-breaking fmt vet lint test ## CI verification pipeline

up: ## start Docker Compose stack
	@docker compose -f deploy/compose/docker-compose.yml up -d

down: ## stop Docker Compose stack
	@docker compose -f deploy/compose/docker-compose.yml down

logs: ## follow compose logs
	@docker compose -f deploy/compose/docker-compose.yml logs -f

cover: ## show coverage
	@go test -cover ./...

fmt-write: ## format with gofumpt
	@gofumpt -w .

help: ## show all targets
	@echo "Use 'make help' to list commands"
