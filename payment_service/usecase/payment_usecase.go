package usecase

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/tomiristapen/banking_service/payment_service/domain/model"
	"github.com/tomiristapen/banking_service/payment_service/domain/repository"
	"github.com/tomiristapen/banking_service/payment_service/infrastructure/grpcclient"
	"github.com/tomiristapen/banking_service/payment_service/infrastructure/mq"
)

type PaymentUsecase struct {
	repo       repository.PaymentRepository
	userClient *grpcclient.UserServiceClient
	mq         *mq.MQPublisher
}

func NewPaymentUsecase(repo repository.PaymentRepository, userClient *grpcclient.UserServiceClient, mq *mq.MQPublisher) *PaymentUsecase {
	return &PaymentUsecase{repo: repo, userClient: userClient, mq: mq}
}

func (uc *PaymentUsecase) PayForService(ctx context.Context, userID, service string, amount float64, category string) (*model.Payment, error) {
	if userID == "" || service == "" || amount <= 0 {
		return nil, errors.New("invalid input")
	}

	// Проверка баланса через UserService
	balance, err := uc.userClient.GetBalance(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to check balance: " + err.Error())
	}
	if balance < amount {
		return nil, errors.New("insufficient funds")
	}

	// Списание баланса
	success, msg, err := uc.userClient.DecreaseBalance(ctx, userID, amount)
	if err != nil {
		return nil, errors.New("failed to decrease balance: " + err.Error())
	}
	if !success {
		return nil, errors.New("balance not decreased: " + msg)
	}

	payment := &model.Payment{
		ID:        generateID(),
		UserID:    userID,
		Service:   service,
		Amount:    amount,
		Category:  normalize(category),
		Type:      "payment",
		Status:    "success",
		CreatedAt: time.Now(),
	}

	if err := uc.repo.Pay(ctx, payment); err != nil {
		return nil, err
	}

	// Публикация события в очередь
	if uc.mq != nil {
		event := &mq.PaymentEvent{
			PaymentID: payment.ID,
			UserID:    payment.UserID,
			Amount:    payment.Amount,
			Service:   payment.Service,
			Category:  payment.Category,
			Status:    payment.Status,
		}
		errPub := uc.mq.PublishPaymentEvent(event)
		if errPub != nil {
			log.Printf("❌ Failed to publish payment event: %v", errPub)
		} else {
			log.Printf("✅ Payment event published to NATS: PaymentID=%s, UserID=%s, Amount=%.2f, Service=%s, Category=%s, Status=%s", event.PaymentID, event.UserID, event.Amount, event.Service, event.Category, event.Status)
		}
	}

	return payment, nil
}

func (uc *PaymentUsecase) ListAvailableServices(ctx context.Context) ([]string, error) {
	return uc.repo.ListAvailableServices(ctx)
}

func (uc *PaymentUsecase) GetPaymentStatus(ctx context.Context, id string) (*model.Payment, error) {
	if id == "" {
		return nil, errors.New("missing id")
	}
	return uc.repo.GetStatus(ctx, id)
}

func normalize(cat string) string {
	return strings.ToLower(strings.TrimSpace(cat))
}

func generateID() string {
	return time.Now().Format("20060102150405")
}
