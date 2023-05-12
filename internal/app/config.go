package app

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/v2"
)

type (
	Config struct {
		Logger   Logger   `koanf:"logger" validate:"required"`
		HTTP     HTTP     `koanf:"http" validate:"required"`
		Postgres Postgres `koanf:"postgres" validate:"required"`
		Redis    Redis    `koanf:"redis" validate:"required"`
	}

	Logger struct {
		Level *int8 `koanf:"level" validate:"required"`
	}

	HTTP struct {
		Host string `koanf:"host" validate:"required"`
		Port string `koanf:"port" validate:"required"`
	}

	Postgres struct {
		ConnString      string        `koanf:"connString" validate:"required"`
		MaxOpenConns    int           `koanf:"maxOpenConns" validate:"required"`
		ConnMaxLifetime time.Duration `koanf:"connMaxLifetime" validate:"required"`
		MaxIdleConns    int           `koanf:"maxIdleConns" validate:"required"`
		ConnMaxIdleTime time.Duration `koanf:"connMaxIdleTime" validate:"required"`
		AutoMigrate     bool          `koanf:"autoMigrate"`
		MigrationsPath  string        `koanf:"migrationsPath" validate:"required"`
	}

	Redis struct {
		Host     string `koanf:"host" validate:"required"`
		Port     string `koanf:"port" validate:"required"`
		Password string `koanf:"password"`
		DB       *int8  `koanf:"db" validate:"required"`
	}
)

func LoadConfig() (*Config, error) {
	k := koanf.New(".")

	defaultLogLevel := int8(-1)
	zero := int8(0)

	cfg := &Config{
		HTTP: HTTP{
			Host: "localhost",
			Port: "8000",
		},
		Logger: Logger{
			Level: &defaultLogLevel,
		},
		Postgres: Postgres{
			ConnString:      "postgresql://root:pass@127.0.0.1:5432/admin?sslmode=disable&application_name=admin-service",
			MaxOpenConns:    10,
			ConnMaxLifetime: 20,
			MaxIdleConns:    15,
			ConnMaxIdleTime: 30,
			AutoMigrate:     true,
			MigrationsPath:  "db/migration",
		},
		Redis: Redis{
			Host:     "127.0.0.1",
			Port:     "6379",
			Password: "",
			DB:       &zero,
		},
	}

	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	err := validator.New().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
