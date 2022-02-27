package repositories

//go:generate mockgen -source=$GOFILE -destination=../mocks/repositories/mock_$GOFILE -package=mockrepositories

import (
	"context"
	"errors"

	"github.com/paulwrubel/moneybags-server/database"
	"github.com/paulwrubel/moneybags-server/models"
)

type IBudgets interface {
	GetAllByUserAccountID(userAccountID string) ([]*models.Budget, error)
	ExistsByID(id string) (bool, error)
	ExistsByUserIDAndName(userID, name string) (bool, error)
	GetByID(id string) (*models.Budget, error)
	GetByUserIDAndName(userAccountID, name string) (*models.Budget, error)
	Create(budget *models.Budget) error
	DeleteByID(id string) error
	Update(budget *models.Budget) error
}

type Budgets struct {
	DB database.IHandler
}

func (b *Budgets) ExistsByID(id string) (bool, error) {
	var count int
	err := b.DB.QueryRow(context.Background(), `
		SELECT count(*) 
		FROM budgets 
		WHERE id = $1`, id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (b *Budgets) ExistsByUserIDAndName(userID, name string) (bool, error) {
	var count int
	err := b.DB.QueryRow(context.Background(), `
		SELECT count(*) 
		FROM budgets 
		WHERE 
			user_account_id = $1 AND
			name = $2`,
		userID,
		name).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (b *Budgets) GetAllByUserAccountID(userAccountID string) ([]*models.Budget, error) {
	rows, err := b.DB.Query(context.Background(), `
		SELECT id, user_account_id, name 
		FROM budgets
		WHERE user_account_id = $1`, userAccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	budgets := []*models.Budget{}
	for rows.Next() {
		budget := &models.Budget{}
		err := rows.Scan(&budget.ID, &budget.UserAccountID, &budget.Name)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}

	return budgets, nil
}

func (b *Budgets) GetByID(id string) (*models.Budget, error) {
	budget := &models.Budget{}
	err := b.DB.QueryRow(context.Background(), `
		SELECT id, user_account_id, name 
		FROM budgets 
		WHERE id = $1`, id).Scan(&budget.ID, &budget.UserAccountID, &budget.Name)
	if err != nil {
		return nil, err
	}

	return budget, nil
}

func (b *Budgets) GetByUserIDAndName(userAccountID, name string) (*models.Budget, error) {
	budget := &models.Budget{}
	err := b.DB.QueryRow(context.Background(), `
		SELECT id, user_account_id, name 
		FROM budgets 
		WHERE 
			user_account_id = $1 AND
			name = $2`,
		userAccountID,
		name,
	).Scan(&budget.ID, &budget.Name)
	if err != nil {
		return nil, err
	}

	return budget, nil
}

func (b *Budgets) Create(budget *models.Budget) error {
	tag, err := b.DB.Exec(context.Background(), `
		INSERT INTO budgets (id, user_account_id, name)
		VALUES ($1, $2, $3)`, budget.ID, budget.UserAccountID, budget.Name)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("failed to create budget: unexpected number of rows affected")
	}

	return nil
}

func (b *Budgets) DeleteByID(id string) error {
	tag, err := b.DB.Exec(context.Background(), `
		DELETE FROM budgets
		WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("failed to delete budget: unexpected number of rows affected")
	}

	return nil
}

func (b *Budgets) Update(budget *models.Budget) error {
	tag, err := b.DB.Exec(context.Background(), `
		UPDATE budgets
		SET 
			user_account_id = $2,
			name = $3
		WHERE id = $1`, budget.ID, budget.UserAccountID, budget.Name)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("failed to update budget: unexpected number of rows affected")
	}

	return nil
}
