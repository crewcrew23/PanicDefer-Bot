package mqproducer

import (
	"log"
	"service-healthz-checker/internal/errs"
	"time"

	"github.com/streadway/amqp"
)

type MQProducer struct {
	conn      *amqp.Connection
	chanel    *amqp.Channel
	queue     *amqp.Queue
	topicName string
}

func New(topicName, addr string) *MQProducer {
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

	return &MQProducer{conn: conn, chanel: ch, queue: &q, topicName: topicName}
}

func (mqp *MQProducer) WriteToTopic(topicName, producereEchange, body string) {
	err := mqp.chanel.Publish(
		producereEchange,
		mqp.topicName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(body),
			Headers:      amqp.Table{"retry_count": 3},
			Timestamp:    time.Now(),
		})
	errs.FailOnError(err, "Failed to publish a message")
	log.Printf("Sent: %s", body)
}

func (mqp *MQProducer) Close() {
	if mqp.chanel != nil {
		if err := mqp.chanel.Close(); err != nil {
			log.Printf("Failed to close channel: %v", err)
		}
	}
	if mqp.conn != nil {
		if err := mqp.conn.Close(); err != nil {
			log.Printf("Failed to close connection: %v", err)
		}
	}
}
