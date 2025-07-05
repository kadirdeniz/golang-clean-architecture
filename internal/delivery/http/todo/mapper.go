package todo

import (
	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/entity"
)

type mapper struct{}

type Mapper interface {
	ToEntity(dto *CreateRequest) *entity.Todo
	ToEntityFromUpdate(existingTodo *entity.Todo, dto *UpdateRequest) *entity.Todo
	ToResponse(entity *entity.Todo) *Response
}

func NewMapper() Mapper {
	return &mapper{}
}

func (m *mapper) ToEntity(dto *CreateRequest) *entity.Todo {
	todo := entity.NewTodo()
	todo.SetTitle(dto.Title)
	todo.SetDescription(*dto.Description)
	return todo
}

func (m *mapper) ToEntityFromUpdate(existingTodo *entity.Todo, dto *UpdateRequest) *entity.Todo {
	// Create a copy of the existing todo to avoid modifying the original
	updatedTodo := &entity.Todo{
		ID:          existingTodo.ID,
		Title:       existingTodo.Title,
		Description: existingTodo.Description,
		Completed:   existingTodo.Completed,
		CreatedAt:   existingTodo.CreatedAt,
		UpdatedAt:   existingTodo.UpdatedAt,
	}

	// Update only the fields that are provided in the request
	if dto.Title != nil {
		updatedTodo.SetTitle(*dto.Title)
	}
	if dto.Description != nil {
		updatedTodo.SetDescription(*dto.Description)
	}
	if dto.Completed != nil {
		if *dto.Completed {
			updatedTodo.MarkAsCompleted()
		} else {
			updatedTodo.MarkAsIncomplete()
		}
	}

	return updatedTodo
}

func (m *mapper) ToResponse(entity *entity.Todo) *Response {
	return &Response{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: &entity.Description,
		Completed:   entity.Completed,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}