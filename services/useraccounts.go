package services

//go:generate mockgen -source=$GOFILE -destination=../mocks/services/mock_$GOFILE -package=mockservices

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/google/uuid"
	"github.com/paulwrubel/moneybags-server/constants"
	"github.com/paulwrubel/moneybags-server/models"
	"github.com/paulwrubel/moneybags-server/repositories"
	"golang.org/x/crypto/bcrypt"
)

type IUserAccounts interface {
	GetInfo(username string) (*models.UserAccount, error)
	ExistsByUsername(username string) (bool, error)
	Create(username, password string, email *string) (*models.UserAccount, error)
	Delete(username string) error
}

type UserAccounts struct {
	Repository repositories.IUserAccounts
}

func (ua *UserAccounts) GetInfo(username string) (*models.UserAccount, error) {
	userAccount, err := ua.Repository.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("error getting user account: %w", err)
	}
	return userAccount, nil
}

func (ua *UserAccounts) ExistsByUsername(username string) (bool, error) {
	exists, err := ua.Repository.ExistsByUsername(username)
	if err != nil {
		return false, fmt.Errorf("error checking if user account exists: %w", err)
	}
	return exists, nil
}

func (ua *UserAccounts) Create(username, password string, email *string) (*models.UserAccount, error) {
	// check username exists
	exists, err := ua.Repository.ExistsByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("error checking if user exists: %w", err)
	}
	if exists {
		return nil, constants.ErrUserExists
	}

	// check password requirement and hash
	if len(password) < 12 {
		return nil, constants.ErrInvalidPassword
	}
	passwordHash, err := getPasswordHash(password)
	if err != nil {
		return nil, err
	}

	// parse and validate email
	emailVal := email
	if emailVal != nil {
		_, err := mail.ParseAddress(*emailVal)
		if err != nil {
			return nil, constants.ErrInvalidEmail
		}
	}

	// make new account
	newUserAccount := &models.UserAccount{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: passwordHash,
		Email:        emailVal,
	}
	err = ua.Repository.Create(newUserAccount)
	if err != nil {
		return nil, fmt.Errorf("error creating user account: %w", err)
	}
	exists, err = ua.Repository.ExistsByID(newUserAccount.ID)
	if err != nil {
		return nil, fmt.Errorf("error validating if user exists: %w", err)
	}
	if !exists {
		return nil, errors.New("user account failed post creation existence check")
	}
	return ua.Repository.GetByID(newUserAccount.ID)

}

func (ua *UserAccounts) Delete(username string) error {
	userAccount, err := ua.Repository.GetByUsername(username)
	if err != nil {
		return fmt.Errorf("error getting user account: %w", err)
	}
	return ua.Repository.DeleteByID(userAccount.ID)
}

func getPasswordHash(password string) (string, error) {
	passHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passHashBytes), nil
}
