package domain

// AddCardRequest is the payload sent by the frontend when a user
// submits their credit/debit card for billing verification.
type AddCardRequest struct {
	CardholderName  string `json:"cardholderName"  binding:"required"`
	CardNumber      string `json:"cardNumber"      binding:"required,len=16"`
	ExpirationMonth string `json:"expirationMonth" binding:"required,len=2"`
	ExpirationYear  string `json:"expirationYear"  binding:"required,len=4"`
	CVV             string `json:"cvv"             binding:"required,min=3,max=4"`
	Address         string `json:"address"         binding:"required"`
	City            string `json:"city"            binding:"required"`
	State           string `json:"state"           binding:"required"`
	ZipCode         string `json:"zipCode"         binding:"required"`
}

// AddCardResponse is returned after a successful card verification.
type AddCardResponse struct {
	MaskedCard string `json:"maskedCard"`
	LastFour   string `json:"lastFour"`
	Brand      string `json:"brand"`
	VerifiedAt string `json:"verifiedAt"`
	Message    string `json:"message"`
}