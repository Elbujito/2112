# Makefile for gqlgen automation
VERSION: 0.0.1

# Variables
PACKAGE_PATH = github.com/Elbujito/2112/graphql-api/model
GQLGEN = github.com/99designs/gqlgen
SCHEMA_FILE = schema.graphqls
REPO = git@github.com:Elbujito/2112.git

# Default target
all: generate

# Generate GraphQL code with gqlgen
generate:
	go run $(GQLGEN) generate

# Initialize gqlgen project
init:
	go run $(GQLGEN) init

# Clean generated files (optional)
clean:
	rm -f generated.go models_gen.go

# Ensure dependencies are updated and vendored
vendor:
	go mod tidy
	go mod vendor

# Run the project
run:
	go run .

# Update gqlgen.yml to include autobind
config:
	@echo "Updating gqlgen.yml to include autobind for $(PACKAGE_PATH)..."
	@sed -i '/autobind:/a\ \ - $(PACKAGE_PATH)' gqlgen.yml

# Help target
help:
	@echo "Usage:"
	@echo "  make generate    - Generate GraphQL code with gqlgen"
	@echo "  make init        - Initialize gqlgen project"
	@echo "  make clean       - Clean generated files"
	@echo "  make vendor      - Update dependencies and vendor"
	@echo "  make run         - Run the project"
	@echo "  make config      - Add autobind path to gqlgen.yml"
	@echo "  make publish     - Publish the package to GitHub with version $(VERSION)"

# Publish the package to GitHub
publish:
	@git add .
	@git commit -m "Release version $(VERSION)"
	@git tag -a v$(VERSION) -m "Version $(VERSION)"
	@git push origin main
	@git push origin v$(VERSION)
	@echo "Package published successfully with version $(VERSION)"
