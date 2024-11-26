
build:
	@go build -ldflags="-w -s -extldflags '-static' -X main.VERSION=$${version:?}" . 
	@chmod +x ./2112

start-dev:
	@docker compose --project-directory ./ -f ./ci/compose/2112-local-dev.yaml up

quick-start-postgres:
	-@mkdir -p ./ci/data/postgres
	@docker compose --project-directory ./ -f ./ci/compose/quick-start-postgres.yaml up --force-recreate --remove-orphans