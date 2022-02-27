package repositories

//go:generate mockgen -source=$GOFILE -destination=../mocks/repositories/mock_$GOFILE -package=mockrepositories

import (
	"context"
	"errors"

	"github.com/paulwrubel/moneybags-server/database"
	"github.com/paulwrubel/moneybags-server/models"
)

type IUserAccounts interface {
	ExistsByID(id string) (bool, error)
	ExistsByUsername(username string) (bool, error)
	GetByID(id string) (*models.UserAccount, error)
	GetByUsername(username string) (*models.UserAccount, error)
	Create(userAccount *models.UserAccount) error
	DeleteByID(id string) error
	Update(userAccount *models.UserAccount) error
}

type UserAccounts struct {
	DB database.IHandler
}

func (ua *UserAccounts) ExistsByID(id string) (bool, error) {
	var count int
	err := ua.DB.QueryRow(context.Background(), `
		SELECT count(*) 
		FROM user_accounts 
		WHERE id = $1`, id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (ua *UserAccounts) ExistsByUsername(username string) (bool, error) {
	var count int
	err := ua.DB.QueryRow(context.Background(), `
		SELECT count(*) 
		FROM user_accounts 
		WHERE username = $1`, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (ua *UserAccounts) GetByID(id string) (*models.UserAccount, error) {
	userAccount := &models.UserAccount{}
	err := ua.DB.QueryRow(context.Background(), `
		SELECT 
			id, 
			username, 
			password_hash, 
			email
		FROM user_accounts
		WHERE id = $1`, id).Scan(
		&userAccount.ID,
		&userAccount.Username,
		&userAccount.PasswordHash,
		&userAccount.Email)
	if err != nil {
		return nil, err
	}

	return userAccount, nil
}

func (ua *UserAccounts) GetByUsername(username string) (*models.UserAccount, error) {
	userAccount := &models.UserAccount{}
	err := ua.DB.QueryRow(context.Background(), `
		SELECT 
			id, 
			username, 
			password_hash, 
			email
		FROM user_accounts
		WHERE username = $1`, username).Scan(
		&userAccount.ID,
		&userAccount.Username,
		&userAccount.PasswordHash,
		&userAccount.Email)
	if err != nil {
		return nil, err
	}

	return userAccount, nil
}

func (ua *UserAccounts) Create(userAccount *models.UserAccount) error {
	tag, err := ua.DB.Exec(context.Background(), `
		INSERT INTO user_accounts (
			id, 
			username, 
			password_hash, 
			email
		) VALUES (
			$1, $2, $3, $4
		)`,
		userAccount.ID,
		userAccount.Username,
		userAccount.PasswordHash,
		userAccount.Email)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("failed to create user account: unexpected number of rows affected")
	}

	return nil
}

func (ua *UserAccounts) DeleteByID(id string) error {
	tag, err := ua.DB.Exec(context.Background(), `
		DELETE FROM user_accounts
		WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("failed to delete user account: unexpected number of rows affected")
	}

	return nil
}

func (ua *UserAccounts) Update(userAccount *models.UserAccount) error {
	tag, err := ua.DB.Exec(context.Background(), `
		UPDATE user_accounts
		SET 
			username = $2, 
			password_hash = $3, 
			email = $4
		WHERE id = $1`,
		userAccount.ID,
		userAccount.Username,
		userAccount.PasswordHash,
		userAccount.Email)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("failed to update user account: unexpected number of rows affected")
	}

	return nil
}
