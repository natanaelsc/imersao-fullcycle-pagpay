package service

import (
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/domain"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/dto"
)

type AccountService struct {
	accountRepository domain.AccountRepository
}

func NewAccountService(accountRepository domain.AccountRepository) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
	}
}

func (s *AccountService) Create(input dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {
	account := dto.ToAccount(input)

	existingAccount, err := s.accountRepository.FindByAPIKey(account.APIKey)
	if err != nil && err != domain.ErrAccountNotFound {
		return nil, err
	}

	if existingAccount != nil {
		return nil, domain.ErrDuplicatedAPIKey
	}

	existingAccount, err = s.accountRepository.FindByEmail(account.Email)
	if err != nil && err != domain.ErrAccountNotFound {
		return nil, err
	}

	if existingAccount.Email == account.Email {
		return nil, domain.ErrAccountAlreadyExists
	}

	err = s.accountRepository.Save(account)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil
}

func (s *AccountService) UpdateBalance(apiKey string, amount float64) (*dto.CreateAccountResponse, error) {
	account, err := s.accountRepository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	account.AddBalance(amount)

	err = s.accountRepository.UpdateBalance(account)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil
}

func (s *AccountService) FindByAPIKey(apiKey string) (*dto.CreateAccountResponse, error) {
	account, err := s.accountRepository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil
}

func (s *AccountService) FindByID(id string) (*dto.CreateAccountResponse, error) {
	account, err := s.accountRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil
}

func (s *AccountService) FindByEmail(email string) (*dto.CreateAccountResponse, error) {
	account, err := s.accountRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil
}
