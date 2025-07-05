package infrastructure

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDatabase is a mock implementation of datasource.Database interface
type MockDatabase struct {
	mock.Mock
}

// NewMockDatabase creates a new mock database instance
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{}
}

// Open mocks the database connection opening
func (m *MockDatabase) Open(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Close mocks the database connection closing
func (m *MockDatabase) Close() error {
	args := m.Called()
	return args.Error(0)
}

// HealthCheck mocks the database health check
func (m *MockDatabase) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// DB mocks the GORM database instance
func (m *MockDatabase) DB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
} 