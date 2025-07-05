// Package mocks contains all generated mocks for testing organized by domain
package mocks

// Health domain mocks
//go:generate mockgen -source=../../internal/delivery/http/health/health_handler.go -destination=health/handler_mock.go -package=health

// Todo domain mocks
//go:generate mockgen -source=../../internal/delivery/http/todo/handler.go -destination=todo/handler_mock.go -package=todo
//go:generate mockgen -source=../../internal/domain/usecase/todo_usecase.go -destination=todo/usecase_mock.go -package=todo
//go:generate mockgen -source=../../internal/delivery/http/todo/mapper.go -destination=todo/mapper_mock.go -package=todo
//go:generate mockgen -source=../../internal/domain/repository/todo_repository.go -destination=todo/repository_mock.go -package=todo
//go:generate mockgen -source=../../internal/infrastructure/datasource/datasource.go -destination=infrastructure/datasource_mock.go -package=infrastructure

// Infrastructure mocks
//go:generate mockgen -source=../../internal/infrastructure/config/config.go -destination=infrastructure/config_mock.go -package=infrastructure
