package domain

import (
	"context"
	"time"
)

// Billing represents a record of service usage and cost.
type Billing struct {
	ID          string    `json:"id" db:"id"`
	ServiceName string    `json:"serviceName" db:"service_name"`
	Amount      float64   `json:"amount" db:"amount"`
	Currency    string    `json:"currency" db:"currency"`
	BillingDate time.Time `json:"billingDate" db:"billing_date"`
	Status      string    `json:"status" db:"status"`
}

// BillingRepository defines the data access layer for billing records.
type BillingRepository interface {
	GetAll(ctx context.Context) ([]Billing, error)
	GetByService(ctx context.Context, serviceName string) ([]Billing, error)
}

// BillingService defines the business logic for billing.
type BillingService interface {
	GetAllBillings(ctx context.Context) ([]Billing, error)
	GetBillingByService(ctx context.Context, serviceName string) ([]Billing, error)
}
