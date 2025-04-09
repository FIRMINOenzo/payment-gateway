package service

import (
	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/dto"
)

type InvoiceService struct {
	invoiceRepository domain.InvoiceRepository
	accountService    *AccountService
}

func NewInvoiceService(invoiceRepository domain.InvoiceRepository, accountService *AccountService) *InvoiceService {
	return &InvoiceService{invoiceRepository: invoiceRepository, accountService: accountService}
}

func (s *InvoiceService) Create(request *dto.CreateInvoiceRequest) (*dto.InvoiceResponse, error) {
	account, err := s.accountService.FindByApiKey(request.ApiKey)
	if err != nil {
		return nil, err
	}

	invoice, err := dto.ToInvoice(request, account.Id)
	if err != nil {
		return nil, err
	}

	if err := invoice.Process(); err != nil {
		return nil, err
	}

	if invoice.Status == domain.InvoiceStatusApproved {
		_, err = s.accountService.UpdateBalance(account.ApiKey, invoice.Amount)
		if err != nil {
			return nil, err
		}
	}

	err = s.invoiceRepository.Save(invoice)
	if err != nil {
		return nil, err
	}

	response := dto.FromDomain(invoice)
	return &response, nil
}

func (s *InvoiceService) GetById(invoiceId, apiKey string) (*dto.InvoiceResponse, error) {
	account, err := s.accountService.FindByApiKey(apiKey)

	if err != nil {
		return nil, err
	}

	invoice, err := s.invoiceRepository.FindById(invoiceId)

	if err != nil {
		return nil, err
	}

	if invoice.AccountId != account.Id {
		return nil, domain.ErrInvoiceNotBelongs
	}

	response := dto.FromDomain(invoice)
	return &response, nil
}

func (s *InvoiceService) ListByAccount(accountId string) ([]*dto.InvoiceResponse, error) {
	_, err := s.accountService.FindById(accountId)
	if err != nil {
		return nil, err
	}

	invoices, err := s.invoiceRepository.FindByAccountId(accountId)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.InvoiceResponse, len(invoices))
	for i, invoice := range invoices {
		formattedInvoice := dto.FromDomain(invoice)
		response[i] = &formattedInvoice
	}

	return response, nil
}

func (s *InvoiceService) ListByAccountApiKey(apiKey string) ([]*dto.InvoiceResponse, error) {
	account, err := s.accountService.FindByApiKey(apiKey)
	if err != nil {
		return nil, err
	}

	return s.ListByAccount(account.Id)
}
