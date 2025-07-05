package logger_test

import (
	"errors"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/config"
	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/logger"
)

func TestLogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logger Suite")
}

var _ = Describe("Logger", func() {
	var testConfig *mockConfig

	BeforeEach(func() {
		testConfig = &mockConfig{
			level:  "info",
			format: "json",
			output: "stdout",
		}
	})

	Describe("NewLogger", func() {
		Context("with valid configuration", func() {
			It("should create logger instance", func() {
				log, err := logger.NewLogger(testConfig)
				Expect(err).To(BeNil())
				Expect(log).NotTo(BeNil())
			})

			It("should implement Logger interface", func() {
				log, err := logger.NewLogger(testConfig)
				Expect(err).To(BeNil())
				_, ok := log.(logger.Logger)
				Expect(ok).To(BeTrue())
			})
		})

		Context("with different log levels", func() {
			It("should create logger with debug level", func() {
				testConfig.level = "debug"
				log, err := logger.NewLogger(testConfig)
				Expect(err).To(BeNil())
				Expect(log).NotTo(BeNil())
			})

			It("should create logger with info level", func() {
				testConfig.level = "info"
				log, err := logger.NewLogger(testConfig)
				Expect(err).To(BeNil())
				Expect(log).NotTo(BeNil())
			})

			It("should create logger with error level", func() {
				testConfig.level = "error"
				log, err := logger.NewLogger(testConfig)
				Expect(err).To(BeNil())
				Expect(log).NotTo(BeNil())
			})
		})

		Context("with different formats", func() {
			It("should create logger with json format", func() {
				testConfig.format = "json"
				log, err := logger.NewLogger(testConfig)
				Expect(err).To(BeNil())
				Expect(log).NotTo(BeNil())
			})

			It("should create logger with console format", func() {
				testConfig.format = "console"
				log, err := logger.NewLogger(testConfig)
				Expect(err).To(BeNil())
				Expect(log).NotTo(BeNil())
			})
		})
	})

	Describe("Logger Methods", func() {
		var log logger.Logger

		BeforeEach(func() {
			var err error
			log, err = logger.NewLogger(testConfig)
			Expect(err).To(BeNil())
		})

		Context("basic logging methods", func() {
			It("should have Debug method", func() {
				Expect(func() {
					log.Debug("Debug message")
				}).NotTo(Panic())
			})

			It("should have Info method", func() {
				Expect(func() {
					log.Info("Info message")
				}).NotTo(Panic())
			})

			It("should have Warn method", func() {
				Expect(func() {
					log.Warn("Warning message")
				}).NotTo(Panic())
			})

			It("should have Error method", func() {
				Expect(func() {
					log.Error("Error message")
				}).NotTo(Panic())
			})

			It("should have Sync method", func() {
				Expect(func() {
					log.Sync()
				}).NotTo(Panic())
			})
		})

		Context("structured logging with fields", func() {
			It("should support string fields", func() {
				Expect(func() {
					log.Info("Test message", logger.String("key", "value"))
				}).NotTo(Panic())
			})

			It("should support int fields", func() {
				Expect(func() {
					log.Info("Test message", logger.Int("count", 42))
				}).NotTo(Panic())
			})

			It("should support error fields", func() {
				testErr := errors.New("test error")
				Expect(func() {
					log.Error("Error message", logger.Error(testErr))
				}).NotTo(Panic())
			})

			It("should support multiple fields", func() {
				Expect(func() {
					log.Info("Test message", 
						logger.String("component", "test"),
						logger.Int("count", 42),
						logger.Bool("success", true),
					)
				}).NotTo(Panic())
			})
		})

		Context("logger with additional fields", func() {
			It("should create logger with fields", func() {
				logWithFields := log.With(logger.String("component", "test"))
				Expect(logWithFields).NotTo(BeNil())
				Expect(logWithFields).NotTo(Equal(log))
			})

			It("should chain With calls", func() {
				logWithFields := log.With(logger.String("component", "test"))
				logWithMoreFields := logWithFields.With(logger.String("operation", "create"))
				Expect(logWithMoreFields).NotTo(BeNil())
			})
		})
	})

	Describe("Convenience Functions", func() {
		var log logger.Logger

		BeforeEach(func() {
			var err error
			log, err = logger.NewLogger(testConfig)
			Expect(err).To(BeNil())
		})

		It("should create logger with request ID", func() {
			logWithRequestID := logger.WithRequestID(log, "req-123")
			Expect(logWithRequestID).NotTo(BeNil())
		})

		It("should create logger with operation", func() {
			logWithOperation := logger.WithOperation(log, "create_todo")
			Expect(logWithOperation).NotTo(BeNil())
		})

		It("should create logger with component", func() {
			logWithComponent := logger.WithComponent(log, "todo_handler")
			Expect(logWithComponent).NotTo(BeNil())
		})

		It("should create logger with error", func() {
			testErr := errors.New("test error")
			logWithError := logger.WithError(log, testErr)
			Expect(logWithError).NotTo(BeNil())
		})
	})
})

// Mock implementations for testing
type mockConfig struct {
	level  string
	format string
	output string
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
		Driver:          "postgres",
		Host:            "localhost",
		Port:            5432,
		Name:            "test_db",
		Username:        "test_user",
		Password:        "test_pass",
		SSLMode:         "disable",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	}
}

func (m *mockConfig) GetLoggingConfig() config.LoggingConfig {
	return config.LoggingConfig{
		Level:  m.level,
		Format: m.format,
		Output: m.output,
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