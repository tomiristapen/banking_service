package mq

import (
	"encoding/json"
	"transaction_service/domain/model"

	"github.com/nats-io/nats.go"
)

type NatsPublisher struct {
	nc    *nats.Conn
	topic string
}

func NewNatsPublisher(nc *nats.Conn, topic string) *NatsPublisher {
	return &NatsPublisher{nc: nc, topic: topic}
}

func (p *NatsPublisher) PublishTransfer(tx *model.Transaction) error {
	data, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	return p.nc.Publish(p.topic, data)
}
