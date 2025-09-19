package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Storage  string `yaml:"storage"`
	Postgres struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Port     int    `yaml:"port"`
	} `yaml:"postgres"`
}

func LoadConfig() (*Config, error) {
	path := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()
	f, err := os.ReadFile(*path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, err
	}
	if env := os.Getenv("DB_HOST"); env != "" {
		cfg.Postgres.Host = env
	}
	return &cfg, nil
}
