package models

import (
	"time"
)

type User struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement;type:integer"`
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Roles        []Role `json:"roles" gorm:"many2many:user_roles;"`
	Password     string
	RegisterDate time.Time
}

type UserWithoutPassword struct {
	ID           int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;NOT NULL"`
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Roles        []Role `json:"roles" gorm:"many2many:user_roles;"`
	RegisterDate time.Time
}

type UserRole struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement;type:integer"`
	UserID int    `json:"user_id" binding:"required"`
	RoleID int    `json:"role_id" binding:"required"`
}

func (UserWithoutPassword) TableName() string {
	return "users"
}
