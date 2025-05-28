package usecase

import (
	"context"

	"github.com/tomiristapen/banking_service/smart_budget_service/domain/model"
	"github.com/tomiristapen/banking_service/smart_budget_service/domain/repository"
)

type BudgetUsecase struct {
	repo repository.BudgetRepository
}

func NewBudgetUsecase(repo repository.BudgetRepository) *BudgetUsecase {
	return &BudgetUsecase{repo: repo}
}

func (u *BudgetUsecase) GetBudgetSummary(ctx context.Context, userID string) ([]*model.BudgetSummary, error) {
	return u.repo.GetBudgetSummary(ctx, userID)
}

func (u *BudgetUsecase) GetCategoryExpenses(ctx context.Context, userID, category string) ([]*model.CategoryExpense, error) {
	return u.repo.GetCategoryExpenses(ctx, userID, category)
}
