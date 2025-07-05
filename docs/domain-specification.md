# Domain Specification - Todo API

## Domain Overview

Basit bir Todo yönetim sistemi. Kullanıcılar todo'ları oluşturabilir, görüntüleyebilir, güncelleyebilir, silebilir ve tamamlandı olarak işaretleyebilir.

## Core Entity

### Todo
```go
type Todo struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

## Use Cases

### 1. Create Todo
- **Input**: Title, Description
- **Output**: Created Todo
- **Business Rules**: Title zorunlu

### 2. Get Todo
- **Input**: Todo ID
- **Output**: Todo or Error
- **Business Rules**: ID geçerli olmalı

### 3. List Todos
- **Input**: Page, Limit (optional)
- **Output**: Paginated Todo listesi
- **Business Rules**: 
  - Default page: 1, Default limit: 10
  - Maximum limit: 100
  - Pagination metadata döner

### 4. Update Todo
- **Input**: Todo ID, Title, Description
- **Output**: Updated Todo
- **Business Rules**: ID geçerli olmalı

### 5. Delete Todo
- **Input**: Todo ID
- **Output**: Success/Error
- **Business Rules**: ID geçerli olmalı

### 6. Complete Todo
- **Input**: Todo ID
- **Output**: Updated Todo
- **Business Rules**: ID geçerli olmalı

### 7. Incomplete Todo
- **Input**: Todo ID
- **Output**: Updated Todo
- **Business Rules**: ID geçerli olmalı

## API Contracts

### Create Todo
```
POST /api/v1/todos
Content-Type: application/json

Request:
{
    "title": "Learn Go",
    "description": "Study Go programming language"
}

Response: 201 Created
{
    "id": 1,
    "title": "Learn Go",
    "description": "Study Go programming language",
    "completed": false,
    "created_at": "2025-01-27T10:00:00Z",
    "updated_at": "2025-01-27T10:00:00Z"
}
```

### Get Todo
```
GET /api/v1/todos/{id}

Response: 200 OK
{
    "id": 1,
    "title": "Learn Go",
    "description": "Study Go programming language",
    "completed": false,
    "created_at": "2025-01-27T10:00:00Z",
    "updated_at": "2025-01-27T10:00:00Z"
}
```

### List Todos
```
GET /api/v1/todos?page=1&limit=10

Response: 200 OK
{
    "data": [
        {
            "id": 1,
            "title": "Learn Go",
            "description": "Study Go programming language",
            "completed": false,
            "created_at": "2025-01-27T10:00:00Z",
            "updated_at": "2025-01-27T10:00:00Z"
        }
    ],
    "pagination": {
        "page": 1,
        "limit": 10,
        "total": 25,
        "total_pages": 3,
        "has_next": true,
        "has_prev": false,
        "next_page": 2,
        "prev_page": null
    }
}
```

### Update Todo
```
PUT /api/v1/todos/{id}
Content-Type: application/json

Request:
{
    "title": "Learn Go Advanced",
    "description": "Study advanced Go concepts"
}

Response: 200 OK
{
    "id": 1,
    "title": "Learn Go Advanced",
    "description": "Study advanced Go concepts",
    "completed": false,
    "created_at": "2025-01-27T10:00:00Z",
    "updated_at": "2025-01-27T10:05:00Z"
}
```

### Delete Todo
```
DELETE /api/v1/todos/{id}

Response: 204 No Content
```

### Complete Todo
```
PATCH /api/v1/todos/{id}/complete

Response: 200 OK
{
    "id": 1,
    "title": "Learn Go Advanced",
    "description": "Study advanced Go concepts",
    "completed": true,
    "created_at": "2025-01-27T10:00:00Z",
    "updated_at": "2025-01-27T10:10:00Z"
}
```

### Incomplete Todo
```
PATCH /api/v1/todos/{id}/incomplete

Response: 200 OK
{
    "id": 1,
    "title": "Learn Go Advanced",
    "description": "Study advanced Go concepts",
    "completed": false,
    "created_at": "2025-01-27T10:00:00Z",
    "updated_at": "2025-01-27T10:15:00Z"
}
```

## Error Responses

### Not Found
```json
{
    "error": "Todo not found",
    "code": "TODO_NOT_FOUND"
}
```

### Validation Error
```json
{
    "error": "Title is required",
    "code": "VALIDATION_ERROR"
}
```

### Internal Error
```json
{
    "error": "Internal server error",
    "code": "INTERNAL_ERROR"
}
```

## Business Rules

1. **Title Required**: Todo oluştururken title zorunlu
2. **ID Validation**: Tüm ID tabanlı işlemlerde ID geçerli olmalı
3. **Timestamps**: CreatedAt ve UpdatedAt otomatik set edilmeli
4. **Completed State**: Complete/Incomplete işlemleri boolean değeri değiştirir
5. **Pagination**: 
   - Default page: 1, Default limit: 10
   - Maximum limit: 100
   - Page ve limit pozitif sayı olmalı
   - Pagination metadata her zaman döner 