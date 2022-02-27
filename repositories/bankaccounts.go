package repositories

//go:generate mockgen -source=$GOFILE -destination=../mocks/repositories/mock_$GOFILE -package=mockrepositories

import (
	"context"
	"errors"

	"github.com/paulwrubel/moneybags-server/database"
	"github.com/paulwrubel/moneybags-server/models"
)

type IBankAccounts interface {
	ExistsByID(id string) (bool, error)
	GetAllByBudgetID(budgetID string) ([]*models.BankAccount, error)
	GetByID(id string) (*models.BankAccount, error)
	Create(bankAccount *models.BankAccount) error
	DeleteByID(id string) error
	Update(bankAccount *models.BankAccount) error
}

type BankAccounts struct {
	DB database.IHandler
}

func (ba *BankAccounts) ExistsByID(id string) (bool, error) {
	var count int
	err := ba.DB.QueryRow(context.Background(), `
		SELECT count(*) 
		FROM bank_accounts 
		WHERE id = $1`, id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (ba *BankAccounts) GetAllByBudgetID(budgetID string) ([]*models.BankAccount, error) {
	rows, err := ba.DB.Query(context.Background(), `
		SELECT id, budget_id, name
		FROM bank_accounts
		WHERE budget_id = $1`, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*models.BankAccount{}
	for rows.Next() {
		account := &models.BankAccount{}
		err := rows.Scan(&account.ID, &account.BudgetID, &account.Name)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (ba *BankAccounts) GetByID(id string) (*models.BankAccount, error) {
	account := &models.BankAccount{}
	err := ba.DB.QueryRow(context.Background(), `
		SELECT id, budget_id, name
		FROM bank_accounts 
		WHERE id = $1`, id).Scan(&account.ID, &account.BudgetID, &account.Name)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (ba *BankAccounts) Create(account *models.BankAccount) error {
	tag, err := ba.DB.Exec(context.Background(), `
		INSERT INTO bank_accounts (id, budget_id, name)
		VALUES ($1, $2, $3)`, account.ID, account.BudgetID, account.Name)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("failed to create bank account: unexpected number of rows affected")
	}

	return nil
}

func (ba *BankAccounts) DeleteByID(id string) error {
	tag, err := ba.DB.Exec(context.Background(), `
		DELETE FROM bank_accounts
		WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("failed to delete bank account: unexpected number of rows affected")
	}

	return nil
}

func (ba *BankAccounts) Update(account *models.BankAccount) error {
	tag, err := ba.DB.Exec(context.Background(), `
		UPDATE bank_accounts 
		SET budget_id = $2, name = $3 
		WHERE id = $1`, account.ID, account.Name, account.BudgetID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("failed to update bank account: unexpected number of rows affected")
	}

	return nil
}
