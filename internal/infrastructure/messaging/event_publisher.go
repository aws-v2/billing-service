package messaging

import (
	"context"
	"time"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/Qarani-m/billing-service/pkg/messaging"
	"github.com/google/uuid"
)

type NATSEventPublisher struct {
	publisher *messaging.NATSPublisher
}

func NewNATSEventPublisher(publisher *messaging.NATSPublisher) *NATSEventPublisher {
	return &NATSEventPublisher{publisher: publisher}
}

type InvoiceCreatedEvent struct {
	CorrelationID string          `json:"correlation_id"`
	InvoiceID     string          `json:"invoice_id"`
	TenantID      string          `json:"tenant_id"`
	Amount        float64         `json:"amount"`
	Timestamp     time.Time       `json:"timestamp"`
}

func (p *NATSEventPublisher) PublishInvoiceCreated(ctx context.Context, invoice *domain.Invoice) error {
	subject := p.publisher.BuildSubject("invoice", "created")
	
	event := InvoiceCreatedEvent{
		CorrelationID: uuid.New().String(),
		InvoiceID:     invoice.ID,
		TenantID:      invoice.TenantID,
		Amount:        invoice.Amount,
		Timestamp:     time.Now(),
	}

	return p.publisher.Publish(subject, event)
}
