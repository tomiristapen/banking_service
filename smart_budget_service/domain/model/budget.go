package model

type CategoryExpense struct {
	PaymentID string
	UserID    string
	Category  string
	Amount    float64
	Service   string
	Status    string
	CreatedAt string // ISO8601
}

type BudgetSummary struct {
	UserID   string
	Category string
	Total    float64
}
