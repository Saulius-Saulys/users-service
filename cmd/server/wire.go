//go:build wireinject

package main

import (
	"context"
	"github.com/Saulius-Saulys/users-service/internal/config"
	"github.com/Saulius-Saulys/users-service/internal/database/postgresql"
	"github.com/Saulius-Saulys/users-service/internal/environment"
	"github.com/Saulius-Saulys/users-service/internal/logger"
	"github.com/Saulius-Saulys/users-service/internal/messaging/rabbitmq"
	"github.com/Saulius-Saulys/users-service/internal/network/http"
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller"
	"github.com/Saulius-Saulys/users-service/internal/repository"
	"github.com/Saulius-Saulys/users-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func inject(ctx context.Context) (userService, error) {
	wire.Build(
		newUserService,
		http.NewHTTPServer,
		http.NewRouter,
		config.NewConfig,
		gin.New,
		environment.NewEnv,
		logger.NewSeparatedLogger,
		repository.NewUser,
		wire.Bind(new(service.UserRepository), new(*repository.User)),
		postgresql.NewUsersDB,
		wire.Bind(new(controller.UserService), new(*service.User)),
		service.NewUser,
		controller.NewUser,
		rabbitmq.NewPublisher,
		wire.Bind(new(service.RabbitMQPublisher), new(*rabbitmq.Publisher)),
	)

	return userService{}, nil
}
