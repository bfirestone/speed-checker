package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Testing  TestingConfig  `mapstructure:"testing"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"db"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`
}

// DSN constructs the PostgreSQL connection string from individual config fields
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Database, d.SSLMode)
}

type TestingConfig struct {
	SpeedTestInterval time.Duration `mapstructure:"speedtest_interval"`
	IperfTestInterval time.Duration `mapstructure:"iperf_interval"`
	IperfTestDuration int           `mapstructure:"iperf_duration"`
}

func Load() (*Config, error) {
	// Use experimental bind struct here to bind mapstructure tags to viper
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())
	v.SetEnvPrefix("SPEED_CHECKER")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Also bind to standard PostgreSQL environment variables
	v.BindEnv("database.host", "POSTGRES_HOST")
	v.BindEnv("database.port", "POSTGRES_PORT")
	v.BindEnv("database.db", "POSTGRES_DB")
	v.BindEnv("database.user", "POSTGRES_USER")
	v.BindEnv("database.password", "POSTGRES_PASSWORD")
	v.BindEnv("database.sslmode", "POSTGRES_SSLMODE")

	// Set defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.host", "localhost")
	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.db", "speedchecker")
	v.SetDefault("database.user", "speedchecker")
	v.SetDefault("database.password", "")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("testing.speedtest_interval", "15m")
	v.SetDefault("testing.iperf_interval", "10m")
	v.SetDefault("testing.iperf_duration", 10)

	// Try to read config file (optional)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/speed-checker/")
	v.AddConfigPath("$HOME/.speed-checker")

	// Read config file if it exists (don't error if it doesn't)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file found but has error
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, continue with env vars and defaults
		log.Println("No config file found, using environment variables and defaults")
	} else {
		log.Printf("Using config file: %s", v.ConfigFileUsed())
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return cfg, nil
}

// Legacy function for backward compatibility
func Default() *Config {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return cfg
}
