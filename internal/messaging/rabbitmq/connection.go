package rabbitmq

import (
	"time"

	"github.com/furdarius/rabbitroutine"
	"go.uber.org/zap"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewConnector(name string, logger *zap.Logger) *rabbitroutine.Connector {
	connectorLogger := logger.With(zap.String("connector", name))
	connector := rabbitroutine.NewConnector(
		rabbitroutine.Config{
			ReconnectAttempts: 30,
			Wait:              time.Second,
		},
	)
	connector.AddAMQPNotifiedListener(func(notified rabbitroutine.AMQPNotified) {
		connectorLogger.Warn("AMQP error encountered", zap.Error(notified.Error))
	})
	connector.AddDialedListener(func(_ rabbitroutine.Dialed) {
		connectorLogger.Info("Connection with RabbitMQ established")
	})
	connector.AddRetriedListener(func(retried rabbitroutine.Retried) {
		connectorLogger.Warn(
			"Failed to establish RabbitMQ connection. Performing retries.",
			zap.Uint("attempt", retried.ReconnectAttempt),
			zap.Error(retried.Error),
		)
	})

	return connector
}

func declareExchange(channel *amqp.Channel, exchangeName string) error {
	exchangeType := "topic"
	durable := true
	autoDelete := false
	internal := false
	noWait := false
	err := channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		durable,
		autoDelete,
		internal,
		noWait,
		nil,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to declare exchange %s", exchangeName)
	}
	return nil
}
