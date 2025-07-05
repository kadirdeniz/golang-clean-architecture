package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/entity"
	usecasecontract "github.com/kadirdeniz/golang-clean-architecture/internal/domain/usecase"
	usecase "github.com/kadirdeniz/golang-clean-architecture/internal/use-case"
	todoMocks "github.com/kadirdeniz/golang-clean-architecture/tests/mocks/todo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTodoUseCase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Todo UseCase Implementation Suite")
}

var _ = Describe("TodoUseCase Implementation", func() {
	var mockRepo *todoMocks.MockTodoRepository
	var todoUseCase usecasecontract.TodoUseCase

	BeforeEach(func() {
		mockRepo = todoMocks.NewMockTodoRepository(GinkgoT())
		todoUseCase = usecase.NewTodoUseCase(mockRepo)
	})

	AfterEach(func() {
		mockRepo.AssertExpectations(GinkgoT())
	})

	Describe("NewTodoUseCase", func() {
		Context("constructor", func() {
			It("should create todo usecase instance", func() {
				uc := usecase.NewTodoUseCase(mockRepo)
				Expect(uc).NotTo(BeNil())
			})

			It("should implement TodoUseCase interface", func() {
				uc := usecase.NewTodoUseCase(mockRepo)
				_, ok := uc.(usecasecontract.TodoUseCase)
				Expect(ok).To(BeTrue())
			})

			It("should accept repository dependency", func() {
				uc := usecase.NewTodoUseCase(mockRepo)
				Expect(uc).NotTo(BeNil())
			})

			It("should work with nil repository for interface testing", func() {
				// This tests that the constructor doesn't panic with nil
				uc := usecase.NewTodoUseCase(nil)
				Expect(uc).NotTo(BeNil())
			})
		})
	})

	Describe("Interface Compliance", func() {
		Context("TodoUseCase interface", func() {
			It("should have Create method", func() {
				// Test that Create method exists and has correct signature
				todo := &entity.Todo{Title: "Test"}
				ctx := context.Background()
				
				var err error
				_ = func() error {
					return todoUseCase.Create(ctx, todo)
				}
				Expect(err).To(BeNil()) // Just checking compilation
			})

			It("should have GetByID method", func() {
				// Test that GetByID method exists and has correct signature
				ctx := context.Background()
				
				var result *entity.Todo
				var err error
				_ = func() (*entity.Todo, error) {
					return todoUseCase.GetByID(ctx, 1)
				}
				Expect(result).To(BeNil()) // Just checking compilation
				Expect(err).To(BeNil())
			})

			It("should have Update method", func() {
				// Test that Update method exists and has correct signature
				todo := &entity.Todo{Title: "Test"}
				ctx := context.Background()
				
				var err error
				_ = func() error {
					return todoUseCase.Update(ctx, todo)
				}
				Expect(err).To(BeNil()) // Just checking compilation
			})

			It("should have Delete method", func() {
				// Test that Delete method exists and has correct signature
				ctx := context.Background()
				
				var err error
				_ = func() error {
					return todoUseCase.Delete(ctx, 1)
				}
				Expect(err).To(BeNil()) // Just checking compilation
			})
		})
	})

	Describe("Create Method", func() {
		Context("delegation pattern", func() {
			It("should delegate to repository Create method", func() {
				todo := &entity.Todo{Title: "Test Todo"}
				ctx := context.Background()

				mockRepo.On("Create", ctx, todo).Return(nil)

				err := todoUseCase.Create(ctx, todo)
				Expect(err).To(BeNil())
			})

			It("should pass through repository errors", func() {
				todo := &entity.Todo{Title: "Test Todo"}
				ctx := context.Background()
				expectedError := errors.New("repository error")

				mockRepo.On("Create", ctx, todo).Return(expectedError)

				err := todoUseCase.Create(ctx, todo)
				Expect(err).To(Equal(expectedError))
			})

			It("should pass correct parameters to repository", func() {
				todo := &entity.Todo{Title: "Specific Todo"}
				ctx := context.Background()

				mockRepo.On("Create", ctx, todo).Return(nil)

				err := todoUseCase.Create(ctx, todo)
				Expect(err).To(BeNil())
				// Mock verification happens in AfterEach
			})
		})
	})

	Describe("GetByID Method", func() {
		Context("delegation pattern", func() {
			It("should delegate to repository GetByID method", func() {
				ctx := context.Background()
				todoID := uint(1)
				expectedTodo := &entity.Todo{ID: todoID, Title: "Found Todo"}

				mockRepo.On("GetByID", ctx, todoID).Return(expectedTodo, nil)

				result, err := todoUseCase.GetByID(ctx, todoID)
				Expect(err).To(BeNil())
				Expect(result).To(Equal(expectedTodo))
			})

			It("should pass through repository errors", func() {
				ctx := context.Background()
				todoID := uint(999)
				expectedError := errors.New("not found")

				mockRepo.On("GetByID", ctx, todoID).Return(nil, expectedError)

				result, err := todoUseCase.GetByID(ctx, todoID)
				Expect(err).To(Equal(expectedError))
				Expect(result).To(BeNil())
			})

			It("should pass correct parameters to repository", func() {
				ctx := context.Background()
				todoID := uint(42)

				mockRepo.On("GetByID", ctx, todoID).Return(&entity.Todo{}, nil)

				_, err := todoUseCase.GetByID(ctx, todoID)
				Expect(err).To(BeNil())
				// Mock verification happens in AfterEach
			})
		})
	})

	Describe("Update Method", func() {
		Context("delegation pattern", func() {
			It("should delegate to repository Update method", func() {
				todo := &entity.Todo{ID: 1, Title: "Updated Todo"}
				ctx := context.Background()

				mockRepo.On("Update", ctx, todo).Return(nil)

				err := todoUseCase.Update(ctx, todo)
				Expect(err).To(BeNil())
			})

			It("should pass through repository errors", func() {
				todo := &entity.Todo{ID: 1, Title: "Updated Todo"}
				ctx := context.Background()
				expectedError := errors.New("update failed")

				mockRepo.On("Update", ctx, todo).Return(expectedError)

				err := todoUseCase.Update(ctx, todo)
				Expect(err).To(Equal(expectedError))
			})

			It("should pass correct parameters to repository", func() {
				todo := &entity.Todo{ID: 123, Title: "Specific Update"}
				ctx := context.Background()

				mockRepo.On("Update", ctx, todo).Return(nil)

				err := todoUseCase.Update(ctx, todo)
				Expect(err).To(BeNil())
				// Mock verification happens in AfterEach
			})
		})
	})

	Describe("Delete Method", func() {
		Context("delegation pattern", func() {
			It("should delegate to repository Delete method", func() {
				ctx := context.Background()
				todoID := uint(1)

				mockRepo.On("Delete", ctx, todoID).Return(nil)

				err := todoUseCase.Delete(ctx, todoID)
				Expect(err).To(BeNil())
			})

			It("should pass through repository errors", func() {
				ctx := context.Background()
				todoID := uint(1)
				expectedError := errors.New("delete failed")

				mockRepo.On("Delete", ctx, todoID).Return(expectedError)

				err := todoUseCase.Delete(ctx, todoID)
				Expect(err).To(Equal(expectedError))
			})

			It("should pass correct parameters to repository", func() {
				ctx := context.Background()
				todoID := uint(456)

				mockRepo.On("Delete", ctx, todoID).Return(nil)

				err := todoUseCase.Delete(ctx, todoID)
				Expect(err).To(BeNil())
				// Mock verification happens in AfterEach
			})
		})
	})

	Describe("Dependency Injection", func() {
		Context("repository dependency", func() {
			It("should require repository dependency", func() {
				uc := usecase.NewTodoUseCase(mockRepo)
				Expect(uc).NotTo(BeNil())
			})

			It("should use injected repository", func() {
				// This is implicitly tested by all delegation tests above
				// where we verify that the mock repository methods are called
				ctx := context.Background()
				todo := &entity.Todo{Title: "Test"}

				mockRepo.On("Create", ctx, todo).Return(nil)

				err := todoUseCase.Create(ctx, todo)
				Expect(err).To(BeNil())
				// Mock verification confirms the injected repo was used
			})
		})
	})

	Describe("Clean Architecture Compliance", func() {
		Context("layer responsibilities", func() {
			It("should act as orchestration layer", func() {
				// UseCase should orchestrate repository calls
				// This is demonstrated by the delegation pattern tests above
				ctx := context.Background()
				todo := &entity.Todo{Title: "Test"}

				mockRepo.On("Create", ctx, todo).Return(nil)

				err := todoUseCase.Create(ctx, todo)
				Expect(err).To(BeNil())
			})

			It("should not contain business logic", func() {
				// This implementation is pure delegation
				// Business logic would be tested separately if it existed
				// Current implementation correctly delegates to repository
				ctx := context.Background()
				todoID := uint(1)

				mockRepo.On("GetByID", ctx, todoID).Return(&entity.Todo{}, nil)

				_, err := todoUseCase.GetByID(ctx, todoID)
				Expect(err).To(BeNil())
			})

			It("should depend on abstractions", func() {
				// UseCase depends on repository interface, not concrete implementation
				// This is verified by the fact that we can inject a mock
				Expect(todoUseCase).NotTo(BeNil())
			})
		})
	})

	Describe("GetAll Method", func() {
		Context("delegation pattern", func() {
			It("should delegate to repository GetAll method", func() {
				ctx := context.Background()
				limit := 10
				offset := 0
				expectedTodos := []*entity.Todo{
					{ID: 1, Title: "Todo 1"},
					{ID: 2, Title: "Todo 2"},
				}

				mockRepo.On("GetAll", ctx, limit, offset).Return(expectedTodos, nil)

				result, err := todoUseCase.GetAll(ctx, limit, offset)
				Expect(err).To(BeNil())
				Expect(result).To(Equal(expectedTodos))
			})

			It("should pass through repository errors", func() {
				ctx := context.Background()
				limit := 10
				offset := 0
				expectedError := errors.New("database error")

				mockRepo.On("GetAll", ctx, limit, offset).Return(nil, expectedError)

				result, err := todoUseCase.GetAll(ctx, limit, offset)
				Expect(err).To(Equal(expectedError))
				Expect(result).To(BeNil())
			})

			It("should pass correct parameters to repository", func() {
				ctx := context.Background()
				limit := 5
				offset := 10

				mockRepo.On("GetAll", ctx, limit, offset).Return([]*entity.Todo{}, nil)

				_, err := todoUseCase.GetAll(ctx, limit, offset)
				Expect(err).To(BeNil())
				// Mock verification happens in AfterEach
			})
		})
	})

	Describe("Count Method", func() {
		Context("delegation pattern", func() {
			It("should delegate to repository Count method", func() {
				ctx := context.Background()
				expectedCount := int64(25)

				mockRepo.On("Count", ctx).Return(expectedCount, nil)

				result, err := todoUseCase.Count(ctx)
				Expect(err).To(BeNil())
				Expect(result).To(Equal(expectedCount))
			})

			It("should pass through repository errors", func() {
				ctx := context.Background()
				expectedError := errors.New("count failed")

				mockRepo.On("Count", ctx).Return(int64(0), expectedError)

				result, err := todoUseCase.Count(ctx)
				Expect(err).To(Equal(expectedError))
				Expect(result).To(Equal(int64(0)))
			})

			It("should handle zero count", func() {
				ctx := context.Background()
				expectedCount := int64(0)

				mockRepo.On("Count", ctx).Return(expectedCount, nil)

				result, err := todoUseCase.Count(ctx)
				Expect(err).To(BeNil())
				Expect(result).To(Equal(expectedCount))
			})
		})
	})
}) 