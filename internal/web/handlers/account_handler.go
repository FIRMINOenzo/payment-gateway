package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/dto"
	"github.com/devfullcycle/imersao22/go-gateway/internal/service"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

func (handler *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, err := handler.accountService.CreateAccount(input)

	if err != nil && err != domain.ErrDuplicatedApiKey {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err == domain.ErrDuplicatedApiKey {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (handler *AccountHandler) GetByApiKey(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")

	if apiKey == "" {
		http.Error(w, "API Key is required", http.StatusUnauthorized)
		return
	}

	account, err := handler.accountService.FindByApiKey(apiKey)

	if err != nil && err != domain.ErrAccountNotFound {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err == domain.ErrAccountNotFound {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}
