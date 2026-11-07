package expense

import "context"

type Service interface {
	GetExpensesByPolitician(ctx context.Context, id string, house *string, limit, offset int, typeCode *string) ([]Expense, error)
	CountExpensesByPolitician(ctx context.Context, id string, house *string, typeCode *string) (int64, error)
	GetExpenseSummaryByPolitician(ctx context.Context, id string, house *string, typeCode *string, monthStart, monthEnd *int) ([]ExpenseCategorySummary, error)
	GetExpenseStatisticsByPolitician(ctx context.Context, id string, house *string, typeCode *string, monthStart, monthEnd *int) (*ExpenseStatistics, error)
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

func (s *service) GetExpensesByPolitician(ctx context.Context, id string, house *string, limit, offset int, typeCode *string) ([]Expense, error) {
	return s.repo.GetExpensesByPolitician(ctx, id, house, limit, offset, typeCode)
}

func (s *service) CountExpensesByPolitician(ctx context.Context, id string, house *string, typeCode *string) (int64, error) {
	return s.repo.CountExpensesByPolitician(ctx, id, house, typeCode)
}

func (s *service) GetExpenseSummaryByPolitician(ctx context.Context, id string, house *string, typeCode *string, monthStart, monthEnd *int) ([]ExpenseCategorySummary, error) {
	return s.repo.GetExpenseSummaryByPolitician(ctx, id, house, typeCode, monthStart, monthEnd)
}

func (s *service) GetExpenseStatisticsByPolitician(ctx context.Context, id string, house *string, typeCode *string, monthStart, monthEnd *int) (*ExpenseStatistics, error) {
	return s.repo.GetExpenseStatisticsByPolitician(ctx, id, house, typeCode, monthStart, monthEnd)
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
