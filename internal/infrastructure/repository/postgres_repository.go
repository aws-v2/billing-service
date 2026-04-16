package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresInvoiceRepository struct {
	db *sqlx.DB
}

func NewPostgresInvoiceRepository(db *sqlx.DB) *PostgresInvoiceRepository {
	return &PostgresInvoiceRepository{db: db}
}

func (r *PostgresInvoiceRepository) Create(ctx context.Context, invoice *domain.Invoice) error {
	query := `
		INSERT INTO invoices (id, tenant_id, amount, status, due_date, created_at, updated_at)
		VALUES (:id, :tenant_id, :amount, :status, :due_date, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, invoice)
	if err != nil {
		return fmt.Errorf("failed to create invoice: %w", err)
	}
	return nil
}

func (r *PostgresInvoiceRepository) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {
	var invoice domain.Invoice
	query := `SELECT * FROM invoices WHERE id = $1`
	err := r.db.GetContext(ctx, &invoice, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}
	return &invoice, nil
}

func (r *PostgresInvoiceRepository) ListByTenant(ctx context.Context, tenantID string) ([]domain.Invoice, error) {
	var invoices []domain.Invoice
	query := `SELECT * FROM invoices WHERE tenant_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &invoices, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}
	return invoices, nil
}
