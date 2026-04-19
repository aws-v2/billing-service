package http

import (
	"net/http"
	"time"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service domain.InvoiceService
}

func NewHandler(service domain.InvoiceService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1/billing")
	{
		v1.POST("/invoices", h.CreateInvoice)
	}
}

type CreateInvoiceRequest struct {
	TenantID string    `json:"tenant_id"`
	Amount   float64   `json:"amount"`
	DueDate  time.Time `json:"due_date"`
}

func (h *Handler) CreateInvoice(c *gin.Context) {
	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	invoice, err := h.service.CreateInvoice(c.Request.Context(), req.TenantID, req.Amount, req.DueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invoice: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": invoice})
}