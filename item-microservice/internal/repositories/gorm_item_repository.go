package repositories

import (
	"login-api/models"

	"gorm.io/gorm"
)

type ItemRepository interface {
	FindAll() ([]models.Item, error)
	FindByID(id string) (*models.Item, error)
	Create(item *models.Item) error
	Update(item *models.Item) error
	Delete(id string) error
}

type GormItemRepository struct {
	db *gorm.DB
}

func NewGormItemRepository(db *gorm.DB) *GormItemRepository {
	return &GormItemRepository{db: db}
}

func (r *GormItemRepository) FindAll() ([]models.Item, error) {
	var items []models.Item
	result := r.db.Find(&items)
	return items, result.Error
}

func (r *GormItemRepository) FindByID(id string) (*models.Item, error) {
	var item models.Item
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *GormItemRepository) Create(item *models.Item) error {
	return r.db.Create(item).Error
}

func (r *GormItemRepository) Update(item *models.Item) error {
	return r.db.Save(item).Error
}

func (r *GormItemRepository) Delete(id string) error {
	return r.db.Delete(&models.Item{}, "id = ?", id).Error
}
