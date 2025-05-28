package repository

import (
	"context"

	"github.com/tomiristapen/banking_service/smart_budget_service/domain/model"
)

type BudgetRepository interface {
	AddExpense(ctx context.Context, expense *model.CategoryExpense) error
	GetBudgetSummary(ctx context.Context, userID string) ([]*model.BudgetSummary, error)
	GetCategoryExpenses(ctx context.Context, userID, category string) ([]*model.CategoryExpense, error)
}
