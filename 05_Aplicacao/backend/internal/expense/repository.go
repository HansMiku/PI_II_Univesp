package expense

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	GetExpensesByPolitician(ctx context.Context, id string, limit, offset int, typeCode *string) ([]Expense, error)
	GetExpenseValue(ctx context.Context, idExpense string) (*ExpenseValue, error)
	GetExpenseType(ctx context.Context, house, typeCode string) (*ExpenseType, error)
	GetSupplier(ctx context.Context, house, taxID string) (*Supplier, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetExpensesByPolitician(ctx context.Context, id string, limit, offset int, typeCode *string) ([]Expense, error) {
	var expenses []Expense
	query := r.db.WithContext(ctx).
		Table("public.despesa").
		Select("public.despesa.*, public.valor_despesa.valor_documento, public.valor_despesa.valor_glosa, public.valor_despesa.valor_liquido").
		Joins("JOIN public.valor_despesa ON public.despesa.id_despesa = public.valor_despesa.id_despesa").
		Where("public.despesa.id_parlamentar = ?", id)

	if typeCode != nil {
		query = query.Where("public.despesa.codigo_tipo = ?", *typeCode)
	}

	err := query.Limit(limit).Offset(offset).Find(&expenses).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching expenses: %w", err)
	}
	return expenses, nil
}

func (r *repository) GetExpenseValue(ctx context.Context, idExpense string) (*ExpenseValue, error) {
	var ev ExpenseValue
	err := r.db.WithContext(ctx).Where("id_despesa = ?", idExpense).First(&ev).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("expense value not found: %w", err)
		}
		return nil, fmt.Errorf("error fetching expense value: %w", err)
	}
	return &ev, nil
}

func (r *repository) GetExpenseType(ctx context.Context, house, typeCode string) (*ExpenseType, error) {
	var et ExpenseType
	err := r.db.WithContext(ctx).Where("casa = ? AND codigo_tipo = ?", house, typeCode).First(&et).Error
	if err != nil {
		return nil, err
	}
	return &et, nil
}

func (r *repository) GetSupplier(ctx context.Context, house, taxID string) (*Supplier, error) {
	var s Supplier
	err := r.db.WithContext(ctx).Where("casa = ? AND cnpj_cpf = ?", house, taxID).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

