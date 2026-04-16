package application

import (
	"context"
	"time"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/google/uuid"
)

type InvoiceServiceImpl struct {
	repo           domain.InvoiceRepository
	eventPublisher domain.EventPublisher
}

func NewInvoiceService(repo domain.InvoiceRepository, eventPublisher domain.EventPublisher) *InvoiceServiceImpl {
	return &InvoiceServiceImpl{
		repo:           repo,
		eventPublisher: eventPublisher,
	}
}

func (s *InvoiceServiceImpl) CreateInvoice(ctx context.Context, tenantID string, amount float64, dueDate time.Time) (*domain.Invoice, error) {
	invoice := &domain.Invoice{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		Amount:    amount,
		Status:    domain.InvoiceStatusPending,
		DueDate:   dueDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, invoice); err != nil {
		return nil, err
	}

	// Publish event
	if err := s.eventPublisher.PublishInvoiceCreated(ctx, invoice); err != nil {
		// Log error but don't fail the request for persistence
		// In a real system, you might want to retry or use an outbox pattern
	}

	return invoice, nil
}
