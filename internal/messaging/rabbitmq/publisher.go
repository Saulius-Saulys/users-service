package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Saulius-Saulys/users-service/internal/config"
	"github.com/Saulius-Saulys/users-service/internal/environment"

	"github.com/furdarius/rabbitroutine"

	"go.uber.org/zap"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	dialogOutputExchangeName = "dialog_messages.output"
)

// Publisher owns the sole rabbit mq connection for message publishing in this service
// (It is recommended to have separate connections for publishing and consuming messages)
// Dialogs are then multiplexed on the Publisher by creating DialogPublishers
type Publisher struct {
	publisher rabbitroutine.Publisher
	logger    *zap.Logger
	Close     func()
	Completed bool
}

func NewPublisher(ctx context.Context, env environment.Env, conf *config.Config, logger *zap.Logger) (*Publisher, error) {
	err := declareOutputExchange(env, conf, logger)
	if err != nil {
		return nil, err
	}
	conn := NewConnector("publisher", logger)
	connectionString := fmt.Sprintf("amqp://%s:%s@%s/", env.RabbitMQUser, env.RabbitMQPassword, conf.RabbitMQ.Address)
	cancelCtx, cancel := context.WithCancel(ctx)
	pool := rabbitroutine.NewLightningPool(conn)
	ensurePub := rabbitroutine.NewFireForgetPublisher(pool)
	pub := rabbitroutine.NewRetryPublisher(
		ensurePub,
		rabbitroutine.PublishMaxAttemptsSetup(3),
		rabbitroutine.PublishDelaySetup(func(_ uint) time.Duration { return time.Second }),
	)
	publisher := NewPublisherImpl(pub, logger, cancel)

	go func() {
		defer func() {
			publisher.Completed = true
		}()
		dialErr := conn.Dial(cancelCtx, connectionString)
		logger.Error("Dial publisher encountered an error", zap.Error(dialErr))
	}()

	return publisher, nil
}

func NewPublisherImpl(pub rabbitroutine.Publisher, logger *zap.Logger, closeFunc func()) *Publisher {
	return &Publisher{
		publisher: pub,
		logger:    logger,
		Close:     closeFunc,
	}
}

type OutputMessage struct {
}

func (p *Publisher) PublishMessage(ctx context.Context, dialogID string, message *OutputMessage) error {
	outputJSON, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "failed to marshal input message for publication on rabbit mq exchange")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	routingKey := fmt.Sprintf("%s.%s.%s", dialogOutputExchangeName, "widget", dialogID)

	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        outputJSON,
	}

	p.logger.Debug("Publishing message", zap.String("message", string(outputJSON)), zap.String("routing_key", routingKey))
	err = p.publisher.Publish(timeoutCtx, dialogOutputExchangeName, routingKey, publishing)
	if err != nil {
		return errors.Wrapf(err, "failed to publish message to exchange %s with routing key %s", dialogOutputExchangeName, routingKey)
	}
	return nil
}

// declareOutputExchange Needs to performed separately on startup to ensure that it exists
func declareOutputExchange(env environment.Env, conf *config.Config, logger *zap.Logger) error {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s/", env.RabbitMQUser, env.RabbitMQPassword, conf.RabbitMQ.Address)
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return errors.Wrap(err, "failed to connect to rabbitmq")
	}
	defer func(conn *amqp.Connection) {
		closeErr := conn.Close()
		if closeErr != nil {
			logger.Warn("Failed to close connection after input exchange declaration")
		}
	}(conn)

	channel, err := conn.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to create rabbitmq channel")
	}
	defer func(channel *amqp.Channel) {
		closeErr := channel.Close()
		if closeErr != nil {
			logger.Warn("Failed to close channel after input exchange declaration")
		}
	}(channel)

	err = declareExchange(channel, dialogOutputExchangeName)
	if err != nil {
		return errors.Wrap(err, "failed to declare dialog input exchange")
	}
	return nil
}
