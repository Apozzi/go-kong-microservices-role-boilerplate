package usecases

import (
	"errors"
	"login-api/internal/repositories"
	"login-api/models"
)

type RoleUseCase struct {
	repo repositories.RoleRepository
}

func NewRoleUseCase(repo repositories.RoleRepository) *RoleUseCase {
	return &RoleUseCase{repo: repo}
}

func (uc *RoleUseCase) GetAll() ([]models.Role, error) {
	return uc.repo.FindAll()
}

func (uc *RoleUseCase) GetByID(id string) (*models.Role, error) {
	return uc.repo.FindByID(id)
}

func (uc *RoleUseCase) Create(role *models.Role) error {
	existingRole, _ := uc.repo.FindByName(role.Name)
	if existingRole != nil {
		return errors.New("role already exists")
	}

	return uc.repo.Create(role)
}

func (uc *RoleUseCase) Update(id string, role *models.Role) error {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return errors.New("role not found")
	}

	return uc.repo.Update(role)
}

func (uc *RoleUseCase) Delete(id string) error {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return errors.New("role not found")
	}

	count, err := uc.repo.CountUsersByRoleID(id)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("role is assigned to users")
	}

	return uc.repo.Delete(id)
}
