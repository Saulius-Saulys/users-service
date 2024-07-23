# Users service

* [General info](#general-info)
* [Linter](#linter)
* [Swagger docs](#swaggerdocs)
* [Mockery](#mockery)


## General info

* To read API specification copy content from `docs/user-service-api_swagger.yaml` and paste it to [Swagger Editor](https://editor.swagger.io/) or any similar tool that can read OpenAPI specification. 

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
