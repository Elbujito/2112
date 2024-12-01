# Environment variables
POSTGRES_COMPOSE_FILE ?= ci/compose/postgres.yaml
WEB_COMPOSE_FILE ?= ci/compose/2112-web.yaml
PYTHON_COMPOSE_FILE ?= ci/compose/2112-python.yaml
SERVER_COMPOSE_FILE ?= ci/compose/2112-server.yaml
VERSION ?= latest

# Targets
.PHONY: all build start-dev quick-start-postgres up down restart logs docker-build web-build web-up web-down web-restart web-logs python-build python-up python-down python-logs python-restart server-build server-up server-down server-logs server-restart clean

# Build the server Docker image
server-build:
	@echo "Building the server application..."
	docker build -f ci/docker/Dockerfile -t server:latest .
	@echo "Server application build complete. Docker image 'server:latest' created."

# Start the server service
server-up:
	@echo "Starting the server service using $(SERVER_COMPOSE_FILE)..."
	@docker-compose -f $(SERVER_COMPOSE_FILE) up -d

# Stop the server service
server-down:
	@echo "Stopping the server service using $(SERVER_COMPOSE_FILE)..."
	@docker-compose -f $(SERVER_COMPOSE_FILE) down

# Restart the server service
server-restart:
	@echo "Restarting the server service using $(SERVER_COMPOSE_FILE)..."
	@docker-compose -f $(SERVER_COMPOSE_FILE) down
	@docker-compose -f $(SERVER_COMPOSE_FILE) up -d

# Show logs for the server service
server-logs:
	@echo "Showing logs for the server service defined in $(SERVER_COMPOSE_FILE)..."
	@docker-compose -f $(SERVER_COMPOSE_FILE) logs -f


# Bring up backend services
db-up:
	@echo "Starting backend services using $(POSTGRES_COMPOSE_FILE)..."
	@docker-compose -f $(POSTGRES_COMPOSE_FILE) up -d

# Bring down backend services
db-down:
	@echo "Stopping backend services using $(POSTGRES_COMPOSE_FILE)..."
	@docker-compose -f $(POSTGRES_COMPOSE_FILE) down

# Restart backend services
db-restart:
	@echo "Restarting backend services using $(POSTGRES_COMPOSE_FILE)..."
	@docker-compose -f $(POSTGRES_COMPOSE_FILE) down
	@docker-compose -f $(POSTGRES_COMPOSE_FILE) up -d

# Show logs for backend services
db-logs:
	@echo "Showing logs for backend services defined in $(POSTGRES_COMPOSE_FILE)..."
	@docker-compose -f $(POSTGRES_COMPOSE_FILE) logs -f

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

# Build the web service
python-build:
	@echo "Building the Python application Docker image..."
	docker build \
		-f ci/docker/Dockerfile.python \
		-t 2112-python .
	@echo "Python application build complete. Docker image '2112-python' created."


# Start Python project container
python-up:
	@echo "Starting Python project container in $(PYTHON_COMPOSE_FILE)..."
	@docker-compose -f $(PYTHON_COMPOSE_FILE) up -d --build

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

# Clean PyInstaller build artifacts
clean:
	@echo "Cleaning up build artifacts..."
	rm -rf build dist *.spec
	@echo "Clean complete."
