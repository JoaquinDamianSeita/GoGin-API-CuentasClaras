package dao

import (
	"time"
)

type Operation struct {
	ID          int `gorm:"column:id; primary_key; not null" json:"id"`
	UserID      uint
	CategoryID  uint
	Type        string
	Amount      float64
	Date        time.Time
	Description string
	BaseModel
}
