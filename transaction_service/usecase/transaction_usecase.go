package usecase

import (
	"transaction_service/domain/model"
	"transaction_service/domain/repository"
)

type NatsPublisher interface {
	PublishTransfer(tx *model.Transaction) error
}

type TransactionUsecase struct {
	repo      repository.TransactionRepository
	publisher NatsPublisher
}

func NewTransactionUsecase(repo repository.TransactionRepository, publisher NatsPublisher) *TransactionUsecase {
	return &TransactionUsecase{repo: repo, publisher: publisher}
}

func (u *TransactionUsecase) CreateTransfer(tx *model.Transaction) error {
	err := u.repo.Create(tx)
	if err != nil {
		return err
	}
	return u.publisher.PublishTransfer(tx)
}

func (u *TransactionUsecase) GetTransferHistory(userID string) ([]*model.Transaction, error) {
	return u.repo.GetHistory(userID)
}

func (u *TransactionUsecase) GetTransferStatus(transferID string) (string, error) {
	return u.repo.GetStatus(transferID)
}
