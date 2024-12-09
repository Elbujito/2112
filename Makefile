# Environment Variables
VERSION_GRAPHQL_API ?= 0.0.6
DOCKER_COMPOSE_FILE = deployments/compose/docker-compose.yaml
DOCKER_PROJECT_NAME = 2112_project

# Service-Specific Settings
GRAPHQL_API = graphql-api
REDIS_SERVICE = redis-service
KEYVAULT_SERVICE = keyvault-service

# Exported Go Variables
export GO111MODULE := on

# Default Target
.DEFAULT_GOAL := help

################################################################################
# Build
################################################################################

.PHONY: build
build: ## Build all services
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

.PHONY: build-graphql-api
build-graphql-api: ## Build the GraphQL API service
	docker-compose -f $(DOCKER_COMPOSE_FILE) build $(GRAPHQL_API)

.PHONY: build-redis
build-redis: ## Build the Redis service
	docker-compose -f $(DOCKER_COMPOSE_FILE) build $(REDIS_SERVICE)

.PHONY: build-keyvault
build-keyvault: ## Build the KeyVault service
	docker-compose -f $(DOCKER_COMPOSE_FILE) build $(KEYVAULT_SERVICE)

################################################################################
# Run
################################################################################

.PHONY: up
up: ## Start all services
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

.PHONY: down
down: ## Stop all services
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

.PHONY: restart
restart: down up ## Restart all services

################################################################################
# Logs
################################################################################

.PHONY: logs
logs: ## View logs for all services
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

.PHONY: logs-graphql-api
logs-graphql-api: ## View logs for the GraphQL API service
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f $(GRAPHQL_API)

.PHONY: logs-redis
logs-redis: ## View logs for the Redis service
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f $(REDIS_SERVICE)

.PHONY: logs-keyvault
logs-keyvault: ## View logs for the KeyVault service
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f $(KEYVAULT_SERVICE)

################################################################################
# Clean
################################################################################

.PHONY: clean
clean: ## Clean up all Docker resources
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --volumes --remove-orphans
	docker system prune -f
	docker volume prune -f

################################################################################
# GraphQL API
################################################################################

.PHONY: gqlgen-generate
gqlgen-generate: ## Generate GraphQL code with gqlgen
	cd graphql-api && go run github.com/99designs/gqlgen generate

.PHONY: gqlgen-init
gqlgen-init: ## Initialize the GraphQL API gqlgen project
	cd graphql-api && go run github.com/99designs/gqlgen init

.PHONY: gqlgen-clean
gqlgen-clean: ## Clean GraphQL generated files
	rm -f graphql-api/generated.go graphql-api/models_gen.go

.PHONY: gqlgen-vendor
gqlgen-vendor: ## Update dependencies for gqlgen
	cd graphql-api && go mod tidy && go mod vendor

.PHONY: gqlgen-run
gqlgen-run: ## Run the GraphQL API project
	cd graphql-api && go run .

.PHONY: gqlgen-publish
gqlgen-publish: ## Publish the GraphQL API to GitHub
	cd graphql-api && git add -A && git commit -m "Release version $(VERSION_GRAPHQL_API)" && \
	git tag -a v$(VERSION_GRAPHQL_API) -m "Version $(VERSION_GRAPHQL_API)" && git push origin main && git push origin v$(VERSION_GRAPHQL_API)
	@echo "GraphQL API published successfully with version $(VERSION_GRAPHQL_API)"

################################################################################
# Help
################################################################################

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
