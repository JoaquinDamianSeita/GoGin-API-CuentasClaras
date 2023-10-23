package repository

import (
	"GoGin-API-Base/dao"

	"gorm.io/gorm"
)

type OperationRepository interface {
}

type OperationRepositoryImpl struct {
	db *gorm.DB
}

func OperationRepositoryInit(db *gorm.DB) *OperationRepositoryImpl {
	db.AutoMigrate(&dao.Operation{})
	return &OperationRepositoryImpl{
		db: db,
	}
}
