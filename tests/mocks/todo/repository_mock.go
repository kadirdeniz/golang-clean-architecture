package todo

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/entity"
	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/repository"
)

// MockTodoRepository is a mock implementation of TodoRepository
type MockTodoRepository struct {
	mock.Mock
}

// NewMockTodoRepository creates a new mock instance
func NewMockTodoRepository(t mock.TestingT) *MockTodoRepository {
	mockRepo := &MockTodoRepository{}
	mockRepo.Test(t)
	return mockRepo
}

// Create mocks the Create method
func (m *MockTodoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *MockTodoRepository) GetByID(ctx context.Context, id uint) (*entity.Todo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Todo), args.Error(1)
}

// Update mocks the Update method
func (m *MockTodoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockTodoRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// List mocks the List method
func (m *MockTodoRepository) List(ctx context.Context, offset, limit int) ([]*entity.Todo, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Todo), args.Error(1)
}

// GetAll mocks the GetAll method
func (m *MockTodoRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Todo, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Todo), args.Error(1)
}

// Count mocks the Count method
func (m *MockTodoRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// Ensure MockTodoRepository implements repository.TodoRepository
var _ repository.TodoRepository = (*MockTodoRepository)(nil) 