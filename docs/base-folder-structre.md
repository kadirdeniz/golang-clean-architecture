# Base Folder Structure - Go Clean Architecture

This document outlines the simplified folder structure for a Go-based Todo API following Clean Architecture principles, without cache and event-driven components.

## Project Structure

```
golang-clean-architecture/
├── .github/
│   └── workflows/
│       ├── ci.yml
│       ├── deploy-dev.yml
│       ├── deploy-staging.yml
│       ├── deploy-prod.yml
│       ├── security.yml
│       └── release.yml
├── cmd/
│   └── server/
│       ├── main.go
│       └── wire.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── domain/
│   │   ├── entity/
│   │   │   └── todo.go
│   │   ├── repository/
│   │   │   ├── base_repository.go
│   │   │   └── todo_repository.go
│   │   └── usecase/
│   │       ├── base_usecase.go
│   │       └── todo_usecase.go
│   ├── delivery/
│   │   ├── http/
│   │   │   ├── handler/
│   │   │   │   └── todo_handler.go
│   │   │   ├── middleware/
│   │   │   │   ├── cors.go
│   │   │   │   ├── logger.go
│   │   │   │   ├── rate_limiter.go
│   │   │   │   ├── recovery.go
│   │   │   │   └── metrics.go
│   │   │   ├── router/
│   │   │   │   └── router.go
│   │   │   ├── dto/
│   │   │   │   ├── todo_dto.go
│   │   │   │   └── response.go
│   │   │   ├── error/
│   │   │   │   └── error_handler.go
│   │   │   └── api/
│   │   │       └── external_service/
│   │   │           ├── dto/
│   │   │           │   ├── request.go
│   │   │           │   └── response.go
│   │   │           └── client.go
│   │   └── grpc/
│   │       ├── handler/
│   │       ├── middleware/
│   │       ├── proto/
│   │       └── service/
│   ├── repository/
│   │   └── postgres/
│   │       ├── base_repository.go
│   │       └── todo_repository.go
│   └── infrastructure/
│       ├── datasource/
│       │   └── postgres/
│       │       ├── connection.go
│       │       └── migration.go
│       ├── config/
│       │   └── config.go
│       ├── logger/
│       │   └── logger.go
│       ├── validator/
│       │   └── validator.go
│       └── metrics/
│           └── prometheus.go
├── configs/
│   ├── config.yaml
│   ├── config.dev.yaml
│   ├── config.staging.yaml
│   ├── config.prod.yaml
│   └── config.test.yaml
├── docker/
│   ├── Dockerfile
│   ├── compose/
│   │   ├── docker-compose.dev.yml
│   │   ├── docker-compose.staging.yml
│   │   └── docker-compose.prod.yml
│   └── .dockerignore
├── migrations/
│   ├── 001_create_todos_table.sql
│   └── 002_create_indexes_and_constraints.sql
├── docs/
│   ├── api/
│   │   └── swagger.json
│   └── architecture/
│       ├── clean-architecture.md
│       └── database-schema.md
├── tests/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── entity_suite_test.go
│   │   │   └── todo_test.go
│   │   ├── repository/
│   │   │   ├── repository_suite_test.go
│   │   │   └── todo_repository_test.go
│   │   └── usecase/
│   │       ├── usecase_suite_test.go
│   │       └── todo_usecase_test.go
│   ├── usecase/
│   │   ├── usecase_suite_test.go
│   │   └── todo_usecase_test.go
│   ├── delivery/
│   │   └── http/
│   │       ├── handler/
│   │       │   ├── handler_suite_test.go
│   │       │   └── todo_handler_test.go
│   │       ├── middleware/
│   │       │   ├── middleware_suite_test.go
│   │       │   ├── cors_test.go
│   │       │   ├── logger_test.go
│   │       │   ├── validator_test.go
│   │       │   ├── rate_limiter_test.go
│   │       │   ├── recovery_test.go
│   │       │   └── metrics_test.go
│   │       └── api/
│   │           └── external_service/
│   │               ├── dto/
│   │               │   ├── request_test.go
│   │               │   └── response_test.go
│   │               └── client_test.go
│   ├── repository/
│   │   └── postgres/
│   │       ├── base_repository_test.go
│   │       └── todo_repository_test.go
│   ├── mocks/
│   │   ├── TodoRepository.go
│   │   └── TodoUseCase.go
│   └── helper/
│       ├── database.go
│       ├── http.go
│       ├── fixtures.go
│       └── metrics.go
├── go.mod
├── go.sum
├── .air.toml
├── go.work
├── Makefile
├── .env.example
├── .gitignore
├── .golangci.yml
├── .pre-commit-config.yaml
└── README.md
```

## Layer Descriptions

### Domain Layer (`internal/domain/`)
- **Entities**: Core business objects (Todo)
- **Repository Interfaces**: Data access contracts
- **Use Case Interfaces**: Business logic contracts

### Application Layer (`internal/usecase/`)
- **Use Case Implementations**: Business logic implementation
- **Orchestration**: Coordination between multiple use cases

### Infrastructure Layer (`internal/repository/`, `internal/delivery/`, `internal/infrastructure/`)
- **Repository Implementations**: Data access implementations (PostgreSQL, external APIs)
- **Delivery**: HTTP handlers, middleware, routing
- **Configuration**: Application configuration management

### Public Packages (`pkg/`)
- **Shared Utilities**: Reusable packages for logging, validation, etc.

## Key Principles

1. **Dependency Inversion**: Inner layers don't depend on outer layers
2. **Separation of Concerns**: Each layer has a specific responsibility
3. **Testability**: Each layer can be tested independently
4. **Maintainability**: Clear boundaries between layers
5. **Scalability**: Easy to add new features without affecting existing code

## Benefits

- **Clean Architecture**: Follows Clean Architecture principles
- **Testability**: Easy to write unit and integration tests
- **Maintainability**: Clear separation of concerns
- **Scalability**: Easy to add new features and modify existing ones
- **Team Collaboration**: Clear structure for team members to understand 