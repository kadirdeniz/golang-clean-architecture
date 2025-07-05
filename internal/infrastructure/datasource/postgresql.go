package datasource

import (
	"context"
	"fmt"

	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/config"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct {
	config config.Config
	gorm   *gorm.DB
}

func NewPostgresDatabase(cfg config.Config) Database {
	return &postgresDatabase{config: cfg}
}

func (p *postgresDatabase) Open(ctx context.Context) error {
	dsn := p.config.GetDatabaseConfig().GetConnectionURL()

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to initialize GORM: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(p.config.GetDatabaseConfig().MaxOpenConns)
	sqlDB.SetMaxIdleConns(p.config.GetDatabaseConfig().MaxIdleConns)
	sqlDB.SetConnMaxLifetime(p.config.GetDatabaseConfig().ConnMaxLifetime)

	p.gorm = gormDB
	return nil
}

func (p *postgresDatabase) Close() error {
	if p.gorm == nil {
		return nil
	}
	
	sqlDB, err := p.gorm.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	return sqlDB.Close()
}

func (p *postgresDatabase) HealthCheck(ctx context.Context) error {
	if p.gorm == nil {
		return fmt.Errorf("database not initialized")
	}
	
	sqlDB, err := p.gorm.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	return sqlDB.PingContext(ctx)
}

func (p *postgresDatabase) DB() *gorm.DB {
	return p.gorm
}
