package todo

import (
	"time"
)

// @Description Request structure for creating a new todo item
type CreateRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=255" example:"Buy groceries" description:"Title of the todo item"`
	Description *string    `json:"description,omitempty" example:"Milk, bread, eggs" description:"Optional description of the todo"`
}

// @Description Request structure for updating an existing todo item (all fields are optional)
type UpdateRequest struct {
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=1,max=255" example:"Updated task title" description:"New title for the todo"`
	Description *string    `json:"description,omitempty" example:"Updated description" description:"New description for the todo"`
	Completed   *bool      `json:"completed,omitempty" example:"true" description:"Completion status of the todo"`
}

// @Description Response structure for todo operations
type Response struct {
	ID          uint       `json:"id" example:"1" description:"Unique identifier for the todo"`
	Title       string     `json:"title" example:"Buy groceries" description:"Title of the todo item"`
	Description *string    `json:"description,omitempty" example:"Milk, bread, eggs" description:"Optional description of the todo"`
	Completed   bool       `json:"completed" example:"false" description:"Completion status of the todo"`
	CreatedAt   time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z" description:"Timestamp when the todo was created"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2024-01-01T00:00:00Z" description:"Timestamp when the todo was last updated"`
}

// @Description Response structure for paginated todo list operations
type ListResponse struct {
	Todos []Response `json:"todos" description:"Array of todo items"`
	Total int64          `json:"total" example:"25" description:"Total number of todos available"`
	Page  int            `json:"page" example:"1" description:"Current page number"`
	Limit int            `json:"limit" example:"10" description:"Number of items per page"`
} 