package mqproducer

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"service-healthz-checker/internal/errs"
	requestmodel "service-healthz-checker/internal/model/requestModel"
	"time"

	"github.com/streadway/amqp"
)

type MQProducer struct {
	conn      *amqp.Connection
	chanel    *amqp.Channel
	queue     *amqp.Queue
	topicName string
	log       *slog.Logger
}

func New(topicName, addr string, log *slog.Logger) *MQProducer {
	conn, err := amqp.Dial(addr)
	errs.FailOnError(err, "failed connect to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
	}
	errs.FailOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		topicName,
		true,
		false,
		false,
		false,
		nil,
	)
	errs.FailOnError(err, "Failed to declare a queue")

	return &MQProducer{conn: conn, chanel: ch, queue: &q, topicName: topicName, log: log}
}

func (mqp *MQProducer) WriteToTopic(producereEchange string, body *requestmodel.RequestCommand) error {

	jsonBody, err := json.Marshal(body)
	if err != nil {
		mqp.log.Debug("Failed to marshal body", slog.Any("Body", body))
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = mqp.chanel.Publish(
		producereEchange,
		mqp.topicName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(jsonBody),
			Headers:      amqp.Table{"retry_count": 3},
			Timestamp:    time.Now(),
		})

	if err != nil {
		mqp.log.Debug("failed to publish a message", slog.String("ERR", err.Error()))
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	mqp.log.Debug("Sent message", slog.String("TO", mqp.topicName), slog.String("MSG", string(jsonBody)))
	return nil
}

func (mqp *MQProducer) Close() {
	if mqp.chanel != nil {
		if err := mqp.chanel.Close(); err != nil {
			mqp.log.Error("Failed to close chanel", slog.String("ERR", err.Error()))
		}
	}
	if mqp.conn != nil {
		if err := mqp.conn.Close(); err != nil {
			mqp.log.Error("Failed to close connection", slog.String("ERR", err.Error()))
		}
	}
}
