# Pagination Implementation Guide

## Overview

Bu doküman, Go'da pagination implementasyonu için en iyi yaklaşımları ve kullanabileceğimiz araçları detaylandırır.

## Pagination Yaklaşımları

### 1. Offset-Based Pagination (Önerilen)
- **Avantajlar**: Basit implementasyon, anlaşılır
- **Dezavantajlar**: Büyük veri setlerinde performans sorunu
- **Kullanım**: Küçük-orta ölçekli uygulamalar için ideal

### 2. Cursor-Based Pagination
- **Avantajlar**: Performanslı, büyük veri setleri için uygun
- **Dezavantajlar**: Karmaşık implementasyon
- **Kullanım**: Büyük ölçekli uygulamalar için ideal

## Go Paketleri

### 1. GORM Pagination
```go
// GORM ile basit pagination
func (r *TodoRepository) GetAll(ctx context.Context, page, limit int) ([]*Todo, int64, error) {
    var todos []*Todo
    var total int64
    
    offset := (page - 1) * limit
    
    // Total count
    if err := r.db.Model(&Todo{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Paginated data
    if err := r.db.Offset(offset).Limit(limit).Find(&todos).Error; err != nil {
        return nil, 0, err
    }
    
    return todos, total, nil
}
```

### 2. Custom Pagination Helper
```go
// Pagination struct'ı
type Pagination struct {
    Page       int   `json:"page"`
    Limit      int   `json:"limit"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
    HasNext    bool  `json:"has_next"`
    HasPrev    bool  `json:"has_prev"`
    NextPage   *int  `json:"next_page"`
    PrevPage   *int  `json:"prev_page"`
}

// Pagination helper
func NewPagination(page, limit int, total int64) *Pagination {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    if limit > 100 {
        limit = 100
    }
    
    totalPages := int(math.Ceil(float64(total) / float64(limit)))
    
    var nextPage, prevPage *int
    
    if page < totalPages {
        next := page + 1
        nextPage = &next
    }
    
    if page > 1 {
        prev := page - 1
        prevPage = &prev
    }
    
    return &Pagination{
        Page:       page,
        Limit:      limit,
        Total:      total,
        TotalPages: totalPages,
        HasNext:    page < totalPages,
        HasPrev:    page > 1,
        NextPage:   nextPage,
        PrevPage:   prevPage,
    }
}
```

### 3. Response Wrapper
```go
// Paginated response wrapper
type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Pagination *Pagination `json:"pagination"`
}

func NewPaginatedResponse(data interface{}, pagination *Pagination) *PaginatedResponse {
    return &PaginatedResponse{
        Data:       data,
        Pagination: pagination,
    }
}
```

## Repository Implementation

```go
// TodoRepository interface güncelleme
type TodoRepository interface {
    Create(ctx context.Context, todo *Todo) error
    GetByID(ctx context.Context, id uint) (*Todo, error)
    GetAll(ctx context.Context, page, limit int) ([]*Todo, int64, error)
    Update(ctx context.Context, todo *Todo) error
    Delete(ctx context.Context, id uint) error
}

// PostgreSQL implementation
func (r *TodoRepository) GetAll(ctx context.Context, page, limit int) ([]*Todo, int64, error) {
    var todos []*Todo
    var total int64
    
    offset := (page - 1) * limit
    
    // Count total records
    if err := r.db.Model(&Todo{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Get paginated data
    if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&todos).Error; err != nil {
        return nil, 0, err
    }
    
    return todos, total, nil
}
```

## Use Case Implementation

```go
// TodoUseCase interface güncelleme
type TodoUseCase interface {
    CreateTodo(ctx context.Context, title, description string) (*Todo, error)
    GetTodo(ctx context.Context, id uint) (*Todo, error)
    GetAllTodos(ctx context.Context, page, limit int) (*PaginatedResponse, error)
    UpdateTodo(ctx context.Context, id uint, title, description string) (*Todo, error)
    DeleteTodo(ctx context.Context, id uint) error
    CompleteTodo(ctx context.Context, id uint) (*Todo, error)
    IncompleteTodo(ctx context.Context, id uint) (*Todo, error)
}

// Use case implementation
func (uc *TodoUseCase) GetAllTodos(ctx context.Context, page, limit int) (*PaginatedResponse, error) {
    todos, total, err := uc.repo.GetAll(ctx, page, limit)
    if err != nil {
        return nil, err
    }
    
    pagination := NewPagination(page, limit, total)
    return NewPaginatedResponse(todos, pagination), nil
}
```

## Handler Implementation

```go
// Handler implementation
func (h *TodoHandler) GetAllTodos(c *gin.Context) {
    // Parse query parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    // Validate parameters
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    if limit > 100 {
        limit = 100
    }
    
    // Call use case
    result, err := h.usecase.GetAllTodos(c.Request.Context(), page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, result)
}
```

## Query Parameters

### Supported Parameters
- `page`: Sayfa numarası (default: 1)
- `limit`: Sayfa başına kayıt sayısı (default: 10, max: 100)

### Example Requests
```
GET /api/v1/todos                    # Default: page=1, limit=10
GET /api/v1/todos?page=2            # Page 2, limit=10
GET /api/v1/todos?limit=20          # Page 1, limit=20
GET /api/v1/todos?page=3&limit=5    # Page 3, limit=5
```

## Performance Considerations

### Database Indexing
```sql
-- Pagination için index
CREATE INDEX idx_todos_created_at ON todos(created_at DESC);
CREATE INDEX idx_todos_completed_created_at ON todos(completed, created_at DESC);
```

### Query Optimization
```go
// Optimized query with proper ordering
func (r *TodoRepository) GetAll(ctx context.Context, page, limit int) ([]*Todo, int64, error) {
    var todos []*Todo
    var total int64
    
    offset := (page - 1) * limit
    
    // Use separate queries for count and data
    countQuery := r.db.Model(&Todo{})
    dataQuery := r.db.Model(&Todo{}).Order("created_at DESC").Offset(offset).Limit(limit)
    
    // Execute count query
    if err := countQuery.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Execute data query
    if err := dataQuery.Find(&todos).Error; err != nil {
        return nil, 0, err
    }
    
    return todos, total, nil
}
```

## Error Handling

```go
// Validation errors
if page < 1 {
    return nil, errors.New("page must be greater than 0")
}
if limit < 1 || limit > 100 {
    return nil, errors.New("limit must be between 1 and 100")
}

// Database errors
if err != nil {
    return nil, fmt.Errorf("failed to fetch todos: %w", err)
}
```

## Testing

```go
func TestTodoRepository_GetAll(t *testing.T) {
    // Setup
    db := setupTestDB()
    repo := NewTodoRepository(db)
    
    // Create test data
    createTestTodos(db, 25)
    
    // Test cases
    tests := []struct {
        name     string
        page     int
        limit    int
        expected int
    }{
        {"first page", 1, 10, 10},
        {"second page", 2, 10, 10},
        {"last page", 3, 10, 5},
        {"custom limit", 1, 5, 5},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            todos, total, err := repo.GetAll(context.Background(), tt.page, tt.limit)
            
            assert.NoError(t, err)
            assert.Len(t, todos, tt.expected)
            assert.Equal(t, int64(25), total)
        })
    }
}
```

Bu implementasyon, [Go pagination best practices](https://medium.easyread.co/how-to-do-pagination-in-postgres-with-golang-in-4-common-ways-12365b9fb528) ve [DEV.to pagination guide](https://dev.to/siddheshk02/how-to-paginate-api-responses-in-go-4cga) doğrultusunda hazırlanmıştır. 