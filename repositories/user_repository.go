package repositories

import (
	"database/sql"

	"balance-tracker/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetUser(id int) (models.User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE id = $1", id)

	user := models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE username = $1", username)

	user := models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) CreateUser(user models.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password, created_at, updated_at) VALUES ($1, $2, $3, $4)", user.Username, user.Password, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) UpdateUser(id int, user models.User) error {
	_, err := r.db.Exec("UPDATE users SET username = $1, password = $2, updated_at = $3 WHERE id = $4", user.Username, user.Password, user.UpdatedAt, id)
	return err
}

func (r *UserRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
