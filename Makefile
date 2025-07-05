# Minimal Makefile for Go project
.PHONY: run test lint mocks clean db-up db-down db-test-up db-test-down db-reset migrate-up migrate-down migrate-force migrate-version migrate-status swagger docker-build docker-run docker-compose-up docker-compose-down

run:
	go run ./cmd/server/main.go

test:
	go test ./...

lint:
	golangci-lint run

mocks:
	cd tests/mocks && go generate ./...

clean:
	rm -rf tmp

# Database commands
db-up:
	docker-compose -f docker/docker-compose.dev.yml up -d postgres

db-down:
	docker-compose -f docker/docker-compose.dev.yml down

db-test-up:
	docker-compose -f docker/docker-compose.test.yml up -d postgres

db-test-down:
	docker-compose -f docker/docker-compose.test.yml down

db-reset:
	docker-compose -f docker/docker-compose.dev.yml down -v
	docker-compose -f docker/docker-compose.dev.yml up -d postgres

# Migration commands
migrate-up:
	./scripts/migrate.sh up

migrate-down:
	./scripts/migrate.sh down

migrate-force:
	./scripts/migrate.sh force $(version)

migrate-version:
	./scripts/migrate.sh version

migrate-status:
	./scripts/migrate.sh status

# Development commands
dev: db-up migrate-up
	go run ./cmd/server/main.go

dev-test: db-test-up
	ENVIRONMENT=test go test ./...
	docker-compose -f docker/docker-compose.test.yml down

# Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs/swagger --parseDependency --parseInternal

# Docker commands
docker-build:
	docker build -t golang-clean-architecture .

docker-run:
	docker run -p 8080:8080 --env-file ./configs/.env golang-clean-architecture

docker-compose-up:
	docker-compose -f docker/docker-compose.dev.yml up -d

docker-compose-down:
	docker-compose -f docker/docker-compose.dev.yml down

docker-compose-build:
	docker-compose -f docker/docker-compose.dev.yml build

docker-compose-logs:
	docker-compose -f docker/docker-compose.dev.yml logs -f

# Production commands
prod-build:
	docker build -t golang-clean-architecture:prod .

prod-run:
	docker run -p 8080:8080 --env-file ./configs/.env.prod golang-clean-architecture:prod 