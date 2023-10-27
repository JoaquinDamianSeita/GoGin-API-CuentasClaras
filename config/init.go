package config

import (
	"GoGin-API-CuentasClaras/api/auth"
	"GoGin-API-CuentasClaras/api/handlers"
	"GoGin-API-CuentasClaras/repository"
	"GoGin-API-CuentasClaras/services"
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
