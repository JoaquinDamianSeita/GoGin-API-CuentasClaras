package dao

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	ID          int `gorm:"column:id; primary_key; not null" json:"id"`
	Name        string
	Description string
	Color       string
}
