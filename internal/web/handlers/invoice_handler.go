package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/dto"
	"github.com/devfullcycle/imersao22/go-gateway/internal/service"
	"github.com/go-chi/chi/v5"
)

const (
	apiKeyHeader = "X-API-Key"
	contentType  = "Content-Type"
	jsonType     = "application/json"
)

type InvoiceHandler struct {
	invoiceService *service.InvoiceService
}

func NewInvoiceHandler(invoiceService *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{invoiceService: invoiceService}
}

func (h *InvoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get(apiKeyHeader)

	var input *dto.CreateInvoiceRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input.ApiKey = apiKey

	invoice, err := h.invoiceService.Create(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)
}

func (h *InvoiceHandler) GetById(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get(apiKeyHeader)

	invoiceId := chi.URLParam(r, "id")
	if invoiceId == "" {
		http.Error(w, "Param id is required", http.StatusBadRequest)
		return
	}

	invoice, err := h.invoiceService.GetById(invoiceId, apiKey)
	if err != nil {
		switch err {
		case domain.ErrInvoiceNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case domain.ErrInvoiceNotBelongs:
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		case domain.ErrUnauthorized:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set(contentType, jsonType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoice)
}

func (h *InvoiceHandler) ListByAccount(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get(apiKeyHeader)

	invoices, err := h.invoiceService.ListByAccountApiKey(apiKey)
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set(contentType, jsonType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoices)
}
