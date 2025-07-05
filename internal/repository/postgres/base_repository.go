package postgres

import (
	"context"

	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/repository"
	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/datasource"
)

type BaseRepository[T any] struct {
	db datasource.Database
}

func NewBaseRepository[T any](db datasource.Database) repository.BaseRepository[T] {
	return &BaseRepository[T]{
		db: db,
	}
}

func (b *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return b.db.DB().WithContext(ctx).Create(entity).Error
}

func (b *BaseRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := b.db.DB().WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (b *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return b.db.DB().WithContext(ctx).Save(entity).Error
}

func (b *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return b.db.DB().WithContext(ctx).Delete(&entity, id).Error
}

func (b *BaseRepository[T]) List(ctx context.Context, offset, limit int) ([]*T, error) {
	var entities []*T
	err := b.db.DB().WithContext(ctx).Offset(offset).Limit(limit).Find(&entities).Error
	return entities, err
}

func (b *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	var entity T
	err := b.db.DB().WithContext(ctx).Model(&entity).Count(&count).Error
	return count, err
}
