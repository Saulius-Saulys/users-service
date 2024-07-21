lint: ## Run linter
	golangci-lint run

swaggerdocs: ## Run swaggerdocs to generate swagger documentation.
	swag init -ot "yaml" -g cmd/server/main.go --instanceName user-service-api

di: ## Generates the go files performing dependency injection
	cd cmd/server && wire
