package entity

import (
	"strings"
	"time"
)

type Todo struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func NewTodo() *Todo {
	now := time.Now()
	return &Todo{
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (t *Todo) SetTitle(title string) {
	t.Title = strings.TrimSpace(title)
	t.UpdatedAt = time.Now()
}

func (t *Todo) SetDescription(description string) {
	t.Description = strings.TrimSpace(description)
	t.UpdatedAt = time.Now()
}

func (t *Todo) MarkAsCompleted() {
	t.Completed = true
	t.UpdatedAt = time.Now()
}

func (t *Todo) MarkAsIncomplete() {
	t.Completed = false
	t.UpdatedAt = time.Now()
}

func (t *Todo) UpdateTitle(title string) {
	t.Title = strings.TrimSpace(title)
	t.UpdatedAt = time.Now()
}

func (t *Todo) UpdateDescription(description string) {
	t.Description = strings.TrimSpace(description)
	t.UpdatedAt = time.Now()
}

func (t *Todo) IsCompleted() bool {
	return t.Completed
} 