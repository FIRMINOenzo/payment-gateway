package dto

import (
	"time"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

const (
	InvoiceStatusPending  = string(domain.InvoiceStatusPending)
	InvoiceStatusApproved = string(domain.InvoiceStatusApproved)
	InvoiceStatusRejected = string(domain.InvoiceStatusRejected)
)

type CreateInvoiceRequest struct {
	ApiKey              string
	Amount              float64 `json:"amount"`
	Description         string  `json:"description"`
	PaymentType         string  `json:"payment_type"`
	CardNumber          string  `json:"card_number"`
	CardHolderName      string  `json:"card_holder_name"`
	CardExpirationMonth int     `json:"card_expiration_month"`
	CardExpirationYear  int     `json:"card_expiration_year"`
	CardCVV             string  `json:"card_cvv"`
}

type InvoiceResponse struct {
	Id             string    `json:"id"`
	AccountId      string    `json:"account_id"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"`
	Description    string    `json:"description"`
	PaymentType    string    `json:"payment_type"`
	CardLastDigits string    `json:"card_last_digits"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func ToInvoice(request *CreateInvoiceRequest, accountId string) (*domain.Invoice, error) {
	card := domain.CreditCard{
		Number:          request.CardNumber,
		HolderName:      request.CardHolderName,
		ExpirationMonth: request.CardExpirationMonth,
		ExpirationYear:  request.CardExpirationYear,
		CVV:             request.CardCVV,
	}

	return domain.NewInvoice(accountId, request.Amount, request.Description, request.PaymentType, card)
}

func FromDomain(invoice *domain.Invoice) InvoiceResponse {
	return InvoiceResponse{
		Id:             invoice.Id,
		AccountId:      invoice.AccountId,
		Amount:         invoice.Amount,
		Status:         string(invoice.Status),
		Description:    invoice.Description,
		PaymentType:    invoice.PaymentType,
		CardLastDigits: invoice.CardLastDigits,
		CreatedAt:      invoice.CreatedAt,
		UpdatedAt:      invoice.UpdatedAt,
	}
}
