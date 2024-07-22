package rabbitmq

import (
	"context"

	"github.com/Saulius-Saulys/users-service/internal/config"
	"github.com/Saulius-Saulys/users-service/internal/environment"
	"go.uber.org/zap"
)

type Client interface {
	GetPublisherClient() *Publisher
	Close()
}

type ClientImpl struct {
	publisher *Publisher
}

func NewClientImpl(ctx context.Context, env environment.Env, conf *config.Config, logger *zap.Logger) (*ClientImpl, error) {
	publisher, err := NewPublisher(ctx, env, conf, logger)
	if err != nil {
		return nil, err
	}

	return &ClientImpl{
		publisher: publisher,
	}, nil
}

func (c *ClientImpl) GetPublisherClient() *Publisher {
	return c.publisher
}

func (c *ClientImpl) Close() {
	c.publisher.Close()
}
