# 🏗️ Go Clean Architecture with AWS CI/CD

A production-ready Go application implementing Clean Architecture principles with comprehensive CI/CD pipeline for AWS deployment using GitHub Actions, ECR, and ECS.

![Go Version](https://img.shields.io/badge/Go-1.23+-blue?logo=go)
![License](https://img.shields.io/github/license/kadirdeniz/golang-clean-architecture)
![Build Status](https://img.shields.io/github/actions/workflow/status/kadirdeniz/golang-clean-architecture/ci.yml?branch=main)
![Security](https://img.shields.io/github/actions/workflow/status/kadirdeniz/golang-clean-architecture/security.yml?label=security&logo=shield)

## 📋 Table of Contents

- [🚀 Features](#-features)
- [🏛️ Architecture](#️-architecture)
- [🛠️ Tech Stack](#️-tech-stack)
- [⚡ Quick Start](#-quick-start)
- [📦 Installation](#-installation)
- [🔧 Development](#-development)
- [🚀 Deployment](#-deployment)
- [🔒 Security](#-security)
- [📊 Monitoring](#-monitoring)
- [🧪 Testing](#-testing)
- [📚 Documentation](#-documentation)
- [🤝 Contributing](#-contributing)

## 🚀 Features

### Core Features
- **Clean Architecture Implementation** - Separation of concerns with clear layer boundaries
- **Dependency Injection** - Using Google Wire for compile-time DI
- **RESTful API** - Built with Fiber framework for high performance
- **Database Integration** - PostgreSQL with GORM ORM
- **Configuration Management** - Environment-specific configs with Viper
- **Structured Logging** - Using Zap for high-performance logging
- **Input Validation** - Request validation with go-playground/validator
- **API Documentation** - Swagger/OpenAPI 3.0 integration
- **Health Checks** - Application health monitoring endpoints

### DevOps & CI/CD Features
- **🔄 Comprehensive CI/CD Pipeline** - GitHub Actions workflows for all environments
- **🐳 Containerization** - Multi-stage Docker builds with security best practices
- **☁️ AWS Deployment** - ECR and ECS integration with auto-scaling
- **🔒 Security Scanning** - Automated vulnerability and secret detection
- **📦 Release Management** - Automated versioning and multi-architecture builds
- **🌍 Multi-Environment Support** - Dev, Staging, Production deployments
- **✅ Quality Gates** - Automated testing, linting, and security checks
- **🔙 Rollback Capabilities** - Safe deployment rollback procedures

### Development Features
- **🧪 Comprehensive Testing** - Unit, integration, and end-to-end tests using Ginkgo/Gomega
- **🔧 Development Tools** - Hot reload, database migrations, mock generation
- **📊 Code Coverage** - Automated coverage reporting and tracking
- **🔍 Code Quality** - golangci-lint integration with custom rules
- **🎯 Pre-commit Hooks** - Automated quality checks before commits

## 🏛️ Architecture

This project follows **Clean Architecture** principles as described by Robert C. Martin, organized into distinct layers:

```
┌─────────────────────────────────────────────────┐
│                   Delivery                      │  ← HTTP Handlers, Routes, Middleware
├─────────────────────────────────────────────────┤
│                   Use Cases                     │  ← Business Logic, Application Services
├─────────────────────────────────────────────────┤
│                   Repository                    │  ← Data Access, Database Operations
├─────────────────────────────────────────────────┤
│                   Domain                        │  ← Entities, Business Rules
└─────────────────────────────────────────────────┘
│                Infrastructure                   │  ← Config, Logger, Database, External APIs
```

### Layer Responsibilities

- **🎯 Domain Layer** (`internal/domain/`) - Core business entities and rules
- **⚙️ Use Case Layer** (`internal/use-case/`) - Application business logic
- **🚪 Delivery Layer** (`internal/delivery/`) - HTTP handlers and routing
- **💾 Repository Layer** (`internal/repository/`) - Data access implementations  
- **🔧 Infrastructure Layer** (`internal/infrastructure/`) - External concerns

## 🛠️ Tech Stack

### Core Technologies
- **Language**: Go 1.23+
- **Web Framework**: [Fiber](https://gofiber.io/) - Express-inspired web framework
- **Database**: PostgreSQL 15+ with [GORM](https://gorm.io/) ORM
- **Dependency Injection**: [Google Wire](https://github.com/google/wire)
- **Configuration**: [Viper](https://github.com/spf13/viper) - Configuration management
- **Logging**: [Zap](https://github.com/uber-go/zap) - Structured, high-performance logging
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)

### DevOps & Infrastructure
- **Containerization**: Docker with multi-stage builds
- **Cloud Platform**: Amazon Web Services (AWS)
- **Container Registry**: Amazon ECR
- **Container Orchestration**: Amazon ECS with Fargate
- **CI/CD**: GitHub Actions
- **Infrastructure as Code**: CloudFormation templates provided
- **Monitoring**: CloudWatch logs and metrics

### Testing & Quality
- **Testing Framework**: [Ginkgo](https://onsi.github.io/ginkgo/) + [Gomega](https://onsi.github.io/gomega/)
- **Mocking**: [Mockery](https://github.com/vektra/mockery)
- **Linting**: [golangci-lint](https://golangci-lint.run/)
- **Security Scanning**: gosec, Trivy, TruffleHog
- **API Documentation**: [Swaggo](https://github.com/swaggo/swag)

## ⚡ Quick Start

### Prerequisites
- Go 1.23+
- Docker & Docker Compose
- AWS CLI v2 (for deployment)
- Make (optional, but recommended)

### 1️⃣ Clone & Setup
```bash
git clone https://github.com/kadirdeniz/golang-clean-architecture.git
cd golang-clean-architecture
make help  # See all available commands
```

### 2️⃣ Start Development Environment
```bash
make dev  # Starts database, runs migrations, generates code, and starts the app
```

### 3️⃣ Verify Installation
```bash
# Check application health
curl http://localhost:8080/health

# View API documentation
open http://localhost:8080/swagger/index.html
```

## 📦 Installation

### Local Development Setup

```bash
# 1. Install dependencies
go mod download && go mod verify

# 2. Generate Wire dependency injection code
make wire-gen

# 3. Start PostgreSQL database
make db-up

# 4. Run database migrations
make migrate-up

# 5. Start the application
make run
```

### Using Docker

```bash
# Build and run with Docker Compose
make docker-compose-up

# Or build and run Docker image directly
make docker-build
make docker-run
```

### Production Setup

See our comprehensive [AWS Deployment Guide](docs/aws-deployment-guide.md) for production deployment instructions.

## 🔧 Development

### Common Development Tasks

```bash
# Start full development environment
make dev

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Generate mocks
make mocks

# Generate API documentation
make swagger

# Run pre-commit checks
make pre-commit
```

### Project Structure

```
golang-clean-architecture/
├── .github/workflows/          # CI/CD pipelines
├── cmd/server/                 # Application entry point
├── internal/                   # Private application code
│   ├── app/                    # Application setup
│   ├── delivery/http/          # HTTP delivery layer
│   ├── domain/                 # Domain entities and interfaces
│   ├── infrastructure/         # Infrastructure concerns
│   ├── repository/             # Repository implementations
│   └── use-case/               # Business logic layer
├── configs/                    # Configuration files
├── docker/                     # Docker configurations
├── docs/                       # Documentation
├── migrations/                 # Database migrations
├── tests/                      # Test files and mocks
└── aws/                        # AWS deployment templates
```

### Configuration

The application supports multiple environments through configuration files:

- `config.local.yaml` - Local development
- `config.dev.yaml` - Development environment  
- `config.staging.yaml` - Staging environment
- `config.prod.yaml` - Production environment

Environment-specific configuration is loaded based on the `CONFIG_FILE_NAME` environment variable.

## 🚀 Deployment

### CI/CD Pipeline Overview

Our GitHub Actions workflows provide a complete CI/CD pipeline:

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| **CI** | Push to `main`/`develop`, PRs | Build, test, lint, security scan, ECR push |
| **Deploy Dev** | Push to `develop` | Automatic deployment to development |
| **Deploy Staging** | Push to `main` | Manual approval deployment to staging |
| **Deploy Prod** | Manual only | Strict approval process for production |
| **Security** | Daily, Push, PRs | Comprehensive security scanning |
| **Release** | Version tags, Manual | Automated release management |

### Deployment Environments

| Environment | Purpose | Deployment | Approval |
|-------------|---------|------------|----------|
| **Development** | Feature testing | Automatic | None |
| **Staging** | Pre-production validation | Manual trigger | Required |
| **Production** | Live application | Manual only | Strict process |

### Quick Deployment Commands

```bash
# Deploy to development (pushes to ECR, triggers auto-deployment)
make deploy-dev

# Deploy to staging (instructions provided)
make deploy-staging

# Deploy to production (instructions provided)  
make deploy-prod

# Create a release
git tag v1.0.0 && git push origin v1.0.0
```

### AWS Infrastructure

The project includes:
- **ECS Clusters** for container orchestration
- **ECR Repository** for Docker images
- **Application Load Balancer** for traffic distribution
- **CloudWatch** for logging and monitoring
- **SSM Parameter Store** for secure configuration

For detailed setup instructions, see [AWS Deployment Guide](docs/aws-deployment-guide.md).

## 🔒 Security

### Security Features

- **🔐 Secrets Management** - AWS SSM Parameter Store integration
- **🛡️ Vulnerability Scanning** - Daily automated scans
- **🔍 Secret Detection** - Prevents credential leaks
- **🐳 Container Security** - Distroless images, non-root users
- **🚪 CORS Configuration** - Environment-specific CORS policies
- **📝 Security Headers** - HTTP security headers implementation
- **🔑 Input Validation** - Comprehensive request validation

### Security Scanning

The security workflow performs:
- Dependency vulnerability checking
- Static code analysis with gosec
- Docker image scanning with Trivy
- Secret detection with TruffleHog
- License compliance checking

```bash
# Run local security scan
make security-scan
```

## 📊 Monitoring

### Application Monitoring

- **Health Checks** - `/health` endpoint for service monitoring
- **Structured Logging** - JSON logs with correlation IDs
- **CloudWatch Integration** - Centralized log aggregation
- **Metrics Collection** - Application and infrastructure metrics

### Observability Stack

- **Logs**: CloudWatch Logs with structured JSON format
- **Metrics**: CloudWatch metrics and custom application metrics
- **Health Checks**: ECS health checks with proper endpoints
- **Alerts**: CloudWatch alarms for critical issues

## 🧪 Testing

### Testing Strategy

- **Unit Tests** - Domain logic and business rules
- **Integration Tests** - Repository and service layer integration
- **End-to-End Tests** - Full HTTP request/response cycles
- **Contract Tests** - API contract validation

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run tests for specific package
go test ./internal/domain/entity -v

# Run integration tests
make dev-test
```

### Testing Tools

- **Framework**: Ginkgo (BDD-style testing)
- **Assertions**: Gomega (fluent assertion library)
- **Mocking**: Mockery (interface mocking)
- **Coverage**: Built-in Go coverage tools
- **Database**: SQLite for fast test execution

## 📚 Documentation

### Available Documentation

- **[AWS Deployment Guide](docs/aws-deployment-guide.md)** - Complete AWS setup and deployment
- **[API Documentation](http://localhost:8080/swagger/)** - Interactive Swagger UI
- **[Architecture Guide](docs/base-folder-structre.md)** - Detailed folder structure
- **[Technology Stack](docs/tech-stack.md)** - Comprehensive tech stack overview

### Generating Documentation

```bash
# Generate Swagger API documentation
make swagger

# View API documentation
open http://localhost:8080/swagger/index.html
```

## 🤝 Contributing

### Development Workflow

1. **Fork** the repository
2. **Create** a feature branch from `develop`
3. **Make** your changes following the coding standards
4. **Run** pre-commit checks: `make pre-commit`
5. **Submit** a pull request to `develop` branch

### Coding Standards

- Follow [Clean Code principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- Use English for all code, comments, and documentation
- Write comprehensive tests for new features
- Ensure all CI/CD checks pass
- Follow the existing code structure and patterns

### Pull Request Guidelines

- Provide clear description of changes
- Include tests for new functionality
- Update documentation if needed
- Ensure CI/CD pipeline passes
- Request review from maintainers

## 📈 Project Status

### Recent Updates

- ✅ Complete CI/CD pipeline with GitHub Actions
- ✅ AWS deployment with ECS and ECR
- ✅ Comprehensive security scanning
- ✅ Multi-environment support
- ✅ Automated release management
- ✅ Docker multi-stage builds
- ✅ Infrastructure as Code templates

### Roadmap

- [ ] Prometheus metrics integration
- [ ] Grafana dashboards
- [ ] Redis caching layer
- [ ] Message queue integration
- [ ] Rate limiting implementation
- [ ] API versioning strategy

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Robert C. Martin
- [go-clean-arch](https://github.com/bxcodec/go-clean-arch) for architectural inspiration
- AWS for cloud infrastructure and services
- The Go community for excellent tooling and libraries

---

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/kadirdeniz/golang-clean-architecture/issues)
- **Discussions**: [GitHub Discussions](https://github.com/kadirdeniz/golang-clean-architecture/discussions)
- **Documentation**: [Project Wiki](https://github.com/kadirdeniz/golang-clean-architecture/wiki)

**Made with ❤️ and ☕ by the Go Clean Architecture team**
