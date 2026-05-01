package domain

// StkPushRequest is the payload sent by the frontend when a user
// initiates an M-Pesa STK push for billing verification.
type StkPushRequest struct {
	// Full international phone number, e.g. "+254712345678"
	Phone       string `json:"phone"       binding:"required"`
	// Name on the Safaricom account
	AccountName string `json:"accountName" binding:"required"`
}

// StkPushResponse is returned after the STK push has been
// dispatched to Safaricom's Daraja API.
type StkPushResponse struct {
	// Daraja CheckoutRequestID — used to poll or confirm the transaction
	CheckoutRequestID string `json:"checkoutRequestId"`
	// Merchant request ID from Daraja
	MerchantRequestID string `json:"merchantRequestId"`
	// Human-readable status, e.g. "STK push sent successfully"
	Message           string `json:"message"`
	// Normalized phone the push was sent to, e.g. "+254712345678"
	Phone             string `json:"phone"`
}

// StkConfirmRequest is the optional body the frontend sends
// once the user has entered their PIN, so the backend can
// query Daraja for the final transaction status.
type StkConfirmRequest struct {
	CheckoutRequestID string `json:"checkoutRequestId" binding:"required"`
}

// StkConfirmResponse is returned after the backend verifies
// the transaction status with Daraja.
type StkConfirmResponse struct {
	// "completed" | "pending" | "failed"
	Status            string `json:"status"`
	// M-Pesa receipt number on success, empty otherwise
	MpesaReceiptNumber string `json:"mpesaReceiptNumber,omitempty"`
	// Amount transacted (should be 1.00 for verification)
	Amount            float64 `json:"amount,omitempty"`
	Message           string  `json:"message"`
}
type DarajaConfig struct {
	ConsumerKey       string
	ConsumerSecret    string
	ShortCode         string
	PassKey           string
	CallbackURL       string
	Environment       string
	DarajaBaseURL     string
}