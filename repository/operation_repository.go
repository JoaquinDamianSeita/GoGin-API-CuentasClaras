package repository

import (
	"GoGin-API-CuentasClaras/dao"

	"gorm.io/gorm"
)

type OperationRepository interface {
	FindOperationsByUser(user dao.User) ([]dao.Operation, error)
	Save(operation *dao.Operation) (dao.Operation, error)
}

type OperationRepositoryImpl struct {
	db *gorm.DB
}

func (u OperationRepositoryImpl) FindOperationsByUser(user dao.User) ([]dao.Operation, error) {
	u.db.Preload("Operations").First(&user)
	return user.Operations, nil
}

func (u OperationRepositoryImpl) Save(operation *dao.Operation) (dao.Operation, error) {
	err := u.db.Create(&operation).Error
	return *operation, err
}

func OperationRepositoryInit(db *gorm.DB) *OperationRepositoryImpl {
	db.AutoMigrate(&dao.Operation{})
	return &OperationRepositoryImpl{
		db: db,
	}
}
