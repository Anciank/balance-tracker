package repositories

import (
	"database/sql"

	"balance-tracker/models"
)

type BalanceRepository struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) *BalanceRepository {
	return &BalanceRepository{db}
}

func (r *BalanceRepository) GetBalances() ([]models.Balance, error) {
	rows, err := r.db.Query("SELECT * FROM balances")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := []models.Balance{}
	for rows.Next() {
		balance := models.Balance{}
		err := rows.Scan(&balance.ID, &balance.UserID, &balance.Amount, &balance.CreatedAt, &balance.UpdatedAt)
		if err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}

	return balances, nil
}

func (r *BalanceRepository) GetBalance(id string) (models.Balance, error) {
	row := r.db.QueryRow("SELECT * FROM balances WHERE id = $1", id)

	balance := models.Balance{}
	err := row.Scan(&balance.ID, &balance.UserID, &balance.Amount, &balance.CreatedAt, &balance.UpdatedAt)
	if err != nil {
		return models.Balance{}, err
	}

	return balance, nil
}

func (r *BalanceRepository) CreateBalance(balance models.Balance) error {
	_, err := r.db.Exec("INSERT INTO balances (user_id, amount) VALUES ($1, $2)", balance.UserID, balance.Amount)
	return err
}

func (r *BalanceRepository) UpdateBalance(id string, balance models.Balance) error {
	_, err := r.db.Exec("UPDATE balances SET user_id = $1, amount = $2, updated_at = $3 WHERE id = $4", balance.UserID, balance.Amount, balance.UpdatedAt, id)
	return err
}

func (r *BalanceRepository) DeleteBalance(id int) error {
	_, err := r.db.Exec("DELETE FROM balances WHERE id = $1", id)
	return err
}

func (r *BalanceRepository) GetLastBalance() (models.Balance, error) {
	var balance models.Balance
	query := "SELECT * FROM balances ORDER BY id DESC LIMIT 1"
	err := r.db.QueryRow(query).Scan(&balance.ID, &balance.UserID, &balance.Amount, &balance.CreatedAt, &balance.UpdatedAt)
	return balance, err
}

func (r *BalanceRepository) GetBalancesByUserID(userID int) ([]models.Balance, error) {
	rows, err := r.db.Query("SELECT * FROM balances WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := []models.Balance{}
	for rows.Next() {
		balance := models.Balance{}
		err := rows.Scan(&balance.ID, &balance.UserID, &balance.Amount, &balance.CreatedAt, &balance.UpdatedAt)
		if err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}

	return balances, nil
}
