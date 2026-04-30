package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CardRepository struct {
	db *sqlx.DB
}

func NewCardRepository(db *sqlx.DB) *CardRepository {
	return &CardRepository{db: db}
}

// SaveCard persists a verified card record.
// Raw card numbers are never passed here — only masked values from the service layer.
func (r *CardRepository) SaveCard(ctx context.Context, record *domain.CardRecord) error {
	query := `
		INSERT INTO billing_cards (
			id,
			user_id,
			cardholder_name,
			masked_card,
			last_four,
			brand,
			expiration_month,
			expiration_year,
			billing_address,
			billing_city,
			billing_state,
			billing_zip,
			verified_at,
			created_at
		) VALUES (
			:id,
			:user_id,
			:cardholder_name,
			:masked_card,
			:last_four,
			:brand,
			:expiration_month,
			:expiration_year,
			:billing_address,
			:billing_city,
			:billing_state,
			:billing_zip,
			:verified_at,
			:created_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, map[string]any{
		"id":               record.ID,
		"user_id":          record.UserID,
		"cardholder_name":  record.CardholderName,
		"masked_card":      record.MaskedCard,
		"last_four":        record.LastFour,
		"brand":            record.Brand,
		"expiration_month": record.ExpirationMonth,
		"expiration_year":  record.ExpirationYear,
		"billing_address":  record.BillingAddress,
		"billing_city":     record.BillingCity,
		"billing_state":    record.BillingState,
		"billing_zip":      record.BillingZip,
		"verified_at":      record.VerifiedAt,
		"created_at":       record.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("cardRepository.SaveCard: %w", err)
	}

	return nil
}

// GetCardByUserID returns the most recently verified card for a user.
func (r *CardRepository) GetCardByUserID(ctx context.Context, userID string) (*domain.CardRecord, error) {
	query := `
		SELECT
			id, user_id, cardholder_name, masked_card, last_four, brand,
			expiration_month, expiration_year, billing_address, billing_city,
			billing_state, billing_zip, verified_at, created_at
		FROM billing_cards
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	var record domain.CardRecord
	if err := r.db.GetContext(ctx, &record, query, userID); err != nil {
		return nil, fmt.Errorf("cardRepository.GetCardByUserID: %w", err)
	}

	return &record, nil
}

// DeleteCard removes a card record by ID, scoped to the owning user
// to prevent cross-user deletion.
func (r *CardRepository) DeleteCard(ctx context.Context, cardID, userID string) error {
	query := `DELETE FROM billing_cards WHERE id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, cardID, userID)
	if err != nil {
		return fmt.Errorf("cardRepository.DeleteCard: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("cardRepository.DeleteCard rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("cardRepository.DeleteCard: card %s not found for user %s", cardID, userID)
	}

	return nil
}

// -- Stubs for the methods required by domain.BillingRepository ----------
// The card repo only owns card-related persistence. The mpesa repo below
// owns STK records. Both are composed in a single service via the facade.
// These stubs prevent the interface from being unsatisfied at compile time
// if you use cardRepository directly. If you only ever use billingRepository
// (which composes both), you can remove these.

func (r *CardRepository) SaveStkTransaction(_ context.Context, _ *domain.StkRecord) error {
	return fmt.Errorf("cardRepository does not own STK records")
}

func (r *CardRepository) UpdateStkTransaction(_ context.Context, _, _, _ string) error {
	return fmt.Errorf("cardRepository does not own STK records")
}

func (r *CardRepository) GetStkTransaction(_ context.Context, _ string) (*domain.StkRecord, error) {
	return nil, fmt.Errorf("cardRepository does not own STK records")
}

// -- DDL (run once in your migrations) ------------------------------------
//
// CREATE TABLE billing_cards (
//     id               UUID         PRIMARY KEY,
//     user_id          UUID         NOT NULL,
//     cardholder_name  VARCHAR(255) NOT NULL,
//     masked_card      VARCHAR(25)  NOT NULL,
//     last_four        CHAR(4)      NOT NULL,
//     brand            VARCHAR(20)  NOT NULL,
//     expiration_month CHAR(2)      NOT NULL,
//     expiration_year  CHAR(4)      NOT NULL,
//     billing_address  TEXT         NOT NULL,
//     billing_city     VARCHAR(100) NOT NULL,
//     billing_state    VARCHAR(100) NOT NULL,
//     billing_zip      VARCHAR(20)  NOT NULL,
//     verified_at      TIMESTAMPTZ  NOT NULL,
//     created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
// );

// CREATE INDEX idx_billing_cards_user_id ON billing_cards(user_id);

// cardRepositoryForTest exposes a helper to seed test data.
// Only used in _test.go files.
type cardRepositoryForTest struct {
	*CardRepository
}

func newCardRepositoryForTest(db *sqlx.DB) *cardRepositoryForTest {
	return &cardRepositoryForTest{CardRepository: &CardRepository{db: db}}
}

func (t *cardRepositoryForTest) seedCard(ctx context.Context, userID string) (*domain.CardRecord, error) {
	record := &domain.CardRecord{
		ID:              "test-card-id",
		UserID:          userID,
		CardholderName:  "TEST USER",
		MaskedCard:      "**** **** **** 4242",
		LastFour:        "4242",
		Brand:           "visa",
		ExpirationMonth: "12",
		ExpirationYear:  "2028",
		BillingAddress:  "123 Test St",
		BillingCity:     "Nairobi",
		BillingState:    "Nairobi",
		BillingZip:      "00100",
		VerifiedAt:      time.Now().UTC(),
		CreatedAt:       time.Now().UTC(),
	}
	if err := t.SaveCard(ctx, record); err != nil {
		return nil, err
	}
	return record, nil
}