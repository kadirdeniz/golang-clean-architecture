# Validator Usage Guide

Bu dokümantasyon, Go Clean Architecture projesinde validator'ın dependency injection ile nasıl kullanılacağını açıklar.

## Genel Bakış

Validator, `go-playground/validator/v10` kütüphanesini kullanarak struct validation sağlar. Interface-based design ile dependency injection için uygun hale getirilmiştir.

## Temel Kullanım

### Validator Oluşturma

```go
package main

import (
    "github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/validator"
)

func main() {
    // Validator instance oluştur
    v := validator.NewValidator()
    // Kullan
    // ...
}
```

### Struct Validation

```go
type Todo struct {
    Title       string `validate:"required,min=1,max=100"`
    Description string `validate:"max=500"`
    Priority    int    `validate:"min=1,max=5"`
    DueDate     time.Time `validate:"required,gt=now"`
}

func createTodo(todo Todo) error {
    v := validator.NewValidator()
    // Basit validation
    if err := v.Validate(todo); err != nil {
        return err
    }
    // Detaylı validation (formatlanmış hatalar)
    if err := v.ValidateStruct(todo); err != nil {
        return err
    }
    return nil
}
```

### Field Validation

```go
func validateTitle(title string) error {
    v := validator.NewValidator()
    return v.ValidateField(title, "required,min=1,max=100")
}

func validatePriority(priority int) error {
    v := validator.NewValidator()
    return v.ValidateField(priority, "min=1,max=5")
}
```

## Dependency Injection ile Kullanım

### Interface Injection

```go
type TodoService struct {
    validator validator.Validator
}

func NewTodoService(v validator.Validator) *TodoService {
    return &TodoService{
        validator: v,
    }
}

func (s *TodoService) CreateTodo(todo Todo) error {
    if err := s.validator.ValidateStruct(todo); err != nil {
        return err
    }
    // Todo creation logic...
    return nil
}
```

### Wire ile Dependency Injection

```go
// wire.go
func InitializeTodoService() *TodoService {
    wire.Build(
        validator.NewValidator,
        NewTodoService,
    )
    return &TodoService{}
}
```

### Test ile Kullanım

```go
func TestTodoService(t *testing.T) {
    // Test validator oluştur
    v := validator.NewValidator()
    // Service'i inject et
    service := NewTodoService(v)
    // Test işlemleri...
    todo := Todo{
        Title:  "Test Todo",
        Priority: 3,
    }
    err := service.CreateTodo(todo)
    Expect(err).To(BeNil())
}
```

## Validation Tags

Validator, go-playground/validator'ın tüm tag'lerini destekler:

```go
type Todo struct {
    ID          string    `validate:"required,uuid"`
    Title       string    `validate:"required,min=1,max=100"`
    Description string    `validate:"max=500"`
    Status      string    `validate:"required,oneof=pending in_progress completed"`
    Priority    int       `validate:"min=1,max=5"`
    DueDate     time.Time `validate:"required,gt=now"`
    Tags        []string  `validate:"max=5,dive,min=1,max=20"`
}
```

## Error Handling

### Validation Error Types

```go
// ValidationError - Tek bir validation hatası
type ValidationError struct {
    Field   string `json:"field"`
    Tag     string `json:"tag"`
    Value   string `json:"value"`
    Message string `json:"message"`
}

// ValidationErrors - Birden fazla validation hatası
type ValidationErrors []ValidationError
```

### Error Response Format

```go
func handleValidationError(err error) map[string]interface{} {
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        return validationErrors.ToMap()
    }
    return map[string]interface{}{
        "error": err.Error(),
    }
}
```

## HTTP Middleware ile Kullanım

```go
func ValidationMiddleware(v validator.Validator) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Request body'yi parse et
        var requestBody Todo
        if err := c.BodyParser(&requestBody); err != nil {
            return c.Status(400).JSON(map[string]string{
                "error": "Invalid request body",
            })
        }
        // Validate et
        if err := v.ValidateStruct(requestBody); err != nil {
            return c.Status(400).JSON(map[string]interface{}{
                "error": "Validation failed",
                "details": err.Error(),
            })
        }
        return c.Next()
    }
}
```

