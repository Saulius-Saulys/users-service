lint: ## Run linter
	golangci-lint run

swaggerdocs: ## Run swaggerdocs to generate swagger documentation.
	swag init -ot "yaml" -g cmd/server/main.go --instanceName user-service-api

di: ## Generates the go files performing dependency injection
	cd cmd/server && wire

up: ## Run the dependencies
	docker-compose -f build/docker-compose.yml up -d

package_migrations: ## Generate go files used to install SQL schemas.
	go-bindata -pkg migrations -ignore bindata -nometadata -prefix internal/database/postgresql/migrations/ -o ./internal/database/postgresql/migrations/bindata.go ./internal/database/postgresql/migrations/

migrate: ## Run migrations
	docker run --rm --network host -v "${PWD}/internal/database/postgresql:/flyway/sql" flyway/flyway -url=jdbc:postgresql://localhost:5435/users-service -user=test_user -baselineOnMigrate="true" migrate -password=dev123 -schemas=users -locations=filesystem:sql/.

run: ## Run the application
	ENVIRONMENT=local RABBITMQ_USER=test_user RABBITMQ_PASSWORD=dev123 POSTGRESQL_DB_USER=test_user POSTGRESQL_DB_PASSWORD=dev123 go run cmd/server/main.go cmd/server/wire_gen.go

