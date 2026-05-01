package domain

import "time"

// CardRecord is the database model for a verified card.
// Raw card numbers are NEVER stored — only masked/tokenized values.
type CardRecord struct {
	ID             string    `json:"id"             db:"id"`
	UserID         string    `json:"userId"         db:"user_id"`
	CardholderName string    `json:"cardholderName" db:"cardholder_name"`
	MaskedCard     string    `json:"maskedCard"     db:"masked_card"`
	LastFour       string    `json:"lastFour"       db:"last_four"`
	Brand          string    `json:"brand"          db:"brand"`
	ExpirationMonth string   `json:"expirationMonth" db:"expiration_month"`
	ExpirationYear  string   `json:"expirationYear"  db:"expiration_year"`
	BillingAddress  string   `json:"billingAddress"  db:"billing_address"`
	BillingCity     string   `json:"billingCity"     db:"billing_city"`
	BillingState    string   `json:"billingState"    db:"billing_state"`
	BillingZip      string   `json:"billingZip"      db:"billing_zip"`
	VerifiedAt      time.Time `json:"verifiedAt"    db:"verified_at"`
	CreatedAt       time.Time `json:"createdAt"     db:"created_at"`
}

// StkRecord is the database model for an M-Pesa STK push transaction.
type StkRecord struct {
	ID                 string    `json:"id"                  db:"id"`
	UserID             string    `json:"userId"              db:"user_id"`
	Phone              string    `json:"phone"               db:"phone"`
	AccountName        string    `json:"accountName"         db:"account_name"`
	CheckoutRequestID  string    `json:"checkoutRequestId"   db:"checkout_request_id"`
	MerchantRequestID  string    `json:"merchantRequestId"   db:"merchant_request_id"`
	// "pending" | "completed" | "failed" | "cancelled"
	Status             string    `json:"status"              db:"status"`
	MpesaReceiptNumber string    `json:"mpesaReceiptNumber"  db:"mpesa_receipt_number"`
	Amount             float64   `json:"amount"              db:"amount"`
	CreatedAt          time.Time `json:"createdAt"           db:"created_at"`
	UpdatedAt          time.Time `json:"updatedAt"           db:"updated_at"`
}

type StkQueryResult struct {
	ResultCode        int    `json:"ResultCode"`
	ResultDesc        string `json:"ResultDesc"`
	CallbackMetadata  struct {
		Item []struct {
			Name  string `json:"Name"`
			Value string `json:"Value"`
		}
	}
}