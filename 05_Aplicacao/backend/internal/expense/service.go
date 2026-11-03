package expense

import "context"

type Service interface {
	GetExpensesByPolitician(ctx context.Context, id string, limit, offset int, typeCode *string) ([]Expense, error)
	GetExpenseValue(ctx context.Context, idExpense string) (*ExpenseValue, error)
	GetExpenseType(ctx context.Context, house, typeCode string) (*ExpenseType, error)
	GetSupplier(ctx context.Context, house, taxID string) (*Supplier, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetExpensesByPolitician(ctx context.Context, id string, limit, offset int, typeCode *string) ([]Expense, error) {
	return s.repo.GetExpensesByPolitician(ctx, id, limit, offset, typeCode)
}

func (s *service) GetExpenseValue(ctx context.Context, idExpense string) (*ExpenseValue, error) {
	return s.repo.GetExpenseValue(ctx, idExpense)
}

func (s *service) GetExpenseType(ctx context.Context, house, typeCode string) (*ExpenseType, error) {
	return s.repo.GetExpenseType(ctx, house, typeCode)
}

func (s *service) GetSupplier(ctx context.Context, house, taxID string) (*Supplier, error) {
	return s.repo.GetSupplier(ctx, house, taxID)
}
