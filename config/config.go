package config

import (
	"fmt"
	"github.com/gookit/slog"
	"mnc-test/util/helper"
	"os"
)

type DbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	DBDriver string
}

type ApiConfig struct {
	Host string
	Port string
}

type Config struct {
	*DbConfig
	*ApiConfig
}

// init db
func (c *Config) ReadDbConfig() error {
	//load env
	if err := helper.Loadenv(); err != nil {
		return err
	}
	//assign struct dbconfig
	c.DbConfig = &DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		DBDriver: os.Getenv("DB_DRIVER"),
	}

	c.ApiConfig = &ApiConfig{
		Host: os.Getenv("API_HOST"),
		Port: os.Getenv("API_PORT"),
	}

	//if env is missing
	if c.DbConfig.Host == "" || c.DbConfig.Port == "" || c.DbConfig.DBName == "" ||
		c.DbConfig.Username == "" || c.DbConfig.Password == "" || c.DbConfig.DBDriver == "" {
		return fmt.Errorf("missing required environment variable ")
	}

	slog.Infof("Database connected to %v", c.DbConfig.Host)
	return nil
}

func NewDbConfig() (*Config, error) {
	cfg := &Config{}

	if err := cfg.ReadDbConfig(); err != nil {
		return nil, fmt.Errorf("Failed to read db config %v", err.Error())
	}

	return cfg, nil
}
