package repository

import "transaction_service/domain/model"

type TransactionRepository interface {
	Create(tx *model.Transaction) error
	GetHistory(userID string) ([]*model.Transaction, error)
	GetStatus(transferID string) (string, error)
}
