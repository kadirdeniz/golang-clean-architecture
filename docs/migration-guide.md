# Migration Guide - Go Clean Architecture Todo API

## Overview

Bu dokümantasyon, Go Clean Architecture Todo API projesinde database migration'larının nasıl yönetileceğini açıklar. Proje, [golang-migrate](https://github.com/golang-migrate/migrate) tool'unu kullanarak migration'ları yönetir.

## Migration Dosya Yapısı

Migration dosyaları `migrations/` klasöründe bulunur ve şu formatta adlandırılır:

```
migrations/
├── 001_create_todos_table.up.sql
├── 001_create_todos_table.down.sql
├── 002_create_indexes_and_constraints.up.sql
└── 002_create_indexes_and_constraints.down.sql
```

### Dosya Naming Convention

- **Up Migration**: `{version}_{description}.up.sql`
- **Down Migration**: `{version}_{description}.down.sql`
- **Version**: 3 haneli sayı (001, 002, 003, ...)
- **Description**: Snake_case formatında açıklama

## Migration Tool Kurulumu

### Go ile Kurulum

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Docker ile Kullanım

```bash
docker run -v $(pwd)/migrations:/migrations migrate/migrate:latest -path /migrations -database "postgres://postgres:postgres@localhost:5432/todo_db_dev?sslmode=disable" up
```

## Migration Komutları

### Makefile Komutları

```bash
# Migration'ları çalıştır
make migrate-up

# Migration'ları geri al
make migrate-down

# Belirli bir versiyona zorla
make migrate-force version=1

# Mevcut migration versiyonunu göster
make migrate-version

# Migration durumunu göster
make migrate-status
```

### Script Komutları

```bash
# Migration'ları çalıştır
./scripts/migrate.sh up

# Migration'ları geri al
./scripts/migrate.sh down

# Belirli bir versiyona zorla
./scripts/migrate.sh force 1

# Mevcut migration versiyonunu göster
./scripts/migrate.sh version

# Migration durumunu göster
./scripts/migrate.sh status

# Yardım göster
./scripts/migrate.sh help
```

### Direkt migrate Komutları

```bash
# Migration'ları çalıştır
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/todo_db_dev?sslmode=disable" up

# Migration'ları geri al
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/todo_db_dev?sslmode=disable" down

# Belirli bir versiyona zorla
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/todo_db_dev?sslmode=disable" force 1

# Mevcut migration versiyonunu göster
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/todo_db_dev?sslmode=disable" version
```

## Environment Variables

Migration script'i aşağıdaki environment variable'ları destekler:

```bash
# Database connection string
DATABASE_URL="postgres://postgres:postgres@localhost:5432/todo_db_dev?sslmode=disable"

# Default değer: postgres://postgres:postgres@localhost:5432/todo_db_dev?sslmode=disable
```

## Docker Compose ile Migration

Docker Compose dosyasında migration service'i otomatik olarak çalışır:

```yaml
# Migration Service
migrate:
  image: migrate/migrate:latest
  container_name: todo_migrate_dev
  environment:
    DATABASE_URL: postgres://postgres:postgres@postgres:5432/todo_db_dev?sslmode=disable
  volumes:
    - ../../migrations:/migrations
  networks:
    - todo_network
  depends_on:
    postgres:
      condition: service_healthy
  command: ["-path", "/migrations", "-database", "postgres://postgres:postgres@postgres:5432/todo_db_dev?sslmode=disable", "up"]
  restart: "no"
```

### Docker Compose Komutları

```bash
# Tüm servisleri başlat (migration dahil)
docker-compose -f docker/docker-compose.dev.yml up -d

# Sadece database'i başlat
docker-compose -f docker/docker-compose.dev.yml up -d postgres

# Migration'ları manuel çalıştır
docker-compose -f docker/docker-compose.dev.yml run --rm migrate
```

## Yeni Migration Oluşturma

### 1. Migration Dosyalarını Oluştur

```bash
# Up migration
touch migrations/003_add_user_id_to_todos.up.sql

# Down migration
touch migrations/003_add_user_id_to_todos.down.sql
```

### 2. Up Migration İçeriği

```sql
-- Migration: Add user_id to todos table
-- Description: Add foreign key relationship to users table
-- Created: 2024-01-27

-- Add user_id column
ALTER TABLE todos ADD COLUMN user_id INTEGER;

-- Add foreign key constraint
ALTER TABLE todos ADD CONSTRAINT fk_todos_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add index for better performance
CREATE INDEX idx_todos_user_id ON todos(user_id);

-- Add comment
COMMENT ON COLUMN todos.user_id IS 'Foreign key reference to users table';
```

### 3. Down Migration İçeriği

