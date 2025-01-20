package repositories

import (
	"login-api/models"

	"gorm.io/gorm"
)

type RoleRepository interface {
	FindAll() ([]models.Role, error)
	FindByID(id string) (*models.Role, error)
	FindByName(name string) (*models.Role, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id string) error
	CountUsersByRoleID(id string) (int64, error)
}

type GormRoleRepository struct {
	db *gorm.DB
}

func NewGormRoleRepository(db *gorm.DB) *GormRoleRepository {
	return &GormRoleRepository{db: db}
}

func (r *GormRoleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	result := r.db.Find(&roles)
	return roles, result.Error
}

func (r *GormRoleRepository) FindByID(id string) (*models.Role, error) {
	var role models.Role
	result := r.db.Where("id = ?", id).First(&role)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

func (r *GormRoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	result := r.db.Where("name = ?", name).First(&role)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

func (r *GormRoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *GormRoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *GormRoleRepository) Delete(id string) error {
	return r.db.Delete(&models.Role{}, "id = ?", id).Error
}

func (r *GormRoleRepository) CountUsersByRoleID(id string) (int64, error) {
	var count int64
	result := r.db.Table("user_roles").Where("role_id = ?", id).Count(&count)
	return count, result.Error
}
