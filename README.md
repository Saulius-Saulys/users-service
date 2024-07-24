# Users service

* [General info](#general-info)
* [Database migrations](#database-migrations)
* [Linter](#linter)
* [Swagger docs](#swaggerdocs)
* [Mockery](#mockery)


## General info

* To read API specification copy content from `docs/user-service-api_swagger.yaml` and paste it to [Swagger Editor](https://editor.swagger.io/) or any similar tool that can read OpenAPI specification.
* To run the application in docker container can be use `make start_docker_app` command that starts postgresql and rabbitmq in docker containers, builds application docker image and runs the application inside docker container.
* To run the application locally can be use `make start_local_app` command that starts postgresql and rabbitmq in docker containers and runs the service.


## Database migrations

To run database migrations need to install tool `go-bindata` that generates go files used to install SQL schemas.

Information on how to install `go-bindata` can be found here - https://github.com/kevinburke/go-bindata

To generate migration files and execute database migrations run following commands:

```
    make package_migrations
    make migrate
```

## Linter

To install linter locally run:

```
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.1

golangci-lint --version
```

To run linter locally run:

```
make lint
```

For more information about `golangci-lint` visit docs: https://golangci-lint.run/

## Swaggerdocs

To generate swagger documentation run

```bash
make swaggerdocs
```

In order to run Swaggerdocs locally, swaggo needs to be installed.

To install swaggo run

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

For more information about `swaggo` visit: <https://github.com/swaggo/swag>

## Mockery

To generate mock interfaces, you need to install 'mockery'

``` bash
brew install mockery
```

After that run
``` bash
mockery;
```
For more information about mockery visit: https://github.com/vektra/mockery
