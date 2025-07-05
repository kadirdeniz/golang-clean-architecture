package datasource_test

import (
	"time"

	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/config"
	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/datasource"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database Factory", func() {
	var testConfig *mockConfig

	BeforeEach(func() {
		// Create a mock config that returns PostgreSQL driver
		testConfig = &mockConfig{
			driver: "postgres",
			host:   "localhost",
			port:   5432,
			name:   "test_db",
			username: "test_user",
			password: "test_pass",
			maxOpenConns: 10,
			maxIdleConns: 5,
			connMaxLifetime: time.Hour,
		}
	})

	Describe("NewDatabase", func() {
		Context("when driver is postgres", func() {
			It("should create PostgreSQL database instance", func() {
				db := datasource.NewDatabase(testConfig)
				Expect(db).NotTo(BeNil())
			})

			It("should return Database interface", func() {
				db := datasource.NewDatabase(testConfig)
				_, ok := db.(datasource.Database)
				Expect(ok).To(BeTrue())
			})
		})

		Context("when driver is postgresql", func() {
			BeforeEach(func() {
				testConfig.driver = "postgresql"
			})

			It("should create PostgreSQL database instance", func() {
				db := datasource.NewDatabase(testConfig)
				Expect(db).NotTo(BeNil())
			})
		})

		Context("when driver is unsupported", func() {
			BeforeEach(func() {
				testConfig.driver = "mysql"
			})

			It("should panic", func() {
				Expect(func() {
					datasource.NewDatabase(testConfig)
				}).To(Panic())
			})
		})

		Context("when driver is empty", func() {
			BeforeEach(func() {
				testConfig.driver = ""
			})

			It("should panic", func() {
				Expect(func() {
					datasource.NewDatabase(testConfig)
				}).To(Panic())
			})
		})
	})
})

// Mock implementations for testing
type mockConfig struct {
	driver          string
	host            string
	port            int
	name            string
	username        string
	password        string
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
}

func (m *mockConfig) GetAppConfig() config.AppConfig {
	return config.AppConfig{
		Name:        "Test App",
		Version:     "1.0.0",
		Environment: "test",
		Port:        8080,
		Host:        "localhost",
	}
}

func (m *mockConfig) GetDatabaseConfig() config.DatabaseConfig {
	return config.DatabaseConfig{
		Driver:          m.driver,
		Host:            m.host,
		Port:            m.port,
		Name:            m.name,
		Username:        m.username,
		Password:        m.password,
		SSLMode:         "disable",
		MaxOpenConns:    m.maxOpenConns,
		MaxIdleConns:    m.maxIdleConns,
		ConnMaxLifetime: m.connMaxLifetime,
	}
}

func (m *mockConfig) GetLoggingConfig() config.LoggingConfig {
	return config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}
}

func (m *mockConfig) GetServerConfig() config.ServerConfig {
	return config.ServerConfig{
		ReadTimeout:             30 * time.Second,
		WriteTimeout:            30 * time.Second,
		IdleTimeout:             60 * time.Second,
		GracefulShutdownTimeout: 10 * time.Second,
	}
}

func (m *mockConfig) GetCORSConfig() config.CORSConfig {
	return config.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}
} 