# Logger Usage Guide

Bu dokümantasyon, Go Clean Architecture projesinde configurable logger'ın nasıl kullanılacağını açıklar.

## Genel Bakış

Logger, `config.Config` struct'ını parametre olarak alarak configurable bir yapıda tasarlanmıştır. Bu sayede farklı ortamlar için farklı logging ayarları kullanılabilir.

## Temel Kullanım

### Logger Oluşturma

```go
package main

import (
    "github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/config"
    "github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/logger"
)

func main() {
    // Config yükle
    cfg, err := config.NewConfig()
    if err != nil {
        panic(err)
    }

    // Logger oluştur
    log, err := logger.NewLogger(cfg)
    if err != nil {
        panic(err)
    }
    defer log.Sync()

    // Logger kullan
    log.Info("Application started")
}
```

### Logging Seviyeleri

```go
// Debug seviyesi (sadece debug modunda görünür)
log.Debug("Debug message", logger.String("key", "value"))

// Info seviyesi
log.Info("Info message", logger.String("key", "value"))

// Warning seviyesi
log.Warn("Warning message", logger.String("key", "value"))

// Error seviyesi
log.Error("Error message", logger.String("key", "value"))

// Fatal seviyesi (uygulamayı sonlandırır)
log.Fatal("Fatal message", logger.String("key", "value"))
```

## Structured Logging

Logger, structured logging için çeşitli field tipleri destekler:

```go
log.Info("Todo action",
    logger.String("todo_id", "123"),
    logger.String("action", "create"),
    logger.String("ip", "192.168.1.1"),
    logger.Int("status_code", 200),
    logger.Float64("duration_ms", 45.67),
    logger.Bool("success", true),
    logger.Time("timestamp", time.Now()),
    logger.Duration("duration", time.Millisecond*100),
    logger.Any("metadata", map[string]interface{}{
        "browser": "Chrome",
        "version": "91.0",
    }),
)
```

## Error Logging

Hataları loglarken `logger.Error` field'ını kullanın:

```go
func processTodo(todoID string) error {
    if err := validateTodo(todoID); err != nil {
        log.Error("Failed to validate todo",
            logger.String("todo_id", todoID),
            logger.Error(err),
        )
        return err
    }
    return nil
}
```

## Logger With Fields

Mevcut logger'a ek alanlar eklemek için `With` metodunu kullanın:

```go
// Component-specific logger
componentLogger := log.With(
    logger.String("component", "todo_service"),
    logger.String("version", "1.0.0"),
)

// Request-specific logger
requestLogger := componentLogger.With(
    logger.String("request_id", "req-123"),
    logger.String("todo_id", "todo-456"),
)

requestLogger.Info("Processing request")
```

## Convenience Functions

Logger, yaygın kullanım senaryoları için convenience fonksiyonlar sağlar:

```go
// HTTP request logging
httpLogger := logger.WithRequestID(log, "req-123")
httpLogger = logger.WithMethod(httpLogger, "POST")
httpLogger = logger.WithPath(httpLogger, "/api/v1/todos")
httpLogger = logger.WithStatusCode(httpLogger, 201)
httpLogger = logger.WithIP(httpLogger, "192.168.1.1")

httpLogger.Info("Request processed successfully")

// Database operation logging
dbLogger := logger.WithComponent(log, "database")
dbLogger = logger.WithOperation(dbLogger, "create_todo")
dbLogger = logger.WithDuration(dbLogger, 45.67)

dbLogger.Info("Database operation completed")
```

## Configuration

Logger, config dosyasından aşağıdaki ayarları alır:

```yaml
logging:
  level: "info"      # debug, info, warn, error, fatal
  format: "json"     # json, console
  output: "stdout"   # stdout, stderr, /path/to/file.log
```

### Environment-Specific Configuration

Farklı ortamlar için farklı logging ayarları:

```yaml
# config.dev.yaml
logging:
  level: "debug"
  format: "console"
  output: "stdout"

# config.prod.yaml
logging:
  level: "info"
  format: "json"
  output: "/var/log/app.log"

# config.test.yaml
logging:
  level: "error"
  format: "json"
  output: "stdout"
```

## Best Practices

### 1. Contextual Information

Her log mesajına yeterli context bilgisi ekleyin:

```go
// İyi örnek
log.Error("Failed to create todo",
    logger.String("todo_id", todoID),
    logger.String("title", title),
    logger.Error(err),
    logger.String("operation", "create_todo"),
)

// Kötü örnek
log.Error("Failed to create todo")
```

### 2. Sensitive Data

Hassas verileri loglamayın:

```go
// Kötü örnek
log.Info("Todo created", logger.String("secret", secretValue))
```

## Testing

Logger'ı test ederken:

```go
func TestTodoService(t *testing.T) {
    // Test config kullan
    os.Setenv("ENVIRONMENT", "test")
    cfg, _ := config.NewConfig()
    log, _ := logger.NewLogger(cfg)
    
    service := NewTodoService(log)
    
    // Test işlemleri...
}
```

## Monitoring ve Alerting

Production ortamında logları izlemek için:

1. **Log Aggregation**: ELK Stack, Splunk, veya benzeri araçlar kullanın
2. **Error Alerting**: Error seviyesindeki loglar için alerting kurun
3. **Performance Monitoring**: Duration field'larını kullanarak performans metrikleri toplayın
4. **Business Metrics**: İş metriklerini loglayın (kullanıcı kayıtları, işlemler, vb.)

Bu yaklaşım, uygulamanızın sağlığını ve performansını etkili bir şekilde izlemenizi sağlar. 