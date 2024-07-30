package services

import (
	"balance-tracker/models"
	"balance-tracker/repositories"
	"log"
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

func (s *BalanceService) GetLastBalanceByUserID(userID int) (models.Balance, error) {
	balance, err := s.balanceRepository.GetLastBalanceByUserID(userID)
	return balance, err
}

func (s *BalanceService) GetBalancesByUserID(userID int) ([]models.Balance, error) {
	balances, err := s.balanceRepository.GetBalancesByUserID(userID)
	return balances, err
}

func (s *BalanceService) CreateNewTransactionByID(userID int, amount int) (models.Balance, error) {
	lastBalance, err := s.balanceRepository.GetLastBalanceByUserID(userID)
	if err != nil {
		log.Println(err)
		return models.Balance{}, err
	}

	newBalance := models.Balance{
		UserID:    lastBalance.UserID,
		Amount:    lastBalance.Amount + float64(amount),
		CreatedAt: lastBalance.CreatedAt,
	}

	err = s.CreateBalance(newBalance)
	if err != nil {
		log.Println(err)
		return models.Balance{}, err
	}

	return newBalance, nil
}
