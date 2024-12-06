# Environment variables
POSTGRES_COMPOSE_FILE ?= ci/compose/test-postgres.yaml
PYTHON_COMPOSE_FILE ?= ci/compose/test-python.yaml
SERVER_COMPOSE_FILE ?= ci/compose/test.yaml
VERSION ?= 0.0.1
SERVER_IMAGE ?= test-go-server
PYTHON_IMAGE ?= test-python

# Build variables
GOARCH := amd64
GOOS := linux
PRODUCT := test
REPO_PATH := templates/go-server
DOCKERNET := test_net
MIGRATIONS_FOLDER := "templates/go-server/assets/migrations"

# Flags
BUILDFLAGS := -mod=vendor
LDFLAGS := "-X main.Version=$(VERSION)"

# Exported Go variables
export GO111MODULE := on
export GOEXPERIMENT := boringcrypto

# Default target
.DEFAULT_GOAL := help

################################################################################
# Build
################################################################################

.PHONY: build
build: ## Build the Go application binary
	cd $(REPO_PATH) && GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 \
	go build -v -ldflags=$(LDFLAGS) $(BUILDFLAGS) -o ../$(PRODUCT) ./cmd/$(PRODUCT)

.PHONY: build-debug
build-debug: ## Build the application with debug flags
	cd $(REPO_PATH) && GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 \
	go build -race -v -ldflags=$(LDFLAGS) $(BUILDFLAGS) -o ../$(PRODUCT) ./cmd/$(PRODUCT)

.PHONY: vendor
vendor: ## Create vendor directory
	cd $(REPO_PATH) && go mod vendor

.PHONY: clean
clean: ## Clean build artifacts
	cd $(REPO_PATH) && go clean -cache -testcache -modcache
	rm -f $(PRODUCT)
	rm -rf $(REPO_PATH)/vendor

################################################################################
# Docker
################################################################################


.PHONY: test-up
test-up: ## Start the database service
	docker-compose -f $(SERVER_COMPOSE_FILE) up --build -d --force-recreate

.PHONY: docker-build
docker-build: ## Build the Docker image
	docker build --build-arg VERSION=$(VERSION) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg GOOS=$(GOOS) \
		-t $(SERVER_IMAGE):$(VERSION) \
		-f ci/docker/Dockerfile .

.PHONY: docker-run
docker-run: ## Run the Docker container
	docker run -d --name $(PRODUCT) \
		--network $(DOCKERNET) \
		-p 8081:8081 -p 8080:8080 -p 8079:8079 \
		$(SERVER_IMAGE):$(VERSION)

.PHONY: docker-stop
docker-stop: ## Stop the Docker container
	docker stop $(PRODUCT) || true && docker rm $(PRODUCT) || true

################################################################################
# Database
################################################################################

.PHONY: db-up
db-up: ## Start the database service
	docker-compose -f $(POSTGRES_COMPOSE_FILE) up -d

.PHONY: db-down
db-down: ## Stop the database service
	docker-compose -f $(POSTGRES_COMPOSE_FILE) down

.PHONY: db-migrate
db-migrate: ## Run database migrations
	cd $(REPO_PATH) && go run ./cmd/migrate


################################################################################
# LINT
################################################################################

LINTER_VERSION = v1.54.2  # Replace with the desired version of golangci-lint

.PHONY: lint-deps
lint-deps:
	go install golang.org/x/lint/golint@latest

.PHONY: lint
lint: lint-deps
	golint ./go-server/...

.PHONY: lint
lint:
	cd go-server && docker run -t --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:$(LINTER_VERSION) golangci-lint run ./... -v --timeout=4m

################################################################################
# Help
################################################################################

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

################################################################################
# Python Service
################################################################################

.PHONY: python-build
python-build: ## Build the Python application Docker image
	docker build -f ci/docker/Dockerfile.python -t $(PYTHON_IMAGE):$(VERSION) .
	@echo "Python application build complete."

.PHONY: python-up
python-up: ## Start Python project container
	docker-compose -f $(PYTHON_COMPOSE_FILE) up -d --build

.PHONY: python-down
python-down: ## Stop Python project container
	docker-compose -f $(PYTHON_COMPOSE_FILE) down

.PHONY: python-logs
python-logs: ## Show logs for Python project container
	docker-compose -f $(PYTHON_COMPOSE_FILE) logs -f

.PHONY: python-restart
python-restart: python-down python-up ## Restart Python project container

################################################################################
# Help
################################################################################

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
