package config

import (
	"GoGin-API-CuentasClaras/api/auth"
	"GoGin-API-CuentasClaras/api/handlers"
	"GoGin-API-CuentasClaras/repository"
	"GoGin-API-CuentasClaras/services"
)

type Initialization struct {
	UserRepo       repository.UserRepository
	operationRepo  repository.OperationRepository
	categoryRepo   repository.CategoryRepository
	userSvc        services.UserService
	operationSvc   services.OperationService
	UserHdler      handlers.UserHandler
	OperationHdler handlers.OperationHandler
	Auth           auth.Auth
}

func NewInitialization(userRepo repository.UserRepository, operationRepo repository.OperationRepository,
	categoryRepo repository.CategoryRepository,
	userService services.UserService, operationSvc services.OperationService,
	UserHdler handlers.UserHandler, OperationHdler handlers.OperationHandler,
	auth auth.Auth) *Initialization {
	return &Initialization{
		UserRepo:       userRepo,
		operationRepo:  operationRepo,
		categoryRepo:   categoryRepo,
		userSvc:        userService,
		operationSvc:   operationSvc,
		UserHdler:      UserHdler,
		OperationHdler: OperationHdler,
		Auth:           auth,
	}
}
