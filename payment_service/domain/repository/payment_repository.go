package repository

import (
	"context"
	"github.com/tomiristapen/banking_service/payment_service/domain/model"
)

type PaymentRepository interface {
	Pay(ctx context.Context, payment *model.Payment) error
	ListAvailableServices(ctx context.Context) ([]string, error)
	GetStatus(ctx context.Context, id string) (*model.Payment, error)
}
