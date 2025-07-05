package postgres

import (
	"context"

	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/entity"
	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/repository"
	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/datasource"
)

type TodoRepository struct {
	BaseRepository[entity.Todo]
	db datasource.Database
}

func NewTodoRepository(db datasource.Database, baseRepository repository.BaseRepository[entity.Todo]) repository.TodoRepository {
	return &TodoRepository{
		BaseRepository: *baseRepository.(*BaseRepository[entity.Todo]),
		db: db,
	}
}

func (r *TodoRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Todo, error) {
	var todos []*entity.Todo
	
	err := r.db.DB().WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&todos).Error
	
	return todos, err
}

