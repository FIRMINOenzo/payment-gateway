package domain

import "errors"

var (
	ErrAccountNotFound  = errors.New("account not found")
	ErrDuplicatedApiKey = errors.New("duplicated api key")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrXApiKeyRequired  = errors.New("x-api-key is required")

	ErrInvoiceNotFound   = errors.New("invoice not found")
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrInvoiceNotBelongs = errors.New("invoice not belongs to account")
)
