package postgres_test

import (
	"context"

	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/entity"
	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/repository"
	"github.com/kadirdeniz/golang-clean-architecture/internal/repository/postgres"
	infrastructureMocks "github.com/kadirdeniz/golang-clean-architecture/tests/mocks/infrastructure"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BaseRepository", func() {
	var mockDB *infrastructureMocks.MockDatabase
	var baseRepo repository.BaseRepository[entity.Todo]

	BeforeEach(func() {
		mockDB = infrastructureMocks.NewMockDatabase()
		baseRepo = postgres.NewBaseRepository[entity.Todo](mockDB)
	})

	Describe("NewBaseRepository", func() {
		Context("constructor", func() {
			It("should create repository instance", func() {
				repo := postgres.NewBaseRepository[entity.Todo](mockDB)
				Expect(repo).NotTo(BeNil())
			})

			It("should implement BaseRepository interface", func() {
				repo := postgres.NewBaseRepository[entity.Todo](mockDB)
				_, ok := repo.(repository.BaseRepository[entity.Todo])
				Expect(ok).To(BeTrue())
			})

			It("should work with different entity types", func() {
				// Test generic nature with different types
				type TestEntity struct {
					ID   uint
					Name string
				}

				testRepo := postgres.NewBaseRepository[TestEntity](mockDB)
				Expect(testRepo).NotTo(BeNil())
				_, ok := testRepo.(repository.BaseRepository[TestEntity])
				Expect(ok).To(BeTrue())
			})
		})
	})

	Describe("Interface Compliance", func() {
		Context("method signatures", func() {
			It("should have Create method", func() {
				// Test that Create method exists and has correct signature
				todo := &entity.Todo{Title: "Test"}
				ctx := context.Background()
				
				// This should compile without errors
				var err error
				_ = func() error {
					return baseRepo.Create(ctx, todo)
				}
				Expect(err).To(BeNil()) // Just checking compilation
			})

			It("should have GetByID method", func() {
				// Test that GetByID method exists and has correct signature
				ctx := context.Background()
				
				// This should compile without errors
				var result *entity.Todo
				var err error
				_ = func() (*entity.Todo, error) {
					return baseRepo.GetByID(ctx, 1)
				}
				Expect(result).To(BeNil()) // Just checking compilation
				Expect(err).To(BeNil())
			})

			It("should have Update method", func() {
				// Test that Update method exists and has correct signature
				todo := &entity.Todo{Title: "Test"}
				ctx := context.Background()
				
				// This should compile without errors
				var err error
				_ = func() error {
					return baseRepo.Update(ctx, todo)
				}
				Expect(err).To(BeNil()) // Just checking compilation
			})

			It("should have Delete method", func() {
				// Test that Delete method exists and has correct signature
				ctx := context.Background()
				
				// This should compile without errors
				var err error
				_ = func() error {
					return baseRepo.Delete(ctx, 1)
				}
				Expect(err).To(BeNil()) // Just checking compilation
			})

			It("should have List method", func() {
				// Test that List method exists and has correct signature
				ctx := context.Background()
				
				// This should compile without errors
				var result []*entity.Todo
				var err error
				_ = func() ([]*entity.Todo, error) {
					return baseRepo.List(ctx, 0, 10)
				}
				Expect(result).To(BeNil()) // Just checking compilation
				Expect(err).To(BeNil())
			})

			It("should have Count method", func() {
				// Test that Count method exists and has correct signature
				ctx := context.Background()
				
				// This should compile without errors
				var result int64
				var err error
				_ = func() (int64, error) {
					return baseRepo.Count(ctx)
				}
				Expect(result).To(Equal(int64(0))) // Just checking compilation
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Generic Type Safety", func() {
		Context("type constraints", func() {
			It("should work with any entity type", func() {
				// Test that the repository works with different entity types
				type AnotherEntity struct {
					ID          uint
					Description string
				}

				anotherRepo := postgres.NewBaseRepository[AnotherEntity](mockDB)
				Expect(anotherRepo).NotTo(BeNil())

				// Verify interface compliance
				_, ok := anotherRepo.(repository.BaseRepository[AnotherEntity])
				Expect(ok).To(BeTrue())
			})

			It("should maintain type safety", func() {
				// This test ensures that the generic types are properly constrained
				// The fact that this compiles means the generic constraints are working
				
				todoRepo := postgres.NewBaseRepository[entity.Todo](mockDB)
				Expect(todoRepo).NotTo(BeNil())

				// Type should be enforced at compile time
				var repo repository.BaseRepository[entity.Todo] = todoRepo
				Expect(repo).NotTo(BeNil())
			})
		})
	})

	Describe("Constructor Parameters", func() {
		Context("database dependency", func() {
			It("should accept database interface", func() {
				repo := postgres.NewBaseRepository[entity.Todo](mockDB)
				Expect(repo).NotTo(BeNil())
			})

			It("should work with nil database for interface testing", func() {
				// This tests that the constructor doesn't panic with nil
				// (though actual operations would fail)
				repo := postgres.NewBaseRepository[entity.Todo](nil)
				Expect(repo).NotTo(BeNil())
			})
		})
	})
})
