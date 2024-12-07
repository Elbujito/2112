# Project settings
DOCKER_COMPOSE_FILE = ci/compose/docker-compose.yaml
DOCKER_PROJECT_NAME = 2112_project

# Service-specific settings
APP_SERVICE = app-service
TERRITORY_SERVICE = territory-service
SATELLITE_POSITION_SERVICE = satellite-position-service
GRAPHQL_GATEWAY = graphql-gateway
REDIS_SERVICE = redis-service
KEYVAULT_SERVICE = keyvault-service

# Commands
.PHONY: all build build-% up down restart logs logs-% clean test test-app test-territory test-satellite

# Default: Build all services
all: build

# Build all Docker images
build:
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

# Build a specific service
build-%:
	docker-compose -f $(DOCKER_COMPOSE_FILE) build $*

# Start all services
up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Stop all services
down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Restart all services
restart: down up

# View logs for all services
logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# View logs for a specific service
logs-%:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f $*

# Run tests
test:
	@echo "Running tests for all services..."
	@$(MAKE) test-app
	@$(MAKE) test-territory
	@$(MAKE) test-satellite

# Test the App Service
test-app:
	@echo "Testing App Service..."
	docker run --rm $(APP_SERVICE) go test ./...

# Test the Territory Service
test-territory:
	@echo "Testing Territory Service..."
	docker run --rm $(TERRITORY_SERVICE) go test ./...

.PHONY: test-satellite
test-satellite: build-satellite-position-service
	@echo "Testing Satellite Position Service..."
	docker run --rm satellite-position-service python -m unittest discover -s tests

# Clean all Docker resources
clean:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --volumes --remove-orphans
	docker system prune -f
	docker volume prune -f
