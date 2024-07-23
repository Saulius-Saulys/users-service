lint: ## Run linter
	golangci-lint run

swaggerdocs: ## Run swaggerdocs to generate swagger documentation.
	swag init -ot "yaml" -g cmd/server/main.go --instanceName user-service-api

di: ## Generates the go files performing dependency injection
	cd cmd/server && wire

up: ## Run the dependencies
	docker-compose -f build/docker-compose.yml up -d
	sleep 5 ## wait until postgresql and rabbitmq are up

package_migrations: ## Generate go files used to install SQL schemas.
	go-bindata -pkg migrations -ignore bindata -nometadata -prefix internal/database/postgresql/migrations/ -o ./internal/database/postgresql/migrations/bindata.go ./internal/database/postgresql/migrations/

migrate: ## Run migrations
	docker run --rm --network host -v "${PWD}/internal/database/postgresql:/flyway/sql" flyway/flyway -url=jdbc:postgresql://localhost:5432/users-service -user=test -baselineOnMigrate="true" migrate -password=dev123 -schemas=users -locations=filesystem:sql/.

local_run: ## Run the application locally
	ENVIRONMENT=local RABBITMQ_HOSTNAME=localhost:5672 RABBITMQ_USER=test_user RABBITMQ_PASSWORD=dev123 POSTGRESQL_DB_HOSTNAME=localhost POSTGRESQL_DB_USER=test POSTGRESQL_DB_PASSWORD=dev123 go run cmd/server/main.go cmd/server/wire_gen.go

docker_build: ## Build the docker image
	docker build --tag 'user_service' -f build/Dockerfile .

docker_run: ## Run the docker container
	docker run -e ENVIRONMENT=local -e RABBITMQ_HOSTNAME=message-broker:5672 -e RABBITMQ_USER=test_user -e RABBITMQ_PASSWORD=dev123 -e POSTGRESQL_DB_HOSTNAME=users-db -e POSTGRESQL_DB_USER=test -e POSTGRESQL_DB_PASSWORD=dev123 --network=build_test-net -p 8080:8080 user_service

start_docker_app: up package_migrations migrate docker_build docker_run

start_local_app: up package_migrations migrate local_run
