package repositories

import (
	"login-api/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(id uint64) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint64) error
	AddRole(userID uint64, roleID string) error
	RemoveRole(userID uint64, roleID string) error
	GetUserWithRoles(id uint64) (*models.User, error)
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Preload("Roles").Find(&users)
	return users, result.Error
}

func (r *GormUserRepository) FindByID(id uint64) (*models.User, error) {
	var user models.User
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *GormUserRepository) GetUserWithRoles(id uint64) (*models.User, error) {
	var user models.User
	result := r.db.Preload("Roles").Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *GormUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *GormUserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *GormUserRepository) Delete(id uint64) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *GormUserRepository) AddRole(userID uint64, roleID string) error {
	return r.db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID).Error
}

func (r *GormUserRepository) RemoveRole(userID uint64, roleID string) error {
	return r.db.Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = ?", userID, roleID).Error
}
