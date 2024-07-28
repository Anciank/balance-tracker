package repositories

import (
	"database/sql"
	"errors"
	"time"

	"balance-tracker/models"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db}
}

func (r *SessionRepository) GetSession(id int) (models.Session, error) {
	row := r.db.QueryRow("SELECT * FROM sessions WHERE id = $1", id)

	session := models.Session{}
	err := row.Scan(&session.ID, &session.CreatedAt, &session.DeletedAt, &session.Token)
	if err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func (r *SessionRepository) GetSessionByToken(token string) (models.Session, error) {
	row := r.db.QueryRow("SELECT id, created_at, deleted_at, token FROM sessions WHERE token = $1", token)

	session := models.Session{}
	err := row.Scan(&session.ID, &session.CreatedAt, &session.DeletedAt, &session.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Session{}, errors.New("session not found")
		}
		return models.Session{}, err
	}

	return session, nil
}

func (r *SessionRepository) CreateSession(session models.Session) error {
	_, err := r.db.Exec("INSERT INTO sessions (created_at, deleted_at, token) VALUES ($1, $2, $3)", session.CreatedAt, session.DeletedAt, session.Token)
	return err
}

func (r *SessionRepository) CreateSessionFromToken(token string) error {
	_, err := r.db.Exec("INSERT INTO sessions (token) VALUES ($1)", token)
	return err
}

func (r *SessionRepository) DeleteSession(id int) error {
	_, err := r.db.Exec("UPDATE sessions SET deleted_at = $1 WHERE id = $2", time.Now(), id)
	return err
}

func (r *SessionRepository) DeleteSessionByToken(token string) error {
	_, err := r.db.Exec("UPDATE sessions SET deleted_at = $1 WHERE token = $2", time.Now(), token)
	return err
}

func (r *SessionRepository) TokenExists(token string) bool {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS (SELECT 1 FROM sessions WHERE token = $1)", token).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
