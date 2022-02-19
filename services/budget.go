package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/paulwrubel/moneybags-server/models"
	"github.com/paulwrubel/moneybags-server/repositories"
)

type IBudget interface {
	Exists(id string) (bool, error)
	Get(id string) (*models.Budget, error)
	GetAll() ([]*models.Budget, error)
	Create(name string) (*models.Budget, error)
	Delete(id string) error
}

type Budget struct {
	Repository repositories.IBudget
}

func (b *Budget) Exists(id string) (bool, error) {
	return b.Repository.ExistsByID(id)
}

func (b *Budget) Get(id string) (*models.Budget, error) {
	return b.Repository.GetByID(id)
}

func (b *Budget) GetAll() ([]*models.Budget, error) {
	return b.Repository.GetAll()
}

func (b *Budget) Create(name string) (*models.Budget, error) {
	newBudget := &models.Budget{
		ID:   uuid.NewString(),
		Name: name,
	}
	err := b.Repository.Create(newBudget)
	if err != nil {
		return nil, err
	}
	exists, err := b.Exists(newBudget.ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("could not create budget")
	}
	return b.Repository.GetByID(newBudget.ID)
}

func (b *Budget) Delete(id string) error {
	return b.Repository.DeleteByID(id)
}
