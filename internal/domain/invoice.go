package domain

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type InvoiceStatus string

const (
	InvoiceStatusPending  InvoiceStatus = "pending"
	InvoiceStatusApproved InvoiceStatus = "approved"
	InvoiceStatusRejected InvoiceStatus = "rejected"
)

type Invoice struct {
	Id             string
	AccountId      string
	Amount         float64
	Status         InvoiceStatus
	Description    string
	PaymentType    string
	CardLastDigits string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreditCard struct {
	Number          string
	CVV             string
	ExpirationMonth int
	ExpirationYear  int
	HolderName      string
}

func NewInvoice(accountId string, amount float64, description string, paymentType string, creditCard CreditCard) (*Invoice, error) {
	if accountId == "" {
		return nil, ErrAccountNotFound
	}

	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	lastDigits := creditCard.Number[len(creditCard.Number)-4:]

	return &Invoice{
		Id:             uuid.New().String(),
		AccountId:      accountId,
		Amount:         amount,
		Status:         InvoiceStatusPending,
		Description:    description,
		PaymentType:    paymentType,
		CardLastDigits: lastDigits,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (i *Invoice) Process() error {
	if i.Amount > 10000 {
		return nil
	}

	randomPercentage := rand.New(rand.NewSource(time.Now().Unix())).Float64()

	var newStatus InvoiceStatus

	if randomPercentage <= 0.7 {
		newStatus = InvoiceStatusApproved
	} else {
		newStatus = InvoiceStatusRejected
	}

	return i.UpdateStatus(newStatus)
}

func (i *Invoice) UpdateStatus(status InvoiceStatus) error {
	if i.Status == status {
		return ErrInvalidStatus
	}

	if status != InvoiceStatusApproved && status != InvoiceStatusRejected {
		return ErrInvalidStatus
	}

	i.Status = status
	i.UpdatedAt = time.Now()

	return nil
}
