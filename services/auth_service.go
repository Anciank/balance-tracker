package services

import (
	"database/sql"
	"errors"
	"log"

	"balance-tracker/models"
	"balance-tracker/repositories"
	"balance-tracker/utils"
)

type AuthService struct {
	userRepository    repositories.UserRepository
	sessionRepository repositories.SessionRepository
}

func NewAuthService(userRepository *repositories.UserRepository, sessionRepository *repositories.SessionRepository) *AuthService {
	return &AuthService{
		userRepository:    *userRepository,
		sessionRepository: *sessionRepository,
	}
}

func (s *AuthService) Login(username string, password string) (string, error) {
	// Get the user
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	// Check the password
	if !utils.ComparePasswords(user.Password, password) {
		return "", errors.New("invalid password")
	}

	// Generate a token
	token, err := utils.GenerateToken(user)
	if err != nil {
		return "", err
	}

	// Create a new session
	err = s.sessionRepository.CreateSessionFromToken(token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(user models.User) error {
	_, err := s.userRepository.GetUserByUsername(user.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		// User doesn't exist, proceed with registration
	} else {
		return errors.New("username is already taken")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	err = s.userRepository.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Logout(token string) error {
	err := s.sessionRepository.DeleteSessionByToken(token)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) TokenValid(token string) bool {
	session, err := s.sessionRepository.GetSessionByToken(token)
	if err != nil {
		log.Println(err)
		return false
	}

	if session.DeletedAt.Valid {
		log.Println(err)
		return false
	}

	return true
}
