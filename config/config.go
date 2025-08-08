package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type LogConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	EnableFile bool   `yaml:"enableFile"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxSize"`
	MaxAge     int    `yaml:"maxAge"`
}

type APIConfig struct {
	Endpoint          string `yaml:"endpoint"`
	IAMEndpoint       string `yaml:"iam_endpoint"`
	Domain            string `yaml:"domain"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	RefreshTokenHours int    `yaml:"refresh_token_hours"`
	PageSize          int    `yaml:"page_size"`
}

type Config struct {
	Log LogConfig `yaml:"log"`
	API APIConfig `yaml:"api"`
}

var Cfg *Config

func LoadConfig(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	// Set defaults if needed
	if cfg.API.PageSize == 0 {
		cfg.API.PageSize = 100
	}
	if cfg.API.RefreshTokenHours == 0 {
		cfg.API.RefreshTokenHours = 3
	}

	Cfg = &cfg
}

func getDefaultConfig() Config {
	return Config{
		Log: LogConfig{
			Level:      "info",
			Format:     "text",
			EnableFile: false,
			Filename:   "./logs/ecs_exporter",
			MaxSize:    100, // 100MB
			MaxAge:     7,   // 7天
		},
		// 其他默认配置...
	}
}
