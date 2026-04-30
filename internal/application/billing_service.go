package application

import (
	"context"

	"github.com/Qarani-m/billing-service/internal/domain"
)

type BillingService struct {
	repo domain.BillingRepository
}

func NewBillingService(repo domain.BillingRepository) *BillingService {
	return &BillingService{repo: repo}
}

func (s *BillingService) GetAllBillings(ctx context.Context) ([]domain.Billing, error) {
	return s.repo.GetAll(ctx)
}

func (s *BillingService) GetBillingByService(ctx context.Context, serviceName string) ([]domain.Billing, error) {
	return s.repo.GetByService(ctx, serviceName)
}
