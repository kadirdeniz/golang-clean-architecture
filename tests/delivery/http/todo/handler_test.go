package todo_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kadirdeniz/golang-clean-architecture/internal/delivery/http/todo"
	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/entity"
	todoMocks "github.com/kadirdeniz/golang-clean-architecture/tests/mocks/todo"
)

var _ = Describe("Todo Handler", func() {
	var (
		handler        todo.Handler
		app            *fiber.App
		mockUseCase    *todoMocks.MockTodoUseCase
		mockMapper     *todoMocks.MockMapper
		ctrl           *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockUseCase = todoMocks.NewMockTodoUseCase(ctrl)
		mockMapper = todoMocks.NewMockMapper(ctrl)
		handler = todo.NewHandler(mockUseCase, mockMapper)
		app = fiber.New()
		app.Post("/todos", handler.Create)
		app.Get("/todos", handler.GetAll)
		app.Get("/todos/:id", handler.GetByID)
		app.Put("/todos/:id", handler.Update)
		app.Delete("/todos/:id", handler.Delete)
	})

	AfterEach(func() {
		if ctrl != nil {
			ctrl.Finish()
		}
		if app != nil {
			app.Shutdown()
		}
	})

	Describe("NewHandler", func() {
		It("should create a new handler instance", func() {
			Expect(handler).ToNot(BeNil())
		})

		It("should implement Handler interface", func() {
			var _ todo.Handler = handler
		})
	})

	Describe("Create Method", func() {
		Context("Success Cases", func() {
			It("should create todo successfully", func() {
				// Prepare request
				createReq := todo.CreateRequest{
					Title:       "Test Todo",
					Description: stringPtr("Test Description"),
				}
				requestBody, _ := json.Marshal(createReq)

				// Prepare mocks
				todoEntity := &entity.Todo{
					ID:          1,
					Title:       "Test Todo",
					Description: "Test Description",
					Completed:   false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				todoResponse := &todo.Response{
					ID:          1,
					Title:       "Test Todo",
					Description: stringPtr("Test Description"),
					Completed:   false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				mockMapper.EXPECT().
					ToEntity(gomock.Any()).
					Return(todoEntity).
					Times(1)

				mockUseCase.EXPECT().
					Create(gomock.Any(), todoEntity).
					Return(nil).
					Times(1)

				mockMapper.EXPECT().
					ToResponse(todoEntity).
					Return(todoResponse).
					Times(1)

				// Make request
				req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(requestBody))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			})

			It("should return correct response format", func() {
				// Prepare request
				createReq := todo.CreateRequest{
					Title:       "Test Todo",
					Description: stringPtr("Test Description"),
				}
				requestBody, _ := json.Marshal(createReq)

				// Prepare mocks
				todoEntity := &entity.Todo{ID: 1, Title: "Test Todo"}
				todoResponse := &todo.Response{
					ID:          1,
					Title:       "Test Todo",
					Description: stringPtr("Test Description"),
					Completed:   false,
				}

				mockMapper.EXPECT().ToEntity(gomock.Any()).Return(todoEntity)
				mockUseCase.EXPECT().Create(gomock.Any(), todoEntity).Return(nil)
				mockMapper.EXPECT().ToResponse(todoEntity).Return(todoResponse)

				// Make request
				req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(requestBody))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				// Check response content type
				Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("application/json"))
			})
		})

		Context("Error Cases", func() {
			It("should return 400 for invalid JSON", func() {
				invalidJSON := []byte(`{"title": "Test Todo", "description":}`)

				req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(invalidJSON))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for empty request body", func() {
				req := httptest.NewRequest(http.MethodPost, "/todos", nil)
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("should return 500 when usecase fails", func() {
				// Prepare request
				createReq := todo.CreateRequest{
					Title:       "Test Todo",
					Description: stringPtr("Test Description"),
				}
				requestBody, _ := json.Marshal(createReq)

				// Prepare mocks
				todoEntity := &entity.Todo{Title: "Test Todo"}

				mockMapper.EXPECT().
					ToEntity(gomock.Any()).
					Return(todoEntity).
					Times(1)

				mockUseCase.EXPECT().
					Create(gomock.Any(), todoEntity).
					Return(errors.New("database error")).
					Times(1)

				// Make request
				req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(requestBody))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("GetByID Method", func() {
		Context("Success Cases", func() {
			It("should get todo by ID successfully", func() {
				// Prepare mocks
				existingTodo := &entity.Todo{
					ID:          1,
					Title:       "Test Todo",
					Description: "Test Description",
					Completed:   false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				todoResponse := &todo.Response{
					ID:          1,
					Title:       "Test Todo",
					Description: stringPtr("Test Description"),
					Completed:   false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				mockUseCase.EXPECT().
					GetByID(gomock.Any(), uint(1)).
					Return(existingTodo, nil).
					Times(1)

				mockMapper.EXPECT().
					ToResponse(existingTodo).
					Return(todoResponse).
					Times(1)

				// Make request
				req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			It("should return correct response format", func() {
				// Prepare mocks
				existingTodo := &entity.Todo{ID: 1, Title: "Test Todo"}
				todoResponse := &todo.Response{
					ID:          1,
					Title:       "Test Todo",
					Description: stringPtr("Test Description"),
					Completed:   false,
				}

				mockUseCase.EXPECT().GetByID(gomock.Any(), uint(1)).Return(existingTodo, nil)
				mockMapper.EXPECT().ToResponse(existingTodo).Return(todoResponse)

				// Make request
				req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				// Check response content type
				Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("application/json"))
			})
		})

		Context("Error Cases", func() {
			It("should return 400 for invalid todo ID", func() {
				req := httptest.NewRequest(http.MethodGet, "/todos/invalid", nil)
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("should return 404 when todo not found", func() {
				mockUseCase.EXPECT().
					GetByID(gomock.Any(), uint(999)).
					Return(nil, errors.New("not found")).
					Times(1)

				// Make request
				req := httptest.NewRequest(http.MethodGet, "/todos/999", nil)
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("Update Method", func() {
		Context("Success Cases", func() {
			It("should update todo successfully", func() {
				// Prepare request
				updateReq := todo.UpdateRequest{
					Title:       stringPtr("Updated Todo"),
					Description: stringPtr("Updated Description"),
					Completed:   boolPtr(true),
				}
				requestBody, _ := json.Marshal(updateReq)

				// Prepare mocks
				existingTodo := &entity.Todo{
					ID:          1,
					Title:       "Old Title",
					Description: "Old Description",
					Completed:   false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				updatedTodo := &entity.Todo{
					ID:          1,
					Title:       "Updated Todo",
					Description: "Updated Description",
					Completed:   true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				todoResponse := &todo.Response{
					ID:          1,
					Title:       "Updated Todo",
					Description: stringPtr("Updated Description"),
					Completed:   true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				mockUseCase.EXPECT().
					GetByID(gomock.Any(), uint(1)).
					Return(existingTodo, nil).
					Times(1)

				mockMapper.EXPECT().
					ToEntityFromUpdate(existingTodo, gomock.Any()).
					Return(updatedTodo).
					Times(1)

				mockUseCase.EXPECT().
					Update(gomock.Any(), updatedTodo).
					Return(nil).
					Times(1)

				mockMapper.EXPECT().
					ToResponse(updatedTodo).
					Return(todoResponse).
					Times(1)

				// Make request
				req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(requestBody))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			It("should return correct response format", func() {
				// Prepare request
				updateReq := todo.UpdateRequest{
					Title: stringPtr("Updated Todo"),
				}
				requestBody, _ := json.Marshal(updateReq)

				// Prepare mocks
				existingTodo := &entity.Todo{ID: 1, Title: "Old Title"}
				updatedTodo := &entity.Todo{ID: 1, Title: "Updated Todo"}
				todoResponse := &todo.Response{
					ID:          1,
					Title:       "Updated Todo",
					Description: stringPtr("Updated Description"),
					Completed:   false,
				}

				mockUseCase.EXPECT().GetByID(gomock.Any(), uint(1)).Return(existingTodo, nil)
				mockMapper.EXPECT().ToEntityFromUpdate(existingTodo, gomock.Any()).Return(updatedTodo)
				mockUseCase.EXPECT().Update(gomock.Any(), updatedTodo).Return(nil)
				mockMapper.EXPECT().ToResponse(updatedTodo).Return(todoResponse)

				// Make request
				req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(requestBody))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				// Check response content type
				Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("application/json"))
			})
		})

		Context("Error Cases", func() {
			It("should return 400 for invalid todo ID", func() {
				req := httptest.NewRequest(http.MethodPut, "/todos/invalid", bytes.NewReader([]byte(`{}`)))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for invalid JSON", func() {
				invalidJSON := []byte(`{"title": "Updated Todo", "completed":}`)

				req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(invalidJSON))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("should return 404 when todo not found", func() {
				// Prepare request
				updateReq := todo.UpdateRequest{
					Title: stringPtr("Updated Todo"),
				}
				requestBody, _ := json.Marshal(updateReq)

				mockUseCase.EXPECT().
					GetByID(gomock.Any(), uint(999)).
					Return(nil, errors.New("not found")).
					Times(1)

				// Make request
				req := httptest.NewRequest(http.MethodPut, "/todos/999", bytes.NewReader(requestBody))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})

			It("should return 500 when usecase update fails", func() {
				// Prepare request
				updateReq := todo.UpdateRequest{
					Title: stringPtr("Updated Todo"),
				}
				requestBody, _ := json.Marshal(updateReq)

				// Prepare mocks
				existingTodo := &entity.Todo{ID: 1, Title: "Old Title"}
				updatedTodo := &entity.Todo{ID: 1, Title: "Updated Todo"}

				mockUseCase.EXPECT().GetByID(gomock.Any(), uint(1)).Return(existingTodo, nil)
				mockMapper.EXPECT().ToEntityFromUpdate(existingTodo, gomock.Any()).Return(updatedTodo)
				mockUseCase.EXPECT().
					Update(gomock.Any(), updatedTodo).
					Return(errors.New("database error")).
					Times(1)

				// Make request
				req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(requestBody))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("Delete Method", func() {
		Context("Success Cases", func() {
			It("should delete todo successfully", func() {
				// Prepare mocks
				existingTodo := &entity.Todo{
					ID:          1,
					Title:       "Todo to Delete",
					Description: "This will be deleted",
					Completed:   false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				mockUseCase.EXPECT().
					GetByID(gomock.Any(), uint(1)).
					Return(existingTodo, nil).
					Times(1)

				mockUseCase.EXPECT().
					Delete(gomock.Any(), uint(1)).
					Return(nil).
					Times(1)

				// Make request
				req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})

			It("should return 204 No Content for successful deletion", func() {
				// Prepare mocks
				existingTodo := &entity.Todo{ID: 1, Title: "Todo to Delete"}

				mockUseCase.EXPECT().GetByID(gomock.Any(), uint(1)).Return(existingTodo, nil)
				mockUseCase.EXPECT().Delete(gomock.Any(), uint(1)).Return(nil)

				// Make request
				req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// Check that response body is empty for 204
				Expect(resp.ContentLength).To(Equal(int64(0)))
			})
		})

		Context("Error Cases", func() {
			It("should return 400 for invalid todo ID", func() {
				req := httptest.NewRequest(http.MethodDelete, "/todos/invalid", nil)
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("should return 404 when todo not found", func() {
				mockUseCase.EXPECT().
					GetByID(gomock.Any(), uint(999)).
					Return(nil, errors.New("not found")).
					Times(1)

				// Make request
				req := httptest.NewRequest(http.MethodDelete, "/todos/999", nil)
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})

			It("should return 500 when usecase delete fails", func() {
				// Prepare mocks
				existingTodo := &entity.Todo{ID: 1, Title: "Todo to Delete"}

				mockUseCase.EXPECT().GetByID(gomock.Any(), uint(1)).Return(existingTodo, nil)
				mockUseCase.EXPECT().
					Delete(gomock.Any(), uint(1)).
					Return(errors.New("database error")).
					Times(1)

				// Make request
				req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
				resp, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("GetAll Method", func() {
		Context("Success Cases", func() {
			It("should return todos with pagination", func() {
				// Arrange
				page := 1
				limit := 10
				offset := 0
				
				todos := []*entity.Todo{
					{ID: 1, Title: "Todo 1", Description: "Description 1", Completed: false},
					{ID: 2, Title: "Todo 2", Description: "Description 2", Completed: true},
				}
				
				totalCount := int64(25)
				
				mockUseCase.EXPECT().
					GetAll(gomock.Any(), limit, offset).
					Return(todos, nil).
					Times(1)
				
				mockUseCase.EXPECT().
					Count(gomock.Any()).
					Return(totalCount, nil).
					Times(1)
				
				mockMapper.EXPECT().
					ToResponse(todos[0]).
					Return(&todo.Response{
						ID:          1,
						Title:       "Todo 1",
						Description: stringPtr("Description 1"),
						Completed:   false,
					}).
					Times(1)
				
				mockMapper.EXPECT().
					ToResponse(todos[1]).
					Return(&todo.Response{
						ID:          2,
						Title:       "Todo 2",
						Description: stringPtr("Description 2"),
						Completed:   true,
					}).
					Times(1)
				
				// Make request
				req := httptest.NewRequest(http.MethodGet, "/todos/", nil)
				resp, err := app.Test(req)
				
				// Assert
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				
				var response todo.ListResponse
				body, _ := io.ReadAll(resp.Body)
				err = json.Unmarshal(body, &response)
				Expect(err).To(BeNil())
				Expect(response.Todos).To(HaveLen(2))
				Expect(response.Total).To(Equal(totalCount))
				Expect(response.Page).To(Equal(page))
				Expect(response.Limit).To(Equal(limit))
			})
			
			It("should handle pagination parameters correctly", func() {
				// Arrange
				expectedLimit := 5
				expectedOffset := 5 // (2-1) * 5
				
				todos := []*entity.Todo{}
				totalCount := int64(0)
				
				mockUseCase.EXPECT().
					GetAll(gomock.Any(), expectedLimit, expectedOffset).
					Return(todos, nil).
					Times(1)
				
				mockUseCase.EXPECT().
					Count(gomock.Any()).
					Return(totalCount, nil).
					Times(1)
				
				// Make request
				req := httptest.NewRequest(http.MethodGet, "/todos/?page=2&limit=5", nil)
				resp, err := app.Test(req)
				
				// Assert
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				
				var response todo.ListResponse
				body, _ := io.ReadAll(resp.Body)
				err = json.Unmarshal(body, &response)
				Expect(err).To(BeNil())
				Expect(response.Page).To(Equal(2))
				Expect(response.Limit).To(Equal(5))
			})
			
			It("should handle invalid pagination parameters", func() {
				// Arrange
				expectedLimit := 10  // default limit when > 100
				expectedOffset := 0  // default page when < 1
				
				todos := []*entity.Todo{}
				totalCount := int64(0)
				
				mockUseCase.EXPECT().
					GetAll(gomock.Any(), expectedLimit, expectedOffset).
					Return(todos, nil).
					Times(1)
				
				mockUseCase.EXPECT().
					Count(gomock.Any()).
					Return(totalCount, nil).
					Times(1)
				
				// Make request
				req := httptest.NewRequest(http.MethodGet, "/todos/?page=-1&limit=200", nil)
				resp, err := app.Test(req)
				
				// Assert
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				
				var response todo.ListResponse
				body, _ := io.ReadAll(resp.Body)
				err = json.Unmarshal(body, &response)
				Expect(err).To(BeNil())
				Expect(response.Page).To(Equal(1))   // corrected to 1
				Expect(response.Limit).To(Equal(10)) // corrected to 10
			})
		})
		
		Context("Error Cases", func() {
			It("should handle usecase GetAll error", func() {
				// Arrange
				mockUseCase.EXPECT().
					GetAll(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
				
				// Make request
				req := httptest.NewRequest(http.MethodGet, "/todos/", nil)
				resp, err := app.Test(req)
				
				// Assert
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
			
			It("should handle usecase Count error", func() {
				// Arrange
				todos := []*entity.Todo{
					{ID: 1, Title: "Todo 1", Description: "Description 1", Completed: false},
				}
				
				mockUseCase.EXPECT().
					GetAll(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(todos, nil).
					Times(1)
				
				mockUseCase.EXPECT().
					Count(gomock.Any()).
					Return(int64(0), errors.New("count error")).
					Times(1)
				
				// Make request
				req := httptest.NewRequest(http.MethodGet, "/todos/", nil)
				resp, err := app.Test(req)
				
				// Assert
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// Helper function to create bool pointer
func boolPtr(b bool) *bool {
	return &b
} 