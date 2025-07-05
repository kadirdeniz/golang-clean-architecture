package todo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/usecase"
)

type Handler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type handler struct {
	todoUseCase usecase.TodoUseCase
	mapper      Mapper
}

func NewHandler(todoUseCase usecase.TodoUseCase, mapper Mapper) Handler {
	return &handler{
		todoUseCase: todoUseCase,
		mapper:      mapper,
	}
}

// @Summary Create a new todo
// @Description Create a new todo item with title, description, and optional
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body dto.TodoCreateRequest true "Todo creation request" example({"title":"Buy groceries","description":"Milk, bread, eggs"})
// @Success 201 {object} dto.TodoResponse "Todo created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input data"
// @Failure 422 {object} map[string]interface{} "Validation error"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /todos [post]
func (h *handler) Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	todo := h.mapper.ToEntity(&req)

	if err := h.todoUseCase.Create(c.Context(), todo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create todo",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(h.mapper.ToResponse(todo))
}

// @Summary Get all todos with pagination
// @Description Retrieve a paginated list of todo items
// @Tags todos
// @Produce json
// @Param page query int false "Page number (default: 1)" example(1)
// @Param limit query int false "Items per page (default: 10, max: 100)" example(10)
// @Success 200 {object} dto.TodoListResponse "Todos retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid pagination parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /todos [get]
func (h *handler) GetAll(c *fiber.Ctx) error {
	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get todos with pagination
	todos, err := h.todoUseCase.GetAll(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve todos",
		})
	}

	// Get total count for pagination metadata
	totalCount, err := h.todoUseCase.Count(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count todos",
		})
	}

	// Map todos to response format
	todoResponses := make([]Response, len(todos))
	for i, todo := range todos {
		todoResponses[i] = *h.mapper.ToResponse(todo)
	}

	// Create paginated response
	response := ListResponse{
		Todos: todoResponses,
		Total: totalCount,
		Page:  page,
		Limit: limit,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// @Summary Get a todo by ID
// @Description Retrieve a specific todo item by its ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID" example(1)
// @Success 200 {object} dto.TodoResponse "Todo retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid todo ID"
// @Failure 404 {object} map[string]interface{} "Todo not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /todos/{id} [get]
func (h *handler) GetByID(c *fiber.Ctx) error {
	// Parse todo ID from URL parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid todo ID",
		})
	}

	// Get todo by ID
	todo, err := h.todoUseCase.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(h.mapper.ToResponse(todo))
}

// @Summary Update an existing todo
// @Description Update an existing todo item by ID with new title, description, or completion status
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID" example(1)
// @Param todo body dto.TodoUpdateRequest true "Todo update request" example({"title":"Updated task","description":"Updated description","completed":true})
// @Success 200 {object} dto.TodoResponse "Todo updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input data"
// @Failure 404 {object} map[string]interface{} "Todo not found"
// @Failure 422 {object} map[string]interface{} "Validation error"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /todos/{id} [put]
func (h *handler) Update(c *fiber.Ctx) error {
	// Parse todo ID from URL parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid todo ID",
		})
	}

	// Parse request body
	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get existing todo to update
	existingTodo, err := h.todoUseCase.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo not found",
		})
	}

	// Update todo with new values
	updatedTodo := h.mapper.ToEntityFromUpdate(existingTodo, &req)

	if err := h.todoUseCase.Update(c.Context(), updatedTodo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update todo",
		})
	}

	return c.Status(fiber.StatusOK).JSON(h.mapper.ToResponse(updatedTodo))
}

// @Summary Delete a todo
// @Description Delete an existing todo item by ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID" example(1)
// @Success 204 "Todo deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid todo ID"
// @Failure 404 {object} map[string]interface{} "Todo not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /todos/{id} [delete]
func (h *handler) Delete(c *fiber.Ctx) error {
	// Parse todo ID from URL parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid todo ID",
		})
	}

	// Check if todo exists before deleting
	_, err = h.todoUseCase.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo not found",
		})
	}

	// Delete the todo
	if err := h.todoUseCase.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete todo",
		})
	}

	// Return 204 No Content for successful deletion
	return c.SendStatus(fiber.StatusNoContent)
}