```sql
-- Migration: Remove user_id from todos table
-- Description: Rollback for adding user_id column
-- Created: 2024-01-27

-- Drop index
DROP INDEX IF EXISTS idx_todos_user_id;

-- Drop foreign key constraint
ALTER TABLE todos DROP CONSTRAINT IF EXISTS fk_todos_user_id;

-- Drop column
ALTER TABLE todos DROP COLUMN IF EXISTS user_id;
```

## Migration Best Practices

### 1. Atomic Operations

Her migration dosyası atomic olmalıdır:

```sql
-- ✅ İyi örnek - Transaction içinde
BEGIN;
ALTER TABLE todos ADD COLUMN user_id INTEGER;
CREATE INDEX idx_todos_user_id ON todos(user_id);
COMMIT;

-- ❌ Kötü örnek - Transaction dışında
ALTER TABLE todos ADD COLUMN user_id INTEGER;
-- Eğer burada hata olursa, migration yarıda kalır
CREATE INDEX idx_todos_user_id ON todos(user_id);
```

### 2. Down Migration'ları Test Et

Her up migration'ından sonra down migration'ını test et:

```bash
# Up migration çalıştır
./scripts/migrate.sh up

# Down migration test et
./scripts/migrate.sh down

# Tekrar up migration çalıştır
./scripts/migrate.sh up
```

### 3. IF NOT EXISTS / IF EXISTS Kullan

```sql
-- ✅ İyi örnek
CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id);
DROP INDEX IF EXISTS idx_todos_user_id;

-- ❌ Kötü örnek
CREATE INDEX idx_todos_user_id ON todos(user_id);
DROP INDEX idx_todos_user_id;
```

### 4. Açıklayıcı Yorumlar

```sql
-- Migration: Add user_id to todos table
-- Description: Add foreign key relationship to users table for user-specific todos
-- Created: 2024-01-27
-- Author: Development Team

-- Add user_id column with proper constraints
ALTER TABLE todos ADD COLUMN user_id INTEGER NOT NULL;

-- Add foreign key constraint with cascade delete
ALTER TABLE todos ADD CONSTRAINT fk_todos_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add index for better query performance
CREATE INDEX idx_todos_user_id ON todos(user_id);

-- Add comments for documentation
COMMENT ON COLUMN todos.user_id IS 'Foreign key reference to users table';
COMMENT ON INDEX idx_todos_user_id IS 'Index for filtering todos by user';
```

### 5. Version Control

Migration dosyaları asla değiştirilmemeli:

```bash
# ✅ İyi - Yeni migration oluştur
touch migrations/004_fix_user_id_constraint.up.sql

# ❌ Kötü - Mevcut migration'ı değiştir
# migrations/003_add_user_id_to_todos.up.sql dosyasını düzenleme
```

## Troubleshooting

### Migration Hatası

```bash
# Migration durumunu kontrol et
./scripts/migrate.sh status

# Migration'ı zorla belirli bir versiyona
./scripts/migrate.sh force 1

# Database'i sıfırla ve migration'ları tekrar çalıştır
make db-reset
make migrate-up
```

### Database Bağlantı Hatası

```bash
# Database'in çalıştığını kontrol et
docker-compose -f docker/docker-compose.dev.yml ps

# Database loglarını kontrol et
docker-compose -f docker/docker-compose.dev.yml logs postgres

# Database'e bağlan
docker-compose -f docker/docker-compose.dev.yml exec postgres psql -U postgres -d todo_db_dev
```

### Migration Tool Hatası

```bash
# Migration tool'un kurulu olduğunu kontrol et
migrate --version

# Migration tool'u yeniden kur
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Production Migration

Production ortamında migration'ları çalıştırırken dikkat edilmesi gerekenler:

### 1. Backup Al

```bash
# Database backup al
pg_dump -h localhost -U postgres -d todo_db_prod > backup_$(date +%Y%m%d_%H%M%S).sql
```

### 2. Maintenance Window

```bash
# Uygulamayı durdur
docker-compose -f docker/docker-compose.prod.yml down

# Migration'ları çalıştır
DATABASE_URL="postgres://user:pass@host:port/db?sslmode=require" ./scripts/migrate.sh up

# Uygulamayı başlat
docker-compose -f docker/docker-compose.prod.yml up -d
```

### 3. Rollback Plan

```bash
# Migration'ları geri al
DATABASE_URL="postgres://user:pass@host:port/db?sslmode=require" ./scripts/migrate.sh down

# Eski versiyona dön
git checkout v1.0.0
```

## Migration Versiyonları

| Version | Description | Date |
|---------|-------------|------|
| 001 | Create todos table | 2024-01-07 |
| 002 | Create indexes and constraints | 2024-01-07 |

## References

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Compose Documentation](https://docs.docker.com/compose/) 