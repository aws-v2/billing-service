package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/Qarani-m/billing-service/internal/infrastructure/repository"
	"github.com/Qarani-m/billing-service/internal/services"

	"github.com/google/uuid"
)

// verificationAmount is the KES amount charged during billing verification.
// Safaricom requires a minimum of KES 1.
const verificationAmount = 1

type MpesaService struct {
	repo   repository.MpesaRepository
	daraja *services.DarajaClient
}

// NewMpesaService creates an MpesaService wired to the Daraja client
// and a persistence repository.
func NewMpesaService(repo *repository.MpesaRepository, cfg domain.DarajaConfig) *MpesaService {
	return &MpesaService{
		repo:   *repo,
		daraja: services.NewDarajaClient(cfg),
	}
}



// InitiateStkPush dispatches an M-Pesa STK push and records
// the pending transaction in the database.
func (s *MpesaService) InitiateStkPush(
	ctx context.Context,
	userID string,
	req domain.StkPushRequest,
) (*domain.StkPushResponse, error) {

	if err := validatePhone(req.Phone); err != nil {
		return nil, err
	}

	result, err := s.daraja.SendStkPush(req.Phone, verificationAmount)
	if err != nil {
		return nil, fmt.Errorf("mpesa: initiate stk push: %w", err)
	}

	record := &domain.StkRecord{
		ID:                uuid.NewString(),
		UserID:            userID,
		Phone:             req.Phone,
		AccountName:       req.AccountName,
		CheckoutRequestID: result.CheckoutRequestID,
		MerchantRequestID: result.MerchantRequestID,
		Status:            "pending",
		Amount:            verificationAmount,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	if err := s.repo.SaveStkTransaction(ctx, record); err != nil {
		// Log but don't fail — the push is already on the user's phone.
		// The callback or confirm step will reconcile the record.
		fmt.Printf("[mpesa] warn: failed to persist stk record %s: %v\n",
			result.CheckoutRequestID, err)
	}

	return &domain.StkPushResponse{
		CheckoutRequestID: result.CheckoutRequestID,
		MerchantRequestID: result.MerchantRequestID,
		Phone:             req.Phone,
		Message:           "STK push sent successfully",
	}, nil
}

// ConfirmStkPush polls Daraja for the transaction status of a previously
// initiated push and updates the database record accordingly.
func (s *MpesaService) ConfirmStkPush(
	ctx context.Context,
	req domain.StkConfirmRequest,
) (*domain.StkConfirmResponse, error) {

	// 1. Verify we know about this checkout request
	record, err := s.repo.GetStkTransaction(ctx, req.CheckoutRequestID)
	if err != nil {
		return nil, fmt.Errorf("mpesa: unknown checkout request: %w", err)
	}
	if record.Status == "completed" {
		return &domain.StkConfirmResponse{
			Status:             "completed",
			MpesaReceiptNumber: record.MpesaReceiptNumber,
			Amount:             record.Amount,
			Message:            "Transaction already confirmed",
		}, nil
	}

	// 2. Query Daraja for current status
	queryResult, err := s.daraja.QueryStkStatus(req.CheckoutRequestID)
	if err != nil {
		return nil, fmt.Errorf("mpesa: daraja query failed: %w", err)
	}

	// 3. Map Daraja result codes to internal status
	//    ResultCode "0" = success, "1032" = cancelled, anything else = failed/pending
	switch queryResult.ResultCode {
	case 0:
		// Success — update DB record
		receiptNumber := extractReceiptNumber(queryResult)
		if err := s.repo.UpdateStkTransaction(ctx,
			req.CheckoutRequestID, "completed", receiptNumber,
		); err != nil {
			fmt.Printf("[mpesa] warn: failed to update stk record to completed: %v\n", err)
		}
		return &domain.StkConfirmResponse{
			Status:             "completed",
			MpesaReceiptNumber: receiptNumber,
			Amount:             verificationAmount,
			Message:            "Payment confirmed successfully",
		}, nil

	case 1032:
		// User cancelled
		_ = s.repo.UpdateStkTransaction(ctx, req.CheckoutRequestID, "cancelled", "")
		return &domain.StkConfirmResponse{
			Status:  "failed",
			Message: "Transaction was cancelled by the user",
		}, nil

	case 1037:
		// DS timeout — user didn't respond
		return &domain.StkConfirmResponse{
			Status:  "pending",
			Message: "No response yet — please check your phone and try again",
		}, nil

	default:
		_ = s.repo.UpdateStkTransaction(ctx, req.CheckoutRequestID, "failed", "")
		return &domain.StkConfirmResponse{
			Status:  "failed",
			Message: fmt.Sprintf("Transaction failed: %s", queryResult.ResultDesc),
		}, nil
	}
}

// validatePhone checks that the phone number is a valid Kenyan Safaricom number.
// Accepts formats: +2547XXXXXXXX, 2547XXXXXXXX, 07XXXXXXXX, 01XXXXXXXX
func validatePhone(phone string) error {
	// Strip leading +
	normalized := phone
	if len(normalized) > 0 && normalized[0] == '+' {
		normalized = normalized[1:]
	}

	// Must be 12 digits starting with 254 7X or 254 1X
	if len(normalized) == 12 && normalized[:3] == "254" {
		prefix := normalized[3:5]
		if prefix == "70" || prefix == "71" || prefix == "72" ||
			prefix == "74" || prefix == "75" || prefix == "76" ||
			prefix == "77" || prefix == "78" || prefix == "79" ||
			prefix == "10" || prefix == "11" {
			return nil
		}
		return fmt.Errorf("phone number is not a supported Safaricom line")
	}

	// Accept local format 07X / 01X (10 digits)
	if len(normalized) == 10 && (normalized[:2] == "07" || normalized[:2] == "01") {
		return nil
	}

	return fmt.Errorf("invalid phone number format — expected +2547XXXXXXXX or 07XXXXXXXX")
}

// extractReceiptNumber pulls the M-Pesa receipt from the query result description.
// In production, this comes from the Daraja callback body instead.
// Here we return the ResultDesc as a fallback — replace with proper callback parsing.
func extractReceiptNumber(result *domain.StkQueryResult) string {
	if result == nil {
		return ""
	}
	// Daraja embeds the receipt in ResultDesc on success for query endpoint.
	// A proper callback handler will have it in CallbackMetadata directly.
	return result.ResultDesc
}