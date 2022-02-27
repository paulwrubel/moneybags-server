package services

//go:generate mockgen -source=$GOFILE -destination=../mocks/services/mock_$GOFILE -package=mockservices

import (
	"errors"

	"github.com/google/uuid"
	"github.com/paulwrubel/moneybags-server/models"
	"github.com/paulwrubel/moneybags-server/repositories"
)

type IBankAccounts interface {
	ExistsByID(id string) (bool, error)
	GetAll(budgetID string) ([]*models.BankAccount, error)
	GetByID(id string) (*models.BankAccount, error)
	Create(budgetID string, name string) (*models.BankAccount, error)
	Delete(id string) error
}

type BankAccounts struct {
	Repository repositories.IBankAccounts
}

func (ba *BankAccounts) ExistsByID(id string) (bool, error) {
	return ba.Repository.ExistsByID(id)
}

func (ba *BankAccounts) GetAll(budgetID string) ([]*models.BankAccount, error) {
	return ba.Repository.GetAllByBudgetID(budgetID)
}

func (ba *BankAccounts) GetByID(id string) (*models.BankAccount, error) {
	return ba.Repository.GetByID(id)
}

func (ba *BankAccounts) Create(budgetID string, name string) (*models.BankAccount, error) {
	newBankAccount := &models.BankAccount{
		ID:       uuid.NewString(),
		BudgetID: budgetID,
		Name:     name,
	}
	err := ba.Repository.Create(newBankAccount)
	if err != nil {
		return nil, err
	}
	exists, err := ba.ExistsByID(newBankAccount.ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("bank account failed post creation existence check")
	}
	return ba.Repository.GetByID(newBankAccount.ID)
}

func (ba *BankAccounts) Delete(id string) error {
	return ba.Repository.DeleteByID(id)
}
