//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kadirdeniz/golang-clean-architecture/internal/app"
	"github.com/kadirdeniz/golang-clean-architecture/internal/delivery/http"
	"github.com/kadirdeniz/golang-clean-architecture/internal/delivery/http/health"
	"github.com/kadirdeniz/golang-clean-architecture/internal/delivery/http/todo"
	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/entity"
	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/repository"
	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/config"
	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/datasource"
	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/validator"
	"github.com/kadirdeniz/golang-clean-architecture/internal/repository/postgres"
	usecase "github.com/kadirdeniz/golang-clean-architecture/internal/use-case"
)

var infrastructureSet = wire.NewSet(
	config.NewConfig,
	datasource.NewDatabase,
	validator.NewValidator,
)

var repositorySet = wire.NewSet(
	postgres.NewTodoRepository,
	provideBaseRepository,
)

var useCaseSet = wire.NewSet(
	usecase.NewTodoUseCase,
)

var todoSet = wire.NewSet(
	todo.NewMapper,
	todo.NewHandler,
	todo.NewRouter,
)

var healthSet = wire.NewSet(
	health.NewHandler,
	health.NewRouter,
)

var appSet = wire.NewSet(
	infrastructureSet,
	repositorySet,
	useCaseSet,
	todoSet,
	healthSet,
	http.NewRouter,
	app.NewApp,
)

func provideBaseRepository(db datasource.Database) repository.BaseRepository[entity.Todo] {
	return postgres.NewBaseRepository[entity.Todo](db)
}

func InitializeApp() (app.App, error) {
	wire.Build(appSet)
	return nil, nil
} 