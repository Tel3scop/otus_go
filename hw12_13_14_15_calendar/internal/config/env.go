package config

import (
	"fmt"
	"log"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

type DataBaseType string

const (
	MemoryDatabaseType   DataBaseType = "memory"
	PostgresDatabaseType DataBaseType = "postgres"
)

// Provider предоставляет интерфейс для получения конфигурации.
type Provider interface {
	Config() *Config
}

// Config общий конфиг.
type Config struct {
	DebugMode   bool         `yaml:"debugMode" envDefault:"false"`
	Environment string       `yaml:"environment" envDefault:"local"`
	Database    DataBaseType `yaml:"database" envDefault:"memory"`
	Postgres    Postgres     `yaml:"postgres"`
	HTTP        HTTP         `yaml:"http"`
	Log         Log          `yaml:"log"`
}

// Config возвращаем сам конфиг.
func (c Config) Config() *Config {
	return &c
}

// Postgres конфиг подключения к БД.
type Postgres struct {
	Host               string `yaml:"host" envDefault:"localhost"`
	Port               string `yaml:"port" envDefault:"5432"`
	User               string `yaml:"user" envDefault:"root"`
	Password           string `yaml:"password" envDefault:"password"`
	DB                 string `yaml:"db" envDefault:"postgres"`
	SslMode            string `yaml:"sslMode" envDefault:"disable"`
	DSN                string
	MaxOpenConnections int `yaml:"maxOpenConnections" envDefault:"100"`
}

// HTTP конфиг подключения к grpc
type HTTP struct {
	Host    string `yaml:"host" envDefault:"localhost"`
	Port    string `env:"port" envDefault:"8080"`
	Address string
}

// Log конфиг для логов.
type Log struct {
	FileName   string `yaml:"fileName" envDefault:"logs/app.log"`
	Level      string `yaml:"level" envDefault:"info"`
	MaxSize    int    `yaml:"maxSize" envDefault:"5"`
	MaxBackups int    `yaml:"maxBackups" envDefault:"3"`
	MaxAge     int    `yaml:"maxAge" envDefault:"10"`
	Compress   bool   `yaml:"compress" envDefault:"false"`
	StdOut     bool   `yaml:"stdOut" envDefault:"false"`
}

// New создаем новый конфиг.
func New(configFile string) (*Config, error) {
	cfg := &Config{}
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to read config file: %s", err.Error())
	}

	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %s", err.Error())
	}
	cfg.HTTP.Address = net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port)
	buildDSN(&cfg.Postgres)

	return cfg, nil
}

func buildDSN(p *Postgres) {
	p.DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.DB, p.SslMode)
}
