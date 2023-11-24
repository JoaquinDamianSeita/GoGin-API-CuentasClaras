package repository

import (
	"GoGin-API-CuentasClaras/dao"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OperationRepository interface {
	FindOperationsByUser(user dao.User) ([]dao.Operation, error)
	Save(operation *dao.Operation) (dao.Operation, error)
	FindOperationByUserAndId(user dao.User, operationID int) (dao.Operation, error)
	Update(operation *dao.Operation) (dao.Operation, error)
}

type OperationRepositoryImpl struct {
	db *gorm.DB
}

func (u OperationRepositoryImpl) FindOperationsByUser(user dao.User) ([]dao.Operation, error) {
	u.db.Preload("Operations").First(&user)
	return user.Operations, nil
}

func (u OperationRepositoryImpl) FindOperationByUserAndId(user dao.User, operationID int) (dao.Operation, error) {
	operation := dao.Operation{
		ID:     operationID,
		UserID: uint(user.ID),
	}
	err := u.db.First(&operation).Error
	if err != nil {
		log.Error("Got and error when find operation by id. Error: ", err)
		return dao.Operation{}, err
	}
	u.db.Preload("Category").First(&operation)
	return operation, nil
}

func (u OperationRepositoryImpl) Save(operation *dao.Operation) (dao.Operation, error) {
	err := u.db.Create(&operation).Error
	return *operation, err
}

func (u OperationRepositoryImpl) Update(operation *dao.Operation) (dao.Operation, error) {
	err := u.db.Save(&operation).Error
	return *operation, err
}

func OperationRepositoryInit(db *gorm.DB) *OperationRepositoryImpl {
	db.AutoMigrate(&dao.Operation{})
	return &OperationRepositoryImpl{
		db: db,
	}
}
