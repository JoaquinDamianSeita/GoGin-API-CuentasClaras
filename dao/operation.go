package dao

import (
	"time"
)

type Operation struct {
	ID          int       `gorm:"column:id; primary_key; not null" json:"id"`
	UserID      uint      `json:"-"`
	Category    Category  `gorm:"foreignKey:ID"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	BaseModel
}
