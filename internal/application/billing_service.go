package application

import (
	"context"
	"log"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/Qarani-m/billing-service/internal/infrastructure/repository"
)

type BillingService struct {
	repo repository.BillingRepository
}

func NewBillingService(repo *repository.BillingRepository) *BillingService {
	return &BillingService{repo: *repo}
}

func (s *BillingService) GetAllBillings(ctx context.Context) ([]domain.Billing, error) {
	return s.repo.GetAll(ctx)
}

func (s *BillingService) GetBillingByService(ctx context.Context, serviceName string) ([]domain.Billing, error) {
	return s.repo.GetByService(ctx, serviceName)
}

func (s *BillingService) RecordMetric(metric domain.Metric) error {
	log.Printf("[billing] Recording metric: service=%s metric=%s value=%.2f %s tenant=%s",
		metric.ServiceID, metric.MetricName, metric.Value, metric.Unit, metric.TenantID)

	if err := s.repo.SaveMetric(metric); err != nil {
		log.Printf("[billing] Failed to save metric: service=%s tenant=%s error=%v",
			metric.ServiceID, metric.TenantID, err)
		return nil
	}

	log.Printf("[billing] Metric saved OK: service=%s tenant=%s", metric.ServiceID, metric.TenantID)
	return nil
}