package config

import (
	"GoGin-API-Base/api/auth"
	"GoGin-API-Base/api/handlers"
	"GoGin-API-Base/repository"
	"GoGin-API-Base/services"
)

type Initialization struct {
	userRepo       repository.UserRepository
	operationRepo  repository.OperationRepository
	categoryRepo   repository.CategoryRepository
	userSvc        services.UserService
	operationSvc   services.OperationService
	UserHdler      handlers.UserHandler
	OperationHdler handlers.OperationHandler
	auth           auth.Auth
}

func NewInitialization(userRepo repository.UserRepository, operationRepo repository.OperationRepository,
	categoryRepo repository.CategoryRepository,
	userService services.UserService, operationSvc services.OperationService,
	UserHdler handlers.UserHandler, OperationHdler handlers.OperationHandler,
	auth auth.Auth) *Initialization {
	return &Initialization{
		userRepo:       userRepo,
		operationRepo:  operationRepo,
		categoryRepo:   categoryRepo,
		userSvc:        userService,
		operationSvc:   operationSvc,
		UserHdler:      UserHdler,
		OperationHdler: OperationHdler,
		auth:           auth,
	}
}
