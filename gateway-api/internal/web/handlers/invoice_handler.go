package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/domain"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/dto"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/service"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/web"
)

type InvoiceHandler struct {
	invoiceService *service.InvoiceService
}

func NewInvoiceHandler(invoiceService *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService: invoiceService,
	}
}

func (h *InvoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateInvoiceRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input.APIKey = r.Header.Get(web.HEADER_X_API_KEY)

	output, err := h.invoiceService.Create(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(web.HEADER_CONTENT_TYPE, web.HEADER_VALUE_APPLICATION_JSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *InvoiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	apiKey := r.Header.Get(web.HEADER_X_API_KEY)
	if apiKey == "" {
		http.Error(w, domain.ErrRequiredAPIKey.Error(), http.StatusBadRequest)
	}

	output, err := h.invoiceService.FindById(id, apiKey)
	if err != nil {
		switch err {
		case domain.ErrInvoiceNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case domain.ErrAccountNotFound:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		case domain.ErrUnauthorizedAccess:
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set(web.HEADER_CONTENT_TYPE, web.HEADER_VALUE_APPLICATION_JSON)
	json.NewEncoder(w).Encode(output)
}

func (h *InvoiceHandler) ListByAccount(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get(web.HEADER_X_API_KEY)
	if apiKey == "" {
		http.Error(w, domain.ErrRequiredAPIKey.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.invoiceService.ListByAccountAPIKey(apiKey)
	if err != nil {
		switch err {
		case domain.ErrAccountNotFound:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set(web.HEADER_CONTENT_TYPE, web.HEADER_VALUE_APPLICATION_JSON)
	json.NewEncoder(w).Encode(output)
}
