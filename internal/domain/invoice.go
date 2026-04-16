package domain

import (
	"context"
	"time"
)

type InvoiceStatus string

const (
	InvoiceStatusPending InvoiceStatus = "Pending"
	InvoiceStatusPaid    InvoiceStatus = "Paid"
)

type Invoice struct {
	ID        string        `json:"id" db:"id"`
	TenantID  string        `json:"tenant_id" db:"tenant_id"`
	Amount    float64       `json:"amount" db:"amount"`
	Status    InvoiceStatus `json:"status" db:"status"`
	DueDate   time.Time     `json:"due_date" db:"due_date"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *Invoice) error
	GetByID(ctx context.Context, id string) (*Invoice, error)
	ListByTenant(ctx context.Context, tenantID string) ([]Invoice, error)
}

type InvoiceService interface {
	CreateInvoice(ctx context.Context, tenantID string, amount float64, dueDate time.Time) (*Invoice, error)
}

type EventPublisher interface {
	PublishInvoiceCreated(ctx context.Context, invoice *Invoice) error
}
