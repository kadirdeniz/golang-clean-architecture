#!/bin/bash

# Migration script for golang-clean-architecture
# Uses golang-migrate to run database migrations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
MIGRATIONS_PATH="./migrations"
DATABASE_URL="postgres://postgres:postgres@localhost:5432/todo_db_dev?sslmode=disable"
COMMAND="up"

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS] [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  up        Run all pending migrations (default)"
    echo "  down      Rollback last migration"
    echo "  down-all  Rollback all migrations"
    echo "  force     Force a specific migration version"
    echo "  version   Show current migration version"
    echo "  create    Create new migration files"
    echo ""
    echo "Options:"
    echo "  -d, --database URL    Database connection string"
    echo "  -p, --path PATH       Path to migrations directory"
    echo "  -h, --help            Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 up"
    echo "  $0 -d 'postgres://user:pass@localhost:5432/db' up"
    echo "  $0 -p ./custom-migrations down"
    echo "  $0 force 20240101000000"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--database)
            DATABASE_URL="$2"
            shift 2
            ;;
        -p|--path)
            MIGRATIONS_PATH="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        up|down|down-all|force|version|create)
            COMMAND="$1"
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_PATH" ]; then
    print_error "Migrations directory not found: $MIGRATIONS_PATH"
    exit 1
fi

# Check if golang-migrate is installed
if ! command -v migrate &> /dev/null; then
    print_error "golang-migrate is not installed. Please install it first:"
    echo "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

print_status "Using migrations path: $MIGRATIONS_PATH"
print_status "Using database URL: $DATABASE_URL"

# Execute migration command
case $COMMAND in
    up)
        print_status "Running migrations up..."
        migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" up
        ;;
    down)
        print_status "Rolling back last migration..."
        migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" down 1
        ;;
    down-all)
        print_status "Rolling back all migrations..."
        migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" down
        ;;
    force)
        if [ -z "$1" ]; then
            print_error "Version number required for force command"
            exit 1
        fi
        print_status "Forcing migration to version: $1"
        migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" force "$1"
        ;;
    version)
        print_status "Current migration version:"
        migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" version
        ;;
    create)
        if [ -z "$1" ]; then
            print_error "Migration name required for create command"
            exit 1
        fi
        print_status "Creating new migration: $1"
        migrate create -ext sql -dir "$MIGRATIONS_PATH" -seq "$1"
        ;;
    *)
        print_error "Unknown command: $COMMAND"
        show_usage
        exit 1
        ;;
esac

print_status "Migration completed successfully!" 