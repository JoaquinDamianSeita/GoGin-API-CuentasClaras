package repository

import (
	"GoGin-API-CuentasClaras/dao"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindCategoryByOperation(operation dao.Operation) (dao.Category, error)
	Save(category *dao.Category) (dao.Category, error)
	FindCategoryById(id int) (dao.Category, error)
}

type CategoryRepositoryImpl struct {
	db *gorm.DB
}

func (u CategoryRepositoryImpl) FindCategoryById(id int) (dao.Category, error) {
	category := dao.Category{
		ID: id,
	}
	err := u.db.First(&category).Error
	if err != nil {
		log.Error("Got and error when find category by id. Error: ", err)
		return dao.Category{}, err
	}
	return category, nil
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
