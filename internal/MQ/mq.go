package mq

import (
	"log/slog"
	mqconsumer "service-healthz-checker/internal/MQ/consumer"
	mqproducer "service-healthz-checker/internal/MQ/producer"
)

type MQ struct {
	producer *mqproducer.MQProducer
	consumer *mqconsumer.MQConsumer
}

func New(topicName string, addr string, log *slog.Logger) *MQ {
	mqp := mqproducer.New(topicName, addr, log)
	mqc := mqconsumer.New(topicName, addr, log)

	return &MQ{producer: mqp, consumer: mqc}
}
