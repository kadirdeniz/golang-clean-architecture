package usecase

import (
	"context"

	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/entity"
	"github.com/kadirdeniz/golang-clean-architecture/internal/domain/repository"
	usecasecontract "github.com/kadirdeniz/golang-clean-architecture/internal/domain/usecase"
)

type todoUseCase struct {
	repo repository.TodoRepository
}

func NewTodoUseCase(repo repository.TodoRepository) usecasecontract.TodoUseCase {
	return &todoUseCase{repo: repo}
}

func (u *todoUseCase) Create(ctx context.Context, todo *entity.Todo) error {
	return u.repo.Create(ctx, todo)
}

func (u *todoUseCase) GetByID(ctx context.Context, id uint) (*entity.Todo, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *todoUseCase) Update(ctx context.Context, todo *entity.Todo) error {
	return u.repo.Update(ctx, todo)
}

func (u *todoUseCase) GetAll(ctx context.Context, limit, offset int) ([]*entity.Todo, error) {
	return u.repo.GetAll(ctx, limit, offset)
}

func (u *todoUseCase) Delete(ctx context.Context, id uint) error {
	return u.repo.Delete(ctx, id)
}

func (u *todoUseCase) Count(ctx context.Context) (int64, error) {
	return u.repo.Count(ctx)
} 