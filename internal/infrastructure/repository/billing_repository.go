package repository

import (
	"context"
	"fmt"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresBillingRepository struct {
	db *sqlx.DB
}

func NewPostgresBillingRepository(db *sqlx.DB) *PostgresBillingRepository {
	return &PostgresBillingRepository{db: db}
}

func (r *PostgresBillingRepository) GetAll(ctx context.Context) ([]domain.Billing, error) {
	var billings []domain.Billing
	query := `SELECT id, service_name, amount, currency, billing_date, status FROM billings ORDER BY billing_date DESC`
	err := r.db.SelectContext(ctx, &billings, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all billings: %w", err)
	}
	return billings, nil
}

func (r *PostgresBillingRepository) GetByService(ctx context.Context, serviceName string) ([]domain.Billing, error) {
	var billings []domain.Billing
	query := `SELECT id, service_name, amount, currency, billing_date, status FROM billings WHERE service_name = $1 ORDER BY billing_date DESC`
	err := r.db.SelectContext(ctx, &billings, query, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch billings for service %s: %w", serviceName, err)
	}
	return billings, nil
}
