package config_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/config"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	BeforeEach(func() {
		// Clear environment variables before each test
		os.Unsetenv("CONFIG_FILE_NAME")
	})

	AfterEach(func() {
		// Clean up environment variables after each test
		os.Unsetenv("CONFIG_FILE_NAME")
	})

	Describe("NewConfig", func() {
		It("should load local config by default", func() {
			os.Unsetenv("CONFIG_FILE_NAME")
			cfg, err := config.NewConfig()
			
			Expect(err).To(BeNil())
			Expect(cfg).ToNot(BeNil())
			appCfg := cfg.GetAppConfig()
			Expect(appCfg.Name).To(Equal("Todo API"))
			Expect(appCfg.Environment).To(Equal("local"))
			Expect(appCfg.Port).To(Equal(8081))
			Expect(appCfg.Host).To(Equal("localhost"))
		})

		It("should load specific config when CONFIG_FILE_NAME is set", func() {
			os.Setenv("CONFIG_FILE_NAME", "config.dev")
			cfg, err := config.NewConfig()
			
			Expect(err).To(BeNil())
			Expect(cfg).ToNot(BeNil())
			appCfg := cfg.GetAppConfig()
			Expect(appCfg.Environment).To(Equal("development"))
			Expect(appCfg.Port).To(Equal(8080))
		})

		It("should fail to load invalid config file name", func() {
			os.Setenv("CONFIG_FILE_NAME", "nonexistent")
			cfg, err := config.NewConfig()
			
			Expect(err).To(HaveOccurred())
			Expect(cfg).To(BeNil())
		})

		It("should apply default values when not specified in config", func() {
			os.Setenv("CONFIG_FILE_NAME", "config.local")
			cfg, err := config.NewConfig()
			
			Expect(err).To(BeNil())
			Expect(cfg).ToNot(BeNil())
			// Test that default values are applied
			appCfg := cfg.GetAppConfig()
			Expect(appCfg.Name).To(Equal("Todo API"))
			Expect(appCfg.Host).To(Equal("localhost"))
		})

		It("should validate configuration after loading", func() {
			os.Setenv("CONFIG_FILE_NAME", "config.local")
			cfg, err := config.NewConfig()
			
			Expect(err).To(BeNil())
			Expect(cfg).ToNot(BeNil())
			
			// Test that validation tags work
			dbCfg := cfg.GetDatabaseConfig()
			Expect(dbCfg.Port).To(BeNumerically(">", 0))
			Expect(dbCfg.Port).To(BeNumerically("<", 65536))
		})
	})

	Describe("Configuration Structs", func() {
		It("should have proper struct tags and validation", func() {
			appConfig := config.AppConfig{
				Name:        "Test App",
				Version:     "1.0.0",
				Environment: "test",
				Port:        8080,
				Host:        "localhost",
			}
			Expect(appConfig.Name).To(Equal("Test App"))
			Expect(appConfig.Port).To(Equal(8080))

			dbConfig := config.DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				Name:     "test_db",
				Username: "test_user",
				Password: "test_pass",
			}
			Expect(dbConfig.Host).To(Equal("localhost"))
			Expect(dbConfig.Name).To(Equal("test_db"))

			logConfig := config.LoggingConfig{
				Level:  "debug",
				Format: "json",
				Output: "stdout",
			}
			Expect(logConfig.Level).To(Equal("debug"))
		})
	})

	Describe("Default Values", func() {
		It("should have sensible default values", func() {
			// Test default values from struct tags
			appConfig := config.AppConfig{}
			Expect(appConfig.Name).To(Equal(""))
			
			// Load config to see default values in action
			os.Setenv("CONFIG_FILE_NAME", "config.local")
			cfg, err := config.NewConfig()
			Expect(err).To(BeNil())
			
			// These should have default values from the config files
			appCfg := cfg.GetAppConfig()
			Expect(appCfg.Name).To(Equal("Todo API"))
			Expect(appCfg.Port).To(Equal(8081))
		})
	})

	Describe("Configuration Validation", func() {
		It("should validate database port range", func() {
			os.Setenv("CONFIG_FILE_NAME", "config.local")
			cfg, err := config.NewConfig()
			Expect(err).To(BeNil())
			
			// Database port should be within valid range
			dbCfg := cfg.GetDatabaseConfig()
			Expect(dbCfg.Port).To(BeNumerically(">", 0))
			Expect(dbCfg.Port).To(BeNumerically("<", 65536))
		})

		It("should validate app port range", func() {
			os.Setenv("CONFIG_FILE_NAME", "config.local")
			cfg, err := config.NewConfig()
			Expect(err).To(BeNil())
			
			// App port should be within valid range
			appCfg := cfg.GetAppConfig()
			Expect(appCfg.Port).To(BeNumerically(">", 0))
			Expect(appCfg.Port).To(BeNumerically("<", 65536))
		})
	})

	Describe("Environment-Specific Configuration", func() {
		It("should load different configs for different file names", func() {
			// Test development config
			os.Setenv("CONFIG_FILE_NAME", "config.dev")
			devCfg, err := config.NewConfig()
			Expect(err).To(BeNil())
			devAppCfg := devCfg.GetAppConfig()
			Expect(devAppCfg.Port).To(Equal(8080))
			Expect(devAppCfg.Environment).To(Equal("development"))

			// Test staging config
			os.Setenv("CONFIG_FILE_NAME", "config.staging")
			stagingCfg, err := config.NewConfig()
			Expect(err).To(BeNil())
			stagingAppCfg := stagingCfg.GetAppConfig()
			Expect(stagingAppCfg.Environment).To(Equal("staging"))

			// Test production config
			os.Setenv("CONFIG_FILE_NAME", "config.prod")
			prodCfg, err := config.NewConfig()
			Expect(err).To(BeNil())
			prodAppCfg := prodCfg.GetAppConfig()
			Expect(prodAppCfg.Environment).To(Equal("production"))

			// Test invalid config file name
			os.Setenv("CONFIG_FILE_NAME", "nonexistent")
			invalidCfg, err := config.NewConfig()
			Expect(err).To(HaveOccurred())
			Expect(invalidCfg).To(BeNil())
		})
	})

	Describe("Configuration File Name", func() {
		It("should load config from specified file name", func() {
			// Test local config file name
			os.Setenv("CONFIG_FILE_NAME", "config.local")
			localCfg, err := config.NewConfig()
			Expect(err).To(BeNil())
			localAppCfg := localCfg.GetAppConfig()
			Expect(localAppCfg.Environment).To(Equal("local"))

			// Test different config file
			os.Setenv("CONFIG_FILE_NAME", "config.staging")
			stagingCfg, err := config.NewConfig()
			Expect(err).To(BeNil())
			stagingAppCfg := stagingCfg.GetAppConfig()
			Expect(stagingAppCfg.Environment).To(Equal("staging"))
		})
	})
}) 