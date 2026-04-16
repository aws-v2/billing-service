package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/gorilla/mux"
)

type Handler struct {
	service domain.InvoiceService
}

func NewHandler(service domain.InvoiceService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/billing/invoices", h.CreateInvoice).Methods(http.MethodPost)
}

type CreateInvoiceRequest struct {
	TenantID string    `json:"tenant_id"`
	Amount   float64   `json:"amount"`
	DueDate  time.Time `json:"due_date"`
}

func (h *Handler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	invoice, err := h.service.CreateInvoice(r.Context(), req.TenantID, req.Amount, req.DueDate)
	if err != nil {
		http.Error(w, "failed to create invoice: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)
}
