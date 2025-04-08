package dto

import (
	"time"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

type AccountResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	ApiKey    string    `json:"api_key,omitempty"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAccountRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ToAccount(data CreateAccountRequest) *domain.Account {
	return domain.NewAccount(data.Name, data.Email)
}

func FromAccount(account *domain.Account) AccountResponse {
	return AccountResponse{
		Id:        account.Id,
		Name:      account.Name,
		Email:     account.Email,
		ApiKey:    account.ApiKey,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}
