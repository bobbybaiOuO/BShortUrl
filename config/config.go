package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config .
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Redis RedisConfig `mapstructure:"redis"`
	Server ServerConfig `mapstructure:"server"`
	App AppConfig `mapstructure:"app"`
	ShortCode ShortCodeConfig `mapstructure:"shortcode"`
}

// LoadConfig .
func LoadConfig(filePath string) (*Config, error) {
	viper.SetConfigFile(filePath)

	viper.SetEnvPrefix("URL_SHORTENER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// DatabaseConfig .
type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Host string `mapstructure:"host"`
	Port int `mapstructure:"port"`
	User string `mapstructure:"user"`
	PassWord string `mapstructure:"password"`
	DBName string `mapstructure:"dbname"`
	SSLMode string `mapstructure:"ssl_mode"`
	// 最大空闲连接数
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// 最大开放连接数
	MaxOpenConns int `mapstructure:"max_open_conns"`
}

// DSN .
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", d.Driver, d.User,
	d.PassWord, d.Host, d.Port, d.DBName, d.SSLMode)
}

// RedisConfig .
type RedisConfig struct {
	Address string `mapstructure:"address"`
	PassWord string `mapstructure:"password"`
	DB int `mapstructure:"db"`
}

// ServerConfig .
type ServerConfig struct {
	Addr string `mapstructure:"addr"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	ReadTimeout time.Duration `mapstructure:"read_timeout"`
}

// AppConfig .
type AppConfig struct {
	BaseURL string `mapstructure:"base_url"`
	DefaultDuration time.Duration `mapstructure:"default_duration"`
	// 清理有效期
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

// ShortCodeConfig .
type ShortCodeConfig struct {
	Length int `mapstructure:"length"`
}