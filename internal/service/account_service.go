package service

import (
	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/dto"
)

type AccountService struct {
	accountRepository domain.AccountRepository
}

func NewAccountService(accountRepository domain.AccountRepository) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
	}
}

func (service *AccountService) CreateAccount(data dto.CreateAccountRequest) (*dto.AccountResponse, error) {
	account := dto.ToAccount(data)

	existingAccount, err := service.accountRepository.FindByApiKey(account.ApiKey)

	if err != nil && err != domain.ErrAccountNotFound {
		return nil, err
	}

	if existingAccount != nil {
		return nil, domain.ErrDuplicatedApiKey
	}

	err = service.accountRepository.Save(account)

	if err != nil {
		return nil, err
	}

	response := dto.FromAccount(account)
	return &response, nil
}

func (service *AccountService) UpdateBalance(apiKey string, amount float64) (*dto.AccountResponse, error) {
	account, err := service.accountRepository.FindByApiKey(apiKey)

	if err != nil {
		return nil, err
	}

	account.UpdateBalance(amount)
	err = service.accountRepository.UpdateBalance(account)

	if err != nil {
		return nil, err
	}

	response := dto.FromAccount(account)
	return &response, nil
}

func (service *AccountService) FindByApiKey(apiKey string) (*dto.AccountResponse, error) {
	account, err := service.accountRepository.FindByApiKey(apiKey)

	if err != nil {
		return nil, err
	}

	response := dto.FromAccount(account)
	return &response, nil
}

func (service *AccountService) FindById(id string) (*dto.AccountResponse, error) {
	account, err := service.accountRepository.FindById(id)

	if err != nil {
		return nil, err
	}

	response := dto.FromAccount(account)
	return &response, nil
}
