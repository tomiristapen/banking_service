package mq

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

type PaymentEvent struct {
	PaymentID string  `json:"payment_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Service   string  `json:"service"`
	Category  string  `json:"category"`
	Status    string  `json:"status"`
}

type MQPublisher struct {
	nc    *nats.Conn
	topic string
}

func NewMQPublisher() *MQPublisher {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	return &MQPublisher{nc: nc, topic: "payment.completed"}
}

func (p *MQPublisher) PublishPaymentEvent(event *PaymentEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.nc.Publish(p.topic, data)
}

func (p *MQPublisher) Close() {
	p.nc.Close()
}
