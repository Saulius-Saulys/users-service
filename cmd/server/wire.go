//go:build wireinject

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/telia-company/convhub-lmm-communication-service/internal/config"
	"github.com/telia-company/convhub-lmm-communication-service/internal/logger"
	"github.com/telia-company/convhub-lmm-communication-service/internal/network/http"
)

func inject(ctx context.Context) (userService, error) {
	wire.Build(
		newUserService,
		http.NewHTTPServer,
		http.NewRouter,
		config.NewConfig,
		gin.New,
		//environment.NewEnv,
		logger.NewSeparatedLogger,
	)

	return userService{}, nil
}
