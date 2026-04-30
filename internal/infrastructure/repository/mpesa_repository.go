package repository

import (
	"context"
	"fmt"
	"time"
	"github.com/Qarani-m/billing-service/internal/domain"

	"github.com/jmoiron/sqlx"

)

type MpesaRepository struct {
	db *sqlx.DB
}

func NewMpesaRepository(db *sqlx.DB) *MpesaRepository {
	return &MpesaRepository{db: db}
}

// SaveStkTransaction persists a newly initiated STK push record with status "pending".
func (r *MpesaRepository) SaveStkTransaction(ctx context.Context, record *domain.StkRecord) error {
	query := `
		INSERT INTO billing_mpesa_transactions (
			id,
			user_id,
			phone,
			account_name,
			checkout_request_id,
			merchant_request_id,
			status,
			mpesa_receipt_number,
			amount,
			created_at,
			updated_at
		) VALUES (
			:id,
			:user_id,
			:phone,
			:account_name,
			:checkout_request_id,
			:merchant_request_id,
			:status,
			:mpesa_receipt_number,
			:amount,
			:created_at,
			:updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, map[string]any{
		"id":                   record.ID,
		"user_id":              record.UserID,
		"phone":                record.Phone,
		"account_name":         record.AccountName,
		"checkout_request_id":  record.CheckoutRequestID,
		"merchant_request_id":  record.MerchantRequestID,
		"status":               record.Status,
		"mpesa_receipt_number": record.MpesaReceiptNumber,
		"amount":               record.Amount,
		"created_at":           record.CreatedAt,
		"updated_at":           record.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("mpesaRepository.SaveStkTransaction: %w", err)
	}

	return nil
}

// UpdateStkTransaction sets the final status and receipt number on a transaction
// once Daraja confirms the outcome.
func (r *MpesaRepository) UpdateStkTransaction(
	ctx context.Context,
	checkoutRequestID string,
	status string,
	receiptNumber string,
) error {
	query := `
		UPDATE billing_mpesa_transactions
		SET
			status               = $1,
			mpesa_receipt_number = $2,
			updated_at           = $3
		WHERE checkout_request_id = $4`

	result, err := r.db.ExecContext(ctx, query, status, receiptNumber, time.Now().UTC(), checkoutRequestID)
	if err != nil {
		return fmt.Errorf("mpesaRepository.UpdateStkTransaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("mpesaRepository.UpdateStkTransaction rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("mpesaRepository.UpdateStkTransaction: no record found for checkout_request_id %s", checkoutRequestID)
	}

	return nil
}

// GetStkTransaction fetches a single STK record by its Daraja CheckoutRequestID.
func (r *MpesaRepository) GetStkTransaction(ctx context.Context, checkoutRequestID string) (*domain.StkRecord, error) {
	query := `
		SELECT
			id, user_id, phone, account_name,
			checkout_request_id, merchant_request_id,
			status, mpesa_receipt_number, amount,
			created_at, updated_at
		FROM billing_mpesa_transactions
		WHERE checkout_request_id = $1
		LIMIT 1`

	var record domain.StkRecord
	if err := r.db.GetContext(ctx, &record, query, checkoutRequestID); err != nil {
		return nil, fmt.Errorf("mpesaRepository.GetStkTransaction: %w", err)
	}

	return &record, nil
}

// GetTransactionsByUserID returns all STK records for a user, newest first.
func (r *MpesaRepository) GetTransactionsByUserID(ctx context.Context, userID string) ([]domain.StkRecord, error) {
	query := `
		SELECT
			id, user_id, phone, account_name,
			checkout_request_id, merchant_request_id,
			status, mpesa_receipt_number, amount,
			created_at, updated_at
		FROM billing_mpesa_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC`

	var records []domain.StkRecord
	if err := r.db.SelectContext(ctx, &records, query, userID); err != nil {
		return nil, fmt.Errorf("mpesaRepository.GetTransactionsByUserID: %w", err)
	}

	return records, nil
}

// GetPendingTransactions returns all STK records still in "pending" status.
// Useful for a background reconciliation job that re-queries Daraja for
// transactions the callback never delivered.
func (r *MpesaRepository) GetPendingTransactions(ctx context.Context) ([]domain.StkRecord, error) {
	query := `
		SELECT
			id, user_id, phone, account_name,
			checkout_request_id, merchant_request_id,
			status, mpesa_receipt_number, amount,
			created_at, updated_at
		FROM billing_mpesa_transactions
		WHERE status = 'pending'
		  AND created_at > NOW() - INTERVAL '1 hour'
		ORDER BY created_at ASC`

	var records []domain.StkRecord
	if err := r.db.SelectContext(ctx, &records, query); err != nil {
		return nil, fmt.Errorf("mpesaRepository.GetPendingTransactions: %w", err)
	}

	return records, nil
}

// -- Stubs for the methods required by domain.BillingRepository ----------
// The mpesa repo only owns STK-related persistence.

func (r *MpesaRepository) SaveCard(_ context.Context, _ *domain.CardRecord) error {
	return fmt.Errorf("mpesaRepository does not own card records")
}

// -- DDL (run once in your migrations) ------------------------------------
//
// CREATE TABLE billing_mpesa_transactions (
//     id                   UUID          PRIMARY KEY,
//     user_id              UUID          NOT NULL,
//     phone                VARCHAR(15)   NOT NULL,
//     account_name         VARCHAR(255)  NOT NULL,
//     checkout_request_id  VARCHAR(100)  NOT NULL UNIQUE,
//     merchant_request_id  VARCHAR(100)  NOT NULL,
//     status               VARCHAR(20)   NOT NULL DEFAULT 'pending',
//     mpesa_receipt_number VARCHAR(50)   NOT NULL DEFAULT '',
//     amount               NUMERIC(10,2) NOT NULL,
//     created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
//     updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW()
// );

// CREATE INDEX idx_mpesa_tx_user_id            ON billing_mpesa_transactions(user_id);
// CREATE INDEX idx_mpesa_tx_checkout_req_id    ON billing_mpesa_transactions(checkout_request_id);
// CREATE INDEX idx_mpesa_tx_status_created_at  ON billing_mpesa_transactions(status, created_at);