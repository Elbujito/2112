# Environment variables
COMPOSE_FILE ?= ci/compose/postgres.yaml
WEB_COMPOSE_FILE ?= ci/compose/2112-web.yaml
PYTHON_COMPOSE_FILE ?= ci/compose/2112-python.yaml
VERSION ?= latest

# Targets
.PHONY: all build start-dev quick-start-postgres up down restart logs docker-build web-build web-up web-down web-restart web-logs

# Build the Go backend binary
build:
	@echo "Building the backend application..."
	@go build -ldflags="-w -s -extldflags '-static' -X main.VERSION=$(VERSION)" -o ./2112 .
	@chmod +x ./2112
	@echo "Backend build complete. Binary is located at ./2112."

# Start local development environment for backend
start-dev:
	@echo "Starting local development environment..."
	@docker compose --project-directory ./ -f ./ci/compose/2112-local-dev.yaml up

# Quick start for PostgreSQL
quick-start-postgres:
	@echo "Starting PostgreSQL environment..."
	@mkdir -p ./ci/data/postgres
	@docker compose --project-directory ./ -f ./ci/compose/quick-start-postgres.yaml up --force-recreate --remove-orphans
	@echo "PostgreSQL environment started."

# Bring up backend services
up:
	@echo "Starting backend services using $(COMPOSE_FILE)..."
	@docker-compose -f $(COMPOSE_FILE) up -d

# Bring down backend services
down:
	@echo "Stopping backend services using $(COMPOSE_FILE)..."
	@docker-compose -f $(COMPOSE_FILE) down

# Restart backend services
restart:
	@echo "Restarting backend services using $(COMPOSE_FILE)..."
	@docker-compose -f $(COMPOSE_FILE) down
	@docker-compose -f $(COMPOSE_FILE) up -d

# Show logs for backend services
logs:
	@echo "Showing logs for backend services defined in $(COMPOSE_FILE)..."
	@docker-compose -f $(COMPOSE_FILE) logs -f

# Build Docker images for backend services
docker-build:
	@echo "Building Docker images for backend services using $(COMPOSE_FILE)..."
	@docker-compose -f $(COMPOSE_FILE) build

# Web service management

# Build the web service
web-build:
	@echo "Building the web application..."
	docker build -f ci/docker/Dockerfile.web -t 2112-web .
	@echo "Web application build complete."


# Bring up the web service
web-up:
	@echo "Starting the web service using $(WEB_COMPOSE_FILE)..."
	@docker-compose -f $(WEB_COMPOSE_FILE) up -d

# Bring down the web service
web-down:
	@echo "Stopping the web service using $(WEB_COMPOSE_FILE)..."
	@docker-compose -f $(WEB_COMPOSE_FILE) down

# Restart the web service
web-restart:
	@echo "Restarting the web service using $(WEB_COMPOSE_FILE)..."
	@docker-compose -f $(WEB_COMPOSE_FILE) down
	@docker-compose -f $(WEB_COMPOSE_FILE) up -d

# Show logs for the web service
web-logs:
	@echo "Showing logs for the web service defined in $(WEB_COMPOSE_FILE)..."
	@docker-compose -f $(WEB_COMPOSE_FILE) logs -f

	# Build Python project Docker image
python-build:
	@echo "Building Python project Docker image..."
	docker build -f ci/docker/Dockerfile.python -t 2112-python .

# Start Python project container
python-up:
	@echo "Starting Python project container..."
	@docker-compose -f $(PYTHON_COMPOSE_FILE) up -d

# Stop Python project container
python-down:
	@echo "Stopping Python project container..."
	@docker-compose -f $(PYTHON_COMPOSE_FILE) down

# Show logs for Python project container
python-logs:
	@echo "Showing logs for Python project container..."
	@docker-compose -f $(PYTHON_COMPOSE_FILE) logs -f

# Restart Python project container
python-restart: python-down python-up
