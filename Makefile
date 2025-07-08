# Go Clean Architecture Makefile with AWS Deployment Support
.PHONY: run test lint mocks clean db-up db-down db-test-up db-test-down db-reset migrate-up migrate-down migrate-force migrate-version migrate-status swagger docker-build docker-run docker-compose-up docker-compose-down wire-gen aws-login ecr-login ecr-push deploy-dev deploy-staging deploy-prod security-scan release

# Application commands
run:
	go run ./cmd/server/main.go

test:
	go test ./... -v -race -coverprofile=coverage.out

test-coverage:
	go test ./... -v -race -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint:
	golangci-lint run --timeout=5m

# Wire dependency injection
wire-gen:
	cd cmd/server && wire && cd ../..

# Mock generation
mocks:
	cd tests/mocks && go generate ./...

clean:
	rm -rf tmp/ dist/ coverage.out coverage.html

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
dev: db-up migrate-up wire-gen
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

# Production Docker commands
prod-build:
	docker build -t golang-clean-architecture:prod .

prod-run:
	docker run -p 8080:8080 --env-file ./configs/.env.prod golang-clean-architecture:prod

# AWS Authentication
aws-login:
	@echo "Logging into AWS..."
	aws sts get-caller-identity

# ECR Authentication and Operations
ecr-login:
	@echo "Logging into Amazon ECR..."
	aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $$(aws sts get-caller-identity --query Account --output text).dkr.ecr.us-east-1.amazonaws.com

ecr-create:
	@echo "Creating ECR repository..."
	aws ecr create-repository --repository-name golang-clean-architecture --region us-east-1 || true

ecr-push: wire-gen docker-build ecr-login
	@echo "Pushing image to ECR..."
	$(eval ACCOUNT_ID=$(shell aws sts get-caller-identity --query Account --output text))
	$(eval ECR_REGISTRY=$(ACCOUNT_ID).dkr.ecr.us-east-1.amazonaws.com)
	$(eval IMAGE_TAG=$(shell git rev-parse --short HEAD))
	docker tag golang-clean-architecture $(ECR_REGISTRY)/golang-clean-architecture:$(IMAGE_TAG)
	docker tag golang-clean-architecture $(ECR_REGISTRY)/golang-clean-architecture:latest
	docker push $(ECR_REGISTRY)/golang-clean-architecture:$(IMAGE_TAG)
	docker push $(ECR_REGISTRY)/golang-clean-architecture:latest
	@echo "Pushed images:"
	@echo "  $(ECR_REGISTRY)/golang-clean-architecture:$(IMAGE_TAG)"
	@echo "  $(ECR_REGISTRY)/golang-clean-architecture:latest"

# GitHub Actions Workflow Triggers
trigger-ci:
	@echo "Triggering CI workflow by pushing to develop branch..."
	git push origin develop

trigger-staging:
	@echo "Triggering staging deployment by pushing to main branch..."
	git push origin main

trigger-prod:
	@echo "Production deployment requires manual trigger via GitHub Actions UI"
	@echo "Go to: https://github.com/$(shell git config --get remote.origin.url | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/actions/workflows/deploy-prod.yml"

# AWS Deployment Commands
deploy-dev: ecr-push
	@echo "Triggering development deployment..."
	@echo "Image pushed to ECR. GitHub Actions will automatically deploy to development."

deploy-staging:
	@echo "To deploy to staging:"
	@echo "1. Ensure your changes are merged to main branch"
	@echo "2. Push to main branch or create a pull request"
	@echo "3. Manual approval will be required for staging deployment"

deploy-prod:
	@echo "To deploy to production:"
	@echo "1. Go to GitHub Actions: https://github.com/$(shell git config --get remote.origin.url | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/actions/workflows/deploy-prod.yml"
	@echo "2. Click 'Run workflow'"
	@echo "3. Provide required inputs (image tag, approver, reason)"
	@echo "4. Wait for security checks and manual approval"

# Security and Quality Assurance
security-scan:
	@echo "Running local security checks..."
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...
	@echo "For comprehensive security scan, check GitHub Actions security workflow"

pre-commit: wire-gen lint test security-scan
	@echo "Pre-commit checks completed successfully!"

# Release Management
release:
	@echo "To create a release:"
	@echo "1. Create and push a version tag: git tag v1.0.0 && git push origin v1.0.0"
	@echo "2. Or use GitHub Actions manual release workflow"
	@echo "3. Release workflow will create binaries and Docker images"

# Infrastructure Management
infra-validate:
	@echo "Validating AWS infrastructure template..."
	aws cloudformation validate-template --template-body file://docs/aws-deployment-guide.md 2>/dev/null || echo "CloudFormation template validation requires separate .yaml file"

# Environment Status
status:
	@echo "=== Go Clean Architecture Status ==="
	@echo "Git Branch: $$(git branch --show-current)"
	@echo "Git Commit: $$(git rev-parse --short HEAD)"
	@echo "Go Version: $$(go version)"
	@echo "Docker Status: $$(docker --version 2>/dev/null || echo 'Docker not available')"
	@echo "AWS CLI: $$(aws --version 2>/dev/null || echo 'AWS CLI not available')"
	@echo ""
	@echo "=== Services Status ==="
	@docker-compose -f docker/docker-compose.dev.yml ps 2>/dev/null || echo "No local services running"

# Help command
help:
	@echo "Go Clean Architecture - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  run              - Run the application locally"
	@echo "  dev              - Start development environment (DB + App)"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  lint             - Run linter"
	@echo "  wire-gen         - Generate Wire dependency injection code"
	@echo "  pre-commit       - Run all pre-commit checks"
	@echo ""
	@echo "Database:"
	@echo "  db-up            - Start local database"
	@echo "  db-down          - Stop local database"
	@echo "  db-reset         - Reset local database"
	@echo "  migrate-up       - Run database migrations"
	@echo "  migrate-down     - Rollback migrations"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Run Docker container"
	@echo "  docker-compose-up - Start all services"
	@echo ""
	@echo "AWS Deployment:"
	@echo "  aws-login        - Check AWS authentication"
	@echo "  ecr-login        - Login to Amazon ECR"
	@echo "  ecr-push         - Build and push image to ECR"
	@echo "  deploy-dev       - Deploy to development"
	@echo "  deploy-staging   - Deploy to staging (instructions)"
	@echo "  deploy-prod      - Deploy to production (instructions)"
	@echo ""
	@echo "Security & Quality:"
	@echo "  security-scan    - Run security vulnerability scan"
	@echo "  swagger          - Generate API documentation"
	@echo ""
	@echo "Release & CI/CD:"
	@echo "  release          - Create new release (instructions)"
	@echo "  trigger-ci       - Trigger CI pipeline"
	@echo "  trigger-staging  - Trigger staging deployment"
	@echo ""
	@echo "Utilities:"
	@echo "  status           - Show environment and services status"
	@echo "  clean            - Clean build artifacts"
	@echo "  help             - Show this help message" 