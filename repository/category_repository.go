package repository

import (
	"GoGin-API-CuentasClaras/dao"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindCategoryByOperation(operation dao.Operation) (dao.Category, error)
	Save(category *dao.Category) (dao.Category, error)
}

type CategoryRepositoryImpl struct {
	db *gorm.DB
}

func (u CategoryRepositoryImpl) FindCategoryByOperation(operation dao.Operation) (dao.Category, error) {
	u.db.Preload("Category").First(&operation)
	return operation.Category, nil
}

func (u CategoryRepositoryImpl) Save(category *dao.Category) (dao.Category, error) {
	err := u.db.Create(&category).Error
	return *category, err
}

func CategoryRepositoryInit(db *gorm.DB) *CategoryRepositoryImpl {
	db.AutoMigrate(&dao.Category{})
	return &CategoryRepositoryImpl{
		db: db,
	}
}
