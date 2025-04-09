package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/domain"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/dto"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/service"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/web"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.accountService.Create(input)
	if err != nil {
		if err == domain.ErrAccountAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(web.HEADER_CONTENT_TYPE, web.HEADER_VALUE_APPLICATION_JSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get(web.HEADER_X_API_KEY)
	if apiKey == "" {
		http.Error(w, domain.ErrRequiredAPIKey.Error(), http.StatusUnauthorized)
		return
	}

	output, err := h.accountService.FindByAPIKey(apiKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(web.HEADER_CONTENT_TYPE, web.HEADER_VALUE_APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
