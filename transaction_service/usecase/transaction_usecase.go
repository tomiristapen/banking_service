package usecase

import (
	"context"
	"transaction_service/domain/model"
	"transaction_service/domain/repository"
	"transaction_service/infrastructure/grpcclient"
)

type NatsPublisher interface {
	PublishTransfer(tx *model.Transaction) error
}

type TransactionUsecase struct {
	repo       repository.TransactionRepository
	publisher  NatsPublisher
	userClient *grpcclient.UserServiceClient
}

func NewTransactionUsecase(repo repository.TransactionRepository, publisher NatsPublisher, userClient *grpcclient.UserServiceClient) *TransactionUsecase {
	return &TransactionUsecase{repo: repo, publisher: publisher, userClient: userClient}
}

func (u *TransactionUsecase) CreateTransfer(tx *model.Transaction) error {
	ctx := context.Background()
	// 1. Проверить баланс отправителя
	balance, err := u.userClient.GetBalance(ctx, tx.FromUserID)
	if err != nil {
		tx.Status = "failed"
		u.repo.Create(tx)
		return err
	}
	if balance < tx.Amount {
		tx.Status = "failed"
		u.repo.Create(tx)
		return nil
	}
	// 2. Списать у отправителя
	success, _, err := u.userClient.DecreaseBalance(ctx, tx.FromUserID, tx.Amount)
	if err != nil || !success {
		tx.Status = "failed"
		u.repo.Create(tx)
		return err
	}
	// 3. Зачислить получателю (если есть метод IncreaseBalance)
	if success, _, err := u.userClient.DecreaseBalance(ctx, tx.ToUserID, -tx.Amount); err != nil || !success {
		tx.Status = "failed"
		u.repo.Create(tx)
		return err
	}
	// 4. Если всё успешно
	tx.Status = "success"
	err = u.repo.Create(tx)
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
