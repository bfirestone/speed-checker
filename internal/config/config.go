package config

import (
	"time"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Testing  TestingConfig  `json:"testing"`
}

type ServerConfig struct {
	Port string `json:"port" default:"8080"`
	Host string `json:"host" default:"localhost"`
}

type DatabaseConfig struct {
	Driver string `json:"driver" default:"sqlite3"`
	DSN    string `json:"dsn" default:"./speedtest_results.db"`
}

type TestingConfig struct {
	SpeedTestInterval time.Duration `json:"speedtest_interval" default:"15m"`
	IperfTestInterval time.Duration `json:"iperf_test_interval" default:"10m"`
	IperfTestDuration int           `json:"iperf_test_duration" default:"10"`
}

func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "localhost",
		},
		Database: DatabaseConfig{
			Driver: "sqlite3",
			DSN:    "./speedtest_results.db?_fk=1",
		},
		Testing: TestingConfig{
			SpeedTestInterval: 15 * time.Minute,
			IperfTestInterval: 10 * time.Minute,
			IperfTestDuration: 10,
		},
	}
}
