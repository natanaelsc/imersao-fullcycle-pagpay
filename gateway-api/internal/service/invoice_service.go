package service

import (
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/domain"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/dto"
)

type InvoiceService struct {
	invoiceRepository domain.InvoiceRepository
	accountService    AccountService
}

func NewInvoiceService(invoiceRepository domain.InvoiceRepository, accountService AccountService) *InvoiceService {
	return &InvoiceService{
		invoiceRepository: invoiceRepository,
		accountService:    accountService,
	}
}

func (s *InvoiceService) Create(input dto.CreateInvoiceRequest) (*dto.InvoiceResponse, error) {
	account, err := s.accountService.FindByAPIKey(input.APIKey)
	if err != nil {
		return nil, domain.ErrAccountNotFound
	}

	invoice, err := dto.ToInvoice(&input, account.ID)
	if err != nil {
		return nil, err
	}

	if err := invoice.Process(); err != nil {
		return nil, err
	}

	if invoice.Status == domain.StatusApproved {
		_, err = s.accountService.UpdateBalance(account.ID, invoice.Amount)
		if err != nil {
			return nil, err
		}
	}

	if err = s.invoiceRepository.Save(invoice); err != nil {
		return nil, err
	}

	return dto.FromInvoice(invoice), nil
}

func (s *InvoiceService) FindById(id, apiKey string) (*dto.InvoiceResponse, error) {
	invoice, err := s.invoiceRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	output, err := s.accountService.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	if invoice.AccountID != output.ID {
		return nil, domain.ErrUnauthorizedAccess
	}

	return dto.FromInvoice(invoice), nil
}

func (s *InvoiceService) ListByAccountID(accountID string) ([]*dto.InvoiceResponse, error) {
	invoices, err := s.invoiceRepository.FindByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	var output []*dto.InvoiceResponse
	for _, invoice := range invoices {
		output = append(output, dto.FromInvoice(invoice))
	}

	return output, nil
}

func (s *InvoiceService) ListByAccountAPIKey(apiKey string) ([]*dto.InvoiceResponse, error) {
	account, err := s.accountService.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	return s.ListByAccountID(account.ID)
}
