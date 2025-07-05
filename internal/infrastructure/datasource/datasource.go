package datasource

import (
	"context"
	"fmt"

	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/config"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type Database interface {
	Open(ctx context.Context) error
	Close() error
	HealthCheck(ctx context.Context) error
	DB() *gorm.DB
} 

func NewDatabase(cfg config.Config) Database {
	switch cfg.GetDatabaseConfig().Driver {
	case "postgres", "postgresql":
		return NewPostgresDatabase(cfg)
	default:
		panic(fmt.Sprintf("unsupported database driver: %s", cfg.GetDatabaseConfig().Driver))
	}
} 