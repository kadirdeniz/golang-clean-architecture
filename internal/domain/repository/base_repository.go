package repository

import (
	"context"
)

type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*T, error)
	Count(ctx context.Context) (int64, error)
} 