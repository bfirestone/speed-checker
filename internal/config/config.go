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
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
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

	// Set defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.host", "localhost")
	v.SetDefault("database.driver", "sqlite3")
	v.SetDefault("database.dsn", "./speedtest_results.db?_fk=1")
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
