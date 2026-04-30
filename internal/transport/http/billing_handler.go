package http

import (
	"net/http"
	"time"

	"github.com/Qarani-m/billing-service/internal/application"
	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/gin-gonic/gin"
)

// BillingHandler wires the HTTP layer to the billing service.
// Inject your BillingService interface here when you have one.
type BillingHandler struct {
	billingService domain.BillingService
	cardService    application.CardService
	mpesaService   application.MpesaService
}

func NewBillingHandler(billingService domain.BillingService, cardService *application.CardService, mpesaService *application.MpesaService) *BillingHandler {
	return &BillingHandler{
		billingService: billingService,
		cardService:    *cardService,
		mpesaService:   *mpesaService,
	}
}


// AddCard handles POST /api/v1/billing/card
func (h *BillingHandler) AddCard(c *gin.Context) {
	var req domain.AddCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: pass req to h.billingService.AddCard(c.Request.Context(), req)
	// Stub response — replace with real service call
	lastFour := req.CardNumber[len(req.CardNumber)-4:]
	resp := domain.AddCardResponse{
		MaskedCard: "**** **** **** " + lastFour,
		LastFour:   lastFour,
		Brand:      detectBrand(req.CardNumber),
		VerifiedAt: time.Now().UTC().Format(time.RFC3339),
		Message:    "Card verified successfully",
	}

	SendSuccessResponse(c, http.StatusOK, "Card verified successfully", resp)
}


// GetAllBillings handles GET /api/v1/billing/billings
func (h *BillingHandler) GetAllBillings(c *gin.Context) {
	billings, err := h.billingService.GetAllBillings(c.Request.Context())
	if err != nil {
		SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendSuccessResponse(c, http.StatusOK, "Billings retrieved successfully", billings)
}


// GetBillingByService handles GET /api/v1/billing/services/:serviceName
func (h *BillingHandler) GetBillingByService(c *gin.Context) {
	serviceName := c.Param("serviceName")
	if serviceName == "" {
		SendErrorResponse(c, http.StatusBadRequest, "serviceName is required")
		return
	}

	billings, err := h.billingService.GetBillingByService(c.Request.Context(), serviceName)
	if err != nil {
		SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendSuccessResponse(c, http.StatusOK, "Billings for service retrieved successfully", billings)
}


// InitiateStkPush handles POST /api/v1/billing/mpesa/stk-push
func (h *BillingHandler) InitiateStkPush(c *gin.Context) {
	var req domain.StkPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: pass req to h.billingService.InitiateStkPush(c.Request.Context(), req)
	// Stub response — replace with real Daraja API integration
	resp := domain.StkPushResponse{
		CheckoutRequestID: "ws_CO_stub_" + time.Now().Format("20060102150405"),
		MerchantRequestID: "stub-merchant-req-id",
		Phone:             req.Phone,
		Message:           "STK push sent successfully",
	}

	SendSuccessResponse(c, http.StatusOK, "STK push initiated successfully", resp)
}


// detectBrand returns a card brand string based on the leading digits.
func detectBrand(number string) string {
	if len(number) == 0 {
		return "unknown"
	}
	switch number[0] {
	case '4':
		return "visa"
	case '5':
		return "mastercard"
	case '3':
		return "amex"
	default:
		return "unknown"
	}
}