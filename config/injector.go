// go:build wireinject
//go:build wireinject
// +build wireinject

package config

import (
	"GoGin-API-CuentasClaras/api/auth"
	"GoGin-API-CuentasClaras/api/handlers"
	"GoGin-API-CuentasClaras/repository"
	"GoGin-API-CuentasClaras/services"

	"github.com/google/wire"
)

var db = wire.NewSet(ConnectToDB)

var userServiceSet = wire.NewSet(services.UserServiceInit,
	wire.Bind(new(services.UserService), new(*services.UserServiceImpl)),
	auth.AuthInit,
	wire.Bind(new(auth.Auth), new(*auth.AuthImpl)),
)

var operationServiceSet = wire.NewSet(services.OperationServiceInit,
	wire.Bind(new(services.OperationService), new(*services.OperationServiceImpl)),
)

var userRepoSet = wire.NewSet(repository.UserRepositoryInit,
	wire.Bind(new(repository.UserRepository), new(*repository.UserRepositoryImpl)),
)

var operationRepoSet = wire.NewSet(repository.OperationRepositoryInit,
	wire.Bind(new(repository.OperationRepository), new(*repository.OperationRepositoryImpl)),
)

var categoryRepoSet = wire.NewSet(repository.CategoryRepositoryInit,
	wire.Bind(new(repository.CategoryRepository), new(*repository.CategoryRepositoryImpl)),
)

var userHdlerSet = wire.NewSet(handlers.UserHandlerInit,
	wire.Bind(new(handlers.UserHandler), new(*handlers.UserHandlerImpl)),
)

var operationHdlerSet = wire.NewSet(handlers.OperationHandlerInit,
	wire.Bind(new(handlers.OperationHandler), new(*handlers.OperationHandlerImpl)),
)

func Init() *Initialization {
	wire.Build(
		NewInitialization, db, userHdlerSet, operationHdlerSet,
		userServiceSet, operationServiceSet, categoryRepoSet,
		userRepoSet, operationRepoSet,
	)
	return nil
}
