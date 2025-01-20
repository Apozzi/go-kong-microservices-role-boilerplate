package usecases

import (
	"errors"
	"login-api/internal/repositories"
	"login-api/models"
	"time"
)

type ItemUseCase struct {
	repo repositories.ItemRepository
}

func NewItemUseCase(repo repositories.ItemRepository) *ItemUseCase {
	return &ItemUseCase{repo: repo}
}

func (uc *ItemUseCase) GetAll() ([]models.Item, error) {
	return uc.repo.FindAll()
}

func (uc *ItemUseCase) GetByID(id string) (*models.Item, error) {
	return uc.repo.FindByID(id)
}

func (uc *ItemUseCase) Create(item *models.Item) error {
	if item.Valor <= 0 {
		return errors.New("valor deve ser maior que zero")
	}

	item.CriadoEm = time.Now()
	item.AtualizadoEm = time.Now()

	return uc.repo.Create(item)
}

func (uc *ItemUseCase) Update(id string, item *models.Item) error {
	if item.Valor <= 0 {
		return errors.New("valor deve ser maior que zero")
	}

	existingItem, err := uc.repo.FindByID(id)
	if err != nil {
		return err
	}

	item.ID = existingItem.ID
	item.CriadoEm = existingItem.CriadoEm
	item.AtualizadoEm = time.Now()

	return uc.repo.Update(item)
}

func (uc *ItemUseCase) Delete(id string) error {
	return uc.repo.Delete(id)
}
