package mqconsumer

import (
	"log"
	"service-healthz-checker/internal/errs"

	"github.com/streadway/amqp"
)

type MQConsumer struct {
	conn      *amqp.Connection
	chanel    *amqp.Channel
	topicName string
}

func New(topicName, addr string) *MQConsumer {

	conn, err := amqp.Dial(addr)
	errs.FailOnError(err, "failed connect to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
	}
	errs.FailOnError(err, "Failed to open a channel")

	return &MQConsumer{conn: conn, chanel: ch, topicName: topicName}
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
	errs.FailOnError(err, "Failed to register a consumer")

	return msgs, nil
}

func (mqc *MQConsumer) Close() {
	if mqc.chanel != nil {
		if err := mqc.chanel.Close(); err != nil {
			log.Printf("Failed to close channel: %v", err)
		}
	}
	if mqc.conn != nil {
		if err := mqc.conn.Close(); err != nil {
			log.Printf("Failed to close connection: %v", err)
		}
	}
}
