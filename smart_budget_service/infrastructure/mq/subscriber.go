package mq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/tomiristapen/banking_service/smart_budget_service/domain/model"
	"github.com/tomiristapen/banking_service/smart_budget_service/domain/repository"
)

type PaymentEvent struct {
	PaymentID string  `json:"payment_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Service   string  `json:"service"`
	Category  string  `json:"category"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}

type MQSubscriber struct {
	nc   *nats.Conn
	repo repository.BudgetRepository
}

func NewMQSubscriber(natsURL string, repo repository.BudgetRepository) *MQSubscriber {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	return &MQSubscriber{nc: nc, repo: repo}
}

func (s *MQSubscriber) SubscribePayments() {
	s.nc.Subscribe("payment.completed", func(m *nats.Msg) {
		log.Printf("[SmartBudget] Получено событие PaymentCompleted: %s", string(m.Data))
		var event PaymentEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Printf("[SmartBudget] Failed to unmarshal payment event: %v", err)
			return
		}
		exp := &model.CategoryExpense{
			PaymentID: event.PaymentID,
			UserID:    event.UserID,
			Category:  event.Category,
			Amount:    event.Amount,
			Service:   event.Service,
			Status:    event.Status,
			CreatedAt: event.CreatedAt,
		}
		log.Printf("[SmartBudget][DEBUG] CategoryExpense type: %+v", exp)
		if err := s.repo.AddExpense(context.Background(), exp); err != nil {
			log.Printf("[SmartBudget] Failed to add expense: %v", err)
		} else {
			log.Printf("[SmartBudget] Expense added: %+v", exp)
		}
	})
	log.Println("[SmartBudget] Subscribed to PaymentCompleted events")
}
