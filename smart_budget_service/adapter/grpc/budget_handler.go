package grpc

import (
	"context"

	smartbudgetpb "github.com/tomiristapen/banking_service/smart_budget_service/proto"
	"github.com/tomiristapen/banking_service/smart_budget_service/usecase"
)

type BudgetHandler struct {
	smartbudgetpb.UnimplementedSmartBudgetServiceServer
	uc *usecase.BudgetUsecase
}

func NewBudgetHandler(uc *usecase.BudgetUsecase) *BudgetHandler {
	return &BudgetHandler{uc: uc}
}

func (h *BudgetHandler) GetBudgetSummary(ctx context.Context, req *smartbudgetpb.BudgetSummaryRequest) (*smartbudgetpb.BudgetSummaryResponse, error) {
	summary, err := h.uc.GetBudgetSummary(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	resp := &smartbudgetpb.BudgetSummaryResponse{}
	for _, cat := range summary {
		resp.Categories = append(resp.Categories, &smartbudgetpb.CategorySummary{
			Category: cat.Category,
			Total:    cat.Total,
		})
	}
	return resp, nil
}

func (h *BudgetHandler) GetCategoryExpenses(ctx context.Context, req *smartbudgetpb.CategoryExpensesRequest) (*smartbudgetpb.CategoryExpensesResponse, error) {
	expenses, err := h.uc.GetCategoryExpenses(ctx, req.UserId, req.Category)
	if err != nil {
		return nil, err
	}
	resp := &smartbudgetpb.CategoryExpensesResponse{Category: req.Category}
	var total float64
	for _, exp := range expenses {
		total += exp.Amount
		resp.Payments = append(resp.Payments, &smartbudgetpb.PaymentInfo{
			PaymentId: exp.PaymentID,
			Amount:    exp.Amount,
			Service:   exp.Service,
			Status:    exp.Status,
			CreatedAt: exp.CreatedAt,
		})
	}
	resp.Total = total
	return resp, nil
}
