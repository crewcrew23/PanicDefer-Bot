package mqconsumer

import (
	"fmt"
	"log/slog"
	"service-healthz-checker/internal/errs"

	"github.com/streadway/amqp"
)

type MQConsumer struct {
	conn      *amqp.Connection
	chanel    *amqp.Channel
	topicName string
	log       *slog.Logger
}

func New(topicName, addr string, log *slog.Logger) *MQConsumer {

	conn, err := amqp.Dial(addr)
	errs.FailOnError(err, "failed connect to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
	}
	errs.FailOnError(err, "Failed to open a channel")

	return &MQConsumer{conn: conn, chanel: ch, topicName: topicName, log: log}
}

func (mqc *MQConsumer) Consume(consumerID string) (<-chan amqp.Delivery, error) {
	msgs, err := mqc.chanel.Consume(
		mqc.topicName,
		consumerID,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		mqc.log.Debug("failed to register consumer", slog.String("ERR", err.Error()))
		return nil, fmt.Errorf("failed to register consumer: %w", err)
	}
	return msgs, nil
}

func (mqc *MQConsumer) Close() {
	if mqc.chanel != nil {
		if err := mqc.chanel.Close(); err != nil {
			mqc.log.Error("Failed to close chanel", slog.String("ERR", err.Error()))
		}
	}
	if mqc.conn != nil {
		if err := mqc.conn.Close(); err != nil {
			mqc.log.Error("Failed to close connection", slog.String("ERR", err.Error()))
		}
	}
}
