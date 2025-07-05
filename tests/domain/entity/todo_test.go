package entity_test

import (
	"testing"
	"time"

	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/entity"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTodo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Todo Entity Suite")
}

var _ = Describe("Todo Entity", func() {
	Describe("NewTodo", func() {
		It("should create a new empty todo with timestamps", func() {
			todo := entity.NewTodo()
			
			Expect(todo).NotTo(BeNil())
			Expect(todo.Title).To(Equal(""))
			Expect(todo.Description).To(Equal(""))
			Expect(todo.Completed).To(BeFalse())
			Expect(todo.CreatedAt).NotTo(BeZero())
			Expect(todo.UpdatedAt).NotTo(BeZero())
			Expect(todo.CreatedAt).To(Equal(todo.UpdatedAt))
		})
	})

	Describe("SetTitle", func() {
		It("should set title with string manipulation", func() {
			todo := entity.NewTodo()
			originalUpdatedAt := todo.UpdatedAt
			
			// Wait a bit to ensure timestamp difference
			time.Sleep(1 * time.Millisecond)
			
			todo.SetTitle("  Test Todo  ")
			
			Expect(todo.Title).To(Equal("Test Todo"))
			Expect(todo.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))
		})

		It("should handle empty title", func() {
			todo := entity.NewTodo()
			originalUpdatedAt := todo.UpdatedAt
			
			time.Sleep(1 * time.Millisecond)
			
			todo.SetTitle("   ")
			
			Expect(todo.Title).To(Equal(""))
			Expect(todo.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))
		})
	})

	Describe("SetDescription", func() {
		It("should set description with string manipulation", func() {
			todo := entity.NewTodo()
			originalUpdatedAt := todo.UpdatedAt
			
			time.Sleep(1 * time.Millisecond)
			
			todo.SetDescription("  Test Description  ")
			
			Expect(todo.Description).To(Equal("Test Description"))
			Expect(todo.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))
		})

		It("should handle empty description", func() {
			todo := entity.NewTodo()
			originalUpdatedAt := todo.UpdatedAt
			
			time.Sleep(1 * time.Millisecond)
			
			todo.SetDescription("   ")
			
			Expect(todo.Description).To(Equal(""))
			Expect(todo.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))
		})
	})

	Describe("MarkAsCompleted", func() {
		It("should mark todo as completed and update timestamp", func() {
			todo := entity.NewTodo()
			todo.SetTitle("Test Todo")
			originalUpdatedAt := todo.UpdatedAt
			
			time.Sleep(1 * time.Millisecond)
			
			todo.MarkAsCompleted()
			
			Expect(todo.Completed).To(BeTrue())
			Expect(todo.IsCompleted()).To(BeTrue())
			Expect(todo.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))
		})
	})

	Describe("MarkAsIncomplete", func() {
		It("should mark todo as incomplete and update timestamp", func() {
			todo := entity.NewTodo()
			todo.SetTitle("Test Todo")
			todo.MarkAsCompleted()
			originalUpdatedAt := todo.UpdatedAt
			
			time.Sleep(1 * time.Millisecond)
			
			todo.MarkAsIncomplete()
			
			Expect(todo.Completed).To(BeFalse())
			Expect(todo.IsCompleted()).To(BeFalse())
			Expect(todo.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))
		})
	})

	Describe("UpdateTitle", func() {
		It("should update title with string manipulation and update timestamp", func() {
			todo := entity.NewTodo()
			todo.SetTitle("Old Title")
			originalUpdatedAt := todo.UpdatedAt
			
			time.Sleep(1 * time.Millisecond)
			
			todo.UpdateTitle("  New Title  ")
			
			Expect(todo.Title).To(Equal("New Title"))
			Expect(todo.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))
		})
	})

	Describe("UpdateDescription", func() {
		It("should update description with string manipulation and update timestamp", func() {
			todo := entity.NewTodo()
			todo.SetDescription("Old Description")
			originalUpdatedAt := todo.UpdatedAt
			
			time.Sleep(1 * time.Millisecond)
			
			todo.UpdateDescription("  New Description  ")
			
			Expect(todo.Description).To(Equal("New Description"))
			Expect(todo.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))
		})
	})
}) 