## Best Practices

### 1. Interface Kullanımı

```go
// Doğru - Interface kullan
type TodoService struct {
    validator validator.Validator
}

// Yanlış - Concrete type kullan
type TodoService struct {
    validator *validator.validatorImpl
}
```

### 2. Error Handling

```go
func (s *TodoService) CreateTodo(todo Todo) error {
    if err := s.validator.ValidateStruct(todo); err != nil {
        // Log the error
        s.logger.Error("Todo validation failed", 
            logger.String("todo_title", todo.Title),
            logger.Error(err),
        )
        return fmt.Errorf("invalid todo data: %w", err)
    }
    
    return nil
}
```

### 3. Custom Validation Rules

```go
func NewValidator() validator.Validator {
    v := validator.New()
    
    // Custom validation rule ekle
    v.RegisterValidation("strong_password", validateStrongPassword)
    
    return &validatorImpl{validate: v}
}

func validateStrongPassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    
    // En az 8 karakter, büyük harf, küçük harf, rakam
    if len(password) < 8 {
        return false
    }
    
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    
    return hasUpper && hasLower && hasNumber
}
```

### 4. Performance Optimization

```go
// Validator'ı singleton olarak kullan
var globalValidator validator.Validator

func init() {
    globalValidator = validator.NewValidator()
}

func GetValidator() validator.Validator {
    return globalValidator
}
```

## Testing

### Unit Testing

```go
func TestTodoValidation(t *testing.T) {
    v := validator.NewValidator()
    
    tests := []struct {
        name    string
        todo    Todo
        wantErr bool
    }{
        {
            name: "valid todo",
            todo: Todo{
                Title:  "Test Todo",
                Priority: 3,
            },
            wantErr: false,
        },
        {
            name: "invalid title",
            todo: Todo{
                Title:  "Test",
                Priority: 3,
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := v.ValidateStruct(tt.todo)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Mock Testing

```go
type MockValidator struct {
    ValidateFunc     func(interface{}) error
    ValidateStructFunc func(interface{}) error
    ValidateFieldFunc  func(interface{}, string) error
}

func (m *MockValidator) Validate(i interface{}) error {
    if m.ValidateFunc != nil {
        return m.ValidateFunc(i)
    }
    return nil
}

func (m *MockValidator) ValidateStruct(i interface{}) error {
    if m.ValidateStructFunc != nil {
        return m.ValidateStructFunc(i)
    }
    return nil
}

func (m *MockValidator) ValidateField(field interface{}, tag string) error {
    if m.ValidateFieldFunc != nil {
        return m.ValidateFieldFunc(field, tag)
    }
    return nil
}

// Test usage
func TestTodoServiceWithMock(t *testing.T) {
    mockValidator := &MockValidator{
        ValidateStructFunc: func(i interface{}) error {
            return fmt.Errorf("validation failed")
        },
    }
    
    service := NewTodoService(mockValidator)
    
    todo := Todo{Title: "Test"}
    err := service.CreateTodo(todo)
    
    Expect(err).ToNot(BeNil())
    Expect(err.Error()).To(ContainSubstring("validation failed"))
}
```

## Integration with Other Components

### Logger Integration

```go
func (s *TodoService) CreateTodo(todo Todo) error {
    if err := s.validator.ValidateStruct(todo); err != nil {
        s.logger.Error("Todo validation failed",
            logger.String("todo_title", todo.Title),
            logger.Error(err),
        )
        return err
    }
    
    s.logger.Info("Todo validation successful",
        logger.String("todo_title", todo.Title),
    )
    
    return nil
}
```

### Config Integration

```go
type ValidatorConfig struct {
    EnableCustomRules bool
    StrictMode        bool
}

func NewValidatorWithConfig(cfg ValidatorConfig) validator.Validator {
    v := validator.New()
    
    if cfg.EnableCustomRules {
        v.RegisterValidation("strong_password", validateStrongPassword)
    }
    
    return &validatorImpl{validate: v}
}
```

Bu yaklaşım, validator'ı Clean Architecture prensiplerine uygun şekilde kullanmanızı ve dependency injection ile entegre etmenizi sağlar. 