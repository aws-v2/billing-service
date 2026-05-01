package repository

import (
	"context"
	"fmt"
 "github.com/google/uuid"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type BillingRepository struct {
	db *sqlx.DB
}
func (r *BillingRepository) SaveMetric(metric domain.Metric) error {
	query := `
		INSERT INTO metrics (id, service_id, metric_name, value, unit, tenant_id, timestamp)
		VALUES (:id, :service_id, :metric_name, :value, :unit, :tenant_id, :timestamp)
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":          generateID(),
		"service_id":  metric.ServiceID,
		"metric_name": metric.MetricName,
		"value":       metric.Value,
		"unit":        metric.Unit,
		"tenant_id":   metric.TenantID,
		"timestamp":   metric.Timestamp,
	})
	if err != nil {
		return fmt.Errorf("failed to save metric: service=%s tenant=%s error=%w", metric.ServiceID, metric.TenantID, err)
	}
	return nil
}

func generateID() string {
	return uuid.New().String()
}
func NewPostgresBillingRepository(db *sqlx.DB) *BillingRepository {
	return &BillingRepository{db: db}
}

func (r *BillingRepository) GetAll(ctx context.Context) ([]domain.Billing, error) {
	var billings []domain.Billing
	query := `SELECT id, service_name, amount, currency, billing_date, status FROM billings ORDER BY billing_date DESC`
	err := r.db.SelectContext(ctx, &billings, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all billings: %w", err)
	}
	return billings, nil
}

func (r *BillingRepository) GetByService(ctx context.Context, serviceName string) ([]domain.Billing, error) {
	var billings []domain.Billing
	query := `SELECT id, service_name, amount, currency, billing_date, status FROM billings WHERE service_name = $1 ORDER BY billing_date DESC`
	err := r.db.SelectContext(ctx, &billings, query, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch billings for service %s: %w", serviceName, err)
	}
	return billings, nil
}
