package mq

import (
	mqconsumer "service-healthz-checker/internal/MQ/consumer"
	mqproducer "service-healthz-checker/internal/MQ/producer"
)

type MQ struct {
	producer *mqproducer.MQProducer
	consumer *mqconsumer.MQConsumer
}

func New(topicName string, addr string) *MQ {
	mqp := mqproducer.New(topicName, addr)
	mqc := mqconsumer.New(topicName, addr)

	return &MQ{producer: mqp, consumer: mqc}
}
