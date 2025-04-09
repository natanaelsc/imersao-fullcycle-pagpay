package domain

import "errors"

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrDuplicatedAPIKey     = errors.New("api key already exists")
	ErrRequiredAPIKey       = errors.New("api key is required")
	ErrInvoiceNotFound      = errors.New("invoice not found")
	ErrUnauthorizedAccess   = errors.New("unauthorized access")

	ErrInvalidAmount = errors.New("invalid amount")
	ErrInvalidStatus = errors.New("invalid status")
)
