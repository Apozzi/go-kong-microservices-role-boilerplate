package models

import (
	"time"
)

// Item model info
// @Description Informações do item
type Item struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;type:integer"`
	Descricao    string    `json:"descricao" gorm:"type:varchar(255);not null"`
	Valor        float64   `json:"valor" gorm:"type:decimal(10,2);not null"`
	CriadoEm     time.Time `json:"criado_em" gorm:"default:CURRENT_TIMESTAMP"`
	AtualizadoEm time.Time `json:"atualizado_em" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (Item) TableName() string {
	return "itens"
}
