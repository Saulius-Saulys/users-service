package main

import (
	"context"
	"github.com/Saulius-Saulys/users-service/internal/logger"
	"github.com/Saulius-Saulys/users-service/internal/messaging/rabbitmq"
	"github.com/Saulius-Saulys/users-service/internal/network/http"
	"go.uber.org/zap"
	"time"
)

type userService struct {
	httpServer     *http.Server
	logger         *zap.Logger
	rabbitMQClient rabbitmq.Client
}

func newUserService(
	httpServer *http.Server,
	logger *zap.Logger,
	rabbitMQClient rabbitmq.Client,
) userService {
	return userService{
		httpServer:     httpServer,
		logger:         logger,
		rabbitMQClient: rabbitMQClient,
	}
}

// @title Users Server API
// @version 1.0.0
// @description this API for endpoints related to user.

// @url http://localhost:8088
func main() {
	cancelableCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wiredServer, err := inject(cancelableCtx)
	if err != nil {
		logger.NewSeparatedLogger().Panic("unable to construct DI", zap.Error(err))
	}
	defer wiredServer.cleanup()

	go func() {
		defer cancel()
		wiredServer.httpServer.Serve()
	}()

	wiredServer.httpServer.GracefulStop()
}

func (us *userService) cleanup() {
	us.logger.Debug("Performing server shutdown cleanup actions")
	us.rabbitMQClient.Close()
	for i := 0; i < 10; i++ {
		if us.rabbitMQClient.GetPublisherClient().Completed {
			return
		}
		time.Sleep(time.Second)
	}
	us.logger.Warn("Failed to close consumer and publisher on server shutdown")
}
