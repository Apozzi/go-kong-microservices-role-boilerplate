package models

type Role struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement;type:integer"`
	Name string `json:"name"`
}
