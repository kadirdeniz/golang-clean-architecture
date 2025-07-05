package usecase

import (
	"context"

	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/entity"
)

type TodoUseCase interface {
	Create(ctx context.Context, todo *entity.Todo) error
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Todo, error)
	GetByID(ctx context.Context, id uint) (*entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}