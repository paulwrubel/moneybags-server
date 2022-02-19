package repositories

import (
	"context"
	"errors"

	"github.com/paulwrubel/moneybags-server/database"
	"github.com/paulwrubel/moneybags-server/models"
)

type IBudget interface {
	ExistsByID(id string) (bool, error)
	GetByID(id string) (*models.Budget, error)
	GetAll() ([]*models.Budget, error)
	Create(budget *models.Budget) error
	DeleteByID(id string) error
	Update(budget *models.Budget) error
}

type Budget struct {
	DB database.IDBHandler
}

func (b *Budget) ExistsByID(id string) (bool, error) {
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

func (b *Budget) GetByID(id string) (*models.Budget, error) {
	budget := &models.Budget{}
	err := b.DB.QueryRow(context.Background(), `
		SELECT id, name 
		FROM budgets 
		WHERE id = $1`, id).Scan(&budget.ID, &budget.Name)
	if err != nil {
		return nil, err
	}

	return budget, nil
}

func (b *Budget) GetAll() ([]*models.Budget, error) {
	rows, err := b.DB.Query(context.Background(), `
		SELECT id, name 
		FROM budgets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	budgets := []*models.Budget{}
	for rows.Next() {
		budget := &models.Budget{}
		err := rows.Scan(&budget.ID, &budget.Name)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}

	return budgets, nil
}

func (b *Budget) Create(budget *models.Budget) error {
	tag, err := b.DB.Exec(context.Background(), `
		INSERT INTO budgets (id, name)
		VALUES ($1, $2)`, budget.ID, budget.Name)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("failed to create budget: unexpected number of rows affected")
	}

	return nil
}

func (b *Budget) DeleteByID(id string) error {
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

func (b *Budget) Update(budget *models.Budget) error {
	tag, err := b.DB.Exec(context.Background(), `
		UPDATE budgets
		SET name = $2
		WHERE id = $1`, budget.ID, budget.Name)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("failed to update budget: unexpected number of rows affected")
	}

	return nil
}
