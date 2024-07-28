package services

import (
	"balance-tracker/models"
	"balance-tracker/repositories"
)

type BalanceService struct {
	balanceRepository repositories.BalanceRepository
}

func NewBalanceService(balanceRepository *repositories.BalanceRepository) *BalanceService {
	return &BalanceService{*balanceRepository}
}

func (s *BalanceService) GetBalances() ([]models.Balance, error) {
	balances, err := s.balanceRepository.GetBalances()
	return balances, err
}

func (s *BalanceService) GetBalance(id string) (models.Balance, error) {
	balance, err := s.balanceRepository.GetBalance(id)
	return balance, err
}

func (s *BalanceService) CreateBalance(balance models.Balance) error {
	err := s.balanceRepository.CreateBalance(balance)
	return err
}

func (s *BalanceService) UpdateBalance(id string, balance models.Balance) error {
	err := s.balanceRepository.UpdateBalance(id, balance)
	return err
}

func (s *BalanceService) DeleteBalance(id int) error {
	err := s.balanceRepository.DeleteBalance(id)
	return err
}

func (s *BalanceService) GetLastBalance() (models.Balance, error) {
	balance, err := s.balanceRepository.GetLastBalance()
	return balance, err
}

func (s *BalanceService) GetBalancesByUserID(userID int) ([]models.Balance, error) {
	balances, err := s.balanceRepository.GetBalancesByUserID(userID)
	return balances, err
}
