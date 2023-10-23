package services

import (
	"GoGin-API-Base/repository"
)

type OperationService interface {
}

type OperationServiceImpl struct {
	userRepository      repository.UserRepository
	operationRepository repository.OperationRepository
	categoryRepository  repository.CategoryRepository
}

func OperationServiceInit(userRepository repository.UserRepository, operationRepository repository.OperationRepository, categoryRepository repository.CategoryRepository) *OperationServiceImpl {
	return &OperationServiceImpl{
		userRepository:      userRepository,
		operationRepository: operationRepository,
		categoryRepository:  categoryRepository,
	}
}
