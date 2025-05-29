package grpc

import (
	"context"
	"transaction_service/domain/model"
	transactionpb "transaction_service/proto"
	"transaction_service/usecase"
)

type TransactionHandler struct {
	transactionpb.UnimplementedTransactionServiceServer
	uc *usecase.TransactionUsecase
}

func NewTransactionHandler(uc *usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{uc: uc}
}

func (h *TransactionHandler) CreateTransfer(ctx context.Context, req *transactionpb.CreateTransferRequest) (*transactionpb.CreateTransferResponse, error) {
	tx := &model.Transaction{
		FromUserID: req.FromUserId,
		ToUserID:   req.ToUserId,
		Amount:     req.Amount,
		Currency:   req.Currency,
		Status:     "pending",
	}
	err := h.uc.CreateTransfer(tx)
	if err != nil {
		return nil, err
	}
	return &transactionpb.CreateTransferResponse{
		TransferId: tx.ID,
		Status:     tx.Status,
	}, nil
}

func (h *TransactionHandler) GetTransferHistory(ctx context.Context, req *transactionpb.GetTransferHistoryRequest) (*transactionpb.GetTransferHistoryResponse, error) {
	txs, err := h.uc.GetTransferHistory(req.UserId)
	if err != nil {
		return nil, err
	}
	var protoTxs []*transactionpb.Transaction
	for _, tx := range txs {
		protoTxs = append(protoTxs, &transactionpb.Transaction{
			Id:         tx.ID,
			FromUserId: tx.FromUserID,
			ToUserId:   tx.ToUserID,
			Amount:     tx.Amount,
			Currency:   tx.Currency,
			Status:     tx.Status,
			CreatedAt:  tx.CreatedAt,
		})
	}
	return &transactionpb.GetTransferHistoryResponse{Transactions: protoTxs}, nil
}

func (h *TransactionHandler) GetTransferStatus(ctx context.Context, req *transactionpb.GetTransferStatusRequest) (*transactionpb.GetTransferStatusResponse, error) {
	status, err := h.uc.GetTransferStatus(req.TransferId)
	if err != nil {
		return nil, err
	}
	return &transactionpb.GetTransferStatusResponse{Status: status}, nil
}
