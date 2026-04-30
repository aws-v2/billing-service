package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/Qarani-m/billing-service/internal/infrastructure/repository"
	"github.com/google/uuid"
)

type CardService struct {
	repo *repository.CardRepository	
}

// NewCardService creates a CardService that satisfies the card-related
// methods of domain.BillingService.
func NewCardService(repo *repository.CardRepository) *CardService {
	return &CardService{repo: repo}
}

// AddCard validates, masks, and persists a card, then returns a safe response.
// Raw card numbers are never written to any storage layer.
func (s *CardService) AddCard(ctx context.Context, userID string, req domain.AddCardRequest) (*domain.AddCardResponse, error) {
	if err := validateCard(req); err != nil {
		return nil, err
	}

	lastFour := req.CardNumber[len(req.CardNumber)-4:]
	brand := detectBrand(req.CardNumber)
	masked := fmt.Sprintf("**** **** **** %s", lastFour)
	now := time.Now().UTC()

	record := &domain.CardRecord{
		ID:              uuid.NewString(),
		UserID:          userID,
		CardholderName:  req.CardholderName,
		MaskedCard:      masked,
		LastFour:        lastFour,
		Brand:           brand,
		ExpirationMonth: req.ExpirationMonth,
		ExpirationYear:  req.ExpirationYear,
		BillingAddress:  req.Address,
		BillingCity:     req.City,
		BillingState:    req.State,
		BillingZip:      req.ZipCode,
		VerifiedAt:      now,
		CreatedAt:       now,
	}

	if err := s.repo.SaveCard(ctx, record); err != nil {
		return nil, fmt.Errorf("failed to persist card record: %w", err)
	}

	return &domain.AddCardResponse{
		MaskedCard: masked,
		LastFour:   lastFour,
		Brand:      brand,
		VerifiedAt: now.Format(time.RFC3339),
		Message:    "Card verified successfully",
	}, nil
}

// validateCard performs business-rule checks beyond what Gin's binding tags cover.
func validateCard(req domain.AddCardRequest) error {
	// Strip spaces just in case the frontend sends them
	number := stripSpaces(req.CardNumber)

	if len(number) < 13 || len(number) > 19 {
		return fmt.Errorf("invalid card number length")
	}
	if !luhnCheck(number) {
		return fmt.Errorf("card number failed Luhn validation")
	}

	month := req.ExpirationMonth
	year := req.ExpirationYear
	now := time.Now()
	if month < "01" || month > "12" {
		return fmt.Errorf("invalid expiration month")
	}
	// Accept YY or YYYY
	fullYear := year
	if len(year) == 2 {
		fullYear = fmt.Sprintf("20%s", year)
	}
	if fullYear < fmt.Sprintf("%d", now.Year()) {
		return fmt.Errorf("card has expired")
	}

	return nil
}

// luhnCheck runs the Luhn algorithm to sanity-check the card number.
func luhnCheck(number string) bool {
	sum := 0
	nDigits := len(number)
	parity := nDigits % 2
	for i := 0; i < nDigits; i++ {
		digit := int(number[i] - '0')
		if digit < 0 || digit > 9 {
			return false
		}
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}

// detectBrand returns a card brand string based on the leading digits.
func detectBrand(number string) string {
	if len(number) == 0 {
		return "unknown"
	}
	switch {
	case number[0] == '4':
		return "visa"
	case number[0] == '5' && number[1] >= '1' && number[1] <= '5':
		return "mastercard"
	case number[:2] == "34" || number[:2] == "37":
		return "amex"
	case number[:4] == "6011":
		return "discover"
	default:
		return "unknown"
	}
}

func stripSpaces(s string) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' {
			result = append(result, s[i])
		}
	}
	return string(result)
}