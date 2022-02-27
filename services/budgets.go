package services

//go:generate mockgen -source=$GOFILE -destination=../mocks/services/mock_$GOFILE -package=mockservices

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/paulwrubel/moneybags-server/models"
	"github.com/paulwrubel/moneybags-server/repositories"
)

type IBudgets interface {
	BelongsTo(userAccountID, budgetID string) (bool, error)
	ExistsByID(id string) (bool, error)
	ExistsByUserIDAndName(userAccountID, name string) (bool, error)
	GetAllByUserAccountID(userAccountID string) ([]*models.Budget, error)
	GetByID(id string) (*models.Budget, error)
	GetByUserIDAndName(userAccountID, name string) (*models.Budget, error)
	Create(userAccountId, name string) (*models.Budget, error)
	Delete(id string) error
}

type Budgets struct {
	RBudgets repositories.IBudgets
}

func (b *Budgets) BelongsTo(userAccountID, budgetID string) (bool, error) {
	budget, err := b.RBudgets.GetByID(budgetID)
	if err != nil {
		return false, fmt.Errorf("failed to get budget by id: %v", err)
	}
	return userAccountID == budget.UserAccountID, nil
}

func (b *Budgets) ExistsByID(id string) (bool, error) {
	return b.RBudgets.ExistsByID(id)
}

func (b *Budgets) ExistsByUserIDAndName(userAccountID, name string) (bool, error) {
	return b.RBudgets.ExistsByUserIDAndName(userAccountID, name)
}

func (b *Budgets) GetAllByUserAccountID(userAccountID string) ([]*models.Budget, error) {
	budgets, err := b.RBudgets.GetAllByUserAccountID(userAccountID)
	if err != nil {
		return nil, err
	}
	return budgets, nil
}

func (b *Budgets) GetByID(id string) (*models.Budget, error) {
	return b.RBudgets.GetByID(id)
}

func (b *Budgets) GetByUserIDAndName(userAccountID, name string) (*models.Budget, error) {
	return b.RBudgets.GetByUserIDAndName(userAccountID, name)
}

func (b *Budgets) Create(userAccountID, name string) (*models.Budget, error) {
	newBudget := &models.Budget{
		ID:            uuid.NewString(),
		UserAccountID: userAccountID,
		Name:          name,
	}
	err := b.RBudgets.Create(newBudget)
	if err != nil {
		return nil, err
	}
	exists, err := b.ExistsByID(newBudget.ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("budget failed post-creation existence check")
	}
	return b.RBudgets.GetByID(newBudget.ID)
}

func (b *Budgets) Delete(id string) error {
	return b.RBudgets.DeleteByID(id)
}
