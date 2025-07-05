package repository

import (
	"context"

	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/entity"
)

type TodoRepository interface {
	BaseRepository[entity.Todo]
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Todo, error)
} 